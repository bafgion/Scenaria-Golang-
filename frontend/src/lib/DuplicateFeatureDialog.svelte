<script lang="ts">
  import { createTranslator, locale } from './i18n'

  export let featurePath = ''
  export let newName = ''
  export let onConfirm: (newName: string) => void = () => {}
  export let onCancel: () => void = () => {}

  $: tr = createTranslator($locale)

  function basename(path: string): string {
    const parts = path.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || path
  }

  function confirm() {
    onConfirm(newName.trim())
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onCancel}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal duplicate-dialog" role="dialog" aria-modal="true" aria-label={tr('dialogs.feature.duplicate.ariaLabel')} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{tr('dialogs.feature.duplicate.title')}</h3>
    <p class="source">{tr('dialogs.feature.duplicate.source', { name: basename(featurePath) })}</p>
    <label>{tr('dialogs.feature.duplicate.copyName')}
      <input bind:value={newName} placeholder={tr('dialogs.feature.duplicate.copyNamePlaceholder')} />
    </label>
    <p class="hint">{tr('dialogs.feature.duplicate.hint')}</p>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={confirm}>{tr('dialogs.feature.duplicate.confirm')}</button>
      <button type="button" on:click={onCancel}>{tr('dialogs.common.cancel')}</button>
    </div>
  </div>
</div>

<style>
  h3 {
    margin: 0 0 8px;
    font-size: 14px;
  }

  .source {
    margin: 0 0 12px;
    font-size: 12px;
    color: var(--color-muted);
  }

  .hint {
    margin: 0 0 12px;
    font-size: 11px;
    color: var(--color-muted);
  }
</style>
