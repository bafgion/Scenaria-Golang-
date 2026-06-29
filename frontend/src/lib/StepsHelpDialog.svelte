<script lang="ts">
  import { onMount } from 'svelte'
  import { SearchSteps } from '../../wailsjs/go/wailsapp/App'

  export type StepEntry = { category: string; template: string; help: string }

  export let onClose: () => void = () => {}
  export let onInsert: ((template: string) => void) | null = null

  let query = ''
  let entries: StepEntry[] = []
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
  $: current = filtered[selected]

  async function onQueryInput() {
    try {
      entries = await SearchSteps(query)
      selected = 0
    } catch {
      /* keep previous */
    }
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
    } else if (e.key === 'Enter' && onInsert && current) {
      e.preventDefault()
      onInsert(current.template)
      onClose()
    }
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="palette-backdrop" role="presentation" on:click={onClose}>
  <div class="palette steps-help" role="dialog" aria-label="Справка по шагам" on:click|stopPropagation>
    <h3>Справка по шагам</h3>
    <input
      class="search"
      bind:value={query}
      placeholder="Поиск по шаблону, категории или описанию…"
      on:input={onQueryInput}
    />
    <div class="body">
      <ul class="list">
        {#if loading}
          <li class="empty">Загрузка…</li>
        {:else if filtered.length === 0}
          <li class="empty">Шаги не найдены</li>
        {:else}
          {#each filtered as entry, i}
            <li>
              <button type="button" class:selected={i === selected} on:click={() => (selected = i)}>
                <span class="cat">{entry.category}</span>
                <span class="tpl">{entry.template}</span>
              </button>
            </li>
          {/each}
        {/if}
      </ul>
      <div class="detail">
        {#if current}
          <div class="detail-cat">{current.category}</div>
          <pre class="detail-template">{current.template}</pre>
          <p class="detail-help">{current.help}</p>
          {#if onInsert}
            <button type="button" class="insert-btn" on:click={() => { onInsert(current.template); onClose() }}>
              Вставить в редактор
            </button>
          {/if}
        {:else}
          <p class="empty">Выберите шаг из списка</p>
        {/if}
      </div>
    </div>
    <div class="actions">
      <button type="button" on:click={onClose}>Закрыть</button>
    </div>
  </div>
</div>

<style>
  .steps-help {
    width: min(920px, 96vw);
    max-height: 86vh;
    display: flex;
    flex-direction: column;
  }

  h3 {
    margin: 0 0 10px;
    font-size: 14px;
  }

  .search {
    width: 100%;
    padding: 6px 8px;
    margin-bottom: 10px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: var(--color-text);
    box-sizing: border-box;
  }

  .body {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
    flex: 1;
    min-height: 320px;
    overflow: hidden;
  }

  .list {
    list-style: none;
    margin: 0;
    padding: 0;
    overflow: auto;
    border: 1px solid var(--color-divider);
    border-radius: 3px;
    background: var(--color-input);
  }

  .list button {
    width: 100%;
    text-align: left;
    padding: 6px 8px;
    border: none;
    border-bottom: 1px solid var(--color-divider);
    background: transparent;
    color: var(--color-text);
    cursor: pointer;
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .list button:hover,
  .list button.selected {
    background: var(--color-selection);
  }

  .cat {
    font-size: 10px;
    color: var(--color-muted);
  }

  .tpl {
    font-size: 12px;
    font-family: var(--font-mono, monospace);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .detail {
    overflow: auto;
    padding: 10px;
    border: 1px solid var(--color-divider);
    border-radius: 3px;
    background: var(--color-toolbar);
  }

  .detail-cat {
    font-size: 11px;
    color: var(--color-muted);
    margin-bottom: 8px;
  }

  .detail-template {
    margin: 0 0 10px;
    padding: 8px;
    background: var(--color-input);
    border-radius: 3px;
    font-size: 12px;
    white-space: pre-wrap;
    word-break: break-word;
  }

  .detail-help {
    margin: 0;
    font-size: 13px;
    line-height: 1.45;
    color: var(--color-text);
  }

  .insert-btn {
    margin-top: 12px;
    padding: 6px 12px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-accent);
    color: var(--color-on-accent, #fff);
    cursor: pointer;
  }

  .empty {
    padding: 12px;
    color: var(--color-muted);
    font-size: 12px;
  }

  .actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 12px;
  }

  .actions button {
    padding: 6px 12px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: var(--color-text);
  }
</style>
