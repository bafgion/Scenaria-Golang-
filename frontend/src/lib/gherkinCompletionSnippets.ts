const URL_RE = /^https?:\/\//i
const SELECTOR_RE = /[#.\[\]=:>]|input|button|select|canvas|textarea/i
const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

/** Turn quoted literals in a step template into Monaco tab stops. */
export function snippetizeInsert(insert: string): string {
  let tabIndex = 0
  const snippet = insert.replace(/"((?:[^"\\]|\\.)*)"/g, (_match, inner: string) => {
    tabIndex++
    const placeholder = quotedPlaceholder(inner, tabIndex)
    return `"$\{${tabIndex}:${placeholder}}"`
  })
  if (tabIndex === 0) {
    return snippetizeHeaderPlaceholder(insert)
  }
  return `${snippet}$0`
}

function snippetizeHeaderPlaceholder(insert: string): string {
  if (insert === 'Функционал: UI сценарий') {
    return 'Функционал: ${1:Название feature}$0'
  }
  if (insert === 'Сценарий: Имя сценария') {
    return 'Сценарий: ${1:Имя сценария}$0'
  }
  return insert
}

function quotedPlaceholder(inner: string, tabIndex: number): string {
  const decoded = inner.replace(/\\"/g, '"').replace(/\\\\/g, '\\')
  if (URL_RE.test(decoded)) return 'https://example.com'
  if (EMAIL_RE.test(decoded)) return 'user@example.com'
  if (SELECTOR_RE.test(decoded)) return 'selector'
  if (/^\d+$/.test(decoded)) return decoded
  if (decoded.includes('{{') && decoded.includes('}}')) return decoded
  if (/^[A-Za-z]:\\/.test(decoded) || decoded.includes('\\')) return 'C:\\\\path\\\\file.pdf'
  if (decoded.length > 36) return decoded.slice(0, 33) + '...'
  if (tabIndex === 1 && /^(текст|text)$/i.test(decoded)) return 'text'
  return decoded || 'text'
}

export function completionSortKey(label: string, typedPrefix: string): string {
  const labelLower = label.toLowerCase()
  const prefix = typedPrefix.trim().toLowerCase()
  if (prefix && labelLower === prefix) return `00_${labelLower}`
  if (prefix && labelLower.startsWith(prefix)) return `01_${labelLower}`
  if (prefix && labelLower.includes(prefix)) return `02_${labelLower}`
  return `03_${labelLower}`
}

export function completionFilterText(item: {
  label: string
  insert: string
  description?: string
}): string {
  return [item.label, item.insert, item.description || ''].join(' ').toLowerCase()
}

export function shouldPreselectCompletion(label: string, typedPrefix: string, index: number): boolean {
  if (index !== 0) return false
  const prefix = typedPrefix.trim().toLowerCase()
  if (!prefix) return false
  return label.toLowerCase().startsWith(prefix)
}

export function usesSnippetTabStops(snippet: string): boolean {
  return /\$\{\d+:/.test(snippet) || snippet.endsWith('$0')
}
