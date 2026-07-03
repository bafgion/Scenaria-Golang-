<script lang="ts">
  import { createTranslator, locale } from './i18n'
  import type { PaletteCommand } from './paletteTypes'

  export let commands: PaletteCommand[] = []
  export let onClose: () => void = () => {}

  $: tr = createTranslator($locale)
  $: sortLocale = $locale === 'en' ? 'en' : 'ru'

  $: hotkeys = commands
    .filter((c) => c.shortcut)
    .sort((a, b) => a.group.localeCompare(b.group, sortLocale) || a.label.localeCompare(b.label, sortLocale))

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
  <div class="modal wide hotkeys-dialog" role="dialog" aria-modal="true" aria-label={tr('menus.hotkeys')} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{tr('menus.hotkeys')}</h3>
    <p class="hotkeys-note">{tr('settings.hotkeysEditorNote')}</p>
    <ul class="hotkeys-list">
      {#each hotkeys as cmd}
        <li>
          <span class="label">{cmd.label}</span>
          <span class="shortcut">{cmd.shortcut}</span>
        </li>
      {/each}
    </ul>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={onClose}>{tr('common.ok')}</button>
    </div>
  </div>
</div>

<style>
  .hotkeys-dialog {
    width: min(480px, 92vw);
    max-height: 70vh;
  }

  .hotkeys-note {
    margin: 0 0 10px;
    font-size: 12px;
    color: var(--color-muted);
  }

  .hotkeys-list {
    list-style: none;
    margin: 0 0 16px;
    padding: 0;
    max-height: 50vh;
    overflow-y: auto;
  }

  .hotkeys-list li {
    display: flex;
    justify-content: space-between;
    gap: 12px;
    padding: 6px 0;
    border-bottom: 1px solid var(--color-border);
    font-size: 13px;
  }

  .hotkeys-list .label {
    flex: 1;
    min-width: 0;
  }

  .hotkeys-list .shortcut {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--color-muted);
    white-space: nowrap;
  }
</style>
