/**
 * E2E coverage for docs/QA-DAILY-USE.md (mock Wails — UI flows, not real Playwright runs).
 * Test titles use checklist IDs: «1.2 …», «4.1 …», etc.
 */
import { expect, test } from '@playwright/test'
import {
  bootApp,
  catalogFeature,
  createNewScenario,
  dirtyTabs,
  mockPath,
  openJournal,
  openMenuItem,
  openResults,
  openSettings,
  openTestProject,
  startRecordingFromDialog,
  statusMessage,
  stopRecording,
  typeInEditor,
} from '../helpers/app'

test.beforeEach(async ({ page }) => {
  await page.addInitScript({ path: mockPath })
})

test.describe('1. Onboarding', () => {
  test('1.1 welcome checklist and quick start auto-opens examples and starts recording', async ({ page }) => {
    await bootApp(page)
    await expect(page.getByRole('button', { name: 'Открыть проект' })).toBeVisible()
    await expect(page.getByText('Записать сценарий')).toBeVisible()
    await page.getByRole('button', { name: 'Быстрый старт' }).click()
    await openJournal(page)
    await expect(page.locator('.panel-body.text-panel')).toContainText('Открыты примеры сценариев')
    await expect(page.locator('.panel-body.text-panel')).toContainText('Запись начата')
    await expect(page.locator('.status-bar .recording-target')).toContainText('Запись')
  })

  test('1.2 Ctrl+B without project auto-opens examples and launches browser', async ({ page }) => {
    await bootApp(page)
    await page.keyboard.press('Control+KeyB')
    await openJournal(page)
    await expect(page.locator('.panel-body.text-panel')).toContainText('Открыты примеры сценариев')
    await expect(page.locator('.browser-overlay')).toBeVisible()
  })

  test('1.3 new project wizard creates feature tab', async ({ page }) => {
    await bootApp(page, '?e2e=new-project')
    await openMenuItem(page, 'Проект', 'Новый проект…')
    const dialog = page.getByRole('dialog', { name: 'Новый проект' })
    await dialog.getByRole('button', { name: 'Обзор…' }).click()
    await dialog.locator('label', { hasText: 'Название feature' }).locator('input').fill('smoke')
    await dialog.getByRole('button', { name: 'Создать' }).click()
    await expect(dialog).toBeHidden({ timeout: 10_000 })
    await expect(page.locator('.editor-tab.file .tab-label', { hasText: 'smoke.feature' })).toBeVisible({
      timeout: 10_000,
    })
    await openJournal(page)
    await expect(page.locator('.panel-body.text-panel')).toContainText('Проект открыт')
  })

  test('1.4 open examples shows catalog hint', async ({ page }) => {
    await bootApp(page, '?e2e=examples')
    await openMenuItem(page, 'Проект', 'Открыть примеры сценариев')
    await expect(page.locator('.catalog-tree')).toBeVisible({ timeout: 10_000 })
    await expect(statusMessage(page)).toContainText('Примеры открыты')
  })

  test('1.5 dismiss checklist persists after reload', async ({ page }) => {
    await bootApp(page)
    await page.getByRole('button', { name: 'Скрыть чеклист' }).click()
    await expect(page.getByRole('button', { name: 'Скрыть чеклист' })).toHaveCount(0)
    await page.reload()
    await expect(page.locator('.ide')).toBeVisible({ timeout: 20_000 })
    await expect(page.getByRole('button', { name: 'Скрыть чеклист' })).toHaveCount(0)
  })

  test('1.6 update modal deferred until project opened', async ({ page }) => {
    await bootApp(page, '?e2e=update-available')
    await expect(page.getByRole('dialog', { name: 'Обновления' })).toHaveCount(0)
    await openMenuItem(page, 'Проект', 'Открыть проект…')
    const openDialog = page.getByRole('dialog', { name: 'Открыть проект' })
    await openDialog.getByRole('textbox').fill('C:/e2e/project')
    await openDialog.getByRole('button', { name: 'Открыть', exact: true }).click()
    await expect(page.getByRole('dialog', { name: 'Обновления' })).toBeVisible({ timeout: 10_000 })
  })
})

test.describe('2. Monaco / editor', () => {
  test('2.1 dirty asterisk only on edited tabs', async ({ page }) => {
    await bootApp(page)
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    await createNewScenario(page)
    await page.locator('.editor-tab.file', { hasText: 'smoke' }).click()
    await expect(page.locator('.editor-tab.file', { hasText: 'smoke' }).locator('.tab-label', { hasText: /\*/ })).toHaveCount(0)
    await typeInEditor(page, '# edited-smoke')
    await expect(page.locator('.editor-tab.file', { hasText: 'smoke' }).locator('.tab-label', { hasText: /\*/ })).toHaveCount(1)
    await page.locator('.editor-tab.file', { hasText: 'novyy-scenariy' }).first().click()
    await expect(dirtyTabs(page)).toHaveCount(2)
    await expect(page.locator('.editor-tab.file.active .tab-label', { hasText: /\*/ })).toHaveCount(1)
  })

  test('2.2 save as uses current editor text', async ({ page }) => {
    await bootApp(page, '?e2e=save-as')
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    const marker = 'SAVE_AS_MARKER'
    await typeInEditor(page, marker)
    await openMenuItem(page, 'Сценарий', 'Сохранить как…')
    await expect(statusMessage(page)).toHaveText('Сохранено', { timeout: 10_000 })
    await expect(page.locator('.editor-tab.file .tab-label', { hasText: 'saved-as.feature' })).toBeVisible()
    await expect(page.locator('.monaco-editor .view-line', { hasText: marker })).toBeVisible()
  })

  test('2.3 settings offers system editor theme', async ({ page }) => {
    await bootApp(page)
    await openSettings(page)
    const dialog = page.getByRole('dialog', { name: /Настройки/ })
    await dialog.getByRole('button', { name: 'Редактор' }).click()
    const themeSelect = dialog.locator('select').filter({ has: page.locator('option[value="system"]') })
    await themeSelect.selectOption('system')
    await dialog.getByRole('button', { name: 'Применить' }).click()
    await expect(dialog).toBeVisible()
    await dialog.getByRole('button', { name: 'Отмена' }).click()
  })

  test('2.4 large file shows simplified mode banner', async ({ page }) => {
    await bootApp(page, '?e2e=large-file')
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    await expect(page.locator('.status-bar .large-file-banner')).toBeVisible({ timeout: 15_000 })
  })

  test('2.5 Ctrl+F and Ctrl+H open find widget', async ({ page }) => {
    await bootApp(page)
    await createNewScenario(page)
    await page.locator('.workspace .monaco-editor .view-lines').click()
    await page.keyboard.press('Control+KeyF')
    await expect(page.locator('.monaco-editor .find-widget.visible')).toBeVisible()
    await page.keyboard.press('Escape')
    await page.keyboard.press('Control+KeyH')
    await expect(page.locator('.monaco-editor .find-widget.visible')).toBeVisible()
  })
})

test.describe('3. Live Recording', () => {
  test('3.1 recording shows target in status bar', async ({ page }) => {
    await bootApp(page, '?e2e=post-record')
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    await startRecordingFromDialog(page)
    await expect(page.locator('.status-bar .recording-target')).toContainText('Запись → smoke.feature')
  })

  test('3.2 tab switch during recording asks confirmation', async ({ page }) => {
    await bootApp(page, '?e2e=post-record')
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    await createNewScenario(page)
    await page.locator('.editor-tab.file', { hasText: 'smoke' }).click()
    await startRecordingFromDialog(page)
    await page.locator('.editor-tab.file', { hasText: 'novyy-scenariy' }).first().click()
    const dialog = page.getByRole('alertdialog', { name: 'Запись активна' })
    await expect(dialog).toBeVisible()
    await expect(dialog.getByText('Больше не спрашивать')).toBeVisible()
    await dialog.getByRole('button', { name: 'Отмена' }).click()
  })

  test('3.3 idle timeout logs to journal', async ({ page }) => {
    await bootApp(page, '?e2e=record-idle')
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    await startRecordingFromDialog(page)
    await openJournal(page)
    await expect(page.locator('.panel-body.text-panel')).toContainText('нет действий', { timeout: 10_000 })
  })

  test('3.4 stop button label during recording', async ({ page }) => {
    await bootApp(page, '?e2e=post-record')
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    await startRecordingFromDialog(page)
    await expect(page.getByRole('button', { name: 'Стоп запись' })).toBeVisible()
    await stopRecording(page)
  })

  test('3.5 headless toggle in recording bar asks confirm', async ({ page }) => {
    await bootApp(page, '?e2e=post-record')
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    await startRecordingFromDialog(page)
    const headless = page.locator('.recording-bar label.check-inline', { hasText: 'Без окна браузера' }).locator('input')
    await headless.click()
    const dialog = page.getByRole('alertdialog', { name: 'Headless' })
    await expect(dialog).toBeVisible()
    await dialog.getByRole('button', { name: 'Отмена' }).click()
  })
})

test.describe('4. Run', () => {
  test('4.1 suite run shows N/M in playing bar', async ({ page }) => {
    await bootApp(page, '?e2e=run-progress')
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    await page.keyboard.press('Control+Enter')
    await expect(page.locator('.playing-bar .play-progress-text')).toContainText(/\d\/3/, { timeout: 10_000 })
  })

  test('4.2 journal streams stdout during run', async ({ page }) => {
    await bootApp(page, '?e2e=run-stream')
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    await page.keyboard.press('Control+Enter')
    await openJournal(page)
    await expect(page.locator('.panel-body.text-panel')).toContainText('E2E_STREAM_LINE', { timeout: 10_000 })
  })

  test('4.3 first Ctrl+Enter opens run dialog', async ({ page }) => {
    await bootApp(page)
    await createNewScenario(page)
    await page.keyboard.press('Control+Enter')
    await expect(page.getByRole('dialog', { name: 'Запуск сценария' })).toBeVisible()
  })

  test('4.4 playing bar cancel stops run', async ({ page }) => {
    await bootApp(page, '?e2e=run-cancel')
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    await page.keyboard.press('Control+Enter')
    await expect(page.locator('.playing-bar')).toBeVisible({ timeout: 10_000 })
    await page.locator('.playing-bar .play-cancel-btn').click()
    await expect(page.locator('.playing-bar .play-cancel', { hasText: 'Останавливаем' })).toBeVisible()
    await expect(page.locator('.playing-bar')).toBeHidden({ timeout: 15_000 })
  })

  test('4.5 dry-run keeps editor available', async ({ page }) => {
    await bootApp(page)
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    await openMenuItem(page, 'Запись и тест', 'Запустить…')
    const dialog = page.getByRole('dialog', { name: 'Запуск сценария' })
    await dialog.locator('label.check-row', { hasText: 'Dry-run' }).locator('input').check()
    await dialog.getByRole('button', { name: 'Запустить' }).click()
    await expect(dialog).toBeHidden({ timeout: 10_000 })
    await expect(page.locator('.playing-bar')).toHaveCount(0)
    await typeInEditor(page, '# dry-run-edit')
    await expect(page.locator('.monaco-editor .view-line', { hasText: 'dry-run-edit' })).toBeVisible()
  })

  test('4.6 run dialog continue-on-fail option', async ({ page }) => {
    await bootApp(page)
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    await openMenuItem(page, 'Запись и тест', 'Запустить…')
    const dialog = page.getByRole('dialog', { name: 'Запуск сценария' })
    await dialog.locator('label.check-row', { hasText: 'Продолжать при ошибке' }).locator('input').check()
    await dialog.getByRole('button', { name: 'Запустить' }).click()
    await expect(dialog).toBeHidden({ timeout: 10_000 })
    const req = await page.evaluate(() => (window as unknown as { __e2eLastRunRequest?: () => { continueOnFail?: boolean } }).__e2eLastRunRequest?.())
    expect(req?.continueOnFail).toBe(true)
  })
})

test.describe('5. Results', () => {
  test('5.1 failed step button navigates in editor', async ({ page }) => {
    await bootApp(page, '?e2e=flaky-run')
    await openTestProject(page)
    await openResults(page)
    await page.getByRole('button', { name: 'Шаг 2' }).first().click()
    await expect(page.locator('.editor-tab.file .tab-label', { hasText: 'smoke.feature' })).toBeVisible({
      timeout: 10_000,
    })
  })

  test('5.2 trace viewer button works', async ({ page }) => {
    await bootApp(page, '?e2e=trace-artifacts')
    await openTestProject(page)
    await openResults(page)
    await page.getByRole('button', { name: 'Trace viewer' }).click()
    await openJournal(page)
    await expect(page.locator('.panel-body.text-panel')).toContainText('Trace viewer:', { timeout: 10_000 })
  })

  test('5.3 flaky badge runs scenario three times', async ({ page }) => {
    await bootApp(page, '?e2e=flaky-run')
    await openTestProject(page)
    await openResults(page)
    await page.getByRole('button', { name: 'Запустить 3×' }).first().click()
    await expect(page.locator('.playing-bar')).toBeHidden({ timeout: 10_000 })
    await openJournal(page)
    await expect(page.locator('.panel-body.text-panel')).toContainText('Прогон 3/3', { timeout: 10_000 })
  })

  test('5.4 allure missing shows install hint', async ({ page }) => {
    await bootApp(page, '?e2e=allure-missing')
    await openTestProject(page)
    await openResults(page)
    await expect(page.getByRole('button', { name: 'Allure не найден' })).toBeVisible()
  })

  test('5.5 double click result row opens feature', async ({ page }) => {
    await bootApp(page, '?e2e=flaky-run')
    await openTestProject(page)
    await openResults(page)
    await page.locator('.results-table tbody tr').first().dblclick()
    await expect(page.locator('.editor-tab.file .tab-label', { hasText: 'smoke.feature' })).toBeVisible({
      timeout: 10_000,
    })
  })
})

test.describe('6. Settings', () => {
  test('6.1 reset defaults restores browser', async ({ page }) => {
    await bootApp(page)
    await openSettings(page)
    const dialog = page.getByRole('dialog', { name: /Настройки/ })
    await dialog.locator('select').first().selectOption('firefox')
    await dialog.getByRole('button', { name: 'Сбросить по умолчанию' }).click()
    await expect(dialog.locator('select').first()).toHaveValue('chromium')
  })

  test('6.2 navWaitUntil visible in settings', async ({ page }) => {
    await bootApp(page)
    await openSettings(page)
    const dialog = page.getByRole('dialog', { name: /Настройки/ })
    const navSelect = dialog.locator('select').filter({ has: page.locator('option[value="networkidle"]') })
    await navSelect.selectOption('networkidle')
    await dialog.getByRole('button', { name: 'Применить' }).click()
    await page.getByRole('button', { name: 'Отмена' }).click()
    await openSettings(page)
    await expect(navSelect).toHaveValue('networkidle')
  })

  test('6.3 extreme workers and slowMo shows warning', async ({ page }) => {
    await bootApp(page)
    await openSettings(page)
    const dialog = page.getByRole('dialog', { name: /Настройки/ })
    await dialog.locator('.setting-card', { hasText: 'Параллельные воркеры' }).locator('input.setting-number').fill('16')
    await dialog.locator('.setting-card', { hasText: 'slow-mo' }).locator('input.setting-number').fill('5000')
    await expect(dialog.getByText('Много воркеров и высокий slow-mo')).toBeVisible()
  })

  test('6.4 apply keeps settings dialog open', async ({ page }) => {
    await bootApp(page)
    await openSettings(page)
    const dialog = page.getByRole('dialog', { name: /Настройки/ })
    await dialog.getByRole('button', { name: 'Применить' }).click()
    await expect(dialog).toBeVisible()
  })

  test('6.5 Ctrl+S ignored while settings open', async ({ page }) => {
    await bootApp(page)
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    await typeInEditor(page, '# unsaved')
    await expect(page.locator('.editor-tab.file .tab-label', { hasText: /\*/ })).toBeVisible()
    await openSettings(page)
    await page.keyboard.press('Control+KeyS')
    await expect(page.getByRole('dialog', { name: /Настройки/ })).toBeVisible()
    await page.getByRole('button', { name: 'Отмена' }).click()
    await expect(page.locator('.editor-tab.file .tab-label', { hasText: 'smoke.feature *' })).toBeVisible()
  })
})
