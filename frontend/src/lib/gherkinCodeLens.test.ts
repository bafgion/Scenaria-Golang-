import { describe, expect, it } from 'vitest'
import { collectRunCodeLenses } from './gherkinCodeLens'

const SAMPLE = `Функционал: Демо
  Сценарий: Вход
    Когда открыт "https://example.com"
    И нажимаю "button.login"
  Сценарий: Выход
    Тогда закрываю браузер
`

describe('collectRunCodeLenses', () => {
  it('adds run lenses on scenario headers and steps', () => {
    const lenses = collectRunCodeLenses(SAMPLE)
    expect(lenses.some((item) => item.line === 2 && item.title === '▶ Запустить сценарий')).toBe(true)
    expect(lenses.some((item) => item.line === 3 && item.title === '▶ с этой строки')).toBe(true)
    expect(lenses.filter((item) => item.scenario === 'Вход').length).toBeGreaterThan(0)
    expect(lenses.filter((item) => item.scenario === 'Выход').length).toBeGreaterThan(0)
  })

  it('includes dry-run actions', () => {
    const lenses = collectRunCodeLenses(SAMPLE)
    expect(lenses.some((item) => item.dryRun && item.title === 'Dry-run')).toBe(true)
  })
})
