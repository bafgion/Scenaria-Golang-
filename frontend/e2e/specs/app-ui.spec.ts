import { expect, test } from '@playwright/test'
import {
  bootApp,
  catalogFeature,
  createNewScenario,
  editorLine,
  flushSessionForE2E,
  mockPath,
  openMenuItem,
  openTestProject,
  startRecordingFromDialog,
  statusMessage,
  stopRecording,
  typeInEditor,
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
  await typeInEditor(page, marker)
  await expect(page.locator('.monaco-editor .view-line', { hasText: marker })).toBeVisible()
  await flushSessionForE2E(page)
  await page.reload()
  await expect(page.locator('.ide')).toBeVisible({ timeout: 20_000 })
  await expect(page.locator('.editor-tab.file .tab-label', { hasText: 'novyy-scenariy.feature' })).toBeVisible({
    timeout: 15_000,
  })
  const untitledTab = page.locator('.editor-tab.file', { hasText: 'novyy-scenariy' })
  await untitledTab.click()
  await expect(page.locator('.feature-workspace:not(.hidden) .monaco-editor')).toBeVisible({ timeout: 15_000 })
  await expect(page.locator('.monaco-editor .view-line', { hasText: marker })).toBeVisible({ timeout: 20_000 })
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

test('closing active tab activates the previous Monaco model', async ({ page }) => {
  await bootApp(page)
  await createNewScenario(page)
  await typeInEditor(page, 'ACTIVE_TAB_ONE_MARKER')
  await createNewScenario(page)
  await typeInEditor(page, 'ACTIVE_TAB_TWO_MARKER')

  await expect(page.locator('.editor-tab.file')).toHaveCount(2)
  await page.locator('.editor-tab.file').nth(1).locator('.tab-close').click()

  const dialog = page.getByRole('dialog', { name: 'Несохранённые изменения' })
  await expect(dialog).toBeVisible()
  await dialog.getByRole('button', { name: 'Не сохранять' }).click()
  await expect(dialog).toBeHidden()

  await expect(page.locator('.editor-tab.file')).toHaveCount(1)
  await expect(page.locator('.editor-tab.file .tab-label')).toContainText('novyy-scenariy')
  await expect(page.locator('.monaco-editor .view-line', { hasText: 'ACTIVE_TAB_ONE_MARKER' })).toBeVisible()
  await expect(page.locator('.monaco-editor .view-line', { hasText: 'ACTIVE_TAB_TWO_MARKER' })).toHaveCount(0)
})

test('closing welcome tab activates the last feature tab', async ({ page }) => {
  await bootApp(page)
  await createNewScenario(page)
  await typeInEditor(page, 'WELCOME_SWITCH_ONE')
  await createNewScenario(page)
  await typeInEditor(page, 'WELCOME_SWITCH_TWO')

  await openMenuItem(page, 'Вид', 'Старт')
  await expect(page.locator('.editor-tab.welcome')).toHaveAttribute('aria-selected', 'true')
  await page.locator('.editor-tab.welcome .tab-close').click()

  await expect(page.locator('.editor-tab.welcome')).toHaveCount(0)
  await expect(page.locator('.monaco-editor .view-line', { hasText: 'WELCOME_SWITCH_TWO' })).toBeVisible()
  await expect(page.locator('.monaco-editor .view-line', { hasText: 'WELCOME_SWITCH_ONE' })).toHaveCount(0)
})

test('undo and redo stay isolated per Monaco tab model', async ({ page }) => {
  await bootApp(page)
  await createNewScenario(page)
  await typeInEditor(page, 'UNDO_REDO_TAB_ONE')
  await createNewScenario(page)
  await typeInEditor(page, 'UNDO_REDO_TAB_TWO')

  const tabs = page.locator('.editor-tab.file')
  await tabs.nth(0).click()
  await expect(page.locator('.monaco-editor .view-line', { hasText: 'UNDO_REDO_TAB_ONE' })).toBeVisible()
  await page.locator('.workspace .monaco-editor .view-lines').click()
  await page.keyboard.press('Control+KeyA')
  await page.keyboard.type('UNDO_REDO_TAB_ONE_REPLACED')
  await expect(page.locator('.monaco-editor .view-line', { hasText: 'UNDO_REDO_TAB_ONE_REPLACED' })).toBeVisible()
  await page.keyboard.press('Control+KeyZ')
  await expect(page.locator('.monaco-editor .view-line', { hasText: 'UNDO_REDO_TAB_ONE' })).toBeVisible()
  await page.keyboard.press('Control+Shift+KeyZ')
  await expect(page.locator('.monaco-editor .view-line', { hasText: 'UNDO_REDO_TAB_ONE_REPLACED' })).toBeVisible()

  await tabs.nth(1).click()
  await expect(page.locator('.monaco-editor .view-line', { hasText: 'UNDO_REDO_TAB_TWO' })).toBeVisible()
  await expect(page.locator('.monaco-editor .view-line', { hasText: 'UNDO_REDO_TAB_ONE_REPLACED' })).toHaveCount(0)
})

test('forceActivate reloads an already active tab without changing selection', async ({ page }) => {
  await bootApp(page, '?e2e=force-activate')
  await openTestProject(page)
  await catalogFeature(page, 'smoke').click()
  const smokeTab = page.locator('.editor-tab.file', { hasText: 'smoke.feature' }).first()
  await expect(smokeTab).toBeVisible({ timeout: 20_000 })
  const smokePath = await smokeTab.getAttribute('title')
  expect(smokePath).toBeTruthy()

  await page.evaluate(async (path) => {
    const hook = (window as unknown as { __e2eLoadFeature?: (path: string, forceActivate?: boolean) => Promise<void> }).__e2eLoadFeature
    if (hook && path) await hook(path, true)
  }, smokePath)

  await expect(page.locator('.editor-tab.file.active .tab-label', { hasText: 'smoke.feature' })).toBeVisible()
  await expect(page.locator('.editor-tab.file')).toHaveCount(1)
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

test('first Ctrl+Enter opens run dialog before options confirmed', async ({ page }) => {
  await bootApp(page)
  await createNewScenario(page)
  await page.keyboard.press('Control+Enter')
  const dialog = page.getByRole('dialog', { name: 'Запуск сценария' })
  await expect(dialog).toBeVisible()
  await dialog.getByRole('button', { name: 'Отмена' }).click()
  await expect(dialog).toBeHidden()
})

test('settings reset defaults restores browser field', async ({ page }) => {
  await bootApp(page)
  await page.keyboard.press('Control+Comma')
  const dialog = page.getByRole('dialog', { name: /Настройки/ })
  await dialog.locator('select').first().selectOption('firefox')
  await dialog.getByRole('button', { name: 'Сбросить по умолчанию' }).click()
  await expect(dialog.locator('select').first()).toHaveValue('chromium')
  await dialog.getByRole('button', { name: 'Отмена' }).click()
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

test('http auth dialog closes over record dialog', async ({ page }) => {
  await bootApp(page)
  await openTestProject(page)
  await page.keyboard.press('Control+KeyR')
  const recordDialog = page.getByRole('dialog', { name: 'Запись сценария' })
  await expect(recordDialog).toBeVisible()
  await recordDialog.getByRole('button', { name: 'HTTP Auth…' }).click()
  const authDialog = page.getByRole('dialog', { name: 'HTTP авторизация' })
  await expect(authDialog).toBeVisible()
  await authDialog.getByRole('button', { name: 'Закрыть' }).click()
  await expect(authDialog).toBeHidden()
  await expect(recordDialog).toBeVisible()
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

test('stale validation response is ignored after text changes', async ({ page }) => {
  await bootApp(page, '?e2e=validation-race')
  await createNewScenario(page)
  await page.locator('.workspace .monaco-editor .view-lines').click()
  await page.keyboard.press('Control+A')
  await page.keyboard.type('VALIDATION_STALE_BAD')
  await page.waitForTimeout(400)
  await page.locator('.workspace .monaco-editor .view-lines').click()
  await page.keyboard.press('Control+A')
  await page.keyboard.type('VALIDATION_STALE_GOOD')
  await page.waitForTimeout(1600)
  await expect(page.locator('.monaco-editor .squiggly-error')).toHaveCount(0)
})

test('stale editor change event does not switch the active tab', async ({ page }) => {
  await bootApp(page, '?e2e=editor-change-race')
  await openTestProject(page)
  await catalogFeature(page, 'smoke').click()

  const smokeTab = page.locator('.editor-tab.file', { hasText: 'smoke.feature' }).first()
  await expect(smokeTab).toBeVisible({ timeout: 20_000 })
  const smokePath = await smokeTab.getAttribute('title')
  expect(smokePath).toBeTruthy()

  await createNewScenario(page)
  await typeInEditor(page, 'EDITOR_CHANGE_KEEP_ACTIVE')

  await page.evaluate((path) => {
    const hook = (window as unknown as { __e2eEmitEditorChange?: (path: string | null, text: string) => void }).__e2eEmitEditorChange
    if (hook) hook(path, 'STALE_EDITOR_CHANGE_EVENT')
  }, smokePath)

  await expect(page.locator('.editor-tab.file.active .tab-label', { hasText: 'novyy-scenariy.feature' })).toBeVisible()
  await expect(page.locator('.monaco-editor .view-line', { hasText: 'EDITOR_CHANGE_KEEP_ACTIVE' })).toBeVisible()
  await expect(page.locator('.monaco-editor .view-line', { hasText: 'STALE_EDITOR_CHANGE_EVENT' })).toHaveCount(0)
})

test('batch run clears stale filters from previous single run', async ({ page }) => {
  await bootApp(page, '?e2e=batch-race')
  await openTestProject(page)
  await page.locator('.catalog-tree .catalog-tree-row.file[title^="smoke.feature"]').click()
  await openMenuItem(page, 'Запись и тест', 'Запустить…')
  const runDialog = page.getByRole('dialog', { name: 'Запуск сценария' })
  await expect(runDialog).toBeVisible()
  await runDialog.getByLabel('Тег').fill('@stale')
  await runDialog.getByLabel('Сценарий (опционально)').fill('stale-single')
  await runDialog.getByRole('button', { name: 'Запустить' }).click()
  await expect(runDialog).toBeHidden({ timeout: 10_000 })

  await page.locator('.explorer-tool-actions .batch-toggle').click()
  await page.keyboard.press('Control+Enter')

  await page.waitForFunction(() => {
    const hook = (window as unknown as { __e2eLastRunRequest?: () => any }).__e2eLastRunRequest
    const req = hook?.()
    return !!req && Array.isArray(req.targets) && req.targets.length > 1
  })

  const request = await page.evaluate(() => (window as unknown as { __e2eLastRunRequest?: () => any }).__e2eLastRunRequest?.())
  expect(request.tag).toBe('')
  expect(request.scenario).toBe('')
  expect(request.startStep).toBe(-1)
  expect(request.endStep).toBe(-1)
  expect(request.targets).toHaveLength(3)
})

test('save completion does not mutate another active tab', async ({ page }) => {
  await bootApp(page, '?e2e=save-race')
  await openTestProject(page)
  await catalogFeature(page, 'smoke').click()
  await expect(editorLine(page, 'тест')).toBeVisible({ timeout: 20_000 })
  await typeInEditor(page, 'SMOKE_SAVE_STALE')
  await createNewScenario(page)
  await typeInEditor(page, 'UNTITLED_SAVE_STALE')
  await page.locator('.editor-tab.file .tab-label', { hasText: 'smoke.feature' }).click()
  await page.keyboard.press('Control+KeyS')
  await page.locator('.editor-tab.file .tab-label', { hasText: 'novyy-scenariy.feature' }).click()
  await expect(statusMessage(page)).toHaveText('Сохранено', { timeout: 10_000 })
  await expect(page.locator('.editor-tab.file .tab-label', { hasText: 'smoke.feature *' })).toHaveCount(0)
  await expect(page.locator('.editor-tab.file .tab-label', { hasText: 'novyy-scenariy.feature *' })).toBeVisible()
})

test('disk reload completion does not mutate another active tab', async ({ page }) => {
  await bootApp(page, '?e2e=reload-race')
  await openTestProject(page)
  await catalogFeature(page, 'smoke').click()
  await page.waitForTimeout(1200)
  await createNewScenario(page)
  await typeInEditor(page, 'UNTITLED_RELOAD_KEEP')
  await page.locator('.editor-tab.file .tab-label', { hasText: 'smoke.feature' }).click()

  await page.evaluate(() => {
    const hook = (window as unknown as { __e2eCheckActiveTabDiskStale?: () => Promise<void> }).__e2eCheckActiveTabDiskStale
    if (hook) void hook()
  })
  await page.waitForTimeout(400)
  await page.locator('.editor-tab.file .tab-label', { hasText: 'novyy-scenariy.feature' }).click()

  const dialog = page.getByRole('alertdialog', { name: 'Файл изменён на диске' })
  await expect(dialog).toBeVisible({ timeout: 10_000 })
  await dialog.getByRole('button', { name: 'Перезагрузить' }).click()
  await expect(dialog).toBeHidden()
  await expect(page.locator('.monaco-editor .view-line', { hasText: 'UNTITLED_RELOAD_KEEP' })).toBeVisible()
  await expect(page.locator('.monaco-editor .view-line', { hasText: 'disk reload updated' })).toHaveCount(0)
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
  await page.getByRole('button', { name: 'Запустить 3×' }).first().click()
  await expect(page.locator('.playing-bar')).toBeHidden({ timeout: 10_000 })
  await openMenuItem(page, 'Вид', 'Журнал')
  await expect(page.locator('.panel-body.text-panel')).toContainText('Прогон 3/3', { timeout: 10_000 })
})

test('new project wizard picks folder and creates project', async ({ page }) => {
  await bootApp(page, '?e2e=new-project')
  await openMenuItem(page, 'Проект', 'Новый проект…')
  const dialog = page.getByRole('dialog', { name: 'Новый проект' })
  await expect(dialog).toBeVisible()
  await dialog.getByRole('button', { name: 'Обзор…' }).click()
  await expect(dialog.locator('input').first()).toHaveValue('C:/e2e/new-project')
  await dialog.getByRole('button', { name: 'Создать' }).click()
  await expect(dialog).toBeHidden({ timeout: 10_000 })
})

test('run progress mock updates playing bar counter', async ({ page }) => {
  await bootApp(page, '?e2e=run-progress')
  await openTestProject(page)
  await catalogFeature(page, 'smoke').click()
  await page.keyboard.press('Control+Enter')
  await expect(page.locator('.playing-bar .play-progress-text')).toContainText(/\d\/3/, { timeout: 10_000 })
})

test('results trace viewer button logs success', async ({ page }) => {
  await bootApp(page, '?e2e=trace-artifacts')
  await openTestProject(page)
  await openMenuItem(page, 'Вид', 'Результаты')
  await page.getByRole('button', { name: 'Trace viewer' }).click()
  await openMenuItem(page, 'Вид', 'Журнал')
  await expect(page.locator('.panel-body.text-panel')).toContainText('Trace viewer:', { timeout: 10_000 })
})
