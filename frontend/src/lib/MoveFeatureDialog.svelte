<script lang="ts">
  export let featurePath = ''
  export let destDirs: string[] = []
  export let destDir = ''
  export let onConfirm: (destDir: string) => void = () => {}
  export let onCancel: () => void = () => {}

  function basename(path: string): string {
    const parts = path.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || path
  }

  function dirLabel(path: string): string {
    const norm = path.replace(/\\/g, '/')
    const parts = norm.split('/').filter(Boolean)
    return parts.length ? parts[parts.length - 1] : norm || '(корень)'
  }

  function confirm() {
    if (destDir) onConfirm(destDir)
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="modal-backdrop" role="presentation" on:click={onCancel}>
  <div class="modal move-dialog" role="dialog" aria-label="Переместить сценарий" on:click|stopPropagation>
    <h3>Переместить сценарий</h3>
    <p class="file">{basename(featurePath)}</p>
    <label>
      Папка назначения
      <select bind:value={destDir}>
        {#each destDirs as dir}
          <option value={dir}>{dirLabel(dir)}</option>
        {/each}
      </select>
    </label>
    <p class="hint">Файл будет перемещён в выбранную папку внутри проекта.</p>
    <div class="modal-actions">
      <button type="button" class="primary" disabled={!destDir} on:click={confirm}>Переместить</button>
      <button type="button" on:click={onCancel}>Отмена</button>
    </div>
  </div>
</div>

<style>
  h3 {
    margin: 0 0 8px;
    font-size: 14px;
  }

  .file {
    margin: 0 0 12px;
    font-size: 12px;
    font-weight: 600;
  }

  label {
    display: grid;
    gap: 4px;
    margin-bottom: 10px;
    font-size: 11px;
    color: var(--color-muted);
  }

  select {
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

  button:disabled {
    opacity: 0.5;
  }
</style>
