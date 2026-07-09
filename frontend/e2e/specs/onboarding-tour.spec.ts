import { expect, test } from '@playwright/test'
import { expectSpotlightOn, expectTourStep } from '../helpers/onboarding'
import { bootApp, mockPath } from '../helpers/app'

test.beforeEach(async ({ page }) => {
  await page.addInitScript({ path: mockPath })
})

test('onboarding tour appears on first launch', async ({ page }) => {
  await bootApp(page, '', { withTour: true })
  await expect(page.getByText('Добро пожаловать в Scenaria')).toBeVisible({ timeout: 15_000 })
  await expect(page.getByRole('button', { name: 'Пропустить обучение' })).toBeVisible()
  await expect(page.getByText('Шаг 1 из 8')).toBeVisible()
})

test('onboarding tour can be skipped', async ({ page }) => {
  await bootApp(page, '', { withTour: true })
  await page.getByRole('button', { name: 'Пропустить обучение' }).click()
  await expect(page.getByText('Добро пожаловать в Scenaria')).toBeHidden()
})

test('onboarding tour advances from step 1 to step 2', async ({ page }) => {
  await bootApp(page, '', { withTour: true })
  await expect(page.getByText('Шаг 1 из 8')).toBeVisible({ timeout: 15_000 })
  await page.getByRole('button', { name: 'Далее' }).click()
  await expect(page.getByText('Шаг 2 из 8')).toBeVisible()
  await expect(page.getByText('Откройте примеры')).toBeVisible()
})

test('onboarding tour step 5 highlights run menu', async ({ page }) => {
  await bootApp(page, '', { withTour: true })
  await page.getByRole('button', { name: 'Далее' }).click()
  await page.getByRole('button', { name: 'Открыть примеры сценариев' }).click()
  await expect(page.locator('.catalog-tree .tree-file-label').first()).toBeVisible({ timeout: 15_000 })
  await page.locator('.catalog-tree .catalog-tree-row.file').first().click()
  await expectTourStep(page, 4, 'Редактор Gherkin')
  await page.getByRole('button', { name: 'Далее' }).click()
  await expectTourStep(page, 5, 'Проверка синтаксиса')
  await expect(page.locator('body.onboarding-tour-active')).toBeVisible()
  await expect(page.locator('.menubar.onboarding-elevated')).toBeVisible()
  await expect(page.locator('.menu-root.open [data-tour="menu-run-dropdown"]')).toBeVisible({ timeout: 15_000 })
  await expectSpotlightOn(page, '[data-tour="menu-run"]')
  await expect(page.locator('.onboarding-ring')).toBeVisible()
  const layers = await page.evaluate(() => {
    const elevated = document.querySelector('.menubar.onboarding-elevated')
    const ring = document.querySelector('.onboarding-ring')
    const z = (el: Element | null) => (el ? Number.parseInt(getComputedStyle(el).zIndex, 10) : 0)
    return {
      blockers: document.querySelectorAll('.onboarding-blocker').length,
      elevatedZ: z(elevated),
      ringZ: z(ring),
    }
  })
  expect(layers.blockers).toBe(0)
  expect(layers.elevatedZ).toBeGreaterThan(2800)
  expect(layers.ringZ).toBeGreaterThan(layers.elevatedZ)
  const overlap = await page.evaluate(() => {
    const ring = document.querySelector('.onboarding-ring')
    const menu = document.querySelector('[data-tour="menu-run"]')
    if (!ring || !menu) return false
    const a = ring.getBoundingClientRect()
    const b = menu.getBoundingClientRect()
    return (
      Math.min(a.right, b.right) - Math.max(a.left, b.left) > 20 &&
      Math.min(a.bottom, b.bottom) - Math.max(a.top, b.top) > 12
    )
  })
  expect(overlap).toBe(true)
  await page.locator('.menu-root.open .menu-dropdown').getByRole('button', { name: 'Проверить…' }).click()
  await expect(page.getByRole('dialog', { name: 'Проверка сценария' })).toBeVisible()
})

test('onboarding tour does not show after skip', async ({ page }) => {
  await bootApp(page, '', { withTour: true })
  await page.getByRole('button', { name: 'Пропустить обучение' }).click()
  await page.reload()
  await expect(page.locator('.ide')).toBeVisible({ timeout: 20_000 })
  await expect(page.getByText('Добро пожаловать в Scenaria')).toHaveCount(0)
})
