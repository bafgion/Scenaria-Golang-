import { describe, expect, it } from 'vitest'
import {
  MAX_RETAINED_CLEAN_TAB_BODIES,
  pathsToRetainModels,
  tabEditorText,
  tabNeedsDiskReload,
  trimRetainedTabBodies,
  type TabBody,
} from './tabMemory'

function tab(path: string, content: string, extra: Partial<TabBody> = {}): TabBody {
  return { path, content, dirty: false, ...extra }
}

describe('tabMemory', () => {
  it('tabEditorText prefers draft for dirty tabs', () => {
    const t = tab('/a.feature', 'saved', { dirty: true, draft: 'edited' })
    expect(tabEditorText(t)).toBe('edited')
  })

  it('tabNeedsDiskReload when body unloaded', () => {
    expect(tabNeedsDiskReload(tab('/a.feature', '', { unloaded: true }))).toBe(true)
    expect(tabNeedsDiskReload(tab('/a.feature', 'x', { dirty: true, unloaded: true }))).toBe(false)
  })

  it('trimRetainedTabBodies unloads oldest clean inactive tabs', () => {
    const tabs: TabBody[] = [
      tab('/1.feature', 'one'),
      tab('/2.feature', 'two'),
      tab('/3.feature', 'three'),
      tab('/active.feature', 'active'),
    ]
    const trimmed = trimRetainedTabBodies(tabs, '/active.feature', 1)
    const unloaded = trimmed.filter((t) => t.unloaded)
    expect(unloaded.length).toBe(2)
    expect(trimmed.find((t) => t.path === '/active.feature')?.unloaded).toBeFalsy()
  })

  it('trimRetainedTabBodies keeps dirty tabs', () => {
    const tabs: TabBody[] = [
      tab('/1.feature', 'one', { dirty: true, draft: 'd1' }),
      tab('/2.feature', 'two'),
      tab('/3.feature', 'three'),
    ]
    const trimmed = trimRetainedTabBodies(tabs, '/x.feature', 0)
    expect(trimmed.find((t) => t.path === '/1.feature')?.unloaded).toBeFalsy()
    expect(trimmed.find((t) => t.path === '/2.feature')?.unloaded).toBe(true)
  })

  it('pathsToRetainModels excludes unloaded clean tabs', () => {
    const tabs: TabBody[] = [
      tab('/a.feature', '', { unloaded: true }),
      tab('/b.feature', 'body', { dirty: true, draft: 'd' }),
      tab('/c.feature', 'keep'),
    ]
    expect(pathsToRetainModels(tabs, '/active.feature')).toEqual(
      expect.arrayContaining(['/b.feature', '/c.feature', '/active.feature']),
    )
    expect(pathsToRetainModels(tabs, '/active.feature')).not.toContain('/a.feature')
  })

  it('default retain limit is reasonable', () => {
    expect(MAX_RETAINED_CLEAN_TAB_BODIES).toBeGreaterThanOrEqual(4)
  })
})
