export type FeatureTemplateOptions = {
  title: string
  scenario: string
  startUrl: string
  tag?: string
}

export function slugifyFileName(title: string): string {
  const raw = title
    .trim()
    .toLowerCase()
    .replace(/[^\p{L}\p{N}]+/gu, '-')
    .replace(/^-+|-+$/g, '')
  return raw || 'scenario'
}

export function buildFeatureTemplate(opts: FeatureTemplateOptions): string {
  const tag = opts.tag?.trim().replace(/^@/, '')
  const tagLine = tag ? `@${tag}\n` : ''
  return `# language: ru
${tagLine}Функционал: ${opts.title.trim() || 'Новый сценарий'}
  Сценарий: ${opts.scenario.trim() || 'Пример'}
    Допустим открыт "${opts.startUrl.trim() || 'https://example.com'}"
    И закрываю браузер
`
}
