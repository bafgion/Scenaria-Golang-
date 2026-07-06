import type { RunForm } from './runTypes'
import { t } from './i18n'

/** Compact label for last run options (status bar / toolbar). */
export function formatLastRunSummary(run: RunForm): string {
  const parts: string[] = []
  if (run.dryRun) {
    parts.push('dry-run')
  } else {
    parts.push(run.headed ? t('statusBar.runSummary.headed') : t('statusBar.runSummary.headless'))
  }
  if (run.html) parts.push(run.htmlLightMode ? 'HTML-light' : 'HTML')
  if (run.trace) parts.push('trace')
  if (run.video) parts.push('video')
  if (run.junit) parts.push('JUnit')
  if (run.continueOnFail) parts.push(t('statusBar.runSummary.allScenarios'))
  if (run.htmlTimestamp) parts.push('HTML+time')
  if (run.workers > 1) parts.push(t('statusBar.runSummary.workers', { count: run.workers }))
  if (run.slowMo > 0) parts.push(t('statusBar.runSummary.slowMo', { ms: run.slowMo }))
  if (run.scenario) parts.push(`«${run.scenario}»`)
  else if (run.tag) parts.push(`@${run.tag}`)
  return parts.join(' · ')
}
