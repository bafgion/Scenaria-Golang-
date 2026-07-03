/**
 * Desktop smoke: интерактивный тур (WebView2 + scenaria-gui.exe).
 *
 * Запуск:
 *   ./scripts/desktop-smoke.ps1
 *   cd frontend && npx playwright test -c e2e/playwright.desktop.config.ts e2e/specs/desktop-tour-onboarding.spec.ts
 */
import { expect, test } from '@playwright/test'
import {
  advanceTourThroughEditor,
  dismissTourIfVisible,
  expectSpotlightOn,
  expectTourStep,
  restartOnboardingTour,
} from '../helpers/onboarding'
import { connectDesktop, disconnectDesktop, dismissBlockingDialogs, waitForAppReady } from '../helpers/desktop'

test.describe.configure({ mode: 'serial', timeout: 240_000 })

test.beforeAll(async () => {
  const page = await connectDesktop()
  await waitForAppReady(page)
})

test.afterAll(async () => {
  const page = await connectDesktop().catch(() => null)
  if (page) await dismissTourIfVisible(page)
  await disconnectDesktop()
})

test('тур: перезапуск из справки', async () => {
  const page = await connectDesktop()
  await dismissBlockingDialogs(page)
  await dismissTourIfVisible(page)
  await restartOnboardingTour(page)
  await expect(page.getByText('Шаг 1 из 8')).toBeVisible({ timeout: 15_000 })
  await expect(page.locator('.onboarding-backdrop, .onboarding-hole')).toHaveCount(1, { timeout: 5000 })
})

test('тур: шаги 1–4 до редактора', async () => {
  const page = await connectDesktop()
  await advanceTourThroughEditor(page)
  await expect(page.locator('.onboarding-ring')).toBeVisible()
})

test('тур: шаг 5 — подсветка меню и проверка синтаксиса', async () => {
  const page = await connectDesktop()
  await advanceTourThroughEditor(page)
  await expectTourStep(page, 5, 'Проверка синтаксиса')
  await expect(page.locator('.menu-root.open [data-tour="menu-run-dropdown"]')).toBeVisible({ timeout: 10_000 })
  await expectSpotlightOn(page, '[data-tour="menu-run"]')

  await page.locator('.menu-root.open .menu-dropdown').getByRole('button', { name: 'Проверить…' }).click()
  const dialog = page.getByRole('dialog', { name: 'Проверка сценария' })
  await expect(dialog).toBeVisible({ timeout: 10_000 })
  await dialog.getByRole('button', { name: 'Проверить' }).click()
  await expect(dialog).toBeHidden({ timeout: 30_000 })
  await expectTourStep(page, 6, /Dry-run/i)
})

test('тур: шаг 6 — dry-run', async () => {
  const page = await connectDesktop()
  await advanceTourThroughEditor(page)
  await expectTourStep(page, 5, 'Проверка синтаксиса')
  await page.locator('.menu-root.open .menu-dropdown').getByRole('button', { name: 'Проверить…' }).click()
  const validateDialog = page.getByRole('dialog', { name: 'Проверка сценария' })
  await validateDialog.getByRole('button', { name: 'Проверить' }).click()
  await expect(validateDialog).toBeHidden({ timeout: 30_000 })

  await expectTourStep(page, 6, /Dry-run/i)
  await expect(page.locator('.menu-root.open')).toBeVisible()
  await page.locator('.menu-root.open .menu-dropdown').getByRole('button', { name: 'Dry-run', exact: true }).click()
  const runDialog = page.getByRole('dialog').filter({ has: page.getByRole('button', { name: 'Запустить' }) })
  if (await runDialog.isVisible({ timeout: 5000 }).catch(() => false)) {
    await runDialog.getByRole('button', { name: 'Запустить' }).click()
    await expect(runDialog).toBeHidden({ timeout: 30_000 })
  }
  await expect(page.locator('footer.status-bar .status-message')).toContainText(/завершён|остановлен/i, {
    timeout: 120_000,
  })
  await expect(page.getByText('Шаг 7 из 8')).toBeVisible({ timeout: 30_000 })
  await page.locator('[data-tour="panel-journal"]').click()
  await expectTourStep(page, 7, 'Журнал')
})
