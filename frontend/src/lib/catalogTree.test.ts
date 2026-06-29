import { describe, expect, it } from 'vitest'
import { buildCatalogTree, countFeatureFiles, type CatalogNode } from './catalogTree'

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
})
