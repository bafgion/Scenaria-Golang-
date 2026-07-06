import type { Page } from '@playwright/test'

export async function showCaption(page: Page, text: string, ms = 2200) {
  await page.evaluate((t) => {
    let el = document.getElementById('demo-caption') as HTMLDivElement | null
    if (!el) {
      el = document.createElement('div')
      el.id = 'demo-caption'
      Object.assign(el.style, {
        position: 'fixed',
        bottom: '56px',
        left: '50%',
        transform: 'translateX(-50%)',
        background: 'rgba(15, 23, 42, 0.88)',
        color: '#f8fafc',
        padding: '14px 28px',
        borderRadius: '10px',
        fontSize: '22px',
        fontWeight: '600',
        zIndex: '99999',
        pointerEvents: 'none',
        fontFamily: 'Segoe UI, system-ui, sans-serif',
        boxShadow: '0 8px 32px rgba(0,0,0,0.35)',
        letterSpacing: '0.01em',
      })
      document.body.appendChild(el)
    }
    el.textContent = t
    el.style.display = 'block'
  }, text)
  await page.waitForTimeout(ms)
  await page.evaluate(() => {
    const el = document.getElementById('demo-caption')
    if (el) el.style.display = 'none'
  })
}

export async function hold(page: Page, ms: number) {
  await page.waitForTimeout(ms)
}

export function reportFileURL(filePath: string) {
  return 'file:///' + filePath.replace(/\\/g, '/')
}
