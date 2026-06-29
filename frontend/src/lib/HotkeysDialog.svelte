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

<div class="palette-backdrop" role="presentation" on:click={onClose}>
  <div class="palette hotkeys" role="dialog" aria-label="Горячие клавиши" on:click|stopPropagation>
    <h3>Горячие клавиши</h3>
    <ul>
      {#each hotkeys as cmd}
        <li>
          <span class="label">{cmd.label}</span>
          <span class="shortcut">{cmd.shortcut}</span>
        </li>
      {/each}
    </ul>
    <div class="actions">
      <button type="button" class="primary" on:click={onClose}>OK</button>
    </div>
  </div>
</div>

<style>
  .hotkeys {
    width: min(480px, 92vw);
    max-height: 70vh;
    overflow: auto;
    padding: 16px;
    background: var(--color-sidebar);
    border: 1px solid var(--color-border);
    border-radius: 4px;
  }

  h3 {
    margin: 0 0 12px;
    font-size: 14px;
    font-weight: 600;
  }

  ul {
    list-style: none;
    margin: 0;
    padding: 0;
  }

  li {
    display: flex;
    justify-content: space-between;
    gap: 16px;
    padding: 6px 0;
    border-bottom: 1px solid var(--color-divider);
    font-size: 12px;
  }

  .label {
    color: var(--color-text);
  }

  .shortcut {
    color: var(--color-muted);
    font-family: var(--font-mono);
    white-space: nowrap;
  }

  .actions {
    margin-top: 12px;
    display: flex;
    justify-content: flex-end;
  }

  .actions button {
    padding: 5px 16px;
    font-size: 12px;
    background: var(--color-primary);
    border: 1px solid var(--color-primary);
    border-radius: 3px;
    color: #fff;
  }
</style>
