export const ONBOARDING_TOUR_VERSION = 1

export type TourStepKind = 'info' | 'action'

export type TourPlacement = 'top' | 'bottom' | 'left' | 'right' | 'center'

export interface OnboardingTourStep {
  id: string
  kind: TourStepKind
  /** Value of [data-tour="…"] on a focusable UI element. */
  target: string
  titleKey: string
  bodyKey: string
  placement?: TourPlacement
  /** Ensure catalog sidebar is visible before highlighting. */
  ensureSidebar?: boolean
}

export const ONBOARDING_TOUR_STEPS: OnboardingTourStep[] = [
  {
    id: 'welcome',
    kind: 'info',
    target: 'welcome-card',
    titleKey: 'onboarding.steps.welcome.title',
    bodyKey: 'onboarding.steps.welcome.body',
    placement: 'center',
  },
  {
    id: 'open-examples',
    kind: 'action',
    target: 'welcome-examples',
    titleKey: 'onboarding.steps.openExamples.title',
    bodyKey: 'onboarding.steps.openExamples.body',
    placement: 'bottom',
  },
  {
    id: 'pick-feature',
    kind: 'action',
    target: 'catalog-tree',
    titleKey: 'onboarding.steps.pickFeature.title',
    bodyKey: 'onboarding.steps.pickFeature.body',
    placement: 'right',
    ensureSidebar: true,
  },
  {
    id: 'editor',
    kind: 'info',
    target: 'editor-workspace',
    titleKey: 'onboarding.steps.editor.title',
    bodyKey: 'onboarding.steps.editor.body',
    placement: 'right',
  },
  {
    id: 'validate',
    kind: 'action',
    target: 'menu-run',
    titleKey: 'onboarding.steps.validate.title',
    bodyKey: 'onboarding.steps.validate.body',
    placement: 'right',
  },
  {
    id: 'dry-run',
    kind: 'action',
    target: 'menu-run',
    titleKey: 'onboarding.steps.dryRun.title',
    bodyKey: 'onboarding.steps.dryRun.body',
    placement: 'right',
  },
  {
    id: 'journal',
    kind: 'action',
    target: 'panel-journal',
    titleKey: 'onboarding.steps.journal.title',
    bodyKey: 'onboarding.steps.journal.body',
    placement: 'top',
  },
  {
    id: 'finish',
    kind: 'info',
    target: 'editor-workspace',
    titleKey: 'onboarding.steps.finish.title',
    bodyKey: 'onboarding.steps.finish.body',
    placement: 'center',
  },
]
