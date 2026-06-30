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

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onCancel}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal move-dialog" role="dialog" aria-modal="true" aria-label="Переместить сценарий" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
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

  .hint {
    margin: 0 0 12px;
    font-size: 11px;
    color: var(--color-muted);
  }
</style>
