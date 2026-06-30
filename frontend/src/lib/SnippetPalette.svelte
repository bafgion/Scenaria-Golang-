<script lang="ts">
  import { onMount } from 'svelte'
  import { SearchSteps } from '../../wailsjs/go/wailsapp/App'
  import { asStepSearchQuery } from './stepSearch'
  import type { SnippetEntry } from './stepTypes'

  export let onClose: () => void
  export let onInsert: (template: string) => void

  let query = ''
  let entries: SnippetEntry[] = []
  let selected = 0
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

  $: filtered = entries.filter((e) => {
    const q = query.trim().toLowerCase()
    if (!q) return true
    return (e.template + e.category + e.help).toLowerCase().includes(q)
  })

  $: if (selected >= filtered.length) selected = Math.max(0, filtered.length - 1)

  async function onQueryInput() {
    try {
      entries = await SearchSteps(asStepSearchQuery(query))
    } catch {
      /* keep previous */
    }
  }

  function runSelected() {
    const item = filtered[selected]
    if (!item) return
    onClose()
    onInsert(item.template)
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      e.preventDefault()
      onClose()
    } else if (e.key === 'ArrowDown') {
      e.preventDefault()
      selected = Math.min(selected + 1, filtered.length - 1)
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      selected = Math.max(selected - 1, 0)
    } else if (e.key === 'Enter') {
      e.preventDefault()
      runSelected()
    }
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="palette-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="palette" role="dialog" aria-modal="true" aria-label="Палитра сниппетов" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <input
      class="palette-input"
      bind:value={query}
      placeholder="Поиск шага…"
      on:input={onQueryInput}
      autofocus
    />
    <ul class="palette-list">
      {#if loading}
        <li class="empty">Загрузка…</li>
      {:else if filtered.length === 0}
        <li class="empty">Шаги не найдены</li>
      {:else}
        {#each filtered as entry, i}
          <li>
            <button class:selected={i === selected} on:click={() => { onClose(); onInsert(entry.template) }}>
              <span class="label">{entry.template}</span>
              <span class="group">{entry.category}</span>
              <span class="help">{entry.help}</span>
            </button>
          </li>
        {/each}
      {/if}
    </ul>
  </div>
</div>

<style>
  .palette-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.45);
    z-index: 1200;
    display: flex;
    align-items: flex-start;
    justify-content: center;
    padding-top: 12vh;
  }

  .palette {
    width: min(640px, 92vw);
    max-height: 70vh;
    background: var(--color-sidebar);
    border: 1px solid var(--color-border);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
    display: flex;
    flex-direction: column;
  }

  .palette-input {
    margin: 12px;
    padding: 8px 10px;
    background: var(--color-input);
    border: 1px solid var(--color-border);
    color: var(--color-text);
    font-size: 14px;
  }

  .palette-list {
    list-style: none;
    margin: 0;
    padding: 0 0 8px;
    overflow-y: auto;
    flex: 1;
  }

  .palette-list button {
    display: grid;
    grid-template-columns: 1fr auto;
    grid-template-rows: auto auto;
    gap: 2px 12px;
    width: 100%;
    text-align: left;
    padding: 8px 14px;
    border: none;
    background: transparent;
    color: var(--color-text);
    cursor: pointer;
    font-size: 12px;
  }

  .palette-list button:hover,
  .palette-list button.selected {
    background: var(--color-input);
  }

  .label {
    font-family: var(--font-mono);
    color: var(--color-text);
    grid-column: 1;
  }

  .group {
    color: var(--color-muted);
    font-size: 10px;
    grid-column: 2;
    grid-row: 1;
  }

  .help {
    color: var(--color-muted);
    font-size: 11px;
    grid-column: 1 / -1;
  }

  .empty {
    padding: 16px;
    color: var(--color-muted);
    font-size: 12px;
  }
</style>
