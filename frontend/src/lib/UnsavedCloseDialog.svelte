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

<div class="palette-backdrop" role="presentation" on:click={onCancel}>
  <div class="palette unsaved-close" role="dialog" aria-label="Несохранённые изменения" on:click|stopPropagation>
    <h3>Несохранённые изменения</h3>
    <p class="hint">Сохранить изменения в «{fileName}» перед закрытием?</p>
    <div class="actions">
      <button type="button" class="primary" on:click={() => onSave()}>Сохранить</button>
      <button type="button" on:click={() => onDiscard()}>Не сохранять</button>
      <button type="button" on:click={() => onCancel()}>Отмена</button>
    </div>
  </div>
</div>

<style>
  .unsaved-close {
    width: min(420px, 92vw);
  }

  .hint {
    margin: 0 0 12px;
    font-size: 12px;
    color: var(--color-muted);
  }

  .actions {
    display: flex;
    gap: 8px;
    justify-content: flex-end;
    flex-wrap: wrap;
  }
</style>
