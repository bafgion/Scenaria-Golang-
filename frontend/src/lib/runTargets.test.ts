import { describe, expect, it } from 'vitest'
import {
  editorTextForRunTarget,
  findTabForRunTarget,
  materializeRunTargetPaths,
  resolveLogicalRunTarget,
} from './runTargets'
import type { TabBody } from './tabMemory'
import { UNTITLED_PREFIX } from './untitled'

const untitled = `${UNTITLED_PREFIX}1/novyy-scenariy.feature`

function tab(path: string, content: string, extra: Partial<TabBody> = {}): TabBody {
  return { path, content, dirty: false, ...extra }
}

describe('findTabForRunTarget', () => {
  const tabs = [tab(untitled, 'body', { dirty: true, draft: 'draft' })]

  it('finds tab by internal untitled path', () => {
    expect(findTabForRunTarget(untitled, tabs, untitled)?.path).toBe(untitled)
  })

  it('finds tab by display label from results panel', () => {
    expect(findTabForRunTarget('novyy-scenariy.feature', tabs, untitled)?.path).toBe(untitled)
  })
})

describe('resolveLogicalRunTarget', () => {
  it('maps display label to untitled id', () => {
    const tabs = [tab(untitled, 'x', { dirty: true })]
    expect(resolveLogicalRunTarget('novyy-scenariy.feature', tabs, untitled)).toBe(untitled)
  })
})

describe('editorTextForRunTarget', () => {
  it('uses live editor text for the active untitled tab', () => {
    const tabs = [tab(untitled, 'saved', { dirty: true, draft: 'stale' })]
    expect(editorTextForRunTarget(untitled, tabs, untitled, 'fallback', 'live')).toBe('live')
  })

  it('uses draft for inactive dirty tabs', () => {
    const saved = tab('C:/proj/a.feature', 'saved')
    const dirty = tab(untitled, 'saved', { dirty: true, draft: 'edited' })
    expect(editorTextForRunTarget(untitled, [saved, dirty], saved.path, 'active')).toBe('edited')
  })
})

describe('materializeRunTargetPaths', () => {
  it('writes temp file for untitled tabs', async () => {
    const temps: string[] = []
    const disk = await materializeRunTargetPaths(
      ['novyy-scenariy.feature'],
      [tab(untitled, 'old', { dirty: true, draft: 'Функционал: X' })],
      untitled,
      'fallback',
      'live body',
      {
        writeTempFeature: async (content) => {
          temps.push(content)
          return 'C:/proj/.scenaria/temp/run-1/scenario.feature'
        },
        readFeature: async () => '',
      },
    )
    expect(temps).toEqual(['live body'])
    expect(disk).toEqual(['C:/proj/.scenaria/temp/run-1/scenario.feature'])
  })

  it('keeps clean saved paths on disk', async () => {
    const disk = await materializeRunTargetPaths(
      ['C:/proj/a.feature'],
      [tab('C:/proj/a.feature', 'content')],
      untitled,
      '',
      undefined,
      {
        writeTempFeature: async () => 'temp',
        readFeature: async () => '',
      },
    )
    expect(disk).toEqual(['C:/proj/a.feature'])
  })
})
