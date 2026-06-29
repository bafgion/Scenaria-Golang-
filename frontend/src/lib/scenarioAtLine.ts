const SCENARIO_LINE =
  /^\s*(?:Сценарий|Структура сценария|Scenario|Scenario Outline)\s*:\s*(.+)$/i

/** Returns scenario titles declared in a .feature file (in source order). */
export function listScenarioTitles(text: string): string[] {
  const titles: string[] = []
  for (const line of text.split('\n')) {
    const match = line.match(SCENARIO_LINE)
    if (match) titles.push(match[1].trim())
  }
  return titles
}

/** Returns the scenario title active at the given 1-based line in a .feature file. */
export function scenarioAtLine(text: string, line: number): string {
  const lines = text.split('\n')
  const limit = Math.min(Math.max(line, 1), lines.length)
  let current = ''
  for (let i = 0; i < limit; i++) {
    const match = lines[i].match(SCENARIO_LINE)
    if (match) current = match[1].trim()
  }
  return current
}
