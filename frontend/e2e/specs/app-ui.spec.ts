import { expect, test } from '@playwright/test'
import {
  bootApp,
  catalogFeature,
  createNewScenario,
  editorLine,
  mockPath,
  openMenuItem,
  openTestProject,
  startRecordingFromDialog,
  statusMessage,
  stopRecording,
} from '../helpers/app'

test.beforeEach(async ({ page }) => {
  await page.addInitScript({ path: mockPath })
})

test('app boots to welcome screen', async ({ page }) => {
  await bootApp(page)
  await expect(page.getByRole('button', { name: /Старт/i }).first()).toBeVisible()
})

test('steps help dialog lists catalog entries', async ({ page }) => {
  await bootApp(page)
  await page.keyboard.press('F1')
  const dialog = page.getByRole('dialog', { name: 'Справка по шагам' })
  await expect(dialog).toBeVisible()
  await expect(dialog.getByText('нажимаю (click)')).toBeVisible()
  await expect(dialog.getByText('Пример', { exact: true })).toBeVisible()
  await expect(dialog.getByText('нажимаю "button.submit"')).toBeVisible()
})

test('command palette opens with Ctrl+Shift+P', async ({ page }) => {
  await bootApp(page)
  await page.keyboard.press('Control+Shift+KeyP')
  await expect(page.getByRole('dialog', { name: 'Палитра команд' })).toBeVisible()
  await expect(page.getByPlaceholder('Введите команду…')).toBeVisible()
})

test('settings dialog opens with Ctrl+Comma', async ({ page }) => {
  await bootApp(page)
  await page.keyboard.press('Control+Comma')
  await expect(page.getByRole('dialog', { name: /Настройки/ })).toBeVisible()
  await expect(page.getByPlaceholder('Поиск настроек')).toBeVisible()
})

test('settings browser engine status updates on selection', async ({ page }) => {
  await bootApp(page)
  await page.keyboard.press('Control+Comma')
  const dialog = page.getByRole('dialog', { name: /Настройки/ })
  await expect(dialog).toBeVisible()
  await expect(dialog.getByText('Chromium: установлен')).toBeVisible()
  await dialog.locator('select').first().selectOption('firefox')
  await expect(dialog.getByText('Firefox: установлен')).toBeVisible()
  await expect(dialog.getByRole('button', { name: 'Переустановить' })).toBeVisible()
})

test('new scenario creates untitled tab', async ({ page }) => {
  await bootApp(page)
  await createNewScenario(page)
})

test('untitled tab restores after reload when project is open', async ({ page }) => {
  await bootApp(page)
  await openTestProject(page)
  await createNewScenario(page)
  const marker = 'E2E_SESSION_RESTORE_MARKER'
  await page.locator('.monaco-editor .view-lines').click()
  await page.keyboard.press('End')
  await page.keyboard.press('Enter')
  await page.keyboard.type(marker)
  await page.waitForTimeout(600)
  await expect(page.locator('.monaco-editor .view-line', { hasText: marker })).toBeVisible()
  await page.reload()
  await expect(page.locator('.ide')).toBeVisible({ timeout: 20_000 })
  await expect(page.locator('.editor-tab.file .tab-label', { hasText: 'novyy-scenariy.feature' })).toBeVisible({
    timeout: 10_000,
  })
  await expect(page.locator('.monaco-editor .view-line', { hasText: marker })).toBeVisible({ timeout: 10_000 })
})

test('hotkeys dialog opens with Shift+F1', async ({ page }) => {
  await bootApp(page)
  await page.keyboard.press('Shift+F1')
  const dialog = page.getByRole('dialog', { name: 'Горячие клавиши' })
  await expect(dialog).toBeVisible()
  await expect(dialog.getByText('Ctrl+Shift+P')).toBeVisible()
  await expect(dialog.getByText('Палитра команд')).toBeVisible()
})

test('monaco find widget opens with Ctrl+H', async ({ page }) => {
  await bootApp(page)
  await createNewScenario(page)
  await page.keyboard.press('Control+KeyH')
  await expect(page.locator('.monaco-editor .find-widget.visible')).toBeVisible()
})

test('unsaved close dialog appears when closing dirty tab', async ({ page }) => {
  await bootApp(page)
  await createNewScenario(page)
  await page.locator('.editor-tab.file .tab-close').click()
  const dialog = page.getByRole('dialog', { name: 'Несохранённые изменения' })
  await expect(dialog).toBeVisible()
  await expect(dialog.getByText('novyy-scenariy.feature')).toBeVisible()
  await dialog.getByRole('button', { name: 'Отмена' }).click()
  await expect(dialog).toBeHidden()
  await expect(page.locator('.editor-tab.file .tab-label', { hasText: 'novyy-scenariy.feature' })).toBeVisible()
})

test('about dialog shows version from mock', async ({ page }) => {
  await bootApp(page)
  await openMenuItem(page, 'Справка', 'О программе')
  const dialog = page.getByRole('dialog', { name: 'О программе' })
  await expect(dialog).toBeVisible()
  await expect(dialog.getByText('e2e-test')).toBeVisible()
})

test('open project enables validate dialog', async ({ page }) => {
  await bootApp(page)
  await openTestProject(page)
  await createNewScenario(page)
  await openMenuItem(page, 'Запись и тест', 'Проверить…')
  const dialog = page.getByRole('dialog', { name: 'Проверка сценария' })
  await expect(dialog).toBeVisible()
  await expect(dialog.getByText('Только синтаксис (без браузера)')).toBeVisible()
  await dialog.getByRole('button', { name: 'Отмена' }).click()
  await expect(dialog).toBeHidden()
})

test('run dialog opens from menu', async ({ page }) => {
  await bootApp(page)
  await createNewScenario(page)
  await openMenuItem(page, 'Запись и тест', 'Запустить…')
  const dialog = page.getByRole('dialog', { name: 'Запуск сценария' })
  await expect(dialog).toBeVisible()
  await expect(dialog.getByText('Dry-run (без браузера)')).toBeVisible()
  await dialog.getByRole('button', { name: 'Отмена' }).click()
  await expect(dialog).toBeHidden()
})

test('snippet palette lists steps from catalog', async ({ page }) => {
  await bootApp(page)
  await createNewScenario(page)
  await openMenuItem(page, 'Сценарий', 'Палитра сниппетов…')
  const dialog = page.getByRole('dialog', { name: 'Палитра сниппетов' })
  await expect(dialog).toBeVisible()
  await expect(dialog.getByPlaceholder('Поиск шага…')).toBeVisible()
  await expect(dialog.getByText('нажимаю "button.submit"')).toBeVisible()
})

test('journal panel opens with Ctrl+Backquote', async ({ page }) => {
  await bootApp(page)
  await page.keyboard.press('Control+Backquote')
  await expect(page.locator('.bottom-panel')).toBeVisible()
  await expect(page.locator('.panel-tab.active', { hasText: 'Журнал' })).toBeVisible()
})

test('validate confirm runs syntax check and logs result', async ({ page }) => {
  await bootApp(page)
  await openTestProject(page)
  await openMenuItem(page, 'Запись и тест', 'Проверить…')
  const dialog = page.getByRole('dialog', { name: 'Проверка сценария' })
  await dialog.getByRole('button', { name: 'Проверить' }).click()
  await expect(dialog).toBeHidden({ timeout: 10_000 })
  await expect(page.locator('.bottom-panel')).toBeVisible()
  await expect(page.locator('.panel-body')).toContainText('Проверка завершена', { timeout: 10_000 })
})

test('plugins dialog opens from menu when project is open', async ({ page }) => {
  await bootApp(page)
  await openTestProject(page)
  await openMenuItem(page, 'Плагины', 'Управление плагинами…')
  const dialog = page.getByRole('dialog', { name: 'Плагины' })
  await expect(dialog).toBeVisible()
  await expect(dialog.getByText('Плагины проекта')).toBeVisible()
  await expect(dialog.getByText('Нет установленных плагинов')).toBeVisible()
  await dialog.getByRole('button', { name: 'Закрыть' }).click()
  await expect(dialog).toBeHidden()
})

test('catalog lists project feature after open', async ({ page }) => {
  await bootApp(page)
  await openTestProject(page)
  await expect(page.locator('.tree-file-label', { hasText: 'smoke' })).toBeVisible()
})

test('settings shows install button when browser is missing', async ({ page }) => {
  await bootApp(page, '?e2e=missing-browser')
  await page.keyboard.press('Control+Comma')
  const dialog = page.getByRole('dialog', { name: /Настройки/ })
  await expect(dialog.getByText('Chromium: не установлен')).toBeVisible()
  await expect(dialog.getByRole('button', { name: 'Установить движок' })).toBeVisible()
})

test('export dialog shows preview for current scenario', async ({ page }) => {
  await bootApp(page)
  await createNewScenario(page)
  await openMenuItem(page, 'Сценарий', 'Экспорт…')
  const dialog = page.getByRole('dialog', { name: 'Экспорт сценария' })
  await expect(dialog).toBeVisible()
  await expect(dialog.getByText('2 шаг(ов)')).toBeVisible()
  await dialog.getByRole('button', { name: 'Отмена' }).click()
  await expect(dialog).toBeHidden()
})

test('import features dialog picks and imports files', async ({ page }) => {
  await bootApp(page, '?e2e=import-pick')
  await openTestProject(page)
  await openMenuItem(page, 'Сценарий', 'Импорт .feature…')
  const dialog = page.getByRole('dialog', { name: 'Импорт feature' })
  await expect(dialog).toBeVisible()
  await dialog.getByRole('button', { name: 'Добавить файлы…' }).click()
  await expect(dialog.getByText('sample.feature')).toBeVisible()
  await dialog.getByRole('button', { name: 'Импортировать' }).click()
  await expect(dialog).toBeHidden({ timeout: 10_000 })
  await page.keyboard.press('Control+Backquote')
  await expect(page.locator('.panel-body')).toContainText('Импортировано файлов: 1', { timeout: 10_000 })
})

test('unsaved close discards dirty tab', async ({ page }) => {
  await bootApp(page)
  await createNewScenario(page)
  await page.locator('.editor-tab.file .tab-close').click()
  const dialog = page.getByRole('dialog', { name: 'Несохранённые изменения' })
  await dialog.getByRole('button', { name: 'Не сохранять' }).click()
  await expect(dialog).toBeHidden()
  await expect(page.locator('.editor-tab.file')).toHaveCount(0)
})

test('export confirm logs success to journal', async ({ page }) => {
  await bootApp(page)
  await createNewScenario(page)
  await openMenuItem(page, 'Сценарий', 'Экспорт…')
  const dialog = page.getByRole('dialog', { name: 'Экспорт сценария' })
  await expect(dialog.getByRole('button', { name: 'Экспорт' })).toBeEnabled()
  await dialog.getByRole('button', { name: 'Экспорт' }).click()
  await expect(dialog).toBeHidden({ timeout: 10_000 })
  await page.keyboard.press('Control+Backquote')
  await expect(page.locator('.panel-body')).toContainText('Экспортировано', { timeout: 10_000 })
})

test('run dialog dry-run logs completion', async ({ page }) => {
  await bootApp(page)
  await openTestProject(page)
  await page.locator('.tree-file-label', { hasText: 'smoke' }).click()
  await openMenuItem(page, 'Запись и тест', 'Запустить…')
  const dialog = page.getByRole('dialog', { name: 'Запуск сценария' })
  await dialog.locator('label.check-row').filter({ hasText: 'Dry-run' }).locator('input').check()
  await dialog.getByRole('button', { name: 'Запустить' }).click()
  await expect(dialog).toBeHidden({ timeout: 10_000 })
  await expect(page.locator('.bottom-panel')).toBeVisible()
  await page.locator('.panel-tab', { hasText: 'Журнал' }).click()
  await expect(page.locator('.panel-body')).toContainText('Dry-run…', { timeout: 10_000 })
  await expect(page.locator('.panel-body')).toContainText('Завершено.', { timeout: 10_000 })
})

test('live recording inserts steps into editor', async ({ page }) => {
  await bootApp(page, '?e2e=post-record')
  await openTestProject(page)
  await catalogFeature(page, 'smoke').click()
  await startRecordingFromDialog(page)
  await expect(editorLine(page, 'нажимаю "#login"')).toBeVisible({
    timeout: 10_000,
  })
  await expect(editorLine(page, 'ввожу "user" в "#email"')).toBeVisible()
  await expect(editorLine(page, 'нажимаю "#submit"')).toBeVisible()

  await editorLine(page, 'нажимаю "#login"').click()
  await expect(page.locator('.monaco-editor .squiggly-warning')).toBeVisible({ timeout: 10_000 })

  await page.keyboard.press('Control+End')
  await page.keyboard.press('Enter')
  await page.keyboard.insertText('\tКогда ')
  await page.keyboard.press('Control+Space')
  await expect(page.locator('.monaco-editor .suggest-widget')).toBeVisible({ timeout: 10_000 })
  await expect(page.locator('.monaco-list-row', { hasText: 'нажимаю' }).first()).toBeVisible()
})

test('Ctrl+S saves feature from editor', async ({ page }) => {
  await bootApp(page)
  await openTestProject(page)
  await catalogFeature(page, 'smoke').click()
  await expect(editorLine(page, 'тест')).toBeVisible({ timeout: 20_000 })
  await page.locator('.workspace .monaco-editor .view-lines').click()
  await page.keyboard.press('End')
  await page.keyboard.press('Enter')
  await page.keyboard.type('    # dirty marker')
  await expect(page.locator('.editor-tab.file .tab-label', { hasText: 'smoke.feature *' })).toBeVisible()
  await page.keyboard.press('Control+KeyS')
  await expect(statusMessage(page)).toHaveText('Сохранено', { timeout: 10_000 })
  await expect(page.locator('.editor-tab.file .tab-label', { hasText: 'smoke.feature *' })).toHaveCount(0)
})

test('Ctrl+Shift+O opens symbol outline in editor', async ({ page }) => {
  await bootApp(page)
  await openTestProject(page)
  await catalogFeature(page, 'smoke').click()
  await expect(editorLine(page, 'тест')).toBeVisible({ timeout: 20_000 })
  await page.locator('.workspace .monaco-editor .view-lines').click()
  await page.keyboard.press('Control+Shift+KeyO')
  await expect(page.locator('.quick-input-widget .monaco-list-row')).not.toHaveCount(0, {
    timeout: 10_000,
  })
})

test('post-record diff compares baseline and recorded steps', async ({ page }) => {
  await bootApp(page, '?e2e=post-record-diff')
  await openTestProject(page)
  await catalogFeature(page, 'smoke').click()
  await startRecordingFromDialog(page)
  await expect(editorLine(page, 'нажимаю "#login"')).toBeVisible({ timeout: 10_000 })
  await stopRecording(page)
  await expect(page.getByRole('status').filter({ hasText: 'Записано шагов' })).toBeVisible({
    timeout: 10_000,
  })
  await page.getByRole('button', { name: 'Сравнить' }).click()
  const dialog = page.getByRole('dialog', { name: 'Изменения после записи' })
  await expect(dialog).toBeVisible()
  await expect(dialog.locator('.diff-host .view-lines').first()).toBeVisible({ timeout: 10_000 })
  expect(await dialog.locator('.diff-host .monaco-editor').count()).toBeGreaterThanOrEqual(2)
})

test('run history flaky filter shows unstable scenarios', async ({ page }) => {
  await bootApp(page, '?e2e=flaky-run')
  await openTestProject(page)
  await openMenuItem(page, 'Запись и тест', 'История запусков…')
  const dialog = page.getByRole('dialog', { name: 'История запусков' })
  await expect(dialog).toBeVisible()
  await expect(dialog.locator('.flaky-tag').first()).toBeVisible()
  await dialog.locator('select').selectOption('flaky')
  await expect(dialog.locator('tbody tr')).toHaveCount(2)
  await expect(dialog.locator('.flaky-tag')).toHaveCount(2)
  await dialog.getByRole('button', { name: 'Закрыть' }).click()
  await expect(dialog).toBeHidden()
  await openMenuItem(page, 'Вид', 'Результаты')
  await expect(page.locator('.results-panel .flaky-tag').first()).toBeVisible()
})
