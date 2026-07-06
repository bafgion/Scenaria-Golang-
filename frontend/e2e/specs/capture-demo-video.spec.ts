/**
 * Демо-видео ~2 мин: GUI (ru) → запись → прогон → HTML-отчёты (успех / падение).
 *
 *   cd frontend && npm run docs:demo-video
 *
 * Результат: docs/videos/scenaria-demo-ru.webm
 */
import { expect, test } from '@playwright/test'
import fs from 'node:fs/promises'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import {
  bootApp,
  catalogFeature,
  mockPath,
  openJournal,
  openMenuItem,
  openRecordDialog,
  openResults,
  stopRecording,
} from '../helpers/app'
import { hold, reportFileURL, showCaption } from '../helpers/demo-video'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const reportDir = path.resolve(__dirname, '../fixtures/report')
const videoOut = path.resolve(__dirname, '../../../docs/videos/scenaria-demo-ru.webm')

test.describe.configure({ mode: 'serial' })

test('запись демо-видео Scenaria (ru)', async ({ page }) => {
  await page.addInitScript({ path: mockPath })
  await bootApp(page, '?e2e=demo-video&locale=ru')

  await showCaption(page, 'Scenaria — BDD-сценарии и автотесты в браузере', 3200)
  await hold(page, 2000)

  await showCaption(page, 'Каталог сценариев и редактор Gherkin', 1800)
  await openMenuItem(page, 'Проект', 'Открыть примеры сценариев')
  await expect(page.locator('.catalog-tree')).toBeVisible({ timeout: 10_000 })
  await hold(page, 2000)
  await catalogFeature(page, '^○ smoke$').click()
  await expect(page.locator('.editor-tab.file .tab-label', { hasText: 'smoke' })).toBeVisible()
  await hold(page, 3500)

  await showCaption(page, 'Проверка синтаксиса Gherkin', 2000)
  await openMenuItem(page, 'Запись и тест', 'Проверить…')
  const validateDialog = page.getByRole('dialog', { name: 'Проверка сценария' })
  await validateDialog.getByRole('button', { name: 'Проверить' }).click()
  await expect(validateDialog).toBeHidden({ timeout: 10_000 })
  await openJournal(page)
  await expect(page.locator('.panel-body.text-panel')).toContainText('Проверка завершена', { timeout: 10_000 })
  await hold(page, 2500)

  await showCaption(page, 'Запись действий пользователя в шаги', 1800)
  await page.keyboard.press('Escape')
  const recordDialog = await openRecordDialog(page)
  await hold(page, 1200)
  await recordDialog.getByRole('button', { name: 'Начать' }).click()
  await expect(recordDialog).toBeHidden({ timeout: 10_000 })
  await expect(page.getByRole('button', { name: 'Стоп запись' })).toBeVisible()
  await expect(page.locator('.status-bar .recording-target')).toContainText('Запись')
  await hold(page, 5500)
  await stopRecording(page)
  await openJournal(page)
  await expect(page.locator('.panel-body.text-panel')).toContainText('Запись', { timeout: 10_000 })
  await hold(page, 3500)

  await showCaption(page, 'Прогон сценариев в Chromium', 1800)
  await openMenuItem(page, 'Запись и тест', 'Запустить…')
  await expect(page.getByRole('dialog')).toBeVisible()
  await hold(page, 2800)
  await page.keyboard.press('Escape')
  await page.keyboard.press('Control+Enter')
  await expect(page.locator('.playing-bar')).toBeVisible({ timeout: 10_000 })
  await expect(page.locator('.playing-bar .play-progress-text')).toContainText(/\d\/2/, { timeout: 10_000 })
  await expect(page.locator('.playing-bar')).toBeHidden({ timeout: 25_000 })
  await openResults(page)
  await expect(page.locator('.results-table tbody tr')).toHaveCount(2)
  await expect(page.locator('.results-table tr.failed')).toHaveCount(1)
  await hold(page, 4500)

  const passedURL = reportFileURL(path.join(reportDir, 'sample-passed.html'))
  await showCaption(page, 'HTML-отчёт: успешный сценарий', 2200)
  await page.goto(passedURL)
  await expect(page.locator('.step-node.passed.active')).toBeVisible({ timeout: 10_000 })
  await hold(page, 3500)
  await page.locator('#tab-actionlog').click()
  await hold(page, 3000)
  await page.locator('#tab-timeline').click()
  await hold(page, 2500)

  const failedURL = reportFileURL(path.join(reportDir, 'sample.html'))
  await showCaption(page, 'HTML-отчёт: падение, trace и inspector', 1800)
  await page.goto(failedURL)
  await expect(page.locator('.step-node.failed.active')).toBeVisible({ timeout: 10_000 })
  await hold(page, 3500)
  await page.locator('#tab-trace').click()
  await expect(page.locator('#traceview .trace-event-row').first()).toBeVisible()
  await hold(page, 3500)
  await page.locator('#traceview .trace-event-row').first().click()
  await hold(page, 3000)
  await page.locator('#traceview .trace-step-row[data-idx="0"]').click()
  await hold(page, 3500)

  await showCaption(page, 'Scenaria — scenaria-go.dev', 3000)
})

test.afterAll(async () => {
  const outDir = path.dirname(videoOut)
  await fs.mkdir(outDir, { recursive: true })
  const artifactsRoot = path.resolve(__dirname, '../../../docs/videos/.playwright-output')
  let newest: { path: string; mtime: number } | null = null
  async function walk(dir: string) {
    let entries: { name: string; isDirectory: () => boolean; isFile: () => boolean }[]
    try {
      entries = await fs.readdir(dir, { withFileTypes: true })
    } catch {
      return
    }
    for (const ent of entries) {
      const full = path.join(dir, ent.name)
      if (ent.isDirectory()) await walk(full)
      else if (ent.isFile() && ent.name.endsWith('.webm')) {
        const st = await fs.stat(full)
        if (!newest || st.mtimeMs > newest.mtime) newest = { path: full, mtime: st.mtimeMs }
      }
    }
  }
  await walk(artifactsRoot)
  if (newest) {
    await fs.copyFile(newest.path, videoOut)
    // eslint-disable-next-line no-console
    console.log('Demo video saved:', videoOut)
  }
})
