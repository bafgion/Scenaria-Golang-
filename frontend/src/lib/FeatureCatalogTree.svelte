<script lang="ts">
  import { type CatalogNode } from './catalogTree'
  import CatalogTreeNode from './CatalogTreeNode.svelte'

  export let tree: CatalogNode
  export let activeFeature = ''
  export let batchSelected: string[] = []
  export let batchMode = false
  export let expandAll = false
  export let collapsed: Set<string> = new Set()

  export let onActivate: (path: string, kind: CatalogNode['kind']) => void = () => {}
  export let onToggleBatch: (path: string) => void = () => {}
  export let onCollapseChange: (key: string, collapsed: boolean) => void = () => {}
  export let onFileContextMenu: (event: MouseEvent, path: string) => void = () => {}
</script>

<div class="catalog-tree" role="tree">
  <CatalogTreeNode
    node={tree}
    depth={0}
    {activeFeature}
    {batchSelected}
    {batchMode}
    {expandAll}
    {collapsed}
    {onActivate}
    {onToggleBatch}
    {onCollapseChange}
    {onFileContextMenu}
  />
</div>

<style>
  .catalog-tree {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;
    padding: 2px 0;
    background: var(--color-sidebar);
  }
</style>
