<script lang="ts">
  export let dryRun = false
  export let tag = ''
  export let excludeTags = ''
  export let scenario = ''
  export let tags: string[] = []
  export let onConfirm: () => void = () => {}
  export let onCancel: () => void = () => {}

  function pickTag(value: string) {
    tag = value
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="modal-backdrop" role="presentation" on:click={onCancel}>
  <div class="modal wide" role="dialog" aria-label="Запуск Vanessa" on:click|stopPropagation>
    <h3>Vanessa {dryRun ? '(dry-run)' : ''}</h3>
    <label>Тег <input bind:value={tag} placeholder="@smoke" /></label>
    {#if tags.length > 0}
      <div class="tag-chips">
        {#each tags as t}
          <button type="button" class="chip" class:active={tag === t} on:click={() => pickTag(t)}>{t}</button>
        {/each}
      </div>
    {/if}
    <label>Исключить теги (через запятую) <input bind:value={excludeTags} placeholder="@wip, @draft" /></label>
    <label>Имя сценария <input bind:value={scenario} placeholder="Необязательно" /></label>
    <p class="hint">Пути к 1C и EPF задаются в <code>.scenaria/vanessa.json</code>.</p>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={onConfirm}>{dryRun ? 'Dry-run' : 'Запустить'}</button>
      <button type="button" on:click={onCancel}>Отмена</button>
    </div>
  </div>
</div>

<style>
  .tag-chips {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin: -4px 0 8px;
  }

  .chip {
    padding: 2px 8px;
    font-size: 11px;
    border: 1px solid var(--color-border);
    border-radius: 10px;
    background: var(--color-input);
    color: var(--color-muted);
  }

  .chip.active {
    border-color: var(--color-primary);
    color: var(--color-text);
    background: var(--color-selected);
  }

  .hint {
    font-size: 11px;
    color: var(--color-muted);
    margin: 0;
  }
</style>
