import { expect, test } from '@playwright/test'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { bootApp, catalogFeature, mockPath, openMenuItem } from '../helpers/app'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const outDir = path.resolve(__dirname, '../../../docs/images')

const labels = {
  ru: {
    project: 'Проект',
    runMenu: 'Запись и тест',
    openExamples: 'Открыть примеры сценариев',
    runDialog: 'Запустить…',
    settings: /Настройки/,
    recordDialog: 'Запись сценария',
    commandPalette: 'Палитра команд',
    palettePlaceholder: 'Введите команду…',
    onboardingNext: 'Далее',
  },
  en: {
    project: 'Project',
    runMenu: 'Run & record',
    openExamples: 'Open example scenarios',
    runDialog: 'Run…',
    settings: /Settings/,
    recordDialog: 'Record scenario',
    commandPalette: 'Command palette',
    palettePlaceholder: 'Type a command…',
    onboardingNext: 'Next',
  },
} as const

async function shotIde(page: import('@playwright/test').Page, name: string) {
  const ide = page.locator('.ide')
  await expect(ide).toBeVisible()
  await page.waitForTimeout(200)
  await ide.screenshot({ path: path.join(outDir, name) })
}

test.describe.configure({ mode: 'serial' })

for (const locale of ['ru', 'en'] as const) {
  test.describe(`docs screenshots (${locale})`, () => {
    test.beforeEach(async ({ page }) => {
      await page.addInitScript({ path: mockPath })
    })

    test('welcome', async ({ page }) => {
      await bootApp(page, `?__fresh=1&locale=${locale}`)
      await shotIde(page, `gui-welcome-${locale}.png`)
    })

    test('main workspace', async ({ page }) => {
      await bootApp(page, `?__fresh=1&locale=${locale}`)
      const l = labels[locale]
      await openMenuItem(page, l.project, l.openExamples)
      await expect(page.locator('.catalog-tree')).toBeVisible({ timeout: 10_000 })
      await catalogFeature(page, '^○ smoke$').click()
      await expect(page.locator('.editor-tab.file .tab-label', { hasText: 'smoke' })).toBeVisible()
      await shotIde(page, `gui-main-${locale}.png`)
    })

    test('settings', async ({ page }) => {
      await bootApp(page, `?__fresh=1&locale=${locale}`)
      await page.keyboard.press('Control+Comma')
      await expect(page.getByRole('dialog', { name: labels[locale].settings })).toBeVisible()
      await shotIde(page, `gui-settings-${locale}.png`)
      await page.keyboard.press('Escape')
    })

    test('run dialog', async ({ page }) => {
      await bootApp(page, `?__fresh=1&locale=${locale}`)
      const l = labels[locale]
      await openMenuItem(page, l.project, l.openExamples)
      await catalogFeature(page, '^○ smoke$').click()
      await openMenuItem(page, l.runMenu, l.runDialog)
      await expect(page.getByRole('dialog')).toBeVisible()
      await shotIde(page, `gui-run-dialog-${locale}.png`)
      await page.keyboard.press('Escape')
    })

    test('record dialog', async ({ page }) => {
      await bootApp(page, `?__fresh=1&locale=${locale}`)
      const l = labels[locale]
      await openMenuItem(page, l.project, l.openExamples)
      await catalogFeature(page, '^○ smoke$').click()
      await page.keyboard.press('Control+KeyR')
      await expect(page.getByRole('dialog', { name: labels[locale].recordDialog })).toBeVisible({
        timeout: 10_000,
      })
      await shotIde(page, `gui-record-dialog-${locale}.png`)
      await page.keyboard.press('Escape')
    })

    test('command palette', async ({ page }) => {
      await bootApp(page, `?__fresh=1&locale=${locale}`)
      await page.keyboard.press('Control+Shift+KeyP')
      const dialog = page.getByRole('dialog', { name: labels[locale].commandPalette })
      await expect(dialog).toBeVisible()
      await expect(page.getByPlaceholder(labels[locale].palettePlaceholder)).toBeVisible()
      await shotIde(page, `gui-command-palette-${locale}.png`)
      await page.keyboard.press('Escape')
    })

    test('onboarding', async ({ page }) => {
      await bootApp(page, `?__fresh=1&locale=${locale}`, { withTour: true })
      await expect(page.getByRole('button', { name: labels[locale].onboardingNext })).toBeVisible({
        timeout: 10_000,
      })
      await shotIde(page, `gui-onboarding-${locale}.png`)
    })
  })
}
