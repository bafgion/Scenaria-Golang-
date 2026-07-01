import { describe, expect, it } from 'vitest'
import {
  buildCatalogTree,
  buildCatalogStructure,
  buildCatalogViewStateFromBase,
  buildRunByPathMap,
  catalogStructureKey,
  collectFeaturePathsUnder,
  countFeatureFiles,
  featurePathsEqual,
  fileTreeLabel,
  mapTreeRunStatus,
  isFeaturePathSelected,
  type CatalogNode,
} from './catalogTree'

describe('catalogTree', () => {
  it('counts feature files in tree', () => {
    const tree: CatalogNode = {
      kind: 'root',
      path: '/proj',
      name: 'proj',
      children: [
        {
          kind: 'dir',
          path: '/proj/a',
          name: 'a',
          children: [
            {
              kind: 'file',
              path: '/proj/a/one.feature',
              name: 'one.feature',
              children: [],
              runSuccess: null,
              runMessage: '',
              runAt: '',
              runRunner: '',
            },
          ],
          runSuccess: null,
          runMessage: '',
          runAt: '',
          runRunner: '',
        },
        {
          kind: 'file',
          path: '/proj/two.feature',
          name: 'two.feature',
          children: [],
          runSuccess: null,
          runMessage: '',
          runAt: '',
          runRunner: '',
        },
      ],
      runSuccess: null,
      runMessage: '',
      runAt: '',
      runRunner: '',
    }
    expect(countFeatureFiles(tree)).toBe(2)
  })

  it('builds tree from flat feature paths', () => {
    const tree = buildCatalogTree('/project', [
      '/project/shop/cart.feature',
      '/project/shop/login.feature',
    ], new Map())
    expect(tree).not.toBeNull()
    expect(countFeatureFiles(tree)).toBe(2)
    expect(tree?.children.some((n) => n.name === 'shop')).toBe(true)
  })

  it('matches feature paths regardless of slash style', () => {
    expect(featurePathsEqual('C:\\proj\\a.feature', 'c:/proj/a.feature')).toBe(true)
    expect(isFeaturePathSelected('C:\\proj\\a.feature', ['c:/proj/a.feature'])).toBe(true)
  })

  it('collects all feature paths under tree node', () => {
    const tree = buildCatalogTree('/project', ['/project/a.feature', '/project/b.feature'], new Map())
    expect(collectFeaturePathsUnder(tree!)).toHaveLength(2)
  })

  it('shows checkbox marks in batch selection mode', () => {
    const node: CatalogNode = {
      kind: 'file',
      path: '/project/a.feature',
      name: 'a.feature',
      children: [],
      runSuccess: null,
      runMessage: '',
      runAt: '',
      runRunner: '',
    }
    expect(fileTreeLabel(node, true, true)).toContain('☑')
    expect(fileTreeLabel(node, true, false)).toContain('☐')
    expect(fileTreeLabel(node, true, true, true)).toBe('☑ a.feature')
  })

  it('catalogStructureKey changes only when feature list changes', () => {
    const a = catalogStructureKey('/p', ['/p/a.feature', '/p/b.feature'])
    const b = catalogStructureKey('/p', ['/p/a.feature', '/p/b.feature'])
    const c = catalogStructureKey('/p', ['/p/c.feature'])
    expect(a).toBe(b)
    expect(a).not.toBe(c)
  })

  it('buildCatalogViewStateFromBase reuses structure and applies run status', () => {
    const base = buildCatalogStructure('/project', ['/project/a.feature', '/project/b.feature'])
    const runMap = buildRunByPathMap([
      { path: '/project/a.feature', success: true, at: 't1', runner: 'pw' },
    ])
    const view = buildCatalogViewStateFromBase('/project', base, '', runMap)
    const fileA = view.tree?.children.find((n) => n.name === 'a')
    const fileB = view.tree?.children.find((n) => n.name === 'b')
    expect(fileA?.runSuccess).toBe(true)
    expect(fileB?.runSuccess).toBe(null)
  })

  it('mapTreeRunStatus does not rebuild directory structure', () => {
    const base = buildCatalogStructure('/project', ['/project/x.feature'])
    const withRun = mapTreeRunStatus(base, buildRunByPathMap([{ path: '/project/x.feature', success: false }]))
    expect(withRun.children[0].runSuccess).toBe(false)
    expect(base.children[0].runSuccess).toBe(null)
  })
})
