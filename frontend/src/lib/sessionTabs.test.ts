import { describe, expect, it } from 'vitest'
import {
  buildSessionTabsSnapshot,
  resolveRestoredActiveTab,
  sessionTabPathsFromSettings,
  untitledContentMap,
} from './sessionTabs'
import { UNTITLED_PREFIX } from './untitled'

describe('buildSessionTabsSnapshot', () => {
  it('includes untitled tabs in open order and bodies', () => {
    const untitled = `${UNTITLED_PREFIX}2/demo.feature`
    const snap = buildSessionTabsSnapshot(
      [
        { path: 'C:/proj/a.feature', content: 'a', dirty: false },
        { path: untitled, content: 'old', dirty: true, draft: 'draft text' },
      ],
      untitled,
      () => 'draft text',
      '__welcome__',
    )
    expect(snap.openTabs).toEqual(['C:/proj/a.feature', untitled])
    expect(snap.untitledTabs).toEqual([{ path: untitled, content: 'draft text' }])
    expect(snap.activeTab).toBe(untitled)
  })

  it('prefers live editor text for active untitled tab when available', () => {
    const untitled = `${UNTITLED_PREFIX}2/demo.feature`
    const snap = buildSessionTabsSnapshot(
      [{ path: untitled, content: 'old', dirty: true, draft: 'draft text' }],
      untitled,
      () => 'live text',
      '__welcome__',
    )
    expect(snap.untitledTabs).toEqual([{ path: untitled, content: 'live text' }])
  })

  it('preserves an intentional empty live edit for the active untitled tab', () => {
    const untitled = `${UNTITLED_PREFIX}2/demo.feature`
    const snap = buildSessionTabsSnapshot(
      [{ path: untitled, content: 'old', dirty: true, draft: 'old' }],
      untitled,
      () => '',
      '__welcome__',
    )
    expect(snap.untitledTabs).toEqual([{ path: untitled, content: '' }])
  })

  it('falls back to stored untitled text when the active Monaco model is for another path', () => {
    const untitled = `${UNTITLED_PREFIX}2/demo.feature`
    const snap = buildSessionTabsSnapshot(
      [{ path: untitled, content: 'old', dirty: true, draft: 'draft text' }],
      untitled,
      () => null,
      '__welcome__',
    )
    expect(snap.untitledTabs).toEqual([{ path: untitled, content: 'draft text' }])
  })
})

describe('sessionTabPathsFromSettings', () => {
  it('falls back to untitled tab paths when openTabs is empty', () => {
    const paths = sessionTabPathsFromSettings([], [{ path: `${UNTITLED_PREFIX}1/x.feature`, content: 'x' }])
    expect(paths).toEqual([`${UNTITLED_PREFIX}1/x.feature`])
  })

  it('merges untitled tab paths when openTabs only has saved feature files', () => {
    const untitled = `${UNTITLED_PREFIX}2/demo.feature`
    const paths = sessionTabPathsFromSettings(
      ['C:/proj/a.feature'],
      [{ path: untitled, content: 'draft' }],
    )
    expect(paths).toEqual(['C:/proj/a.feature', untitled])
  })

  it('deduplicates open and untitled paths deterministically', () => {
    const untitled = `${UNTITLED_PREFIX}2/demo.feature`
    const paths = sessionTabPathsFromSettings(
      [' C:/proj/a.feature ', untitled, 'C:/proj/a.feature'],
      [
        { path: untitled, content: 'draft' },
        { path: `${UNTITLED_PREFIX}3/other.feature`, content: 'other' },
      ],
    )
    expect(paths).toEqual(['C:/proj/a.feature', untitled, `${UNTITLED_PREFIX}3/other.feature`])
  })
})

describe('resolveRestoredActiveTab', () => {
  const welcomeKey = '__welcome__'
  const tabs = [
    { path: 'C:/proj/a.feature', content: 'a', dirty: false },
    { path: 'C:/proj/b.feature', content: 'b', dirty: false },
  ]

  it('ignores saved welcome active tab when feature tabs exist', () => {
    expect(resolveRestoredActiveTab(welcomeKey, ['C:/proj/a.feature'], tabs, welcomeKey)).toBe('C:/proj/a.feature')
    expect(resolveRestoredActiveTab('', ['C:/proj/a.feature'], tabs, welcomeKey)).toBe('C:/proj/a.feature')
    expect(resolveRestoredActiveTab('C:/proj/a.feature', ['C:/proj/a.feature', 'C:/proj/b.feature'], tabs, welcomeKey)).toBe('C:/proj/a.feature')
  })

  it('restores focus to saved untitled tab', () => {
    const untitled = `${UNTITLED_PREFIX}3/demo.feature`
    const restored = [
      { path: 'C:/proj/a.feature', content: 'a', dirty: false },
      { path: untitled, content: 'draft', dirty: true },
    ]
    expect(
      resolveRestoredActiveTab(untitled, ['C:/proj/a.feature', untitled], restored, welcomeKey),
    ).toBe(untitled)
  })
})

describe('untitledContentMap', () => {
  it('maps path to content', () => {
    const map = untitledContentMap([{ path: 'p', content: 'body' }])
    expect(map.get('p')).toBe('body')
  })

  it('keeps the latest non-empty content for duplicate untitled paths', () => {
    const map = untitledContentMap([
      { path: 'p', content: 'first' },
      { path: 'p', content: '' },
      { path: 'p', content: 'second' },
    ])
    expect(map.get('p')).toBe('second')
  })

  it('does not replace non-empty duplicate content with a later empty entry', () => {
    const map = untitledContentMap([
      { path: 'p', content: 'first' },
      { path: 'p', content: '   ' },
    ])
    expect(map.get('p')).toBe('first')
  })
})
