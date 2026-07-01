import { describe, expect, it } from 'vitest'
import {
  rebuildLiveRecordStepLines,
  upsertRecordedStepInText,
} from './recordedStepEditor'

const sample = `Функционал: Demo
Сценарий: Main
\tДопустим шаблон
\tТогда шаблон
`

function applyRecordedSteps(
  text: string,
  steps: string[],
  lineByIndex: Record<number, number>,
): { text: string; lineByIndex: Record<number, number> } {
  let nextText = text
  let nextMap = lineByIndex
  for (let i = 0; i < steps.length; i++) {
    const result = upsertRecordedStepInText(nextText, i, steps[i], nextMap)
    nextText = result.text
    nextMap = result.lineByIndex
  }
  return { text: nextText, lineByIndex: nextMap }
}

describe('record stop starts fresh segment', () => {
  const firstSegment = [
    'открыт "https://example.com"',
    'нажимаю "#one"',
    'нажимаю "#two"',
  ]

  it('new segment after stop appends from index 0 without duplicating prior steps', () => {
    const first = applyRecordedSteps(sample, firstSegment, {})
    const afterStopMap: Record<number, number> = {}

    const secondSegment = ['нажимаю "#three"', 'нажимаю "#four"']
    const second = applyRecordedSteps(first.text, secondSegment, afterStopMap)

    expect(second.text).toContain('нажимаю "#one"')
    expect(second.text).toContain('нажимаю "#three"')
    expect(second.text.split('нажимаю "#one"').length - 1).toBe(1)
    expect(second.text.split('нажимаю "#three"').length - 1).toBe(1)
    expect(second.lineByIndex[0]).toBeGreaterThan(first.lineByIndex[firstSegment.length - 1]!)
  })

  it('pause keeps line map so same indices update in place', () => {
    let text = sample
    let map: Record<number, number> = {}
    const r1 = upsertRecordedStepInText(text, 0, 'нажимаю "#a"', map)
    text = r1.text
    map = r1.lineByIndex
    const r2 = upsertRecordedStepInText(text, 0, 'нажимаю "#a-updated"', map)
    expect(r2.text.split('нажимаю "#a"').length - 1).toBe(0)
    expect(r2.text).toContain('нажимаю "#a-updated"')
    expect(Object.keys(r2.lineByIndex)).toHaveLength(1)
  })
})

describe('Chaos record stop segment editor', () => {
  const bodies = [
    'нажимаю "#a"',
    'ввожу "x" в "#email"',
    'открыт "https://site.test"',
    'нажимаю "button.submit"',
  ]

  for (let seed = 0; seed < 30; seed++) {
    it(`seed-${seed} stop resets session index without duplicating editor`, () => {
      let text = sample
      let map: Record<number, number> = {}
      const segmentA: string[] = []

      const roundsA = 1 + (seed % 4)
      for (let i = 0; i < roundsA; i++) {
        segmentA.push(`${bodies[(seed + i) % bodies.length]} /*a${i}*/`)
      }
      const appliedA = applyRecordedSteps(text, segmentA, map)
      text = appliedA.text
      const linesAfterA = text.split('\n').length

      map = {}
      const segmentB: string[] = []
      const roundsB = 1 + ((seed + 1) % 3)
      for (let i = 0; i < roundsB; i++) {
        segmentB.push(`${bodies[(seed + roundsA + i) % bodies.length]} /*b${i}*/`)
      }
      const appliedB = applyRecordedSteps(text, segmentB, map)
      expect(appliedB.text.split('\n').length).toBeGreaterThan(linesAfterA)
      expect(appliedB.text).toContain(segmentA[0])
      for (const step of segmentA) {
        expect(appliedB.text.split(step).length - 1).toBe(1)
      }
    })
  }
})

describe('rebuildLiveRecordStepLines', () => {
  it('still maps trailing steps for pause-time updates', () => {
    const text = `Функционал: X
Сценарий: Y
\tДопустим шаблон
\tКогда записан один
\tИ записан два
`
    expect(rebuildLiveRecordStepLines(text, 2)).toEqual({ 0: 3, 1: 4 })
  })
})
