import { writable } from 'svelte/store'
import type { gui } from '../../wailsjs/go/models'
import { ONBOARDING_TOUR_VERSION } from '../lib/onboarding/tourSteps'

export type UiPrefsState = {
  toolbarCompact: boolean
  stepsPanelVisible: boolean
  stepsPanelHeight: number
  sidebarWidth: number
  checklistDismissed: boolean
  welcomePlayedSuccess: boolean
  onboardingCompleted: boolean
  onboardingDismissed: boolean
  runDialogConfirmed: boolean
  pickerDuringRecording: boolean
  stepsPanelCollapsed: boolean
}

export const defaultUiPrefsState: UiPrefsState = {
  toolbarCompact: false,
  stepsPanelVisible: true,
  stepsPanelHeight: 160,
  sidebarWidth: 260,
  checklistDismissed: false,
  welcomePlayedSuccess: false,
  onboardingCompleted: false,
  onboardingDismissed: false,
  runDialogConfirmed: false,
  pickerDuringRecording: false,
  stepsPanelCollapsed: true,
}

export function createUiPrefsStore(initial: UiPrefsState = defaultUiPrefsState) {
  const store = writable<UiPrefsState>(initial)
  return {
    subscribe: store.subscribe,
    patch(partial: Partial<UiPrefsState>) {
      store.update((state) => ({ ...state, ...partial }))
    },
    patchLocal(partial: Partial<UiPrefsState>) {
      this.patch(partial)
    },
    setStepsPanelCollapsed(stepsPanelCollapsed: boolean) {
      store.update((state) => ({ ...state, stepsPanelCollapsed }))
    },
    toggleStepsPanelCollapsed() {
      store.update((state) => ({ ...state, stepsPanelCollapsed: !state.stepsPanelCollapsed }))
    },
    applyFromDTO(dto: gui.AppSettingsDTO) {
      store.update((state) => {
        let onboardingCompleted = !!dto.onboardingCompleted
        let onboardingDismissed = !!dto.onboardingDismissed
        if (!onboardingCompleted && dto.onboardingVersion !== ONBOARDING_TOUR_VERSION) {
          onboardingCompleted = false
          onboardingDismissed = false
        }
        return {
          ...state,
          toolbarCompact: !!dto.toolbarCompact,
          stepsPanelVisible: dto.stepsPanelVisible !== false,
          stepsPanelHeight: dto.stepsPanelHeight || 160,
          sidebarWidth: dto.sidebarWidth || state.sidebarWidth,
          checklistDismissed: !!dto.checklistDismissed,
          welcomePlayedSuccess: !!dto.welcomePlayedSuccess,
          onboardingCompleted,
          onboardingDismissed,
          runDialogConfirmed: !!dto.runDialogConfirmed,
          pickerDuringRecording: !!dto.pickerDuringRecording,
        }
      })
    },
    dtoFields(state: UiPrefsState): Partial<gui.AppSettingsDTO> {
      return {
        toolbarCompact: state.toolbarCompact,
        stepsPanelVisible: state.stepsPanelVisible,
        stepsPanelHeight: state.stepsPanelHeight,
        sidebarWidth: state.sidebarWidth,
        checklistDismissed: state.checklistDismissed,
        welcomePlayedSuccess: state.welcomePlayedSuccess,
        onboardingCompleted: state.onboardingCompleted,
        onboardingDismissed: state.onboardingDismissed,
        onboardingVersion: ONBOARDING_TOUR_VERSION,
        runDialogConfirmed: state.runDialogConfirmed,
        pickerDuringRecording: state.pickerDuringRecording,
      }
    },
    snapshot(): UiPrefsState {
      let state = defaultUiPrefsState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultUiPrefsState)
    },
  }
}
