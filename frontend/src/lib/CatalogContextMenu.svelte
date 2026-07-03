<script lang="ts">
  import { createTranslator, locale } from './i18n'

  export let x = 0
  export let y = 0
  export let onRun: () => void = () => {}
  export let onDryRun: () => void = () => {}
  export let onOpen: () => void = () => {}
  export let onDuplicate: () => void = () => {}
  export let onRename: () => void = () => {}
  export let onMove: () => void = () => {}
  export let onDelete: () => void = () => {}
  export let onReveal: () => void = () => {}
  export let onClose: () => void = () => {}

  $: tr = createTranslator($locale)

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="context-menu-backdrop" role="presentation" on:click={onClose} on:contextmenu|preventDefault={onClose}>
  <div class="context-menu" style="left: {x}px; top: {y}px" role="menu" tabindex="0" on:click|stopPropagation on:keydown|stopPropagation>
    <button type="button" on:click={onRun}>{tr('dialogs.contextMenu.catalog.run')}</button>
    <button type="button" on:click={onDryRun}>{tr('dialogs.contextMenu.catalog.dryRun')}</button>
    <button type="button" on:click={onOpen}>{tr('dialogs.contextMenu.catalog.open')}</button>
    <button type="button" on:click={onDuplicate}>{tr('dialogs.contextMenu.catalog.duplicate')}</button>
    <button type="button" on:click={onRename}>{tr('dialogs.contextMenu.catalog.rename')}</button>
    <button type="button" on:click={onMove}>{tr('dialogs.contextMenu.catalog.move')}</button>
    <button type="button" on:click={onReveal}>{tr('dialogs.contextMenu.catalog.reveal')}</button>
    <button type="button" class="danger" on:click={onDelete}>{tr('dialogs.contextMenu.catalog.delete')}</button>
  </div>
</div>
