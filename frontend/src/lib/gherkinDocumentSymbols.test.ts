import { describe, expect, it } from 'vitest'
import { flattenFeatureSymbols, parseFeatureSymbols } from './gherkinDocumentSymbols'

const SAMPLE = `Функционал: Демо
  @smoke
  Сценарий: Вход
    Когда открыт "https://example.com"
    И нажимаю "button.login"
    Если вижу ".banner"
      И нажимаю "button.accept"
`

describe('parseFeatureSymbols', () => {
  it('builds feature → scenario → steps hierarchy', () => {
    const roots = parseFeatureSymbols(SAMPLE)
    expect(roots).toHaveLength(1)
    expect(roots[0].kind).toBe('feature')
    expect(roots[0].name).toBe('Демо')

    const scenario = roots[0].children.find((node) => node.kind === 'scenario')
    expect(scenario?.name).toBe('Вход')
    expect(scenario?.children.some((node) => node.kind === 'tag')).toBe(true)
    expect(scenario?.children.filter((node) => node.kind === 'step')).toHaveLength(2)
    const block = scenario?.children.find((node) => node.kind === 'block')
    expect(block?.children).toHaveLength(1)
  })

  it('flattens tree with depth', () => {
    const flat = flattenFeatureSymbols(parseFeatureSymbols(SAMPLE))
    expect(flat[0].depth).toBe(0)
    expect(flat.some((row) => row.kind === 'step' && row.depth >= 2)).toBe(true)
  })
})
