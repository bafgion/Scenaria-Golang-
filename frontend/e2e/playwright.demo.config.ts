import { defineConfig, devices } from '@playwright/test'

/** Запись демо-видео для docs/videos/ — не часть CI e2e. */
export default defineConfig({
  testDir: './specs',
  testMatch: ['capture-demo-video.spec.ts'],
  timeout: 300_000,
  fullyParallel: false,
  workers: 1,
  retries: 0,
  reporter: 'list',
  outputDir: '../../docs/videos/.playwright-output',
  use: {
    baseURL: 'http://127.0.0.1:4173',
    viewport: { width: 1440, height: 900 },
    deviceScaleFactor: 1,
    locale: 'ru-RU',
    video: { mode: 'on', size: { width: 1440, height: 900 } },
  },
  webServer: {
    command: 'npm run build && npm run preview -- --host 127.0.0.1 --port 4173',
    url: 'http://127.0.0.1:4173',
    reuseExistingServer: !process.env.CI,
    timeout: 180_000,
  },
  projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'] } }],
})
