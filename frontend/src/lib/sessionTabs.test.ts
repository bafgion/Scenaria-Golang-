import { describe, expect, it } from 'vitest'
import {
  buildSessionTabsSnapshot,
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
      'draft text',
      '__welcome__',
    )
    expect(snap.openTabs).toEqual(['C:/proj/a.feature', untitled])
    expect(snap.untitledTabs).toEqual([{ path: untitled, content: 'draft text' }])
    expect(snap.activeTab).toBe(untitled)
  })
})

describe('sessionTabPathsFromSettings', () => {
  it('falls back to untitled tab paths when openTabs is empty', () => {
    const paths = sessionTabPathsFromSettings([], [{ path: `${UNTITLED_PREFIX}1/x.feature`, content: 'x' }])
    expect(paths).toEqual([`${UNTITLED_PREFIX}1/x.feature`])
  })
})

describe('untitledContentMap', () => {
  it('maps path to content', () => {
    const map = untitledContentMap([{ path: 'p', content: 'body' }])
    expect(map.get('p')).toBe('body')
  })
})
