<script lang="ts">
  export let featurePath = ''
  export let newName = ''
  export let onConfirm: (newName: string) => void = () => {}
  export let onCancel: () => void = () => {}

  function basename(path: string): string {
    const parts = path.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || path
  }

  function stem(path: string): string {
    const name = basename(path)
    return name.replace(/\.feature$/i, '')
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
  <div class="modal duplicate-dialog" role="dialog" aria-modal="true" aria-label="Дублировать сценарий" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>Дублировать сценарий</h3>
    <p class="source">Источник: {basename(featurePath)}</p>
    <label>Имя копии (без .feature)
      <input bind:value={newName} placeholder="scenario-copy" />
    </label>
    <p class="hint">Оставьте пустым для автоматического имени <code>*-copy.feature</code>.</p>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={confirm}>Дублировать</button>
      <button type="button" on:click={onCancel}>Отмена</button>
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

  label {
    display: grid;
    gap: 4px;
    margin-bottom: 8px;
    font-size: 11px;
    color: var(--color-muted);
  }

  input {
    padding: 6px 8px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: var(--color-text);
    font-size: 12px;
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
