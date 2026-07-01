import { expect, test } from '@playwright/test'
import {
  bootApp,
  createNewScenario,
  editorLine,
  expectRunMenuItemDisabled,
  expectRunMenuItemEnabled,
  openTestProject,
  pauseRecording,
  startRecordingFromDialog,
  statusMessage,
  stopRecording,
} from '../helpers/app'

test('record resume does not duplicate steps', async ({ page }) => {
  await bootApp(page, '?e2e=record-resume')
  await openTestProject(page)
  await createNewScenario(page)

  await startRecordingFromDialog(page)
  await expect(editorLine(page, 'https://iana.org/domains/example')).toHaveCount(1, { timeout: 10_000 })

  await stopRecording(page)
  await expectRunMenuItemDisabled(page, 'Пауза')

  await startRecordingFromDialog(page)
  await expect(editorLine(page, 'https://iana.org/domains/example')).toHaveCount(1, { timeout: 10_000 })
  await expectRunMenuItemEnabled(page, 'Пауза')
})

test('second record segment appends new steps from index zero', async ({ page }) => {
  await bootApp(page, '?e2e=record-resume')
  await openTestProject(page)
  await createNewScenario(page)

  await startRecordingFromDialog(page)
  await stopRecording(page)
  await startRecordingFromDialog(page)

  await page.evaluate(() => {
    const emit = window.__e2eEmit
    if (!emit) return
    emit('record-step', { index: 0, line: 'нажимаю "#after-stop"' })
    emit('record-step', { index: 1, line: 'нажимаю "#after-stop-2"' })
  })

  await expect(editorLine(page, 'нажимаю "#after-stop"')).toHaveCount(1, { timeout: 10_000 })
  await expect(editorLine(page, 'нажимаю "#after-stop-2"')).toHaveCount(1, { timeout: 10_000 })
  await expect(editorLine(page, 'https://iana.org/domains/example')).toHaveCount(1)
})

test('pause and stop work after record resume', async ({ page }) => {
  await bootApp(page, '?e2e=record-resume')
  await openTestProject(page)
  await createNewScenario(page)

  await startRecordingFromDialog(page)
  await stopRecording(page)
  await startRecordingFromDialog(page)

  await expectRunMenuItemEnabled(page, 'Пауза')
  await pauseRecording(page)
  await expect(statusMessage(page)).toContainText('Пауза', { timeout: 10_000 })

  await stopRecording(page)
  await expectRunMenuItemDisabled(page, 'Пауза')
})
