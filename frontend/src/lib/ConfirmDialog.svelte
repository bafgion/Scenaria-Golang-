<script lang="ts">
  import { createTranslator, locale } from './i18n'

  export let title = ''
  export let message = ''
  export let confirmLabel = ''
  export let danger = false
  export let dontAskAgainLabel = ''
  export let onConfirm: (dontAskAgain?: boolean) => void = () => {}
  export let onCancel: () => void = () => {}

  $: tr = createTranslator($locale)
  $: dialogTitle = title || tr('common.confirm')
  $: dialogConfirmLabel = confirmLabel || tr('common.ok')

  let dontAskAgain = false

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel()
    if (e.key === 'Enter') onConfirm()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop modal-layer-confirm" role="presentation" on:click={onCancel}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal confirm-dialog" role="alertdialog" aria-modal="true" aria-label={dialogTitle} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{dialogTitle}</h3>
    <p class="message">{message}</p>
    {#if dontAskAgainLabel}
      <label class="dont-ask">
        <input type="checkbox" bind:checked={dontAskAgain} />
        {dontAskAgainLabel}
      </label>
    {/if}
    <div class="modal-actions">
      <button type="button" class="primary" class:danger on:click={() => onConfirm(dontAskAgain)}>{dialogConfirmLabel}</button>
      <button type="button" on:click={onCancel}>{tr('common.cancel')}</button>
    </div>
  </div>
</div>

<style>
  .confirm-dialog {
    width: min(420px, 92vw);
  }

  .message {
    margin: 0 0 12px;
    font-size: 13px;
    color: var(--color-text);
  }

  .dont-ask {
    display: flex;
    align-items: center;
    gap: 8px;
    margin: 0 0 12px;
    font-size: 12px;
    color: var(--color-muted);
    cursor: pointer;
  }
</style>
