import { describe, expect, it } from 'vitest'
import { createTranslator } from './i18n'
import { formatCloseAppMessage, normalizeCloseGuardReasons } from './closeAppGuard'

describe('closeAppGuard', () => {
  it('formats localized close-app message from reason codes', () => {
    const tr = createTranslator('ru')
    const message = formatCloseAppMessage(tr, ['unsaved_tabs', 'active_recorder'])
    expect(message).toContain('несохранённые вкладки')
    expect(message).toContain('активная запись/браузер')
    expect(message).toContain('Всё равно закрыть приложение?')
  })

  it('ignores unknown reason codes', () => {
    expect(normalizeCloseGuardReasons(['unsaved_tabs', 'legacy reason'])).toEqual(['unsaved_tabs'])
  })
})
