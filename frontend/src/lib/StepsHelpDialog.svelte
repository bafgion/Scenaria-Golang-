<script lang="ts">
  import { onMount, tick } from 'svelte'
  import { SearchSteps } from '../../wailsjs/go/wailsapp/App'
  import { asStepSearchQuery } from './stepSearch'
  import type { StepHelpEntry } from './stepTypes'

  function entryText(entry: StepHelpEntry): string {
    return [
      entry.label,
      entry.action,
      entry.category,
      entry.description,
      entry.template,
      entry.example,
      entry.help,
      ...(entry.parameters ?? []),
    ].join(' ')
  }

  function displayLabel(entry: StepHelpEntry): string {
    if (entry.label) {
      return entry.action ? `${entry.label} (${entry.action})` : entry.label
    }
    return entry.template
  }

  function detailDescription(entry: StepHelpEntry): string {
    return entry.description || entry.help || ''
  }

  function detailExample(entry: StepHelpEntry): string {
    return entry.example || entry.template || ''
  }

  export let onClose: () => void = () => {}
  export let onInsert: ((template: string) => void) | null = null
  export let initialQuery = ''

  let query = ''
  let entries: StepHelpEntry[] = []
  let selected = 0
  let loading = true
  let searchInput: HTMLInputElement

  async function loadEntries(q: unknown) {
    const searchQ = asStepSearchQuery(q)
    loading = true
    try {
      entries = await SearchSteps(searchQ)
    } catch {
      entries = []
    } finally {
      loading = false
    }
  }

  onMount(async () => {
    query = asStepSearchQuery(initialQuery)
    await loadEntries(query)
    await tick()
    searchInput?.focus()
  })

  $: filtered = entries.filter((e) => {
    const q = query.trim().toLowerCase()
    if (!q) return true
    return entryText(e).toLowerCase().includes(q)
  })

  $: if (selected >= filtered.length) selected = Math.max(0, filtered.length - 1)
  $: current = filtered[selected]

  async function onQueryInput() {
    await loadEntries(query)
    selected = 0
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

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="palette-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="palette steps-help" role="dialog" aria-modal="true" aria-label="Справка по шагам" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>Справка по шагам</h3>
    <input
      bind:this={searchInput}
      class="search"
      bind:value={query}
      placeholder="Поиск по шаблону, категории или описанию… (params — наборы параметров)"
      on:input={onQueryInput}
    />
    <details class="params-help" open={!query.trim() || query.toLowerCase().includes('param')}>
      <summary>Наборы параметров (.params.json)</summary>
      <p>
        Для сценариев-шаблонов без таблицы «Примеры» положите рядом с <code>имя.feature</code> файл
        <code>имя.params.json</code>:
      </p>
      <pre class="params-example">{`{
  "scenarios": {
    "Название сценария": [
      { "url": "/catalog", "title": "Items" },
      { "url": "/offers", "title": "Offers" }
    ]
  }
}`}</pre>
      <p class="params-note">
        Ключи совпадают с плейсхолдерами в шагах (<code>&lt;url&gt;</code>). Если в feature уже есть таблица
        «Примеры», используется она. При запуске runner разворачивает каждую строку в отдельный прогон.
      </p>
    </details>
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
                <span class="tpl">{displayLabel(entry)}</span>
              </button>
            </li>
          {/each}
        {/if}
      </ul>
      <div class="detail">
        {#if current}
          <header class="detail-header">
            <h4 class="detail-title">{current.label || displayLabel(current)}</h4>
            <p class="detail-meta">
              <span class="badge">{current.category}</span>
              {#if current.action}<code class="action">{current.action}</code>{/if}
            </p>
          </header>
          {#if detailDescription(current)}
            <p class="detail-desc">{detailDescription(current)}</p>
          {/if}
          {#if current.parameters?.length}
            <section class="detail-block">
              <div class="block-title">Параметры</div>
              {#each current.parameters as param}
                <p class="param-line">{param}</p>
              {/each}
            </section>
          {/if}
          <section class="detail-block">
            <div class="block-title">Пример</div>
            <pre class="detail-example">{detailExample(current)}</pre>
          </section>
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
    padding: 16px;
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

  .params-help {
    margin: 0 0 10px;
    padding: 8px 10px;
    border: 1px solid var(--color-divider);
    border-radius: 3px;
    background: rgba(0, 122, 204, 0.06);
    font-size: 11px;
    color: var(--color-muted);
  }

  .params-help summary {
    cursor: pointer;
    color: var(--color-text);
    font-weight: 600;
    margin-bottom: 6px;
  }

  .params-example {
    margin: 6px 0;
    padding: 8px;
    background: var(--color-input);
    border-radius: 3px;
    font-size: 10px;
    overflow: auto;
  }

  .params-note {
    margin: 0;
    line-height: 1.45;
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
    background: var(--color-selected);
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
    padding: 12px;
    border: 1px solid var(--color-divider);
    border-radius: 3px;
    background: var(--color-toolbar);
  }

  .detail-header {
    margin-bottom: 12px;
    padding-bottom: 12px;
    border-bottom: 1px solid var(--color-divider);
  }

  .detail-title {
    margin: 0 0 8px;
    font-size: 16px;
    font-weight: 600;
    color: var(--color-text);
  }

  .detail-meta {
    margin: 0;
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 8px;
  }

  .badge {
    display: inline-block;
    padding: 1px 8px;
    border: 1px solid var(--color-border);
    border-radius: 4px;
    font-size: 10px;
    font-weight: 600;
    color: var(--color-muted);
  }

  .action {
    font-family: var(--font-mono, monospace);
    font-size: 11px;
    color: #9cdc8a;
    background: transparent;
  }

  .detail-desc {
    margin: 0 0 12px;
    font-size: 13px;
    line-height: 1.45;
    color: var(--color-text);
  }

  .detail-block {
    margin-top: 12px;
  }

  .block-title {
    margin-bottom: 8px;
    font-size: 10px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--color-muted);
  }

  .param-line {
    margin: 0 0 6px;
    font-size: 12px;
    line-height: 1.45;
    color: var(--color-text);
  }

  .detail-example {
    margin: 0;
    padding: 0 0 0 10px;
    border-left: 2px solid var(--color-border);
    background: transparent;
    font-family: var(--font-mono, monospace);
    font-size: 12px;
    line-height: 1.5;
    color: #ce9178;
    white-space: pre-wrap;
    word-break: break-word;
  }

  .insert-btn {
    margin-top: 12px;
    padding: 6px 12px;
    border: 1px solid var(--color-primary);
    border-radius: 3px;
    background: var(--color-primary);
    color: #fff;
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
