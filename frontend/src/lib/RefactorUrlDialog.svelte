<script lang="ts">
  export let initialUrl = 'https://example.com'
  export let onConfirm: (url: string) => void = () => {}
  export let onClose: () => void = () => {}

  let url = initialUrl

  function submit() {
    const value = url.trim()
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
  <div class="modal" role="dialog" aria-modal="true" aria-label="Обновить стартовый URL" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>Обновить стартовый URL</h3>
    <p class="hint">Заменяет URL во всех шагах «открыт» в текущем файле.</p>
    <label>
      Новый URL
      <input bind:value={url} placeholder="https://example.com" />
    </label>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={submit} disabled={!url.trim()}>Применить</button>
      <button type="button" on:click={onClose}>Отмена</button>
    </div>
  </div>
</div>

<style>
  .hint {
    font-size: 12px;
    color: var(--color-muted);
    margin: 0 0 8px;
  }
</style>
