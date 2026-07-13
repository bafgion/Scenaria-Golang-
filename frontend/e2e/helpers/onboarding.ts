import { expect, type Page } from '@playwright/test'
import { dismissTourIfVisible, openMenuItem } from './app'
import { clickCatalogFeature, dismissBlockingDialogs } from './desktop'

export { dismissTourIfVisible }

export const TOUR_TOTAL = 8

export async function restartOnboardingTour(page: Page): Promise<void> {
  await dismissBlockingDialogs(page)
  await dismissTourIfVisible(page)
  await page.keyboard.press('Escape')
  await openMenuItem(page, 'Справка', 'Обучение…')
  await expect(page.getByText('Шаг 1 из 8')).toBeVisible({ timeout: 20_000 })
}

export async function expectTourStep(page: Page, step: number, title?: string | RegExp): Promise<void> {
  await expect(page.getByText(`Шаг ${step} из ${TOUR_TOTAL}`)).toBeVisible({ timeout: 15_000 })
  if (title) {
    await expect(page.getByRole('heading', { name: title })).toBeVisible({ timeout: 15_000 })
  }
}

export async function tourNext(page: Page): Promise<void> {
  await page.getByRole('button', { name: 'Далее' }).click()
}

/** Spotlight ring should overlap the target element (same screen region). */
export async function expectSpotlightOn(page: Page, targetSelector: string): Promise<void> {
  const ring = page.locator('.onboarding-ring')
  await expect(ring).toBeVisible({ timeout: 10_000 })
  const overlap = await page.evaluate((selector) => {
    const ringEl = document.querySelector('.onboarding-ring')
    const target = document.querySelector(selector)
    if (!ringEl || !target) return false
    const a = ringEl.getBoundingClientRect()
    const b = target.getBoundingClientRect()
    const overlapW = Math.max(0, Math.min(a.right, b.right) - Math.max(a.left, b.left))
    const overlapH = Math.max(0, Math.min(a.bottom, b.bottom) - Math.max(a.top, b.top))
    return overlapW > 20 && overlapH > 12
  }, targetSelector)
  expect(overlap).toBe(true)
}

export async function advanceTourThroughEditor(page: Page): Promise<void> {
  await restartOnboardingTour(page)
  await tourNext(page)
  await expectTourStep(page, 2, 'Откройте примеры')
  await page.getByRole('button', { name: 'Открыть примеры сценариев' }).click()
  await expect(page.locator('.catalog-tree .tree-file-label').first()).toBeVisible({ timeout: 30_000 })
  await expectTourStep(page, 3, 'Выберите сценарий')
  await clickCatalogFeature(page, 0)
  await expect(page.locator('.monaco-editor')).toBeVisible({ timeout: 20_000 })
  await expectTourStep(page, 4, 'Редактор Gherkin')
  await tourNext(page)
}
