import { defineConfig } from '@playwright/test'

/**
 * Desktop smoke — Playwright подключается к уже запущенному scenaria-gui.exe
 * через WebView2 CDP (порт задаётся в scripts/desktop-smoke.ps1).
 * webServer не используется.
 */
export default defineConfig({
  testDir: './specs',
  testMatch: ['desktop-smoke.spec.ts', 'desktop-tour-onboarding.spec.ts'],
  timeout: 180_000,
  fullyParallel: false,
  workers: 1,
  retries: 0,
  reporter: 'list',
  use: {
    trace: 'retain-on-failure',
  },
  projects: [{ name: 'desktop', use: {} }],
})
