import { describe, expect, it } from 'vitest'
import { formatStepInlayLabel } from './gherkinInlayHints'

describe('formatStepInlayLabel', () => {
  it('formats click and fill hints', () => {
    expect(
      formatStepInlayLabel({
        line: 1,
        kind: 'click',
        action: 'Нажать',
        element: 'button.submit',
        value: '',
        error: '',
      }),
    ).toBe('click → button.submit')

    expect(
      formatStepInlayLabel({
        line: 2,
        kind: 'fill',
        action: 'Ввести',
        element: 'input[name=email]',
        value: 'text',
        error: '',
      }),
    ).toBe('fill → input[name=email] / "text"')
  })

  it('shows parse errors', () => {
    expect(
      formatStepInlayLabel({
        line: 3,
        kind: '',
        action: 'шаг',
        element: '',
        value: '',
        error: 'unknown step',
      }),
    ).toBe('⚠ unknown step')
  })
})
