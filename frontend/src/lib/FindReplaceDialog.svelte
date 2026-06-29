<script lang="ts">
  export let findText = ''
  export let replaceText = ''
  export let caseSensitive = false
  export let onFindNext: () => boolean = () => false
  export let onReplace: () => boolean = () => false
  export let onReplaceAll: () => number = () => 0
  export let onClose: () => void = () => {}

  let status = ''

  function findNext() {
    status = onFindNext() ? '' : 'Не найдено'
  }

  function replaceOne() {
    const ok = onReplace()
    status = ok ? 'Заменено' : 'Не найдено'
  }

  function replaceAll() {
    const n = onReplaceAll()
    status = n > 0 ? `Заменено: ${n}` : 'Не найдено'
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      e.preventDefault()
      onClose()
    } else if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) {
      e.preventDefault()
      replaceAll()
    } else if (e.key === 'Enter') {
      e.preventDefault()
      findNext()
    }
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="palette-backdrop" role="presentation" on:click={onClose}>
  <div class="palette find-replace" role="dialog" aria-label="Найти и заменить" on:click|stopPropagation>
    <h3>Найти и заменить</h3>
    <label>Найти <input bind:value={findText} autofocus /></label>
    <label>Заменить <input bind:value={replaceText} /></label>
    <label class="check-row">
      <input type="checkbox" bind:checked={caseSensitive} /> Учитывать регистр
    </label>
    {#if status}<p class="status">{status}</p>{/if}
    <div class="actions">
      <button type="button" on:click={findNext}>Найти далее</button>
      <button type="button" on:click={replaceOne}>Заменить</button>
      <button type="button" class="primary" on:click={replaceAll}>Заменить все</button>
      <button type="button" on:click={onClose}>Закрыть</button>
    </div>
  </div>
</div>

<style>
  .find-replace {
    width: min(420px, 92vw);
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

  .status {
    margin: 0;
    font-size: 12px;
    color: var(--color-muted);
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
</style>
