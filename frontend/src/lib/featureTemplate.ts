import { getLocale, t } from './i18n'

export type FeatureTemplateOptions = {
  title: string
  scenario: string
  startUrl: string
  tag?: string
}

const STEP = '\t'

export function slugifyFileName(title: string): string {
  const raw = title
    .trim()
    .toLowerCase()
    .replace(/[^\p{L}\p{N}]+/gu, '-')
    .replace(/^-+|-+$/g, '')
  return raw || 'scenario'
}

/** Locale-aware Gherkin template (shop-ui-recorder/examples parity). */
export function buildFeatureTemplate(opts: FeatureTemplateOptions): string {
  const locale = getLocale()
  const tag = opts.tag?.trim().replace(/^@/, '')
  const tagLine = tag ? `@${tag}\n` : ''
  const url = opts.startUrl.trim() || 'https://example.com'
  const title = opts.title.trim() || t('dialogs.project.beginnerExamples.title')
  const scenario = opts.scenario.trim() || t('dialogs.project.beginnerExamples.scenario')
  const featureKw = locale === 'en' ? 'Feature' : 'Функционал'
  const scenarioKw = locale === 'en' ? 'Scenario' : 'Сценарий'
  const given = t('dialogs.record.gherkinGiven')
  const then = t('dialogs.record.gherkinThen')
  const and = t('dialogs.record.gherkinAnd')
  const openStep = t('dialogs.record.openStepTemplate', { url })
  const seeStep = t('dialogs.record.seeStepTemplate')
  const checkTextStep = t('dialogs.record.checkTextStepTemplate')
  return `# language: ${locale}
${tagLine}${featureKw}: ${title}
${scenarioKw}: ${scenario}
${STEP}${given} ${openStep}
${STEP}${then} ${seeStep}
${STEP}${and} ${checkTextStep}
`
}
