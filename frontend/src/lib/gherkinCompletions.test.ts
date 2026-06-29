import { describe, expect, it } from 'vitest'
import { formatInsertText } from './gherkinCompletions'

describe('formatInsertText', () => {
  it('prefixes keyword when line has no step keyword', () => {
    const insert = formatInsertText('Функционал: demo', {
      label: 'нажимаю',
      insert: 'нажимаю "button.submit"',
      description: '',
    })
    expect(insert).toBe('Когда нажимаю "button.submit"')
  })

  it('keeps indent and inserts snippet body after keyword', () => {
    const insert = formatInsertText('\tКогда на', {
      label: 'нажимаю',
      insert: 'нажимаю "button.submit"',
      description: '',
    })
    expect(insert).toBe('нажимаю "button.submit"')
  })

  it('adds trailing space for bare keyword snippet on empty line', () => {
    const insert = formatInsertText('\t', {
      label: 'Когда',
      insert: 'Когда',
      description: '',
    })
    expect(insert).toBe('\tКогда ')
  })
})
