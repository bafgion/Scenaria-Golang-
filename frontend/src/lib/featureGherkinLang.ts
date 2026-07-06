/** Detect `# language: ru|en` from feature file text (mirrors backend gherkin.ParseLanguageTag). */
export type FeatureGherkinLanguage = 'ru' | 'en'

const LANGUAGE_TAG_RE = /^#\s*language\s*:\s*(ru|en)\s*$/i

export function detectFeatureGherkinLanguage(text: string): FeatureGherkinLanguage {
  for (const raw of text.split('\n')) {
    const line = raw.trim()
    if (!line) continue
    if (!line.startsWith('#')) break
    const match = LANGUAGE_TAG_RE.exec(line)
    if (match) {
      return match[1].toLowerCase() === 'en' ? 'en' : 'ru'
    }
  }
  return 'ru'
}

export function defaultStepKeyword(lang: FeatureGherkinLanguage): string {
  return lang === 'en' ? 'When' : 'Когда'
}
