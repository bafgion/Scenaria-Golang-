<script lang="ts">
  import { onMount } from 'svelte'
  import { SearchSteps } from '../../wailsjs/go/wailsapp/App'
  import { asStepSearchQuery } from './stepSearch'

  export type StepEntry = { category: string; template: string; help: string }

  export let onInsert: (template: string) => void = () => {}
  export let onClose: () => void = () => {}

  let query = ''
  let entries: StepEntry[] = []
  let loading = true

  onMount(async () => {
    try {
      entries = await SearchSteps('')
    } catch {
      entries = []
    } finally {
      loading = false
    }
  })

  async function onQueryInput() {
    loading = true
    try {
      entries = await SearchSteps(asStepSearchQuery(query))
    } catch {
      entries = []
    } finally {
      loading = false
    }
  }

  function pick(template: string) {
    onInsert(template)
    onClose()
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal wide tall" role="dialog" aria-modal="true" aria-label="Вставить шаг" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>Вставить шаг</h3>
    <input bind:value={query} placeholder="Поиск шага…" on:input={onQueryInput} />
    <div class="step-list">
      {#if loading}
        <p class="empty">Загрузка…</p>
      {:else if entries.length === 0}
        <p class="empty">Шаги не найдены</p>
      {:else}
        {#each entries as step}
          <button type="button" class="step-row" on:click={() => pick(step.template)}>
            <div class="template">{step.template}</div>
            <div class="meta">{step.category} — {step.help}</div>
          </button>
        {/each}
      {/if}
    </div>
    <div class="modal-actions">
      <button type="button" on:click={onClose}>Закрыть</button>
    </div>
  </div>
</div>

<style>
  .step-list {
    max-height: 360px;
    overflow: auto;
    margin: 8px 0;
    border: 1px solid var(--color-border);
    border-radius: 3px;
  }

  .step-row {
    display: block;
    width: 100%;
    text-align: left;
    padding: 8px 10px;
    border: none;
    border-bottom: 1px solid var(--color-divider);
    background: transparent;
    color: var(--color-text);
    cursor: pointer;
  }

  .step-row:hover {
    background: var(--color-selected);
  }

  .step-row .template {
    font-family: var(--font-mono, monospace);
    font-size: 12px;
  }

  .step-row .meta {
    font-size: 11px;
    color: var(--color-muted);
    margin-top: 2px;
  }

  .empty {
    padding: 12px;
    font-size: 12px;
    color: var(--color-muted);
    margin: 0;
  }
</style>
