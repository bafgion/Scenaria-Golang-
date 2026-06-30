<script lang="ts">
  import type { gui } from '../../wailsjs/go/models'
  import { BRAND_NAME } from './brand'

  export let currentVersion = ''
  export let info: gui.UpdateInfoDTO | null = null
  export let message = ''
  export let hasUpdate = false
  export let downloading = false
  export let canAutoApply = false
  export let onClose: () => void = () => {}
  export let onOpenRelease: () => void = () => {}
  export let onDownload: () => void = () => {}
  export let onApply: () => void = () => {}

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal update-dialog" role="dialog" aria-modal="true" aria-label="Обновления" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>Обновления {BRAND_NAME}</h3>
    <p class="version">Текущая версия: {info?.currentVersion || currentVersion || '—'}</p>
    {#if hasUpdate}
      <p class="update-available">Доступна версия {info?.latestVersion || 'новее'}</p>
    {:else}
      <p class="up-to-date">Установлена актуальная версия</p>
    {/if}
    {#if message}
      <pre class="details">{message.trim()}</pre>
    {/if}
    <div class="modal-actions">
      {#if hasUpdate && info?.htmlUrl}
        <button type="button" on:click={onOpenRelease}>Страница релиза</button>
      {/if}
      {#if hasUpdate && canAutoApply && info?.downloadUrl}
        <button type="button" class="primary" disabled={downloading} on:click={onApply}>
          {downloading ? 'Обновление…' : 'Установить обновление'}
        </button>
      {/if}
      {#if hasUpdate && info?.downloadUrl}
        <button type="button" disabled={downloading} on:click={onDownload}>
          {downloading ? 'Скачивание…' : 'Скачать вручную'}
        </button>
      {/if}
      <button type="button" on:click={onClose}>Закрыть</button>
    </div>
  </div>
</div>

<style>
  h3 {
    margin: 0 0 8px;
    font-size: 14px;
  }

  .version {
    margin: 0 0 8px;
    font-size: 12px;
    color: var(--color-muted);
  }

  .update-available {
    margin: 0 0 8px;
    font-size: 13px;
    color: var(--color-accent);
  }

  .up-to-date {
    margin: 0 0 8px;
    font-size: 13px;
    color: var(--color-success, #4ec9b0);
  }

  .details {
    margin: 0 0 12px;
    padding: 8px;
    font-size: 11px;
    line-height: 1.4;
    white-space: pre-wrap;
    background: var(--color-input);
    border: 1px solid var(--color-border);
    border-radius: 3px;
    max-height: 160px;
    overflow: auto;
  }

  .modal-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    justify-content: flex-end;
  }

  button {
    padding: 6px 12px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: var(--color-text);
    font-size: 12px;
  }

  button.primary {
    background: var(--color-accent);
    color: var(--color-on-accent, #fff);
    border-color: var(--color-accent);
  }
</style>
