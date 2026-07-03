import { describe, expect, it } from 'vitest'

/**
 * Documented z-index scale from style.css :root.
 * Onboarding targets (2810) must stay above dimming/blockers (2750–2751)
 * and below the tour card (2820).
 */
const Z = {
  contextMenu: 1400,
  menubar: 2200,
  menuDropdown: 2201,
  browserOverlay: 2400,
  modal: 2500,
  modalTop: 2600,
  modalConfirm: 2700,
  onboardingDim: 2750,
  onboardingBlock: 2751,
  onboardingTarget: 2810,
  onboardingRing: 2815,
  onboardingCard: 2820,
  splash: 3000,
} as const

describe('z-index layer scale', () => {
  it('shell chrome stays below modals', () => {
    expect(Z.menubar).toBeLessThan(Z.modal)
    expect(Z.menuDropdown).toBeLessThan(Z.modal)
    expect(Z.browserOverlay).toBeLessThan(Z.modal)
  })

  it('modals stay below onboarding tour', () => {
    expect(Z.modalConfirm).toBeLessThan(Z.onboardingDim)
  })

  it('onboarding interactive targets are above dimming and blockers', () => {
    expect(Z.onboardingTarget).toBeGreaterThan(Z.onboardingBlock)
    expect(Z.onboardingTarget).toBeGreaterThan(Z.onboardingDim)
  })

  it('onboarding ring and card stack above interactive targets', () => {
    expect(Z.onboardingRing).toBeGreaterThan(Z.onboardingTarget)
    expect(Z.onboardingCard).toBeGreaterThan(Z.onboardingRing)
  })

  it('splash is above everything else', () => {
    const values = Object.values(Z)
    expect(Z.splash).toBe(Math.max(...values))
  })
})
