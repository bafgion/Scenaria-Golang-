<script lang="ts">
  export let pluginName = ''
  export let pluginTitle = ''
  export let dryRun = false
  export let tag = ''
  export let scenario = ''
  export let scenarios: string[] = []
  export let tags: string[] = []
  export let onConfirm: (payload: { tag: string; scenario: string; dryRun: boolean }) => void = () => {}
  export let onCancel: () => void = () => {}

  function pickTag(value: string) {
    tag = value
  }

  function confirm() {
    onConfirm({ tag: tag.trim(), scenario: scenario.trim(), dryRun })
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onCancel}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal plugin-run-dialog" role="dialog" aria-modal="true" aria-label="Запуск плагина" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>Запуск: {pluginTitle || pluginName}</h3>
    <label>Тег (опционально)
      <input bind:value={tag} placeholder="@smoke" />
    </label>
    {#if tags.length > 0}
      <div class="tag-chips">
        {#each tags as t}
          <button type="button" class="chip" class:active={tag === t} on:click={() => pickTag(t)}>{t}</button>
        {/each}
      </div>
    {/if}
    <label>Сценарий (опционально)
      <input bind:value={scenario} placeholder="Название сценария" list="plugin-scenario-list" />
    </label>
    {#if scenarios.length > 0}
      <datalist id="plugin-scenario-list">
        {#each scenarios as name}
          <option value={name}></option>
        {/each}
      </datalist>
      <div class="tag-chips">
        {#each scenarios as name}
          <button type="button" class="chip" class:active={scenario === name} on:click={() => (scenario = name)}>{name}</button>
        {/each}
      </div>
    {/if}
    <label class="check-row">
      <input type="checkbox" bind:checked={dryRun} />
      Dry-run (без выполнения)
    </label>
    <p class="hint">Для Vanessa используйте пункты меню «Vanessa run…» / «Vanessa (dry)…».</p>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={confirm}>Запустить</button>
      <button type="button" on:click={onCancel}>Отмена</button>
    </div>
  </div>
</div>

<style>
  h3 {
    margin: 0 0 12px;
    font-size: 14px;
  }

  label {
    display: grid;
    gap: 4px;
    margin-bottom: 10px;
    font-size: 11px;
    color: var(--color-muted);
  }

  .check-row {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 12px;
    color: var(--color-text);
  }

  .hint {
    margin: 0 0 12px;
    font-size: 11px;
    color: var(--color-muted);
  }

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
</style>
