import { writable } from 'svelte/store'
import { gui } from '../../wailsjs/go/models'

export type VanessaRunFormState = {
  dry: boolean
  preferRerun: boolean
  tag: string
  excludeTags: string
  scenario: string
  rerunDir: string
  installEpf: boolean
  epfUrl: string
  epfDest: string
  platformExe: string
  epfPath: string
  ib: string
  reportAllure: boolean
  vaDir: string
  vaFiles: string
}

export type VanessaRunState = VanessaRunFormState & {
  dialogScenarios: string[]
  running: boolean
  watchDir: string
  plannedTotal: number
  snapshot: gui.VanessaRunSnapshotDTO
}

export const defaultVanessaRunFormState: VanessaRunFormState = {
  dry: false,
  preferRerun: false,
  tag: '',
  excludeTags: '',
  scenario: '',
  rerunDir: '',
  installEpf: false,
  epfUrl: '',
  epfDest: '',
  platformExe: '',
  epfPath: '',
  ib: '',
  reportAllure: false,
  vaDir: '',
  vaFiles: '',
}

export const defaultVanessaRunState: VanessaRunState = {
  ...defaultVanessaRunFormState,
  dialogScenarios: [],
  running: false,
  watchDir: '',
  plannedTotal: 1,
  snapshot: new gui.VanessaRunSnapshotDTO(),
}

export function createVanessaRunStore(initial: VanessaRunState = defaultVanessaRunState) {
  const store = writable<VanessaRunState>(initial)
  let pollTimer: ReturnType<typeof setInterval> | null = null
  return {
    subscribe: store.subscribe,
    patch(partial: Partial<VanessaRunState>) {
      store.update((s) => ({ ...s, ...partial }))
    },
    prepareDialog(opts: {
      dry: boolean
      preferRerun?: boolean
      scenario?: string
      dialogScenarios?: string[]
      vaDir?: string
    }) {
      store.set({
        ...defaultVanessaRunState,
        dry: opts.dry,
        preferRerun: !!opts.preferRerun,
        scenario: opts.scenario ?? '',
        dialogScenarios: opts.dialogScenarios ?? [],
        vaDir: opts.vaDir ?? '',
      })
    },
    setRunning(running: boolean) {
      store.update((s) => ({ ...s, running }))
    },
    setSnapshot(snapshot: gui.VanessaRunSnapshotDTO) {
      store.update((s) => ({ ...s, snapshot }))
    },
    resetRuntime() {
      store.update((s) => ({
        ...s,
        running: false,
        watchDir: '',
        plannedTotal: 1,
        snapshot: new gui.VanessaRunSnapshotDTO(),
      }))
    },
    setPollTimer(timer: ReturnType<typeof setInterval> | null) {
      if (pollTimer) clearInterval(pollTimer)
      pollTimer = timer
    },
    clearPollTimer() {
      if (pollTimer) {
        clearInterval(pollTimer)
        pollTimer = null
      }
    },
    buildPluginRequest(state: VanessaRunState): gui.PluginRunRequest {
      const exclude = state.excludeTags
        .split(',')
        .map((t) => t.trim())
        .filter(Boolean)
      return {
        name: 'vanessa',
        dryRun: state.dry,
        tag: state.tag.trim(),
        excludeTags: exclude,
        scenario: state.scenario.trim(),
        rerunFailedRunDir: state.rerunDir.trim(),
        installEpf: state.installEpf,
        epfUrl: state.epfUrl.trim(),
        epfDest: state.epfDest.trim(),
        platformExe: state.platformExe.trim(),
        epfPath: state.epfPath.trim(),
        ibConnection: state.ib.trim(),
        reportAllure: state.reportAllure,
        vaDir: state.vaDir.trim(),
        vaFiles: state.vaFiles.trim(),
      }
    },
    snapshot(): VanessaRunState {
      let state = defaultVanessaRunState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      this.clearPollTimer()
      store.set(defaultVanessaRunState)
    },
  }
}
