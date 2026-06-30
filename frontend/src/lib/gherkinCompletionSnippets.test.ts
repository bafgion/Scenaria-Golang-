import { describe, expect, it } from 'vitest'
import {
  completionFilterText,
  completionSortKey,
  shouldPreselectCompletion,
  snippetizeInsert,
  usesSnippetTabStops,
} from './gherkinCompletionSnippets'

describe('snippetizeInsert', () => {
  it('creates tab stop for selector', () => {
    const out = snippetizeInsert('нажимаю "button.submit"')
    expect(out).toBe('нажимаю "${1:selector}"$0')
    expect(usesSnippetTabStops(out)).toBe(true)
  })

  it('creates tab stops for fill step', () => {
    const out = snippetizeInsert('ввожу "текст" в "input[name=email]"')
    expect(out).toBe('ввожу "${1:text}" в "${2:selector}"$0')
  })

  it('uses example url placeholder', () => {
    const out = snippetizeInsert('открыт "https://site.com"')
    expect(out).toBe('открыт "${1:https://example.com}"$0')
  })

  it('leaves plain steps unchanged', () => {
    expect(snippetizeInsert('обновляю страницу')).toBe('обновляю страницу')
  })
})

describe('completion ranking helpers', () => {
  it('sorts exact prefix first', () => {
    expect(completionSortKey('нажимаю', 'наж')).toBe('01_нажимаю')
    expect(completionSortKey('нажимаю', 'нажимаю')).toBe('00_нажимаю')
  })

  it('builds filter text from label and insert', () => {
    expect(
      completionFilterText({
        label: 'нажимаю',
        insert: 'нажимаю "x"',
        description: 'Клик',
      }),
    ).toContain('клик')
  })

  it('preselects only first exact-prefix match', () => {
    expect(shouldPreselectCompletion('нажимаю', 'наж', 0)).toBe(true)
    expect(shouldPreselectCompletion('нажимаю', 'наж', 1)).toBe(false)
  })
})
