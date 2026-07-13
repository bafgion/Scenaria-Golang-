<script lang="ts">
  import { createTranslator, locale } from './i18n'
  import { PickOpenFile, PickOpenFiles } from '../../wailsjs/go/wailsapp/App'

  export let destDirs: string[] = []
  export let destDir = ''
  export let busy = false
  export let onImport: (payload: { destDir: string; paths: string[] }) => void = () => {}
  export let onClose: () => void = () => {}

  $: tr = createTranslator($locale)

  let paths: string[] = []
  let error = ''

  function basename(path: string): string {
    const parts = path.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || path
  }

  function dirLabel(path: string): string {
    const norm = path.replace(/\\/g, '/')
    const parts = norm.split('/').filter(Boolean)
    return parts.length ? parts[parts.length - 1] : tr('dialogs.feature.move.root')
  }

  async function addFiles() {
    error = ''
    const picked = await PickOpenFiles(tr('dialogs.import.features.pickFilesTitle'))
    if (!picked?.length) return
    const valid = picked.filter((p) => p.toLowerCase().endsWith('.feature'))
    if (valid.length === 0) {
      error = tr('dialogs.import.features.errorPickFeatures')
      return
    }
    const merged = [...paths]
    for (const p of valid) {
      if (!merged.includes(p)) merged.push(p)
    }
    paths = merged
  }

  async function addFile() {
    error = ''
    const picked = await PickOpenFile(tr('dialogs.import.features.pickFileTitle'))
    if (!picked) return
    if (!picked.toLowerCase().endsWith('.feature')) {
      error = tr('dialogs.import.features.errorPickFeature')
      return
    }
    if (!paths.includes(picked)) paths = [...paths, picked]
  }

  function removePath(index: number) {
    paths = paths.filter((_, i) => i !== index)
  }

  function confirm() {
    if (!destDir || paths.length === 0) {
      error = tr('dialogs.import.features.errorDestAndFiles')
      return
    }
    onImport({ destDir, paths })
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal wide import-features-dialog" role="dialog" aria-modal="true" aria-label={tr('dialogs.import.features.ariaLabel')} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{tr('dialogs.import.features.title')}</h3>
    <p class="hint">{tr('dialogs.import.features.hint')}</p>
    <label>
      {tr('dialogs.import.features.destDir')}
      <select bind:value={destDir} disabled={busy}>
        {#each destDirs as dir}
          <option value={dir}>{dirLabel(dir)}</option>
        {/each}
      </select>
    </label>
    <div class="files-header">
      <span>{tr('dialogs.import.features.files', { count: paths.length })}</span>
      <div class="file-actions">
        <button type="button" disabled={busy} on:click={addFiles}>{tr('dialogs.import.features.addFiles')}</button>
        <button type="button" disabled={busy} on:click={addFile}>{tr('dialogs.import.features.addFile')}</button>
      </div>
    </div>
    <ul class="file-list">
      {#if paths.length === 0}
        <li class="empty">{tr('dialogs.import.features.noFiles')}</li>
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
      <button type="button" class="primary" disabled={busy || !paths.length || !destDir} on:click={confirm}>{tr('dialogs.import.features.confirm')}</button>
      <button type="button" disabled={busy} on:click={onClose}>{tr('dialogs.common.cancel')}</button>
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

  .files-header {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    align-items: center;
    justify-content: space-between;
    margin: 12px 0 6px;
    font-size: 12px;
  }

  .file-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }

  .file-list {
    margin: 0 0 10px;
    padding: 0;
    list-style: none;
    max-height: 180px;
    overflow: auto;
    border: 1px solid var(--color-border);
    border-radius: 4px;
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
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .error {
    margin: 0 0 10px;
    font-size: 12px;
    color: var(--color-error);
  }

  @media (max-width: 560px) {
    .files-header {
      align-items: stretch;
      flex-direction: column;
    }

    .file-actions button {
      flex: 1 1 140px;
    }
  }
</style>
