import type { gui } from '../../wailsjs/go/models'

export type StepInlayRow = Pick<gui.EditorStepRow, 'line' | 'kind' | 'action' | 'element' | 'value' | 'error'>

function formatHintValue(value: string): string {
  if (!value) return ''
  if (/^https?:\/\//i.test(value)) return value
  if (/^[#.\[{@]/.test(value) || /[=>\[\].]/.test(value)) return value
  if (value.length > 28) return `"${value.slice(0, 25)}…"`
  return `"${value}"`
}

/** Builds gray inlay hint text for a parsed editor step row. */
export function formatStepInlayLabel(step: StepInlayRow): string | null {
  if (step.error) {
    const message = step.error.length > 36 ? `${step.error.slice(0, 33)}…` : step.error
    return `⚠ ${message}`
  }
  const kind = (step.kind || step.action || '').trim()
  if (!kind) return null

  if (step.element && step.value) {
    return `${kind} → ${formatHintValue(step.element)} / ${formatHintValue(step.value)}`
  }
  if (step.element) return `${kind} → ${formatHintValue(step.element)}`
  if (step.value) return `${kind} → ${formatHintValue(step.value)}`
  return kind
}

export function inlayHintColumn(lineText: string): number {
  const trimmedEnd = lineText.trimEnd()
  return Math.max(1, trimmedEnd.length + 1)
}
