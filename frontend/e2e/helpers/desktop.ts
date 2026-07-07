import { chromium, expect, type Browser, type Page } from '@playwright/test'
import { openMenuItem } from './app'
import { dismissTourIfVisible } from './onboarding'

let browser: Browser | null = null
let page: Page | null = null

export function desktopCdpUrl(): string {
  return process.env.DESKTOP_CDP_URL || 'http://127.0.0.1:9333'
}

function isAppPage(p: Page): boolean {
  const url = p.url()
  return url.includes('wails.localhost') || url === 'about:blank'
}

async function pickAppPage(deadlineMs = 90_000): Promise<Page> {
  if (!browser) {
    throw new Error('connectDesktop: browser not connected')
  }
  const context = browser.contexts()[0]
  if (!context) {
    throw new Error('CDP: нет browser context')
  }

  const deadline = Date.now() + deadlineMs
  while (Date.now() < deadline) {
    const pages = context.pages().filter(isAppPage)
    if (pages.length > 0) {
      return pages[0]
    }
    const remaining = Math.max(500, deadline - Date.now())
    const newPage = await context
      .waitForEvent('page', { timeout: Math.min(remaining, 2000) })
      .catch(() => null)
    if (newPage && isAppPage(newPage)) {
      return newPage
    }
  }
  throw new Error('CDP: страница scenaria-gui не найдена')
}

/** Подключиться к WebView2 scenaria-gui через CDP (приложение уже запущено). */
export async function connectDesktop(): Promise<Page> {
  if (page) return page
  const url = desktopCdpUrl()
  browser = await chromium.connectOverCDP(url, { timeout: 60_000 })
  page = await pickAppPage()
  await page.waitForFunction(() => document.querySelector('.ide'), null, { timeout: 90_000 })
  return page
}

export async function disconnectDesktop(): Promise<void> {
  page = null
  if (browser) {
    await browser.close()
    browser = null
  }
}

export async function resetRunDialogConfirmed(page: Page): Promise<void> {
  await page.evaluate(async () => {
    const api = (window as unknown as {
      go?: { wailsapp?: { App?: { LoadSettings: () => Promise<Record<string, unknown>>; SaveSettings: (s: Record<string, unknown>) => Promise<void> } } }
    }).go?.wailsapp?.App
    if (!api?.LoadSettings || !api?.SaveSettings) return
    const settings = await api.LoadSettings()
    settings.runDialogConfirmed = false
    await api.SaveSettings(settings)
  })
}

export async function ensureRussianLocale(p: Page): Promise<void> {
  const projectMenu = p.locator('.menubar .menu-trigger').first()
  await expect(projectMenu).toBeVisible({ timeout: 30_000 })
  const label = ((await projectMenu.textContent()) ?? '').trim()
  if (label.includes('Проект')) return

  const needsReload = await p.evaluate(async () => {
    const api = (
      window as unknown as {
        go?: {
          wailsapp?: {
            App?: {
              LoadSettings: () => Promise<Record<string, unknown>>
              SaveSettings: (s: Record<string, unknown>) => Promise<void>
            }
          }
        }
      }
    ).go?.wailsapp?.App
    if (!api?.LoadSettings || !api?.SaveSettings) return false
    const settings = await api.LoadSettings()
    settings.uiLocale = 'ru'
    await api.SaveSettings(settings)
    return true
  })
  if (!needsReload) return
  await p.reload({ waitUntil: 'domcontentloaded' })
  await p.waitForFunction(() => document.querySelector('.ide'), null, { timeout: 90_000 })
  await expect(p.locator('.menubar .menu-trigger', { hasText: 'Проект' })).toBeVisible({ timeout: 30_000 })
}

export async function dismissBlockingDialogs(p: Page): Promise<void> {
  for (let i = 0; i < 4; i++) {
    const update = p.getByRole('dialog', { name: /обновлени|update/i })
    if (await update.isVisible().catch(() => false)) {
      await update.getByRole('button', { name: /Закрыть|Close/i }).click()
      await p.waitForTimeout(250)
      continue
    }
    const confirm = p.locator('.modal-backdrop.modal-layer-confirm')
    if ((await confirm.count()) === 0) return
    const openOther = p.getByRole('alertdialog', { name: 'Открыть другой проект' })
    if (await openOther.isVisible().catch(() => false)) {
      await openOther.getByRole('button', { name: 'Открыть' }).click()
      await p.waitForTimeout(250)
      continue
    }
    const cancel = p.getByRole('button', { name: 'Отмена' })
    if (await cancel.isVisible().catch(() => false)) {
      await cancel.click()
    } else {
      await p.keyboard.press('Escape')
    }
    await p.waitForTimeout(250)
  }
}

export async function clickCatalogFeature(p: Page, index = 0): Promise<void> {
  await dismissBlockingDialogs(p)
  const row = p.locator('.catalog-tree .catalog-tree-row.file').nth(index)
  await expect(row).toBeVisible({ timeout: 15_000 })
  await row.click()
}

export async function waitForAppReady(p: Page): Promise<void> {
  await dismissBlockingDialogs(p)
  await ensureRussianLocale(p)
  await dismissBlockingDialogs(p)
  await dismissTourIfVisible(p)
  await expect(p.locator('.ide')).toBeVisible({ timeout: 90_000 })
  await p
    .waitForFunction(
      () => {
        const welcome = document.querySelector('.welcome-panel, [role="button"]')
        const menubar = document.querySelector('.menubar')
        const catalog = document.querySelector('.catalog-tree')
        return Boolean(welcome || menubar || catalog)
      },
      { timeout: 30_000 },
    )
    .catch(() => {
      /* menubar should exist once .ide is visible */
    })
}

/** Открыть встроенные примеры, если каталог пуст (состояние serial-прогона). */
export async function ensureExamplesProject(p: Page): Promise<void> {
  await dismissBlockingDialogs(p)
  if ((await p.locator('.catalog-tree .tree-file-label').count()) > 0) {
    return
  }
  await openMenuItem(p, 'Проект', 'Открыть примеры сценариев')
  await expect(p.locator('.catalog-tree')).toBeVisible({ timeout: 30_000 })
  await dismissBlockingDialogs(p)
}
