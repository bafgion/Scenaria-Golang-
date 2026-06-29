import { describe, expect, it } from 'vitest'
import { listScenarioTitles, mergeScenarioNames, scenarioAtLine } from './scenarioAtLine'

const sample = `Функция: Demo
  Сценарий: Первый
    открыт "https://a.com"
  Сценарий: Второй
    нажимаю "#ok"
`

describe('scenarioAtLine', () => {
  it('returns active scenario at line', () => {
    expect(scenarioAtLine(sample, 1)).toBe('')
    expect(scenarioAtLine(sample, 3)).toBe('Первый')
    expect(scenarioAtLine(sample, 5)).toBe('Второй')
    expect(scenarioAtLine(sample, 99)).toBe('Второй')
  })

  it('lists scenario titles in order', () => {
    expect(listScenarioTitles(sample)).toEqual(['Первый', 'Второй'])
  })

  it('merges scenario name groups without duplicates', () => {
    expect(mergeScenarioNames(['A', 'B'], ['B', 'C'])).toEqual(['A', 'B', 'C'])
  })
})
