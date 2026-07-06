/**
 * User journey E2E — цепочки действий, эмулирующие типичный рабочий день QA-инженера.
 *
 * В отличие от qa-daily-use.spec.ts (пункты чеклиста 1:1), здесь — сквозные сценарии:
 * открытие проекта → правка → прогон → разбор результатов.
 *
 * Локальный запуск:
 *   cd frontend && npm run test:e2e:journeys
 *
 * Полный E2E (все спеки):
 *   cd frontend && npm run test:e2e
 */
import { expect, test } from '@playwright/test'
import {
  bootApp,
  catalogFeature,
  confirmRunFromDialog,
  createNewScenario,
  focusBrowserFromToolbar,
  journalPanel,
  mockPath,
  openJournal,
  openMenuItem,
  openRecordDialog,
  openResults,
  openSettings,
  openTestProject,
  runCurrentScenario,
  saveScenario,
  startRecordingFromDialog,
  statusMessage,
  stopRecording,
  typeInEditor,
} from '../helpers/app'

test.beforeEach(async ({ page }) => {
  await page.addInitScript({ path: mockPath })
})

test.describe('Сценарий: знакомство с примерами', () => {
  test('примеры → smoke → проверка → dry-run → журнал', async ({ page }) => {
    await bootApp(page, '?e2e=examples')
    await openMenuItem(page, 'Проект', 'Открыть примеры сценариев')
    await expect(page.locator('.catalog-tree')).toBeVisible({ timeout: 10_000 })
    await catalogFeature(page, '^○ smoke$').click()
    await expect(page.locator('.editor-tab.file .tab-label', { hasText: 'smoke' })).toBeVisible()

    await openMenuItem(page, 'Запись и тест', 'Проверить…')
    const validateDialog = page.getByRole('dialog', { name: 'Проверка сценария' })
    await validateDialog.getByRole('button', { name: 'Проверить' }).click()
    await expect(validateDialog).toBeHidden({ timeout: 10_000 })
    await openJournal(page)
    await expect(journalPanel(page)).toContainText('Проверка завершена', { timeout: 10_000 })

    await openMenuItem(page, 'Запись и тест', 'Запустить…')
    await confirmRunFromDialog(page, { dryRun: true })
    await typeInEditor(page, '# journey-dry-run-edit')
    await expect(page.locator('.monaco-editor .view-line', { hasText: 'journey-dry-run-edit' })).toBeVisible()
    await expect(page.locator('.playing-bar')).toHaveCount(0)
  })
})

test.describe('Сценарий: новый проект и сценарий', () => {
  test('мастер проекта → правка → сохранение → проверка', async ({ page }) => {
    await bootApp(page, '?e2e=new-project')
    await openMenuItem(page, 'Проект', 'Новый проект…')
    const wizard = page.getByRole('dialog', { name: 'Новый проект' })
    await wizard.getByRole('button', { name: 'Обзор…' }).click()
    await wizard.locator('label', { hasText: 'Название feature' }).locator('input').fill('checkout')
    await wizard.getByRole('button', { name: 'Создать' }).click()
    await expect(wizard).toBeHidden({ timeout: 10_000 })
    await expect(page.locator('.editor-tab.file .tab-label', { hasText: 'checkout.feature' })).toBeVisible()

    const marker = 'JOURNEY_CHECKOUT_STEP'
    await typeInEditor(page, marker)
    await saveScenario(page)
    await expect(statusMessage(page)).toHaveText('Сохранено', { timeout: 10_000 })
    await expect(page.locator('.monaco-editor .view-line', { hasText: marker })).toBeVisible()

    await openMenuItem(page, 'Запись и тест', 'Проверить…')
    await page.getByRole('dialog', { name: 'Проверка сценария' }).getByRole('button', { name: 'Проверить' }).click()
    await openJournal(page)
    await expect(journalPanel(page)).toContainText('Проверка завершена', { timeout: 10_000 })
  })
})

test.describe('Сценарий: live-запись', () => {
  test('диалог записи → HTTP Auth → старт → стоп → журнал', async ({ page }) => {
    await bootApp(page, '?e2e=post-record')
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()

    const recordDialog = await openRecordDialog(page)
    await recordDialog.getByRole('button', { name: 'HTTP Auth…' }).click()
    const authDialog = page.getByRole('dialog', { name: 'HTTP авторизация' })
    await expect(authDialog).toBeVisible()
    await authDialog.locator('label', { hasText: 'Хост' }).locator('input').fill('example.com')
    await authDialog.locator('label', { hasText: 'Логин' }).locator('input').fill('user')
    await authDialog.getByRole('button', { name: 'Закрыть' }).click()
    await expect(authDialog).toBeHidden()
    await expect(recordDialog).toBeVisible()

    await recordDialog.getByRole('button', { name: 'Начать' }).click()
    await expect(recordDialog).toBeHidden({ timeout: 10_000 })
    await expect(page.locator('.status-bar .recording-target')).toContainText('Запись → smoke.feature')
    await expect(page.getByRole('button', { name: 'Стоп запись' })).toBeVisible()

    await stopRecording(page)
    await openJournal(page)
    await expect(journalPanel(page)).toContainText('Запись', { timeout: 10_000 })
  })

  test('запись → пауза → переключение вкладки с confirm → отмена', async ({ page }) => {
    await bootApp(page, '?e2e=post-record')
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    await createNewScenario(page)
    await page.locator('.editor-tab.file', { hasText: 'smoke' }).click()
    await startRecordingFromDialog(page)

    await page.locator('.editor-tab.file', { hasText: 'novyy-scenariy' }).first().click()
    const confirm = page.getByRole('alertdialog', { name: 'Запись активна' })
    await expect(confirm).toBeVisible()
    await confirm.getByRole('button', { name: 'Отмена' }).click()
    await expect(page.locator('.editor-tab.file.active', { hasText: 'smoke' })).toBeVisible()
  })
})

test.describe('Сценарий: прогон suite', () => {
  test('suite 3 сценария → playing bar → журнал', async ({ page }) => {
    await bootApp(page, '?e2e=run-progress')
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    await runCurrentScenario(page)
    await expect(page.locator('.playing-bar .play-progress-text')).toContainText(/\d\/3/, { timeout: 10_000 })
    await expect(page.locator('.playing-bar')).toBeHidden({ timeout: 15_000 })
    await openJournal(page)
    await expect(journalPanel(page)).toContainText('Завершено', { timeout: 10_000 })
  })

  test('прогон со стримом в журнал → отмена', async ({ page }) => {
    await bootApp(page, '?e2e=run-cancel')
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    await runCurrentScenario(page)
    await expect(page.locator('.playing-bar')).toBeVisible({ timeout: 10_000 })
    await openJournal(page)
    await expect(journalPanel(page)).toContainText('progress', { timeout: 10_000 })
    await page.locator('.playing-bar .play-cancel-btn').click()
    await expect(page.locator('.playing-bar')).toBeHidden({ timeout: 15_000 })
  })
})

test.describe('Сценарий: разбор упавшего теста', () => {
  test('результаты → шаг ошибки → flaky 3× → trace', async ({ page }) => {
    await bootApp(page, '?e2e=flaky-run')
    await openTestProject(page)
    await openResults(page)

    await page.getByRole('button', { name: 'Шаг 2' }).first().click()
    await expect(page.locator('.editor-tab.file.active .tab-label', { hasText: 'smoke.feature' })).toBeVisible()

    await page.getByRole('button', { name: 'Запустить 3×' }).first().click()
    await expect(page.locator('.playing-bar')).toBeHidden({ timeout: 10_000 })
    await openJournal(page)
    await expect(journalPanel(page)).toContainText('Прогон 3/3', { timeout: 10_000 })

    await openResults(page)
    await page.getByRole('button', { name: 'Trace viewer' }).click()
    await openJournal(page)
    await expect(journalPanel(page)).toContainText('Trace viewer:', { timeout: 10_000 })
  })
})

test.describe('Сценарий: браузер и фокус', () => {
  test('открыть браузер → показать браузер → закрыть', async ({ page }) => {
    await bootApp(page)
    await openTestProject(page)
    await page.locator('button.tool-btn.primary[title="Браузер (Ctrl+B)"]').click()
    await expect(statusMessage(page)).toContainText('Браузер открыт', { timeout: 10_000 })
    await expect(page.getByRole('button', { name: 'Закрыть браузер' })).toBeVisible()

    await focusBrowserFromToolbar(page)
    const calls = await page.evaluate(() => (window as unknown as { __e2eFocusBrowserCalls?: () => number }).__e2eFocusBrowserCalls?.())
    expect(calls).toBe(1)
    await openJournal(page)
    await expect(journalPanel(page)).toContainText('Окно браузера выведено на передний план')

    await page.getByRole('button', { name: 'Закрыть браузер' }).click()
    await expect(page.getByRole('button', { name: 'Закрыть браузер' })).toHaveCount(0)
  })
})

test.describe('Сценарий: настройки без потери работы', () => {
  test('правка сценария → настройки → apply → Ctrl+S не сохраняет → закрыть', async ({ page }) => {
    await bootApp(page)
    await openTestProject(page)
    await catalogFeature(page, 'smoke').click()
    await typeInEditor(page, '# before-settings')
    await expect(page.locator('.editor-tab.file .tab-label', { hasText: /\*/ })).toBeVisible()

    await openSettings(page)
    const settings = page.getByRole('dialog', { name: /Настройки/ })
    await settings.getByRole('button', { name: 'Сбросить по умолчанию' }).click()
    await settings.getByRole('button', { name: 'Применить' }).click()
    await expect(settings).toBeVisible()
    await page.keyboard.press('Control+KeyS')
    await expect(settings).toBeVisible()
    await settings.getByRole('button', { name: 'Отмена' }).click()

    await expect(page.locator('.editor-tab.file .tab-label', { hasText: /\*/ })).toBeVisible()
    await expect(page.locator('.monaco-editor .view-line', { hasText: 'before-settings' })).toBeVisible()
  })
})

test.describe('Сценарий: плагины', () => {
  test('управление плагинами → запуск → отмена → закрыть', async ({ page }) => {
    await bootApp(page, '?e2e=with-plugins')
    await openTestProject(page)
    await openMenuItem(page, 'Плагины', 'Управление плагинами…')
    const plugins = page.getByRole('dialog', { name: 'Плагины' })
    await expect(plugins.locator('table td div', { hasText: /^demo$/ })).toBeVisible()
    await plugins.getByRole('button', { name: 'Запуск…' }).click()
    await expect(plugins).toBeHidden()

    const runDialog = page.getByRole('dialog', { name: 'Запуск плагина' })
    await expect(runDialog).toBeVisible()
    await runDialog.getByRole('button', { name: 'Отмена' }).click()
    await expect(runDialog).toBeHidden()
  })
})

test.describe('Сценарий: модальные окна (регрессия)', () => {
  test('Escape закрывает только верхний диалог', async ({ page }) => {
    await bootApp(page)
    await openTestProject(page)
    await openRecordDialog(page)
    await page.getByRole('dialog', { name: 'Запись сценария' }).getByRole('button', { name: 'HTTP Auth…' }).click()
    await expect(page.getByRole('dialog', { name: 'HTTP авторизация' })).toBeVisible()
    await page.keyboard.press('Escape')
    await expect(page.getByRole('dialog', { name: 'HTTP авторизация' })).toBeHidden()
    await expect(page.getByRole('dialog', { name: 'Запись сценария' })).toBeVisible()
    await page.keyboard.press('Escape')
    await expect(page.getByRole('dialog', { name: 'Запись сценария' })).toBeHidden()
  })
})
