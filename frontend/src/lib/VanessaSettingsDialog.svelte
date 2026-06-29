<script lang="ts">
  import { onMount } from 'svelte'
  import { ReadVanessaSettingsJSON, SaveVanessaSettingsJSON } from '../../wailsjs/go/wailsapp/App'

  export let onClose: () => void = () => {}
  export let onLog: (message: string) => void = () => {}

  let jsonText = ''
  let busy = false
  let error = ''

  onMount(async () => {
    busy = true
    try {
      jsonText = await ReadVanessaSettingsJSON()
    } catch (e: any) {
      error = String(e)
    } finally {
      busy = false
    }
  })

  async function save() {
    busy = true
    error = ''
    try {
      JSON.parse(jsonText)
      await SaveVanessaSettingsJSON(jsonText)
      onLog('Настройки Vanessa сохранены в .scenaria/vanessa.json')
      onClose()
    } catch (e: any) {
      error = String(e)
    } finally {
      busy = false
    }
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal wide tall" role="dialog" aria-modal="true" aria-label="Настройки Vanessa" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>Настройки Vanessa (1C)</h3>
    <p class="hint">Файл <code>.scenaria/vanessa.json</code> — пути к платформе 1C, EPF Vanessa и каталогу запусков.</p>
    <textarea bind:value={jsonText} spellcheck="false" disabled={busy}></textarea>
    {#if error}<p class="error">{error}</p>{/if}
    <div class="modal-actions">
      <button type="button" class="primary" on:click={save} disabled={busy || !jsonText.trim()}>Сохранить</button>
      <button type="button" on:click={onClose} disabled={busy}>Отмена</button>
    </div>
  </div>
</div>

<style>
  textarea {
    width: 100%;
    min-height: 320px;
    font-family: var(--font-mono, monospace);
    font-size: 12px;
    resize: vertical;
    margin: 8px 0;
  }

  .hint {
    font-size: 11px;
    color: var(--color-muted);
    margin: 0 0 4px;
  }

  .error {
    color: var(--color-error, #c62828);
    font-size: 12px;
    margin: 0;
  }
</style>
