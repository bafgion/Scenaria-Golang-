import { describe, expect, it } from 'vitest'
import { batchRunFormFrom, defaultRunForm, runFormFromMode, type RunForm } from './runTypes'

describe('batchRunFormFrom', () => {
  it('clears stale scenario, tag and step range while preserving other run settings', () => {
    const lastRun: RunForm = defaultRunForm({
      tag: '@smoke',
      scenario: 'Login flow',
      dryRun: false,
      headed: true,
      installPW: true,
      html: true,
      browser: 'firefox',
      baseUrl: 'https://example.com',
      startStep: 2,
      endStep: 4,
    })

    const batch = batchRunFormFrom(lastRun, true)

    expect(batch.dryRun).toBe(true)
    expect(batch.headed).toBe(true)
    expect(batch.installPW).toBe(true)
    expect(batch.html).toBe(true)
    expect(batch.browser).toBe('firefox')
    expect(batch.baseUrl).toBe('https://example.com')
    expect(batch.tag).toBe('')
    expect(batch.scenario).toBe('')
    expect(batch.startStep).toBe(-1)
    expect(batch.endStep).toBe(-1)
    expect(lastRun.tag).toBe('@smoke')
    expect(lastRun.scenario).toBe('Login flow')
    expect(lastRun.startStep).toBe(2)
    expect(lastRun.endStep).toBe(4)
  })
})

describe('runFormFromMode', () => {
  const lastRun: RunForm = defaultRunForm({
    tag: '@smoke',
    scenario: 'Login flow',
    dryRun: false,
    headed: true,
    installPW: true,
    html: true,
    browser: 'firefox',
    baseUrl: 'https://example.com',
    startStep: 2,
    endStep: 4,
  })

  it('prepares single-mode run forms without stale batch filters', () => {
    const single = runFormFromMode(lastRun, 'single', { dryRun: true, scenario: 'Checkout flow' })

    expect(single.dryRun).toBe(true)
    expect(single.scenario).toBe('Checkout flow')
    expect(single.tag).toBe('')
    expect(single.startStep).toBe(-1)
    expect(single.endStep).toBe(-1)
    expect(single.browser).toBe('firefox')
  })

  it('prepares tag-mode run forms by clearing scenario and step range', () => {
    const tag = runFormFromMode(lastRun, 'tag', { dryRun: true, tag: '@api' })

    expect(tag.dryRun).toBe(true)
    expect(tag.tag).toBe('@api')
    expect(tag.scenario).toBe('')
    expect(tag.startStep).toBe(-1)
    expect(tag.endStep).toBe(-1)
  })

  it('prepares step-range run forms by clearing stale tag filters', () => {
    const range = runFormFromMode(lastRun, 'step-range', {
      dryRun: false,
      scenario: 'Login flow',
      startStep: 1,
      endStep: 3,
    })

    expect(range.tag).toBe('')
    expect(range.scenario).toBe('Login flow')
    expect(range.startStep).toBe(1)
    expect(range.endStep).toBe(3)
  })
})
