import type { gui } from '../../wailsjs/go/models'
import type { EditorSettings } from './editorOptions'

export function filterScenarioHints(
  hints: gui.ScenarioHintDTO[],
  settings: Pick<EditorSettings, 'scenarioHintsShowWarning' | 'scenarioHintsShowInfo'>,
): gui.ScenarioHintDTO[] {
  return hints.filter((hint) => {
    if (hint.severity === 'warning') return settings.scenarioHintsShowWarning
    return settings.scenarioHintsShowInfo
  })
}

export async function applyAutoFixableScenarioHints(
  text: string,
  hints: gui.ScenarioHintDTO[],
  applyFix: (hint: gui.ScenarioHintDTO, currentText: string) => Promise<string | null>,
): Promise<string> {
  let current = text
  const pending = hints.filter((h) => h.autoFixable)
  for (const hint of pending) {
    const next = await applyFix(hint, current)
    if (next && next !== current) {
      current = next
    }
  }
  return current
}
