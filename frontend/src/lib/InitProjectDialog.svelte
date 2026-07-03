<script lang="ts">
  import { createTranslator, locale } from './i18n'

  export let projectPath = ''
  export let onConfirm: () => void = () => {}
  export let onCancel: () => void = () => {}

  $: tr = createTranslator($locale)

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onCancel}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal init-dialog" role="dialog" aria-modal="true" aria-label={tr('dialogs.project.init.ariaLabel')} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{tr('dialogs.project.init.title')}</h3>
    <p class="hint">
      {tr('dialogs.project.init.hintCreate')}
      {#if projectPath}
        <br /><span class="path">{projectPath}</span>
      {/if}
    </p>
    <p class="hint">{tr('dialogs.project.init.hintFiles')}</p>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={onConfirm}>{tr('dialogs.project.init.confirm')}</button>
      <button type="button" on:click={onCancel}>{tr('dialogs.common.cancel')}</button>
    </div>
  </div>
</div>

<style>
  h3 {
    margin: 0 0 12px;
    font-size: 14px;
  }

  .hint {
    margin: 0 0 10px;
    font-size: 12px;
    color: var(--color-muted);
    line-height: 1.45;
  }

  .path {
    word-break: break-all;
    color: var(--color-text);
  }
</style>
