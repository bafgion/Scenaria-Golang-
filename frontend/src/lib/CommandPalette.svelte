<script lang="ts">
  import type { PaletteCommand } from './paletteTypes'

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
  <div class="palette command-palette" role="dialog" aria-modal="true" aria-label="Палитра команд" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
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
