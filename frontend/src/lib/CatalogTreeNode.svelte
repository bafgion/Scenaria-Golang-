<script lang="ts">
  import { icons } from './icons'
  import { type CatalogNode, fileTreeLabel, fileTreeTooltip } from './catalogTree'
  import Self from './CatalogTreeNode.svelte'

  const DRAG_TYPE = 'application/x-scenaria-feature'

  export let node: CatalogNode
  export let depth = 0
  export let activeFeature = ''
  export let batchSelected: string[] = []
  export let batchMode = false
  export let expandAll = false
  export let collapsed: Set<string> = new Set()
  export let dropTarget = ''

  export let onActivate: (path: string, kind: CatalogNode['kind']) => void = () => {}
  export let onToggleBatch: (path: string) => void = () => {}
  export let onCollapseChange: (key: string, collapsed: boolean) => void = () => {}
  export let onFileContextMenu: (event: MouseEvent, path: string) => void = () => {}
  export let onMoveFeature: (src: string, destDir: string) => void = () => {}
  export let onDropTarget: (path: string) => void = () => {}

  $: expanded = node.kind === 'file' ? false : expandAll || !collapsed.has(node.path)
  $: isFolder = node.kind !== 'file'
  $: isDropTarget = isFolder && dropTarget === node.path

  function toggleExpand() {
    if (node.kind === 'file') return
    onCollapseChange(node.path, expanded)
  }

  function isSelected(path: string): boolean {
    const norm = path.replace(/\\/g, '/').toLowerCase()
    return batchSelected.some((p) => p.replace(/\\/g, '/').toLowerCase() === norm)
  }

  function onRowClick(e: MouseEvent) {
    if (node.kind === 'file') {
      if (batchMode || e.ctrlKey || e.metaKey) {
        onToggleBatch(node.path)
        return
      }
      onActivate(node.path, node.kind)
      return
    }
    toggleExpand()
  }

  function onRowDblClick() {
    if (node.kind !== 'file') onActivate(node.path, node.kind)
  }

  function onDragStart(e: DragEvent) {
    if (node.kind !== 'file' || !e.dataTransfer) return
    e.dataTransfer.setData(DRAG_TYPE, node.path)
    e.dataTransfer.effectAllowed = 'move'
  }

  function onDragEnd() {
    onDropTarget('')
  }

  function onDragOver(e: DragEvent) {
    if (!isFolder || !e.dataTransfer?.types.includes(DRAG_TYPE)) return
    e.preventDefault()
    e.dataTransfer.dropEffect = 'move'
    onDropTarget(node.path)
  }

  function onDragLeave() {
    if (isDropTarget) onDropTarget('')
  }

  function onDrop(e: DragEvent) {
    if (!isFolder) return
    e.preventDefault()
    onDropTarget('')
    const src = e.dataTransfer?.getData(DRAG_TYPE)
    if (src) onMoveFeature(src, node.path)
  }

  $: rowClass = [
    'catalog-tree-row',
    node.kind === 'file' ? 'file' : 'folder',
    node.kind === 'file' && node.path === activeFeature ? 'active' : '',
    node.kind === 'file' && isSelected(node.path) ? 'batch-selected' : '',
    node.kind === 'file' && node.runSuccess === true ? 'run-ok' : '',
    node.kind === 'file' && node.runSuccess === false ? 'run-fail' : '',
    node.kind === 'file' && node.runSuccess === null ? 'run-none' : '',
    isDropTarget ? 'drop-target' : '',
  ]
    .filter(Boolean)
    .join(' ')

  $: indent = `${depth * 16}px`
  $: selected = isSelected(node.path)
</script>

<div class="catalog-tree-item">
  <button
    type="button"
    class={rowClass}
    style="padding-left: calc({indent} + 4px)"
    title={node.kind === 'file' ? fileTreeTooltip(node) : node.path}
    draggable={node.kind === 'file'}
    on:click={onRowClick}
    on:dblclick={onRowDblClick}
    on:contextmenu={(e) => node.kind === 'file' && onFileContextMenu(e, node.path)}
    on:dragstart={onDragStart}
    on:dragend={onDragEnd}
    on:dragover={onDragOver}
    on:dragleave={onDragLeave}
    on:drop={onDrop}
  >
    {#if node.kind !== 'file'}
      <span class="tree-chevron" aria-hidden="true">
        {#if expanded}
          {@html icons.chevronDown}
        {:else}
          {@html icons.chevronRight}
        {/if}
      </span>
      <span class="tree-folder">📁 {node.name}</span>
    {:else}
      <span class="tree-chevron spacer" aria-hidden="true"></span>
      <span class="tree-file-label">{fileTreeLabel(node, batchMode, selected)}</span>
    {/if}
  </button>

  {#if node.kind !== 'file' && expanded}
    {#each node.children as child (child.path)}
      <Self
        node={child}
        depth={depth + 1}
        {activeFeature}
        {batchSelected}
        {batchMode}
        {expandAll}
        {collapsed}
        {dropTarget}
        {onActivate}
        {onToggleBatch}
        {onCollapseChange}
        {onFileContextMenu}
        {onMoveFeature}
        {onDropTarget}
      />
    {/each}
  {/if}
</div>

<style>
  .catalog-tree-item {
    display: block;
  }

  .catalog-tree-row {
    display: flex;
    align-items: center;
    width: 100%;
    min-height: 20px;
    padding: 0 4px 0 4px;
    border: none;
    background: transparent;
    color: var(--color-text);
    font-size: 12px;
    line-height: 20px;
    text-align: left;
    cursor: default;
    box-sizing: border-box;
  }

  .catalog-tree-row:hover {
    background: var(--color-input);
  }

  .catalog-tree-row.active {
    background: #37373d;
    color: #fff;
  }

  .catalog-tree-row.batch-selected {
    background: rgba(9, 71, 113, 0.35);
  }

  .catalog-tree-row.drop-target {
    outline: 1px dashed var(--color-accent);
    background: rgba(9, 71, 113, 0.25);
  }

  .catalog-tree-row.run-ok .tree-file-label {
    color: var(--color-success);
  }

  .catalog-tree-row.run-fail .tree-file-label {
    color: var(--color-error);
  }

  .catalog-tree-row.run-none .tree-file-label {
    color: var(--color-muted);
  }

  .catalog-tree-row.active .tree-file-label {
    color: inherit;
  }

  .tree-chevron {
    width: 14px;
    height: 14px;
    flex-shrink: 0;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    margin-right: 2px;
    color: var(--color-muted);
  }

  .tree-chevron.spacer {
    visibility: hidden;
  }

  .tree-chevron :global(svg) {
    width: 12px;
    height: 12px;
  }

  .tree-folder,
  .tree-file-label {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    min-width: 0;
    flex: 1;
  }

  .tree-folder {
    color: var(--color-text);
  }
</style>
