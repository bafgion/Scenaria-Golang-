<script lang="ts">
  export let pluginName = ''
  export let pluginTitle = ''
  export let dryRun = false
  export let tag = ''
  export let scenario = ''
  export let onConfirm: (payload: { tag: string; scenario: string; dryRun: boolean }) => void = () => {}
  export let onCancel: () => void = () => {}

  function confirm() {
    onConfirm({ tag: tag.trim(), scenario: scenario.trim(), dryRun })
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="modal-backdrop" role="presentation" on:click={onCancel}>
  <div class="modal plugin-run-dialog" role="dialog" aria-label="Запуск плагина" on:click|stopPropagation>
    <h3>Запуск: {pluginTitle || pluginName}</h3>
    <label>Тег (опционально)
      <input bind:value={tag} placeholder="@smoke" />
    </label>
    <label>Сценарий (опционально)
      <input bind:value={scenario} placeholder="Название сценария" />
    </label>
    <label class="check-row">
      <input type="checkbox" bind:checked={dryRun} />
      Dry-run (без выполнения)
    </label>
    <p class="hint">Для Vanessa используйте пункты меню «Vanessa run…» / «Vanessa (dry)…».</p>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={confirm}>Запустить</button>
      <button type="button" on:click={onCancel}>Отмена</button>
    </div>
  </div>
</div>

<style>
  h3 {
    margin: 0 0 12px;
    font-size: 14px;
  }

  label {
    display: grid;
    gap: 4px;
    margin-bottom: 10px;
    font-size: 11px;
    color: var(--color-muted);
  }

  input[type='text'],
  input:not([type]) {
    padding: 6px 8px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: var(--color-text);
    font-size: 12px;
  }

  .check-row {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 12px;
    color: var(--color-text);
  }

  .hint {
    margin: 0 0 12px;
    font-size: 11px;
    color: var(--color-muted);
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
