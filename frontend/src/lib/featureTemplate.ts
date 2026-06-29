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

/** Шаблон как в shop-ui-recorder/examples: без # language, шаги с табуляцией. */
export function buildFeatureTemplate(opts: FeatureTemplateOptions): string {
  const tag = opts.tag?.trim().replace(/^@/, '')
  const tagLine = tag ? `@${tag}\n` : ''
  const url = opts.startUrl.trim() || 'https://example.com'
  const title = opts.title.trim() || 'Примеры для новичков'
  const scenario = opts.scenario.trim() || 'Первая проверка страницы'
  return `${tagLine}Функционал: ${title}
Сценарий: ${scenario}
${STEP}Допустим открыт "${url}"
${STEP}Тогда вижу "h1"
${STEP}И проверяю текст "Example Domain" в "h1"
${STEP}И закрываю браузер
`
}
