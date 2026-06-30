<script lang="ts">
  import { onDestroy, onMount } from 'svelte'
  import { SearchSteps } from '../../wailsjs/go/wailsapp/App'
  import { asStepSearchQuery } from './stepSearch'
  import type { SnippetEntry } from './stepTypes'

  export let onClose: () => void
  export let onInsert: (template: string) => void

  let query = ''
  let entries: SnippetEntry[] = []
  let selected = 0
  let loading = true
  let searchDebounceTimer: ReturnType<typeof setTimeout> | null = null

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
    if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
    searchDebounceTimer = setTimeout(async () => {
      searchDebounceTimer = null
      try {
        entries = await SearchSteps(asStepSearchQuery(query))
      } catch {
        /* keep previous */
      }
    }, 200)
  }

  onDestroy(() => {
    if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
  })

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
  <div class="palette snippet-palette" role="dialog" aria-modal="true" aria-label="Палитра сниппетов" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <input
      class="palette-input"
      bind:value={query}
      placeholder="Поиск шага…"
      on:input={onQueryInput}
      autofocus
    />
    <ul class="palette-list snippet-grid">
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
