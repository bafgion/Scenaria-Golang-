import { expect, type Page } from '@playwright/test'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
export const mockPath = path.join(__dirname, '..', 'fixtures', 'wails-mock.js')
export const E2E_PROJECT = 'C:/e2e/project'

export async function bootApp(page: Page, search = '') {
  await page.addInitScript({ path: mockPath })
  await page.goto(`/${search}`)
  await expect(page.locator('.ide')).toBeVisible({ timeout: 20_000 })
}

export async function createNewScenario(page: Page) {
  await page.locator('.menubar .menu-trigger', { hasText: 'Сценарий' }).click()
  await page.getByRole('button', { name: 'Новый', exact: true }).click()
  await expect(page.locator('.editor-tab.file .tab-label', { hasText: 'novyy-scenariy.feature' })).toBeVisible({
    timeout: 10_000,
  })
}

export async function openTestProject(page: Page, projectPath = E2E_PROJECT) {
  await page.locator('.menubar .menu-trigger', { hasText: 'Проект' }).click()
  await page.getByRole('button', { name: 'Открыть проект…' }).click()
  const dialog = page.getByRole('dialog', { name: 'Открыть проект' })
  await expect(dialog).toBeVisible()
  await dialog.locator('input').fill(projectPath)
  await dialog.getByRole('button', { name: 'Открыть', exact: true }).click()
  await expect(dialog).toBeHidden({ timeout: 10_000 })
  await page.locator('.menubar .menu-trigger', { hasText: 'Запись и тест' }).click()
  await expect(page.getByRole('button', { name: 'Проверить…' })).toBeEnabled({ timeout: 10_000 })
  await page.keyboard.press('Escape')
}

export async function openMenuItem(page: Page, menu: string, item: string | RegExp) {
  await page.locator('.menubar .menu-trigger', { hasText: menu }).click()
  await page.getByRole('button', { name: item }).click()
}
