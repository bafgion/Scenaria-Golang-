import { writable } from 'svelte/store'
import type { gui } from '../../wailsjs/go/models'
import type { Locale } from '../lib/i18n'
import {
  DEFAULT_EDITOR_SETTINGS,
  editorSettingsFromDTO,
  editorSettingsToDTO,
  type EditorSettings,
} from '../lib/editorOptions'

export type AppSettingsState = {
  browser: string
  headless: boolean
  parallelWorkers: number
  slowMo: number
  scrollBeforeClick: boolean
  disableRecordUrlWait: boolean
  hoverRecordMinMs: number
  maxLoopIterations: number
  checkUpdatesOnStartup: boolean
  selectorClickStrategies: string[]
  selectorInputStrategies: string[]
  libraryHeuristicsMui: boolean
  libraryHeuristicsAnt: boolean
  navWaitUntil: string
  editor: EditorSettings
  startUrl: string
  uiLocale: Locale
}

export const defaultAppSettingsState: AppSettingsState = {
  browser: 'chromium',
  headless: false,
  parallelWorkers: 1,
  slowMo: 0,
  scrollBeforeClick: false,
  disableRecordUrlWait: false,
  hoverRecordMinMs: 600,
  maxLoopIterations: 100,
  checkUpdatesOnStartup: true,
  selectorClickStrategies: ['text', 'contextual', 'aria', 'title', 'testid', 'id'],
  selectorInputStrategies: ['label', 'placeholder', 'aria', 'name', 'testid', 'id'],
  libraryHeuristicsMui: true,
  libraryHeuristicsAnt: true,
  navWaitUntil: 'domcontentloaded',
  editor: { ...DEFAULT_EDITOR_SETTINGS },
  startUrl: '',
  uiLocale: 'ru',
}

export function createSettingsStore(initial: AppSettingsState = defaultAppSettingsState) {
  const store = writable<AppSettingsState>(initial)
  return {
    subscribe: store.subscribe,
    patch(partial: Partial<AppSettingsState>) {
      store.update((s) => ({ ...s, ...partial }))
    },
    setEditor(editor: EditorSettings) {
      store.update((s) => ({ ...s, editor }))
    },
    applyFromDTO(dto: gui.AppSettingsDTO) {
      store.update((s) => ({
        ...s,
        browser: dto.browser || 'chromium',
        headless: dto.headless,
        parallelWorkers: dto.parallelWorkers || 1,
        slowMo: dto.slowMo ?? 0,
        scrollBeforeClick: dto.scrollBeforeClick ?? false,
        disableRecordUrlWait: dto.disableRecordUrlWait ?? false,
        hoverRecordMinMs: dto.hoverRecordMinMs || 600,
        maxLoopIterations: dto.maxLoopIterations || 100,
        checkUpdatesOnStartup: dto.checkUpdatesOnStartup !== false,
        selectorClickStrategies: dto.selectorClickStrategies?.length
          ? [...dto.selectorClickStrategies]
          : s.selectorClickStrategies,
        selectorInputStrategies: dto.selectorInputStrategies?.length
          ? [...dto.selectorInputStrategies]
          : s.selectorInputStrategies,
        libraryHeuristicsMui: dto.libraryHeuristicsMui !== false,
        libraryHeuristicsAnt: dto.libraryHeuristicsAnt !== false,
        navWaitUntil: dto.navWaitUntil || 'domcontentloaded',
        editor: editorSettingsFromDTO(dto.editor),
        startUrl: dto.startUrl || '',
        uiLocale: dto.uiLocale === 'en' ? 'en' : 'ru',
      }))
    },
    reset() {
      store.set(defaultAppSettingsState)
    },
    snapshot(): AppSettingsState {
      let state = defaultAppSettingsState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    dtoFields(state: AppSettingsState): Partial<gui.AppSettingsDTO> {
      return {
        browser: state.browser,
        headless: state.headless,
        parallelWorkers: state.parallelWorkers,
        slowMo: state.slowMo,
        maxLoopIterations: state.maxLoopIterations,
        scrollBeforeClick: state.scrollBeforeClick,
        disableRecordUrlWait: state.disableRecordUrlWait,
        hoverRecordMinMs: state.hoverRecordMinMs,
        checkUpdatesOnStartup: state.checkUpdatesOnStartup,
        selectorClickStrategies: state.selectorClickStrategies,
        selectorInputStrategies: state.selectorInputStrategies,
        libraryHeuristicsMui: state.libraryHeuristicsMui,
        libraryHeuristicsAnt: state.libraryHeuristicsAnt,
        navWaitUntil: state.navWaitUntil,
        editor: editorSettingsToDTO(state.editor),
        startUrl: state.startUrl,
        uiLocale: state.uiLocale,
      }
    },
  }
}
