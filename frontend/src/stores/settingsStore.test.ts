import { describe, expect, it } from 'vitest'
import { gui } from '../../wailsjs/go/models'
import { createSettingsStore, defaultAppSettingsState } from './settingsStore'

describe('settingsStore', () => {
  it('applies persisted DTO fields', () => {
    const store = createSettingsStore()
    store.applyFromDTO(
      gui.AppSettingsDTO.createFrom({
        browser: 'firefox',
        headless: true,
        parallelWorkers: 3,
        slowMo: 120,
        maxLoopIterations: 50,
        navWaitUntil: 'load',
        selectorClickStrategies: ['id'],
        selectorInputStrategies: ['name'],
      }),
    )

    let snapshot = defaultAppSettingsState
    const unsub = store.subscribe((s) => {
      snapshot = s
    })
    unsub()

    expect(snapshot.browser).toBe('firefox')
    expect(snapshot.headless).toBe(true)
    expect(snapshot.parallelWorkers).toBe(3)
    expect(snapshot.slowMo).toBe(120)
    expect(snapshot.maxLoopIterations).toBe(50)
    expect(snapshot.navWaitUntil).toBe('load')
    expect(snapshot.selectorClickStrategies).toEqual(['id'])
    expect(snapshot.selectorInputStrategies).toEqual(['name'])
  })

  it('patches editor settings independently', () => {
    const store = createSettingsStore()
    store.setEditor({ ...defaultAppSettingsState.editor, fontSize: 16, theme: 'scenaria-light' })
    let fontSize = 0
    const unsub = store.subscribe((s) => {
      fontSize = s.editor.fontSize
    })
    unsub()
    expect(fontSize).toBe(16)
  })

  it('exports dto fields from snapshot', () => {
    const store = createSettingsStore()
    store.patch({ browser: 'webkit', parallelWorkers: 2 })
    const fields = store.dtoFields(store.snapshot())
    expect(fields.browser).toBe('webkit')
    expect(fields.parallelWorkers).toBe(2)
  })
})
