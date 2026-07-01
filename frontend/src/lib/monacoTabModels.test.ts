import { describe, expect, it } from 'vitest'
import { MonacoTabModelStore, featureTabUri } from './monacoTabModels'

function mockMonaco() {
  const models = new Map<string, { value: string; disposed: boolean; ref: ReturnType<typeof makeModel> | null }>()

  function makeModel(key: string) {
    return {
      getValue: () => models.get(key)!.value,
      setValue: (v: string) => {
        models.get(key)!.value = v
      },
      isDisposed: () => models.get(key)!.disposed,
      dispose: () => {
        models.get(key)!.disposed = true
      },
    }
  }

  return {
    Uri: {
      parse: (uri: string) => ({ toString: () => uri }),
    },
    editor: {
      getModel: (uri: { toString: () => string }) => {
        const entry = models.get(uri.toString())
        if (!entry || entry.disposed) return null
        if (!entry.ref) entry.ref = makeModel(uri.toString())
        return entry.ref
      },
      createModel: (value: string, _lang: string, uri: { toString: () => string }) => {
        const key = uri.toString()
        const ref = makeModel(key)
        models.set(key, { value, disposed: false, ref })
        return ref
      },
    },
  }
}

describe('monacoTabModels', () => {
  it('featureTabUri encodes path consistently', () => {
    const monaco = mockMonaco()
    const a = featureTabUri(monaco as never, 'C:\\proj\\a.feature').toString()
    const b = featureTabUri(monaco as never, 'C:/proj/a.feature').toString()
    expect(a).toBe(b)
    expect(a).toContain('inmemory://scenaria/feature/')
  })

  it('store reuses model for same path', () => {
    const monaco = mockMonaco()
    const store = new MonacoTabModelStore()
    const m1 = store.getOrCreate(monaco as never, '/p/a.feature', 'hello')
    const m2 = store.getOrCreate(monaco as never, '/p/a.feature', 'changed')
    expect(m1).toBe(m2)
    expect(m1.getValue()).toBe('hello')
    expect(store.trackedPaths()).toHaveLength(1)
  })

  it('release disposes model and allows recreate', () => {
    const monaco = mockMonaco()
    const store = new MonacoTabModelStore()
    const m1 = store.getOrCreate(monaco as never, '/p/a.feature', 'v1')
    store.release(monaco as never, '/p/a.feature')
    expect(m1.isDisposed()).toBe(true)
    const m2 = store.getOrCreate(monaco as never, '/p/a.feature', 'v2')
    expect(m2).not.toBe(m1)
    expect(m2.getValue()).toBe('v2')
  })

  it('releaseExcept keeps only requested paths', () => {
    const monaco = mockMonaco()
    const store = new MonacoTabModelStore()
    store.getOrCreate(monaco as never, '/a.feature', 'a')
    store.getOrCreate(monaco as never, '/b.feature', 'b')
    store.releaseExcept(monaco as never, ['/b.feature'])
    expect(store.trackedPaths()).toEqual(['/b.feature'])
  })
})
