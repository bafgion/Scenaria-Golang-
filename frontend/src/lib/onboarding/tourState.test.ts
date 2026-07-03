import { describe, expect, it } from 'vitest'
import { canAdvanceTourStep, isTourStepComplete, maxValidTourStepIndex } from './tourState'
import { ONBOARDING_TOUR_STEPS } from './tourSteps'

const welcomeKey = '__welcome__'

function ctx(overrides: Partial<import('./tourState').TourContext> = {}) {
  return {
    projectPath: '',
    featuresCount: 0,
    isWelcome: true,
    activeTab: welcomeKey,
    welcomeKey,
    validateDone: false,
    dryRunDone: false,
    bottomPanelOpen: false,
    bottomTab: 'journal',
    ...overrides,
  }
}

describe('tourState', () => {
  it('open-examples completes when project has features', () => {
    expect(isTourStepComplete('open-examples', ctx())).toBe(false)
    expect(isTourStepComplete('open-examples', ctx({ projectPath: '/examples', featuresCount: 3 }))).toBe(true)
  })

  it('pick-feature completes when editor tab is active', () => {
    expect(
      isTourStepComplete('pick-feature', ctx({ projectPath: '/examples', isWelcome: false, activeTab: '/a.feature' })),
    ).toBe(true)
  })

  it('info steps always allow advance', () => {
    const welcome = ONBOARDING_TOUR_STEPS[0]
    expect(canAdvanceTourStep(welcome, ctx())).toBe(true)
  })

  it('action steps require completion', () => {
    const validate = ONBOARDING_TOUR_STEPS.find((s) => s.id === 'validate')!
    expect(canAdvanceTourStep(validate, ctx())).toBe(false)
    expect(canAdvanceTourStep(validate, ctx({ validateDone: true }))).toBe(true)
  })

  it('maxValidTourStepIndex allows open-examples before project is open', () => {
    const openExamples = ONBOARDING_TOUR_STEPS.findIndex((s) => s.id === 'open-examples')
    expect(maxValidTourStepIndex(ctx())).toBe(openExamples)
  })

  it('maxValidTourStepIndex rewinds when editor tab is closed', () => {
    const pick = ONBOARDING_TOUR_STEPS.findIndex((s) => s.id === 'pick-feature')
    const validate = ONBOARDING_TOUR_STEPS.findIndex((s) => s.id === 'validate')
    expect(
      maxValidTourStepIndex(
        ctx({ projectPath: '/examples', featuresCount: 3, isWelcome: true, activeTab: welcomeKey }),
      ),
    ).toBe(pick)
    expect(
      maxValidTourStepIndex(
        ctx({
          projectPath: '/examples',
          featuresCount: 3,
          isWelcome: false,
          activeTab: '/a.feature',
          validateDone: false,
        }),
      ),
    ).toBe(validate)
  })
})
