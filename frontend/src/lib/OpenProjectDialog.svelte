<script lang="ts">
  export let initialPath = ''
  export let recentProjects: string[] = []
  export let onConfirm: (path: string) => void = () => {}
  export let onClose: () => void = () => {}

  let path = initialPath

  function submit() {
    const value = path.trim()
    if (!value) return
    onConfirm(value)
    onClose()
  }

  function pickRecent(p: string) {
    path = p
    submit()
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
    if (e.key === 'Enter') submit()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="modal-backdrop" role="presentation" on:click={onClose}>
  <div class="modal" role="dialog" aria-label="Открыть проект" on:click|stopPropagation>
    <h3>Открыть проект</h3>
    <p class="hint">Укажите папку с .feature сценариями (если диалог выбора папки недоступен).</p>
    <label>
      Путь к папке
      <input bind:value={path} placeholder="C:\Projects\my-scenarios" />
    </label>
    {#if recentProjects.length > 0}
      <div class="recents">
        <p class="recents-title">Недавние проекты</p>
        {#each recentProjects.slice(0, 8) as project}
          <button type="button" class="recent-item" on:click={() => pickRecent(project)}>{project}</button>
        {/each}
      </div>
    {/if}
    <div class="modal-actions">
      <button type="button" class="primary" on:click={submit} disabled={!path.trim()}>Открыть</button>
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

  .recents {
    margin-top: 12px;
  }

  .recents-title {
    font-size: 11px;
    color: var(--color-muted);
    margin: 0 0 6px;
  }

  .recent-item {
    display: block;
    width: 100%;
    text-align: left;
    padding: 6px 8px;
    margin-bottom: 4px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: var(--color-text);
    font-size: 12px;
    cursor: pointer;
  }

  .recent-item:hover {
    background: var(--color-selected);
  }
</style>
