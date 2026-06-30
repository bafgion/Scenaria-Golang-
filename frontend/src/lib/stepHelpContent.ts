import type { StepHelpEntry } from './stepTypes'

export function stepHelpDescription(entry: StepHelpEntry): string {
  return entry.description || entry.help || ''
}

export function stepHelpExample(entry: StepHelpEntry): string {
  return entry.example || entry.template || ''
}

export function stepHelpTitle(entry: StepHelpEntry): string {
  if (entry.label) {
    return entry.action ? `${entry.label} (${entry.action})` : entry.label
  }
  return entry.template
}

export function formatStepHoverMarkdown(entry: StepHelpEntry): string {
  const lines: string[] = []
  lines.push(`**${stepHelpTitle(entry)}**`)
  if (entry.category) {
    lines.push(`*${entry.category}*`)
  }
  const desc = stepHelpDescription(entry)
  if (desc) {
    lines.push('', desc)
  }
  if (entry.parameters?.length) {
    lines.push('', '**Параметры**')
    for (const param of entry.parameters) {
      lines.push(`- ${param}`)
    }
  }
  const example = stepHelpExample(entry)
  if (example) {
    lines.push('', '**Пример**', '```', example, '```')
  }
  lines.push('', '_F1 — полная справка по шагам_')
  return lines.join('\n')
}

export function hasStepHelp(entry: StepHelpEntry | null | undefined): entry is StepHelpEntry {
  if (!entry) return false
  return Boolean(entry.label || entry.template || entry.action)
}
