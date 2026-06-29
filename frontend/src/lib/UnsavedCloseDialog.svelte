<script lang="ts">
  export let fileName = ''
  export let onSave: () => void | Promise<void> = () => {}
  export let onDiscard: () => void = () => {}
  export let onCancel: () => void = () => {}

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      e.preventDefault()
      onCancel()
    }
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onCancel}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal unsaved-close" role="dialog" aria-modal="true" aria-label="Несохранённые изменения" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>Несохранённые изменения</h3>
    <p class="message">Сохранить изменения в «{fileName}» перед закрытием?</p>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={() => onSave()}>Сохранить</button>
      <button type="button" on:click={() => onDiscard()}>Не сохранять</button>
      <button type="button" on:click={() => onCancel()}>Отмена</button>
    </div>
  </div>
</div>

<style>
  .message {
    margin: 0;
    font-size: 13px;
    color: var(--color-text);
  }
</style>
