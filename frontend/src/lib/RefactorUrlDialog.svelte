<script lang="ts">
  import { createTranslator, locale } from './i18n'

  export let initialUrl = 'https://example.com'
  export let onConfirm: (url: string) => void = () => {}
  export let onClose: () => void = () => {}

  $: tr = createTranslator($locale)

  let url = initialUrl

  function submit() {
    const value = url.trim()
    if (!value) return
    onConfirm(value)
    onClose()
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
    if (e.key === 'Enter') submit()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal" role="dialog" aria-modal="true" aria-label={tr('dialogs.feature.refactorUrl.ariaLabel')} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{tr('dialogs.feature.refactorUrl.title')}</h3>
    <p class="hint">{tr('dialogs.feature.refactorUrl.hint')}</p>
    <label>
      {tr('dialogs.feature.refactorUrl.newUrl')}
      <input bind:value={url} placeholder={tr('dialogs.feature.refactorUrl.newUrlPlaceholder')} />
    </label>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={submit} disabled={!url.trim()}>{tr('dialogs.common.apply')}</button>
      <button type="button" on:click={onClose}>{tr('dialogs.common.cancel')}</button>
    </div>
  </div>
</div>

<style>
  .hint {
    font-size: 12px;
    color: var(--color-muted);
    margin: 0 0 8px;
  }
</style>
