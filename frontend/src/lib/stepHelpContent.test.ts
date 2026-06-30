import { describe, expect, it } from 'vitest'
import { formatStepHoverMarkdown } from './stepHelpContent'
import type { StepHelpEntry } from './stepTypes'

const clickEntry: StepHelpEntry = {
  label: 'нажимаю',
  action: 'click',
  category: 'Формы и ввод',
  description: 'Клик по элементу',
  template: 'нажимаю "button.submit"',
  example: 'нажимаю "button.submit"',
  parameters: ['selector — CSS/XPath селектор элемента'],
  help: 'Клик по элементу',
}

describe('stepHelpContent', () => {
  it('formatStepHoverMarkdown includes title, params and example', () => {
    const md = formatStepHoverMarkdown(clickEntry)
    expect(md).toContain('**нажимаю (click)**')
    expect(md).toContain('Формы и ввод')
    expect(md).toContain('**Параметры**')
    expect(md).toContain('selector')
    expect(md).toContain('button.submit')
    expect(md).toContain('F1')
  })
})
