import { expect, type Page } from '@playwright/test'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
export const mockPath = path.join(__dirname, '..', 'fixtures', 'wails-mock.js')
export const E2E_PROJECT = 'C:/e2e/project'

const RUN_MENU = 'Запись и тест'

function withFreshSessionQuery(search: string): string {
  const path = search.startsWith('/') ? search : `/${search}`
  return path.includes('?') ? `${path}&__fresh=1` : `${path}?__fresh=1`
}

export async function bootApp(page: Page, search = '', opts?: { keepSession?: boolean }) {
  await page.addInitScript({ path: mockPath })
  const url = opts?.keepSession ? (search.startsWith('/') ? search : `/${search}`) : withFreshSessionQuery(search)
  await page.goto(url)
  await expect(page.locator('.ide')).toBeVisible({ timeout: 20_000 })
}

export async function createNewScenario(page: Page) {
  await openMenuItem(page, 'Сценарий', 'Новый', { exact: true })
  await expect(page.getByText('novyy-scenariy.feature')).toBeVisible({ timeout: 10_000 })
}

export async function openTestProject(page: Page, projectPath = E2E_PROJECT) {
  await openMenuItem(page, 'Проект', 'Открыть проект…')
  const dialog = page.getByRole('dialog', { name: 'Открыть проект' })
  await expect(dialog).toBeVisible()
  await dialog.getByRole('textbox').fill(projectPath)
  await dialog.getByRole('button', { name: 'Открыть', exact: true }).click()
  await expect(dialog).toBeHidden({ timeout: 10_000 })
  await openRunMenu(page)
  await expect(page.locator('.menu-root.open .menu-dropdown').getByRole('button', { name: 'Проверить…' })).toBeEnabled({
    timeout: 10_000,
  })
  await page.keyboard.press('Escape')
}

export async function openRunMenu(page: Page) {
  await page.locator('.menubar .menu-trigger', { hasText: RUN_MENU }).click()
  await expect(page.locator('.menu-root.open .menu-dropdown')).toBeVisible()
}

export async function clickRunMenuItem(page: Page, label: string | RegExp) {
  await openRunMenu(page)
  await page.locator('.menu-root.open .menu-dropdown').getByRole('button', { name: label }).click()
}

export async function startRecordingFromDialog(page: Page) {
  await clickRunMenuItem(page, /Запись\s+Ctrl\+R/i)
  const recordDialog = page.getByRole('dialog', { name: 'Запись сценария' })
  await expect(recordDialog).toBeVisible()
  await recordDialog.getByRole('button', { name: 'Начать' }).click()
  await expect(recordDialog).toBeHidden({ timeout: 10_000 })
}

export async function stopRecording(page: Page) {
  await clickRunMenuItem(page, 'Стоп')
}

export async function pauseRecording(page: Page) {
  await clickRunMenuItem(page, 'Пауза')
}

export async function expectRunMenuItemDisabled(page: Page, label: string | RegExp) {
  await openRunMenu(page)
  await expect(page.locator('.menu-root.open .menu-dropdown').getByRole('button', { name: label })).toBeDisabled()
  await page.keyboard.press('Escape')
}

export async function expectRunMenuItemEnabled(page: Page, label: string | RegExp) {
  await openRunMenu(page)
  await expect(page.locator('.menu-root.open .menu-dropdown').getByRole('button', { name: label })).toBeEnabled()
  await page.keyboard.press('Escape')
}

export function statusMessage(page: Page) {
  return page.locator('footer.status-bar .status-message')
}

export function editorLine(page: Page, text: string | RegExp) {
  return page.locator('.workspace .monaco-editor').getByText(text)
}

export async function openMenuItem(page: Page, menu: string, item: string | RegExp, opts?: { exact?: boolean }) {
  await page.locator('.menubar .menu-trigger', { hasText: menu }).click()
  await page
    .locator('.menu-root.open .menu-dropdown')
    .getByRole('button', { name: item, exact: opts?.exact })
    .click()
}

export function catalogFeature(page: Page, featureName: string) {
  return page.locator('.catalog-tree').getByRole('button', { name: new RegExp(featureName) })
}
