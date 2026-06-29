<script lang="ts">
  export let projectPath = ''
  export let onConfirm: () => void = () => {}
  export let onCancel: () => void = () => {}

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="modal-backdrop" role="presentation" on:click={onCancel}>
  <div class="modal init-dialog" role="dialog" aria-label="Init проекта" on:click|stopPropagation>
    <h3>Init проекта</h3>
    <p class="hint">
      Создать или обновить каталог <code>.scenaria/</code> в проекте
      {#if projectPath}
        <br /><span class="path">{projectPath}</span>
      {/if}
    </p>
    <p class="hint">Будут добавлены <code>project.json</code> и пример TestClient, если их ещё нет.</p>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={onConfirm}>Init</button>
      <button type="button" on:click={onCancel}>Отмена</button>
    </div>
  </div>
</div>

<style>
  h3 {
    margin: 0 0 12px;
    font-size: 14px;
  }

  .hint {
    margin: 0 0 10px;
    font-size: 12px;
    color: var(--color-muted);
    line-height: 1.45;
  }

  .path {
    word-break: break-all;
    color: var(--color-text);
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
