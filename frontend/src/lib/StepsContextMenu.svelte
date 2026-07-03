<script lang="ts">
  import { createTranslator, locale } from './i18n'

  export let x = 0
  export let y = 0
  export let onRunFrom: () => void = () => {}
  export let onRunTo: () => void = () => {}
  export let onDryRunFrom: () => void = () => {}
  export let onDryRunTo: () => void = () => {}
  export let onGoto: () => void = () => {}
  export let onHelp: () => void = () => {}
  export let onClose: () => void = () => {}

  $: tr = createTranslator($locale)

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="context-menu-backdrop" role="presentation" on:click={onClose} on:contextmenu|preventDefault={onClose}>
  <div class="context-menu" style="left: {x}px; top: {y}px" role="menu" tabindex="0" on:click|stopPropagation on:keydown|stopPropagation>
    <button type="button" on:click={onRunFrom}>{tr('dialogs.steps.contextMenu.runFrom')}</button>
    <button type="button" on:click={onRunTo}>{tr('dialogs.steps.contextMenu.runTo')}</button>
    <button type="button" on:click={onDryRunFrom}>{tr('dialogs.steps.contextMenu.dryRunFrom')}</button>
    <button type="button" on:click={onDryRunTo}>{tr('dialogs.steps.contextMenu.dryRunTo')}</button>
    <button type="button" on:click={onGoto}>{tr('dialogs.steps.contextMenu.goto')}</button>
    <button type="button" on:click={onHelp}>{tr('dialogs.steps.contextMenu.help')}</button>
  </div>
</div>
