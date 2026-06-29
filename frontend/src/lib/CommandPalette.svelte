<script lang="ts">
  export type PaletteCommand = {
    id: string
    label: string
    group: string
    shortcut?: string
    run: () => void
  }

  export let commands: PaletteCommand[] = []
  export let onClose: () => void

  let query = ''
  let selected = 0

  $: filtered = commands.filter((c) => {
    const q = query.trim().toLowerCase()
    if (!q) return true
    return (c.label + c.group).toLowerCase().includes(q)
  })

  $: if (selected >= filtered.length) selected = Math.max(0, filtered.length - 1)

  function runSelected() {
    const cmd = filtered[selected]
    if (cmd) {
      onClose()
      cmd.run()
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
  <div class="palette" role="dialog" aria-modal="true" aria-label="Палитра команд" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <input class="palette-input" bind:value={query} placeholder="Введите команду…" autofocus />
    <ul class="palette-list">
      {#each filtered as cmd, i}
        <li>
          <button class:selected={i === selected} on:click={() => { onClose(); cmd.run() }}>
            <span class="label">{cmd.label}</span>
            <span class="group">{cmd.group}</span>
            {#if cmd.shortcut}<span class="shortcut">{cmd.shortcut}</span>{/if}
          </button>
        </li>
      {:else}
        <li class="empty">Команды не найдены</li>
      {/each}
    </ul>
  </div>
</div>

<style>
  .palette-backdrop {
    position: fixed;
    inset: 0;
    z-index: 1500;
    background: rgba(0, 0, 0, 0.45);
    display: flex;
    align-items: flex-start;
    justify-content: center;
    padding-top: 12vh;
  }

  .palette {
    width: min(560px, 92vw);
    background: var(--color-sidebar);
    border: 1px solid var(--color-border);
    border-radius: 4px;
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.5);
    overflow: hidden;
  }

  .palette-input {
    width: 100%;
    padding: 12px 14px;
    border: none;
    border-bottom: 1px solid var(--color-border);
    background: var(--color-toolbar);
    color: var(--color-text);
    font-size: 14px;
    outline: none;
  }

  .palette-list {
    list-style: none;
    margin: 0;
    padding: 4px 0;
    max-height: 360px;
    overflow-y: auto;
  }

  .palette-list button {
    display: grid;
    grid-template-columns: 1fr auto auto;
    gap: 12px;
    width: 100%;
    text-align: left;
    padding: 8px 14px;
    border: none;
    background: transparent;
    color: var(--color-text);
    font-size: 13px;
  }

  .palette-list button.selected,
  .palette-list button:hover {
    background: var(--color-selected);
  }

  .group {
    color: var(--color-muted);
    font-size: 11px;
  }

  .shortcut {
    color: var(--color-muted);
    font-size: 11px;
    font-family: var(--font-mono);
  }

  .empty {
    padding: 12px 14px;
    color: var(--color-muted);
    font-size: 12px;
  }
</style>
