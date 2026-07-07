<script lang="ts">
  import { createTranslator, locale } from './i18n'

  export let x = 0
  export let y = 0
  export let featureCount = 0
  export let onRun: () => void = () => {}
  export let onDryRun: () => void = () => {}
  export let onVanessa: () => void = () => {}
  export let onVanessaDry: () => void = () => {}
  export let onSelectBatch: () => void = () => {}
  export let onRefresh: () => void = () => {}
  export let onClose: () => void = () => {}

  $: tr = createTranslator($locale)

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="context-menu-backdrop" role="presentation" on:click={onClose} on:contextmenu|preventDefault={onClose}>
  <div class="context-menu" style="left: {x}px; top: {y}px" role="menu" tabindex="0" on:click|stopPropagation on:keydown|stopPropagation>
    <button type="button" on:click={onRefresh}>{tr('dialogs.contextMenu.folder.refresh')}</button>
    <button type="button" on:click={onRun} disabled={featureCount === 0}>{tr('dialogs.contextMenu.folder.runAll', { count: featureCount })}</button>
    <button type="button" on:click={onDryRun} disabled={featureCount === 0}>{tr('dialogs.contextMenu.folder.dryRun')}</button>
    <button type="button" on:click={onVanessa} disabled={featureCount === 0}>{tr('dialogs.contextMenu.folder.vanessa')}</button>
    <button type="button" on:click={onVanessaDry} disabled={featureCount === 0}>{tr('dialogs.contextMenu.folder.vanessaDry')}</button>
    <button type="button" on:click={onSelectBatch} disabled={featureCount === 0}>{tr('dialogs.contextMenu.folder.selectBatch')}</button>
  </div>
</div>
