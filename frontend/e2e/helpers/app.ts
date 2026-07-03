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

export async function dismissTourIfVisible(page: Page): Promise<void> {
  const skip = page.getByRole('button', { name: 'Пропустить обучение' })
  if (await skip.isVisible().catch(() => false)) {
    await skip.click()
    await expect(skip).toBeHidden({ timeout: 10_000 })
  }
}

export async function bootApp(page: Page, search = '', opts?: { keepSession?: boolean; withTour?: boolean }) {
  await page.addInitScript({ path: mockPath })
  const url = opts?.keepSession ? (search.startsWith('/') ? search : `/${search}`) : withFreshSessionQuery(search)
  await page.goto(url)
  await expect(page.locator('.ide')).toBeVisible({ timeout: 20_000 })
  if (!opts?.withTour) {
    await dismissTourIfVisible(page)
  }
}

export async function createNewScenario(page: Page) {
  await openMenuItem(page, 'Сценарий', 'Новый', { exact: true })
  await expect(page.locator('.editor-tab.file .tab-label').filter({ hasText: 'novyy-scenariy' }).first()).toBeVisible({
    timeout: 10_000,
  })
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
  await page.locator('.menubar .menu-trigger', { hasText: RUN_MENU }).click({ force: true })
  await expect(page.locator('.menu-root.open .menu-dropdown')).toBeVisible()
}

export async function clickRunMenuItem(page: Page, label: string | RegExp) {
  await openRunMenu(page)
  await page.locator('.menu-root.open .menu-dropdown').getByRole('button', { name: label }).click()
}

export async function startRecordingFromDialog(page: Page) {
  await page.keyboard.press('Control+KeyR')
  const recordDialog = page.getByRole('dialog', { name: 'Запись сценария' })
  await expect(recordDialog).toBeVisible()
  await recordDialog.getByRole('button', { name: 'Начать' }).click()
  await expect(recordDialog).toBeHidden({ timeout: 10_000 })
}

export async function stopRecording(page: Page) {
  await page.keyboard.press('Control+Shift+KeyR')
  await page.waitForTimeout(300)
}

export async function pauseRecording(page: Page) {
  await page.keyboard.press('Alt+KeyP')
  await page.waitForTimeout(300)
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

export async function openJournal(page: Page) {
  await openMenuItem(page, 'Вид', 'Журнал')
  await expect(page.locator('.panel-tab.active', { hasText: 'Журнал' })).toBeVisible()
}

export async function openResults(page: Page) {
  await openMenuItem(page, 'Вид', 'Результаты')
  await expect(page.locator('.panel-tab.active', { hasText: 'Результаты' })).toBeVisible()
}

export async function openSettings(page: Page) {
  await page.keyboard.press('Control+Comma')
  await expect(page.getByRole('dialog', { name: /Настройки/ })).toBeVisible()
}

export async function typeInEditor(page: Page, text: string) {
  await page.locator('.workspace .monaco-editor .view-lines').click()
  await page.keyboard.press('End')
  await page.keyboard.press('Enter')
  await page.keyboard.type(text)
}

export async function openRecordDialog(page: Page) {
  await page.keyboard.press('Control+KeyR')
  const dialog = page.getByRole('dialog', { name: 'Запись сценария' })
  await expect(dialog).toBeVisible()
  return dialog
}

export async function confirmRunFromDialog(page: Page, opts?: { dryRun?: boolean; continueOnFail?: boolean }) {
  const dialog = page.getByRole('dialog', { name: 'Запуск сценария' })
  await expect(dialog).toBeVisible()
  if (opts?.dryRun) {
    await dialog.locator('label.check-row', { hasText: 'Dry-run' }).locator('input').check()
  }
  if (opts?.continueOnFail) {
    await dialog.locator('label.check-row', { hasText: 'Продолжать при ошибке' }).locator('input').check()
  }
  await dialog.getByRole('button', { name: 'Запустить' }).click()
  await expect(dialog).toBeHidden({ timeout: 10_000 })
}

/** First Ctrl+Enter may open RunDialog; confirms if needed. */
export async function runCurrentScenario(page: Page, opts?: { dryRun?: boolean }) {
  await page.keyboard.press('Control+Enter')
  const dialog = page.getByRole('dialog', { name: 'Запуск сценария' })
  if (await dialog.isVisible().catch(() => false)) {
    await confirmRunFromDialog(page, opts)
  }
}

export async function saveScenario(page: Page) {
  await page.keyboard.press('Control+KeyS')
}

export async function focusBrowserFromToolbar(page: Page) {
  await page.getByRole('button', { name: 'Показать браузер' }).click()
}

export function journalPanel(page: Page) {
  return page.locator('.panel-body.text-panel')
}

export function dirtyTabs(page: Page) {
  return page.locator('.editor-tab.file .tab-label', { hasText: /\*/ })
}
