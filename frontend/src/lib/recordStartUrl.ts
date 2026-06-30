const GOTO_LINE =
  /^\s*(?:(?:Допустим|Дано|Когда|Тогда|И|Но|Given|When|Then|And|But)\s+)?(?:открыт|открываю|перехожу на)\s+"((?:[^"\\]|\\.)*)"/i

/** First «открыт/открываю» URL in a .feature buffer. */
export function firstGotoURLFromText(text: string): string {
  for (const line of text.split('\n')) {
    const match = line.match(GOTO_LINE)
    if (match?.[1]) {
      return match[1].replace(/\\"/g, '"')
    }
  }
  return ''
}

export function resolveRecordStartURL(opts: {
  editorText?: string
  startURL?: string
  recordURL?: string
  lastRunBaseUrl?: string
}): string {
  const fromEditor = firstGotoURLFromText(opts.editorText ?? '')
  if (fromEditor) return fromEditor

  const manual = (opts.startURL ?? '').trim() || (opts.recordURL ?? '').trim()
  if (manual) return manual

  const fromRun = (opts.lastRunBaseUrl ?? '').trim()
  if (fromRun) return fromRun

  return ''
}
