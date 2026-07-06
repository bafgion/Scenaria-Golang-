import { t } from './i18n'
import { ru } from './i18n/locales/ru'

export type FeatureTemplateOptions = {
  title: string
  scenario: string
  startUrl: string
  tag?: string
}

export const DEFAULT_FEATURE_GHERKIN_LANGUAGE = 'ru' as const

const STEP = '\t'

function ruT(key: string, params?: Record<string, string | number>): string {
  return t(key, params, ru)
}

export function slugifyFileName(title: string): string {
  const raw = title
    .trim()
    .toLowerCase()
    .replace(/[^\p{L}\p{N}]+/gu, '-')
    .replace(/^-+|-+$/g, '')
  return raw || 'scenario'
}

/** New feature scaffold — always Russian Gherkin (`# language: ru`), regardless of UI locale. */
export function buildFeatureTemplate(opts: FeatureTemplateOptions): string {
  const tag = opts.tag?.trim().replace(/^@/, '')
  const tagLine = tag ? `@${tag}\n` : ''
  const url = opts.startUrl.trim() || 'https://example.com'
  const title = opts.title.trim() || ruT('dialogs.project.beginnerExamples.title')
  const scenario = opts.scenario.trim() || ruT('dialogs.project.beginnerExamples.scenario')
  const given = ruT('dialogs.record.gherkinGiven')
  const then = ruT('dialogs.record.gherkinThen')
  const and = ruT('dialogs.record.gherkinAnd')
  const openStep = ruT('dialogs.record.openStepTemplate', { url })
  const seeStep = ruT('dialogs.record.seeStepTemplate')
  const checkTextStep = ruT('dialogs.record.checkTextStepTemplate')
  return `# language: ru
${tagLine}Функционал: ${title}
Сценарий: ${scenario}
${STEP}${given} ${openStep}
${STEP}${then} ${seeStep}
${STEP}${and} ${checkTextStep}
`
}

/** Gherkin preview for baseline recording — same default dialect as new features. */
export function buildGherkinPreview(feature: string, scenario: string, stepBodies: string[]): string {
  const title = feature.trim() || ruT('dialogs.record.featureDefault')
  const scen = scenario.trim() || ruT('dialogs.record.baselineScenarioDefault')
  const lines = stepBodies
    .map((s, i) => `    ${i === 0 ? ruT('dialogs.record.gherkinGiven') : ruT('dialogs.record.gherkinAnd')} ${s.trim()}`)
    .join('\n')
  return `# language: ru
Функционал: ${title}
  Сценарий: ${scen}
${lines || `    ${ruT('dialogs.record.gherkinFallbackStep')}`}`
}

export function defaultOpenStepTemplate(url: string): string {
  return ruT('dialogs.record.openStepTemplate', { url })
}
