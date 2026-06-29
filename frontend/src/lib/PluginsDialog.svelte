<script lang="ts">
  import { onMount } from 'svelte'
  import { ListPlugins, InstallPlugin, UninstallPlugin } from '../../wailsjs/go/wailsapp/App'
  import type { gui } from '../../wailsjs/go/models'

  export let onClose: () => void = () => {}
  export let onRunPlugin: (name: string, dryRun: boolean) => void = () => {}
  export let onAskConfirm: (message: string) => Promise<boolean> = (message) =>
    Promise.resolve(window.confirm(message))

  let entries: gui.PluginEntryDTO[] = []
  let name = 'vanessa'
  let source = ''
  let busy = false
  let error = ''
  let loading = true

  onMount(() => void refresh())

  async function refresh() {
    loading = true
    error = ''
    try {
      entries = await ListPlugins()
    } catch (e: any) {
      entries = []
      error = String(e)
    } finally {
      loading = false
    }
  }

  async function install() {
    if (!name.trim() || !source.trim()) {
      error = 'Укажите имя и источник (URL или путь к .zip)'
      return
    }
    busy = true
    error = ''
    try {
      await InstallPlugin(name.trim(), source.trim())
      source = ''
      await refresh()
    } catch (e: any) {
      error = String(e)
    } finally {
      busy = false
    }
  }

  async function uninstall(entry: gui.PluginEntryDTO) {
    if (!(await onAskConfirm(`Удалить плагин «${entry.name}» из манифеста?`))) return
    busy = true
    error = ''
    try {
      await UninstallPlugin(entry.name)
      await refresh()
    } catch (e: any) {
      error = String(e)
    } finally {
      busy = false
    }
  }

  function pluginTitle(entry: gui.PluginEntryDTO): string {
    return entry.description || entry.id || entry.name
  }

  function runnableEntries(): gui.PluginEntryDTO[] {
    return entries.filter((e) => e.runnable)
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="palette-backdrop" role="presentation" on:click={onClose}>
  <div class="palette plugins-dialog" role="dialog" aria-label="Плагины" on:click|stopPropagation>
    <h3>Плагины проекта</h3>
    <p class="hint">Плагины устанавливаются в <code>addons/&lt;имя&gt;/</code> и регистрируются в <code>.scenaria/plugins.json</code>.</p>

    {#if loading}
      <p class="empty">Загрузка…</p>
    {:else if entries.length === 0}
      <p class="empty">Нет установленных плагинов</p>
    {:else}
      <table>
        <thead>
          <tr><th>Имя</th><th>Источник</th><th></th></tr>
        </thead>
        <tbody>
          {#each entries as entry}
            <tr>
              <td>
                <div>{entry.name}</div>
                {#if entry.description}<div class="meta">{entry.description}</div>{/if}
              </td>
              <td class="source" title={entry.source}>{entry.source}</td>
              <td><button type="button" class="danger" disabled={busy} on:click={() => uninstall(entry)}>Удалить</button></td>
            </tr>
          {/each}
        </tbody>
      </table>
    {/if}

    <div class="install-form">
      <label>Имя <input bind:value={name} placeholder="vanessa" disabled={busy} /></label>
      <label>Источник (URL или .zip)
        <input bind:value={source} placeholder="https://…/plugin.zip" disabled={busy} />
      </label>
      <button type="button" class="primary" disabled={busy} on:click={install}>Установить</button>
    </div>

    {#each runnableEntries() as entry (entry.name)}
      <div class="runners">
        <span>{pluginTitle(entry)}:</span>
        {#if entry.vanessa}
          <button type="button" disabled={busy} on:click={() => { onClose(); onRunPlugin(entry.name, true) }}>Dry-run</button>
          <button type="button" disabled={busy} on:click={() => { onClose(); onRunPlugin(entry.name, false) }}>Запуск</button>
        {:else}
          <button type="button" disabled={busy} on:click={() => { onClose(); onRunPlugin(entry.name, true) }}>Dry-run…</button>
          <button type="button" disabled={busy} on:click={() => { onClose(); onRunPlugin(entry.name, false) }}>Запуск…</button>
        {/if}
      </div>
    {/each}

    {#if error}<p class="error">{error}</p>{/if}

    <div class="actions">
      <button type="button" on:click={onClose}>Закрыть</button>
    </div>
  </div>
</div>

<style>
  .plugins-dialog {
    width: min(720px, 96vw);
    max-height: 86vh;
    overflow: auto;
  }

  h3 {
    margin: 0 0 8px;
    font-size: 14px;
  }

  .hint {
    margin: 0 0 12px;
    font-size: 11px;
    color: var(--color-muted);
    line-height: 1.4;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-size: 12px;
    margin-bottom: 12px;
  }

  th, td {
    padding: 6px 8px;
    border-bottom: 1px solid var(--color-divider);
    text-align: left;
  }

  .source {
    max-width: 320px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: var(--color-muted);
  }

  .meta {
    font-size: 10px;
    color: var(--color-muted);
    margin-top: 2px;
  }

  .install-form {
    display: grid;
    gap: 8px;
    margin-bottom: 12px;
  }

  label {
    display: grid;
    gap: 4px;
    font-size: 11px;
    color: var(--color-muted);
  }

  input {
    padding: 6px 8px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: var(--color-text);
    font-size: 12px;
  }

  .runners {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 10px;
    font-size: 12px;
  }

  .empty, .error {
    font-size: 12px;
    margin: 8px 0;
  }

  .error {
    color: var(--color-error);
  }

  .actions {
    display: flex;
    justify-content: flex-end;
  }

  button {
    padding: 6px 12px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: var(--color-text);
    font-size: 12px;
  }

  button.primary {
    background: var(--color-accent);
    color: var(--color-on-accent, #fff);
    border-color: var(--color-accent);
    justify-self: start;
  }

  button.danger {
    color: var(--color-error);
  }
</style>
