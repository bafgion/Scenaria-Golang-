<script lang="ts">
  export let findText = ''
  export let replaceText = ''
  export let caseSensitive = false
  export let busy = false
  export let onConfirm: () => void | Promise<void> = () => {}
  export let onClose: () => void = () => {}

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      e.preventDefault()
      onClose()
    }
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="palette-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="palette find-replace" role="dialog" aria-modal="true" aria-label="Замена по проекту" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>Замена по проекту</h3>
    <p class="hint">Заменяет текст во всех .feature файлах открытого проекта.</p>
    <label>Найти <input bind:value={findText} autofocus /></label>
    <label>Заменить <input bind:value={replaceText} /></label>
    <label class="check-row">
      <input type="checkbox" bind:checked={caseSensitive} /> Учитывать регистр
    </label>
    <div class="actions">
      <button type="button" class="primary" disabled={busy || !findText} on:click={() => onConfirm()}>
        {busy ? 'Замена…' : 'Заменить во всех файлах'}
      </button>
      <button type="button" on:click={onClose}>Отмена</button>
    </div>
  </div>
</div>

<style>
  .find-replace {
    width: min(440px, 92vw);
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 10px;
    background: var(--color-sidebar);
    border: 1px solid var(--color-border);
    border-radius: 4px;
  }

  h3 {
    margin: 0;
    font-size: 14px;
    font-weight: 600;
  }

  .hint {
    margin: 0;
    font-size: 12px;
    color: var(--color-muted);
    line-height: 1.4;
  }

  label {
    display: flex;
    flex-direction: column;
    gap: 4px;
    font-size: 12px;
    color: var(--color-text);
  }

  input[type='text'],
  input:not([type]) {
    padding: 6px 8px;
    background: var(--color-input);
    border: 1px solid var(--color-border);
    color: var(--color-text);
    border-radius: 3px;
    font-family: inherit;
    font-size: 13px;
  }

  .check-row {
    flex-direction: row;
    align-items: center;
    gap: 8px;
  }

  .actions {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    justify-content: flex-end;
  }

  .actions button {
    padding: 5px 12px;
    font-size: 12px;
    background: var(--color-input);
    border: 1px solid var(--color-border);
    border-radius: 3px;
    color: var(--color-text);
  }

  .actions button.primary {
    background: var(--color-primary);
    border-color: var(--color-primary);
    color: #fff;
  }

  .actions button:disabled {
    opacity: 0.5;
    cursor: default;
  }
</style>
