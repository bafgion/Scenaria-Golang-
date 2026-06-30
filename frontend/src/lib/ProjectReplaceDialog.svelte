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
  .find-replace .actions {
    margin-top: 4px;
  }
</style>
