import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './specs',
  testIgnore: ['**/desktop-smoke.spec.ts', '**/desktop-tour-onboarding.spec.ts', '**/capture-docs-screenshots.spec.ts', '**/capture-demo-video.spec.ts'],
  timeout: 60_000,
  fullyParallel: false,
  retries: process.env.CI ? 1 : 0,
  use: {
    baseURL: 'http://127.0.0.1:4173',
    trace: 'on-first-retry',
  },
  webServer: {
    command: 'npm run build && npm run preview -- --host 127.0.0.1 --port 4173',
    url: 'http://127.0.0.1:4173',
    reuseExistingServer: !process.env.CI,
    timeout: 120_000,
  },
  projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'] } }],
})
