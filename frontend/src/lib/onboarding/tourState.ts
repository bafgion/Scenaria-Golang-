import type { OnboardingTourStep } from './tourSteps'
import { ONBOARDING_TOUR_STEPS } from './tourSteps'

export interface TourContext {
  projectPath: string
  featuresCount: number
  isWelcome: boolean
  activeTab: string
  welcomeKey: string
  validateDone: boolean
  dryRunDone: boolean
  bottomPanelOpen: boolean
  bottomTab: string
}

export function isTourStepComplete(stepId: string, ctx: TourContext): boolean {
  switch (stepId) {
    case 'welcome':
    case 'editor':
    case 'finish':
      return true
    case 'open-examples':
      return !!ctx.projectPath && ctx.featuresCount > 0
    case 'pick-feature':
      return !ctx.isWelcome && ctx.activeTab !== ctx.welcomeKey
    case 'validate':
      return ctx.validateDone
    case 'dry-run':
      return ctx.dryRunDone
    case 'journal':
      return ctx.bottomPanelOpen && ctx.bottomTab === 'journal'
    default:
      return false
  }
}

export function canAdvanceTourStep(step: OnboardingTourStep, ctx: TourContext): boolean {
  if (step.kind === 'info') return true
  return isTourStepComplete(step.id, ctx)
}

/** Highest tour step that matches the current UI state (rewind when user breaks the flow). */
export function maxValidTourStepIndex(ctx: TourContext): number {
  const idx = (id: string) => ONBOARDING_TOUR_STEPS.findIndex((s) => s.id === id)
  const hasProject = !!ctx.projectPath && ctx.featuresCount > 0
  const hasFeatureTab = !ctx.isWelcome && ctx.activeTab !== ctx.welcomeKey

  if (!hasProject) return idx('open-examples')
  if (!hasFeatureTab) return idx('pick-feature')
  if (!ctx.validateDone) return idx('validate')
  if (!ctx.dryRunDone) return idx('dry-run')
  if (!(ctx.bottomPanelOpen && ctx.bottomTab === 'journal')) return idx('journal')
  return idx('finish')
}

export function tourTargetSelector(target: string): string {
  return `[data-tour="${target}"]`
}
