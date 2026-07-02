import type { RunForm } from './runTypes'

/** Compact label for last run options (status bar / toolbar). */
export function formatLastRunSummary(run: RunForm): string {
  const parts: string[] = []
  if (run.dryRun) {
    parts.push('dry-run')
  } else {
    parts.push(run.headed ? 'с окном' : 'headless')
  }
  if (run.html) parts.push('HTML')
  if (run.trace) parts.push('trace')
  if (run.video) parts.push('video')
  if (run.junit) parts.push('JUnit')
  if (run.continueOnFail) parts.push('все сценарии')
  if (run.htmlTimestamp) parts.push('HTML+time')
  if (run.workers > 1) parts.push(`${run.workers} ворк.`)
  if (run.slowMo > 0) parts.push(`slow ${run.slowMo}мс`)
  if (run.scenario) parts.push(`«${run.scenario}»`)
  else if (run.tag) parts.push(`@${run.tag}`)
  return parts.join(' · ')
}
