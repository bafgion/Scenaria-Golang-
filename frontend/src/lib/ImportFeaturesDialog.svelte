<script lang="ts">
  import { PickOpenFile } from '../../wailsjs/go/wailsapp/App'

  export let projectPath = ''
  export let destDirs: string[] = []
  export let destDir = ''
  export let busy = false
  export let onImport: (payload: { destDir: string; paths: string[] }) => void = () => {}
  export let onClose: () => void = () => {}

  let paths: string[] = []
  let error = ''

  function basename(path: string): string {
    const parts = path.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || path
  }

  function dirLabel(path: string): string {
    const norm = path.replace(/\\/g, '/')
    const parts = norm.split('/').filter(Boolean)
    return parts.length ? parts[parts.length - 1] : norm || '(корень)'
  }

  async function addFile() {
    error = ''
    const picked = await PickOpenFile('Импорт .feature')
    if (!picked) return
    if (!picked.toLowerCase().endsWith('.feature')) {
      error = 'Выберите файл .feature'
      return
    }
    if (!paths.includes(picked)) paths = [...paths, picked]
  }

  function removePath(index: number) {
    paths = paths.filter((_, i) => i !== index)
  }

  function confirm() {
    if (!destDir || paths.length === 0) {
      error = 'Укажите папку и хотя бы один файл'
      return
    }
    onImport({ destDir, paths })
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="modal-backdrop" role="presentation" on:click={onClose}>
  <div class="modal wide import-features-dialog" role="dialog" aria-label="Импорт feature" on:click|stopPropagation>
    <h3>Импорт .feature в проект</h3>
    <p class="hint">Копирует внешние сценарии в папку проекта (аналог drag-and-drop на каталог).</p>
    <label>
      Папка назначения
      <select bind:value={destDir} disabled={busy}>
        {#each destDirs as dir}
          <option value={dir}>{dirLabel(dir)}</option>
        {/each}
      </select>
    </label>
    <div class="files-header">
      <span>Файлы ({paths.length})</span>
      <button type="button" disabled={busy} on:click={addFile}>Добавить файл…</button>
    </div>
    <ul class="file-list">
      {#if paths.length === 0}
        <li class="empty">Нет выбранных файлов</li>
      {:else}
        {#each paths as path, index}
          <li>
            <span title={path}>{basename(path)}</span>
            <button type="button" disabled={busy} on:click={() => removePath(index)}>×</button>
          </li>
        {/each}
      {/if}
    </ul>
    {#if error}<p class="error">{error}</p>{/if}
    <div class="modal-actions">
      <button type="button" class="primary" disabled={busy || !paths.length || !destDir} on:click={confirm}>Импортировать</button>
      <button type="button" disabled={busy} on:click={onClose}>Отмена</button>
    </div>
  </div>
</div>

<style>
  h3 {
    margin: 0 0 8px;
    font-size: 14px;
  }

  .hint {
    margin: 0 0 12px;
    font-size: 11px;
    color: var(--color-muted);
    line-height: 1.4;
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

  .files-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin: 12px 0 6px;
    font-size: 12px;
  }

  .file-list {
    margin: 0 0 10px;
    padding: 0;
    list-style: none;
    max-height: 180px;
    overflow: auto;
    border: 1px solid var(--color-border);
    border-radius: 3px;
  }

  li {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    padding: 6px 8px;
    font-size: 12px;
    border-bottom: 1px solid var(--color-divider);
  }

  li:last-child {
    border-bottom: none;
  }

  li.empty {
    color: var(--color-muted);
    justify-content: center;
  }

  li span {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .error {
    margin: 0 0 10px;
    font-size: 12px;
    color: var(--color-error);
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
