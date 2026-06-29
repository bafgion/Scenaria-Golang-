<script lang="ts">
  export let currentVersion = ''
  export let message = ''
  export let hasUpdate = false
  export let onClose: () => void = () => {}

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal update-dialog" role="dialog" aria-modal="true" aria-label="Обновления" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>Обновления Scenaria</h3>
    <p class="version">Текущая версия: {currentVersion || '—'}</p>
    {#if hasUpdate}
      <p class="update-available">Доступна новая версия</p>
    {:else}
      <p class="up-to-date">Установлена актуальная версия</p>
    {/if}
    {#if message}
      <pre class="details">{message.trim()}</pre>
    {/if}
    <div class="modal-actions">
      <button type="button" class="primary" on:click={onClose}>OK</button>
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
