type TranslateFn = (key: string, params?: Record<string, string | number>) => string

export type CloseGuardReason = 'unsaved_tabs' | 'active_run' | 'active_recorder'

const KNOWN_REASONS = new Set<CloseGuardReason>(['unsaved_tabs', 'active_run', 'active_recorder'])

export function normalizeCloseGuardReasons(reasons: string[]): CloseGuardReason[] {
  const out: CloseGuardReason[] = []
  for (const reason of reasons) {
    if (KNOWN_REASONS.has(reason as CloseGuardReason)) {
      out.push(reason as CloseGuardReason)
    }
  }
  return out
}

export function formatCloseAppMessage(tr: TranslateFn, reasons: string[]): string {
  const normalized = normalizeCloseGuardReasons(reasons)
  const bullets = normalized
    .map((reason) => `• ${tr(`confirm.closeApp.reason.${reason}`)}`)
    .join('\n')
  return `${tr('confirm.closeApp.messageIntro')}\n${bullets}\n\n${tr('confirm.closeApp.messageOutro')}`
}
