import { describe, expect, it } from 'vitest'
import { redactSecrets } from './redaction'

describe('redactSecrets', () => {
  it('redacts common secret forms while keeping useful context', () => {
    const got = redactSecrets(
      'url=https://user:pass@example.com/path password="open sesame" Authorization: Bearer abc.def Cookie: sid=123',
    )
    expect(got).not.toContain('user:pass@')
    expect(got).not.toContain('open sesame')
    expect(got).not.toContain('abc.def')
    expect(got).not.toContain('sid=123')
    expect(got).toContain('https://example.com/path')
    expect(got).toContain('[REDACTED]')
  })
})
