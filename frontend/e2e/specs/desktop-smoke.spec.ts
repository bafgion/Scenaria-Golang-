/**
 * Desktop smoke (Windows + WebView2) — реальный scenaria-gui.exe, без wails-mock.
 *
 * Запускается из scripts/desktop-smoke.ps1 (CDP :9333).
 * Локально:
 *   wails build -platform windows/amd64
 *   ./scripts/desktop-smoke.ps1
 *
 * Покрывает UI-сценарии из app-ui.spec.ts / qa-daily-use, которые не требуют
 * mock-режимов (?e2e=…), эмуляции прогона или нативных диалогов ОС.
 */
import { expect, test } from '@playwright/test'
import {
  catalogFeature,
  createNewScenario,
  openJournal,
  openMenuItem,
  openRecordDialog,
  openSettings,
  statusMessage,
  typeInEditor,
} from '../helpers/app'
import {
  clickCatalogFeature,
  connectDesktop,
  disconnectDesktop,
  dismissBlockingDialogs,
  ensureExamplesProject,
  waitForAppReady,
  resetRunDialogConfirmed,
} from '../helpers/desktop'

test.describe.configure({ mode: 'serial', timeout: 180_000 })

test.beforeAll(async () => {
  const page = await connectDesktop()
  await waitForAppReady(page)
})

test.afterAll(async () => {
  await disconnectDesktop()
})

test('приложение загружается после splash', async () => {
  const page = await connectDesktop()
  await expect(page.locator('.ide')).toBeVisible()
  await expect(page.locator('.menubar')).toBeVisible({ timeout: 30_000 })
})

test('открыть примеры → каталог → сценарий в редакторе', async () => {
  const page = await connectDesktop()
  await ensureExamplesProject(page)
  await expect(statusMessage(page)).toContainText(/examples|Примеры/i, { timeout: 15_000 })

  const feature = page.locator('.catalog-tree .tree-file-label').first()
  await expect(feature).toBeVisible()
  const name = (await feature.textContent())?.replace(/^○\s*/, '').trim() || ''
  await clickCatalogFeature(page, 0)
  await expect(page.locator('.editor-tab.file .tab-label').first()).toBeVisible({ timeout: 15_000 })
  await expect(page.locator('.monaco-editor .view-lines')).toBeVisible({ timeout: 15_000 })
  if (name) {
    await expect(page.locator('.editor-tab.file .tab-label', { hasText: name })).toBeVisible()
  }
})

test('каталог показывает feature после открытия примеров', async () => {
  const page = await connectDesktop()
  await ensureExamplesProject(page)
  await expect(page.locator('.catalog-tree .tree-file-label').first()).toBeVisible()
})

test('проверка сценария пишет в журнал', async () => {
  const page = await connectDesktop()
  await ensureExamplesProject(page)
  await clickCatalogFeature(page, 0)
  await expect(page.locator('.monaco-editor')).toBeVisible({ timeout: 15_000 })

  await openMenuItem(page, 'Запись и тест', 'Проверить…')
  const dialog = page.getByRole('dialog', { name: 'Проверка сценария' })
  await expect(dialog).toBeVisible()
  await dialog.getByRole('button', { name: 'Проверить' }).click()
  await expect(dialog).toBeHidden({ timeout: 60_000 })
  await openJournal(page)
  await expect(page.locator('.panel-body.text-panel')).toContainText(/Проверка|завершен/i, {
    timeout: 60_000,
  })
})

test('диалог проверки: опция «Только синтаксис» и отмена', async () => {
  const page = await connectDesktop()
  await ensureExamplesProject(page)
  await clickCatalogFeature(page, 0)
  await openMenuItem(page, 'Запись и тест', 'Проверить…')
  const dialog = page.getByRole('dialog', { name: 'Проверка сценария' })
  await expect(dialog).toBeVisible()
  await expect(dialog.getByText('Только синтаксис (без браузера)')).toBeVisible()
  await dialog.getByRole('button', { name: 'Отмена' }).click()
  await expect(dialog).toBeHidden()
})

test('настройки открываются и закрываются', async () => {
  const page = await connectDesktop()
  await dismissBlockingDialogs(page)
  await openSettings(page)
  const dialog = page.getByRole('dialog', { name: /Настройки/ })
  await expect(dialog.getByPlaceholder('Поиск настроек')).toBeVisible()
  await dialog.getByRole('button', { name: 'Отмена' }).click()
  await expect(dialog).toBeHidden()
})

test('настройки: сброс по умолчанию восстанавливает браузер', async () => {
  const page = await connectDesktop()
  await openSettings(page)
  const dialog = page.getByRole('dialog', { name: /Настройки/ })
  await dialog.locator('select').first().selectOption('firefox')
  await dialog.getByRole('button', { name: 'Сбросить по умолчанию' }).click()
  await expect(dialog.locator('select').first()).toHaveValue('chromium')
  await dialog.getByRole('button', { name: 'Отмена' }).click()
})

test('настройки: Apply оставляет диалог открытым', async () => {
  const page = await connectDesktop()
  await openSettings(page)
  const dialog = page.getByRole('dialog', { name: /Настройки/ })
  await dialog.getByRole('button', { name: 'Применить' }).click()
  await expect(dialog).toBeVisible()
  await dialog.getByRole('button', { name: 'Отмена' }).click()
})

test('настройки: предупреждение при экстремальных workers и slowMo', async () => {
  const page = await connectDesktop()
  await openSettings(page)
  const dialog = page.getByRole('dialog', { name: /Настройки/ })
  await dialog.locator('.setting-card', { hasText: 'Параллельные воркеры' }).locator('input.setting-number').fill('16')
  await dialog.locator('.setting-card', { hasText: 'slow-mo' }).locator('input.setting-number').fill('5000')
  await expect(dialog.getByText('Много воркеров и высокий slow-mo')).toBeVisible()
  await dialog.getByRole('button', { name: 'Отмена' }).click()
})

test('палитра команд и справка по шагам', async () => {
  const page = await connectDesktop()
  await dismissBlockingDialogs(page)
  await page.keyboard.press('Control+Shift+KeyP')
  const palette = page.getByRole('dialog', { name: 'Палитра команд' })
  await expect(palette).toBeVisible()
  await expect(page.getByPlaceholder('Введите команду…')).toBeVisible()
  await page.keyboard.press('Escape')

  await page.keyboard.press('F1')
  const help = page.getByRole('dialog', { name: 'Справка по шагам' })
  await expect(help).toBeVisible()
  await expect(help.getByText('нажимаю (click)', { exact: false }).first()).toBeVisible({ timeout: 15_000 })
  await expect(help.locator('.list button').first()).toBeVisible({ timeout: 15_000 })
  await help.getByRole('button', { name: 'Закрыть' }).click()
})

test('горячие клавиши: Shift+F1', async () => {
  const page = await connectDesktop()
  await page.keyboard.press('Shift+F1')
  const dialog = page.getByRole('dialog', { name: 'Горячие клавиши' })
  await expect(dialog).toBeVisible()
  await expect(dialog.getByText('Ctrl+Shift+P')).toBeVisible()
  await expect(dialog.getByText('Палитра команд')).toBeVisible()
  await dialog.getByRole('button', { name: 'OK' }).click()
  await expect(dialog).toBeHidden()
})

test('журнал: Ctrl+` переключает нижнюю панель', async () => {
  const page = await connectDesktop()
  const ide = page.locator('.ide')
  if (!(await ide.evaluate((el) => el.classList.contains('panel-open')))) {
    await page.keyboard.press('Control+Backquote')
  }
  await expect(ide).toHaveClass(/panel-open/)
  await expect(page.locator('.panel-tab.active', { hasText: 'Журнал' })).toBeVisible()
  await page.keyboard.press('Control+Backquote')
  await expect(ide).not.toHaveClass(/panel-open/)
})

test('новый сценарий создаёт вкладку', async () => {
  const page = await connectDesktop()
  await ensureExamplesProject(page)
  await createNewScenario(page)
})

test('палитра сниппетов показывает шаги каталога', async () => {
  const page = await connectDesktop()
  await ensureExamplesProject(page)
  await createNewScenario(page)
  await openMenuItem(page, 'Сценарий', 'Палитра сниппетов…')
  const dialog = page.getByRole('dialog', { name: 'Палитра сниппетов' })
  await expect(dialog).toBeVisible()
  await expect(dialog.getByPlaceholder('Поиск шага…')).toBeVisible()
  await expect(dialog.locator('.snippet-grid .label', { hasText: /нажимаю/i }).first()).toBeVisible({
    timeout: 15_000,
  })
  await page.keyboard.press('Escape')
})

test('диалог запуска открывается из меню', async () => {
  const page = await connectDesktop()
  await ensureExamplesProject(page)
  await createNewScenario(page)
  await openMenuItem(page, 'Запись и тест', 'Запустить…')
  const dialog = page.getByRole('dialog', { name: 'Запуск сценария' })
  await expect(dialog).toBeVisible()
  await expect(dialog.getByText('Dry-run (без браузера)')).toBeVisible()
  await dialog.getByRole('button', { name: 'Отмена' }).click()
  await expect(dialog).toBeHidden()
})

test('первый Ctrl+Enter открывает диалог запуска', async () => {
  const page = await connectDesktop()
  await ensureExamplesProject(page)
  await dismissBlockingDialogs(page)
  await resetRunDialogConfirmed(page)
  await page.reload()
  await waitForAppReady(page)
  await ensureExamplesProject(page)
  await createNewScenario(page)
  await page.locator('.menubar').first().click()
  await page.locator('.tool-btn.primary-run').click()
  const dialog = page.getByRole('dialog', { name: 'Запуск сценария' })
  await expect(dialog).toBeVisible()
  await dialog.getByRole('button', { name: 'Отмена' }).click()
  await expect(dialog).toBeHidden()
})

test('экспорт: превью текущего сценария', async () => {
  const page = await connectDesktop()
  await ensureExamplesProject(page)
  await clickCatalogFeature(page, 0)
  await openMenuItem(page, 'Сценарий', 'Экспорт…')
  const dialog = page.getByRole('dialog', { name: 'Экспорт сценария' })
  await expect(dialog).toBeVisible()
  await expect(dialog.getByText(/\d+ шаг\(ов\)/)).toBeVisible()
  await dialog.getByRole('button', { name: 'Отмена' }).click()
  await expect(dialog).toBeHidden()
})

test('несохранённые изменения при закрытии вкладки', async () => {
  const page = await connectDesktop()
  await createNewScenario(page)
  const activeTab = page.locator('.editor-tab.file.active')
  await typeInEditor(page, '# desktop-smoke-dirty')
  await expect(activeTab.locator('.tab-label', { hasText: /\*/ })).toBeVisible()
  await activeTab.locator('.tab-close').click()
  const dialog = page.getByRole('dialog', { name: 'Несохранённые изменения' })
  await expect(dialog).toBeVisible()
  await dialog.getByRole('button', { name: 'Отмена' }).click()
  await expect(dialog).toBeHidden()
  await expect(activeTab).toBeVisible()
})

test('плагины: диалог управления', async () => {
  const page = await connectDesktop()
  await ensureExamplesProject(page)
  await openMenuItem(page, 'Плагины', 'Управление плагинами…')
  const dialog = page.getByRole('dialog', { name: 'Плагины' })
  await expect(dialog).toBeVisible()
  await expect(dialog.getByText('Плагины проекта')).toBeVisible()
  await dialog.getByRole('button', { name: 'Закрыть' }).click()
  await expect(dialog).toBeHidden()
})

test('о программе показывает версию', async () => {
  const page = await connectDesktop()
  await openMenuItem(page, 'Справка', 'О программе')
  const dialog = page.getByRole('dialog', { name: 'О программе' })
  await expect(dialog).toBeVisible()
  await expect(dialog.locator('.about-version')).not.toBeEmpty()
  await dialog.getByRole('button', { name: 'OK' }).click()
  await expect(dialog).toBeHidden()
})

test('диалог записи открывается и закрывается', async () => {
  const page = await connectDesktop()
  await ensureExamplesProject(page)
  await catalogFeature(page, 'pervaya').click().catch(async () => {
    await clickCatalogFeature(page, 0)
  })
  const recordDialog = await openRecordDialog(page)
  await expect(recordDialog.getByRole('button', { name: 'Начать' })).toBeVisible()
  await page.keyboard.press('Escape')
  await expect(recordDialog).toBeHidden()
})

test('запись: HTTP Auth закрывается поверх диалога записи', async () => {
  const page = await connectDesktop()
  await ensureExamplesProject(page)
  await catalogFeature(page, 'pervaya').click().catch(async () => {
    await clickCatalogFeature(page, 0)
  })

  const recordDialog = await openRecordDialog(page)
  await recordDialog.getByRole('button', { name: 'HTTP Auth…' }).click()
  const auth = page.getByRole('dialog', { name: 'HTTP авторизация' })
  await expect(auth).toBeVisible()
  await expect(recordDialog).toBeVisible()

  await auth.getByRole('button', { name: 'Закрыть' }).click()
  await expect(auth).toBeHidden()
  await expect(recordDialog).toBeVisible()

  await recordDialog.getByRole('button', { name: 'HTTP Auth…' }).click()
  await expect(auth).toBeVisible()
  await page.keyboard.press('Escape')
  await expect(auth).toBeHidden()
  await expect(recordDialog).toBeVisible()
  await page.keyboard.press('Escape')
})
