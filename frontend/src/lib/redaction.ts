const REDACTION = '[REDACTED]'

const urlPattern = /\b[a-z][a-z0-9+.-]*:\/\/[^\s"'<>]+/gi
const bearerPattern = /\b(Bearer|Basic)\s+[A-Za-z0-9._~+/=-]+/gi
const keyValuePattern =
  /\b(password|passwd|pwd|secret|token|authorization|api[_-]?key|cookie|credential|scenaria_email_code|email_code|http_auth)(\s*[:=]\s*)("[^"]*"|'[^']*'|[^\s,;]+)/gi

export function redactSecrets(text: string): string {
  if (!text) return ''
  return text
    .replace(urlPattern, (raw) => stripUrlCredentials(raw))
    .replace(bearerPattern, (_match, scheme) => `${scheme} ${REDACTION}`)
    .replace(keyValuePattern, (_match, key, separator) => `${key}${separator}${REDACTION}`)
}

function stripUrlCredentials(raw: string): string {
  try {
    const parsed = new URL(raw)
    if (!parsed.username && !parsed.password) return raw
    parsed.username = ''
    parsed.password = ''
    return parsed.toString()
  } catch {
    return raw
  }
}
