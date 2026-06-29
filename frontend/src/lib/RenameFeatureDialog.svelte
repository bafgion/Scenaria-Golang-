<script lang="ts">
  export let currentPath = ''
  export let onConfirm: (newName: string) => void = () => {}
  export let onClose: () => void = () => {}

  function basename(path: string): string {
    const parts = path.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || path
  }

  let name = basename(currentPath).replace(/\.feature$/i, '')

  function submit() {
    const value = name.trim()
    if (!value) return
    onConfirm(value)
    onClose()
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
    if (e.key === 'Enter') submit()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal" role="dialog" aria-modal="true" aria-label="Переименовать сценарий" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>Переименовать сценарий</h3>
    <label>
      Новое имя файла
      <input bind:value={name} placeholder="scenario-name" />
    </label>
    <p class="hint">Расширение .feature добавится автоматически.</p>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={submit} disabled={!name.trim()}>Переименовать</button>
      <button type="button" on:click={onClose}>Отмена</button>
    </div>
  </div>
</div>

<style>
  .hint {
    font-size: 12px;
    color: var(--color-muted);
    margin: 0;
  }
</style>
