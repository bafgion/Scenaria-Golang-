<script lang="ts">
  import { createTranslator, locale } from './i18n'

  export let fileName = ''
  export let onSave: () => void | Promise<void> = () => {}
  export let onDiscard: () => void = () => {}
  export let onCancel: () => void = () => {}

  $: tr = createTranslator($locale)

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      e.preventDefault()
      onCancel()
    }
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop modal-layer-top" role="presentation" on:click={onCancel}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal unsaved-close" role="dialog" aria-modal="true" aria-label={tr('dialogs.unsavedClose.ariaLabel')} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{tr('dialogs.unsavedClose.title')}</h3>
    <p class="message">{tr('dialogs.unsavedClose.message', { fileName })}</p>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={() => onSave()}>{tr('dialogs.common.save')}</button>
      <button type="button" on:click={() => onDiscard()}>{tr('dialogs.unsavedClose.discard')}</button>
      <button type="button" on:click={() => onCancel()}>{tr('dialogs.common.cancel')}</button>
    </div>
  </div>
</div>

<style>
  .message {
    margin: 0;
    font-size: 13px;
    color: var(--color-text);
  }
</style>
