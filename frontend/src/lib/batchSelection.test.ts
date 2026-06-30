import { describe, expect, it } from 'vitest'
import { buildCatalogTree } from './catalogTree'
import {
  buildBatchSelectedSet,
  isPathBatchSelected,
  resolveCatalogFileClickAction,
  selectAllFeaturesUnder,
  toggleBatchModeState,
  toggleBatchPath,
} from './batchSelection'

describe('batchSelection', () => {
  const tree = buildCatalogTree(
    '/project',
    ['/project/a.feature', '/project/shop/b.feature', '/project/shop/c.feature'],
    new Map(),
  )!

  it('toggleBatchModeState enables mode and selects all features in tree', () => {
    const off = toggleBatchModeState(false, tree)
    expect(off.batchMode).toBe(true)
    expect(off.batchSelected).toHaveLength(3)
    expect(off.batchSelected).toContain('/project/a.feature')

    const on = toggleBatchModeState(true, tree)
    expect(on.batchMode).toBe(false)
    expect(on.batchSelected).toEqual([])
  })

  it('toggleBatchModeState with empty tree selects nothing', () => {
    const result = toggleBatchModeState(false, null)
    expect(result.batchMode).toBe(true)
    expect(result.batchSelected).toEqual([])
  })

  it('selectAllFeaturesUnder returns every feature path', () => {
    expect(selectAllFeaturesUnder(tree)).toEqual([
      '/project/a.feature',
      '/project/shop/b.feature',
      '/project/shop/c.feature',
    ])
  })

  it('toggleBatchPath adds and removes with normalized paths', () => {
    let selected = toggleBatchPath([], 'C:\\project\\a.feature')
    expect(selected).toEqual(['C:\\project\\a.feature'])

    selected = toggleBatchPath(selected, 'c:/project/a.feature')
    expect(selected).toEqual([])

    selected = toggleBatchPath(selected, '/project/b.feature')
    expect(selected).toEqual(['/project/b.feature'])
  })

  it('buildBatchSelectedSet and isPathBatchSelected use O(1) lookup', () => {
    const set = buildBatchSelectedSet(['C:\\proj\\a.feature', '/proj/b.feature'])
    expect(set.size).toBe(2)
    expect(isPathBatchSelected('c:/proj/a.feature', set)).toBe(true)
    expect(isPathBatchSelected('/proj/c.feature', set)).toBe(false)
  })

  it('re-toggling batch mode off clears selection without touching tree', () => {
    const enabled = toggleBatchModeState(false, tree)
    const disabled = toggleBatchModeState(enabled.batchMode, tree)
    expect(disabled.batchSelected).toEqual([])
    expect(selectAllFeaturesUnder(tree)).toHaveLength(3)
  })

  it('resolveCatalogFileClickAction: batch mode toggles without opening file', () => {
    expect(resolveCatalogFileClickAction(true, false, false)).toBe('toggle-batch')
    expect(resolveCatalogFileClickAction(false, true, false)).toBe('toggle-batch')
    expect(resolveCatalogFileClickAction(false, false, true)).toBe('toggle-batch')
    expect(resolveCatalogFileClickAction(false, false, false)).toBe('open')
  })
})
