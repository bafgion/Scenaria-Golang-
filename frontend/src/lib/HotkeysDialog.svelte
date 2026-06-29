<script lang="ts">
  import type { PaletteCommand } from './CommandPalette.svelte'

  export let commands: PaletteCommand[] = []
  export let onClose: () => void = () => {}

  $: hotkeys = commands
    .filter((c) => c.shortcut)
    .sort((a, b) => a.group.localeCompare(b.group, 'ru') || a.label.localeCompare(b.label, 'ru'))

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      e.preventDefault()
      onClose()
    }
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal wide hotkeys-dialog" role="dialog" aria-modal="true" aria-label="Горячие клавиши" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>Горячие клавиши</h3>
    <ul class="hotkeys-list">
      {#each hotkeys as cmd}
        <li>
          <span class="label">{cmd.label}</span>
          <span class="shortcut">{cmd.shortcut}</span>
        </li>
      {/each}
    </ul>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={onClose}>OK</button>
    </div>
  </div>
</div>

<style>
  .hotkeys-dialog {
    width: min(480px, 92vw);
    max-height: 70vh;
  }

  .hotkeys-list {
    list-style: none;
    margin: 0;
    padding: 0;
    overflow: auto;
    flex: 1;
    border: 1px solid var(--color-border);
    border-radius: 4px;
    background: var(--color-bg);
  }

  li {
    display: flex;
    justify-content: space-between;
    gap: 16px;
    padding: 8px 12px;
    border-bottom: 1px solid var(--color-divider);
    font-size: 12px;
  }

  li:last-child {
    border-bottom: none;
  }

  .label {
    color: var(--color-text);
  }

  .shortcut {
    color: var(--color-muted);
    font-family: var(--font-mono);
    white-space: nowrap;
  }
</style>
