<script lang="ts">
  import { onMount } from 'svelte'
  import { createTranslator, locale } from './i18n'
  import { ListPlugins, InstallPlugin, UninstallPlugin } from '../../wailsjs/go/wailsapp/App'
  import type { gui } from '../../wailsjs/go/models'

  export let onClose: () => void = () => {}
  export let onRunPlugin: (name: string, dryRun: boolean) => void = () => {}
  export let childModalOpen = false
  export let onAskConfirm: (message: string) => Promise<boolean> = (message) =>
    Promise.resolve(window.confirm(message))

  $: tr = createTranslator($locale)

  let entries: gui.PluginEntryDTO[] = []
  let name = 'vanessa'
  let source = ''
  let busy = false
  let error = ''
  let loading = true

  $: runnablePluginEntries = entries.filter((e) => e.runnable)

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
      error = tr('dialogs.plugins.errorNameSource')
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
    if (!(await onAskConfirm(tr('dialogs.plugins.confirmUninstall', { name: entry.name })))) return
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

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
  function onBackdropClose() {
    if (!childModalOpen) onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onBackdropClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal wide tall plugins-dialog" role="dialog" aria-modal="true" aria-label={tr('dialogs.plugins.ariaLabel')} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{tr('dialogs.plugins.title')}</h3>
    <p class="hint">{tr('dialogs.plugins.hint')}</p>

    {#if loading}
      <p class="empty">{tr('dialogs.common.loading')}</p>
    {:else if entries.length === 0}
      <p class="empty">{tr('dialogs.plugins.empty')}</p>
    {:else}
      <table>
        <thead>
          <tr><th>{tr('dialogs.plugins.colName')}</th><th>{tr('dialogs.plugins.colSource')}</th><th></th></tr>
        </thead>
        <tbody>
          {#each entries as entry}
            <tr>
              <td>
                <div>{entry.name}</div>
                {#if entry.description}<div class="meta">{entry.description}</div>{/if}
              </td>
              <td class="source" title={entry.source}>{entry.source}</td>
              <td><button type="button" class="danger" disabled={busy} on:click={() => uninstall(entry)}>{tr('dialogs.plugins.uninstall')}</button></td>
            </tr>
          {/each}
        </tbody>
      </table>
    {/if}

    <div class="install-form">
      <label>{tr('dialogs.plugins.name')} <input bind:value={name} placeholder={tr('dialogs.plugins.namePlaceholder')} disabled={busy} /></label>
      <label>{tr('dialogs.plugins.source')}
        <input bind:value={source} placeholder={tr('dialogs.plugins.sourcePlaceholder')} disabled={busy} />
      </label>
      <button type="button" class="primary" disabled={busy} on:click={install}>{tr('dialogs.plugins.install')}</button>
    </div>

    {#each runnablePluginEntries as entry (entry.name)}
      <div class="runners">
        <span>{pluginTitle(entry)}:</span>
        {#if entry.vanessa}
          <button type="button" disabled={busy} on:click={() => { onClose(); onRunPlugin(entry.name, true) }}>{tr('dialogs.plugins.dryRun')}</button>
          <button type="button" disabled={busy} on:click={() => { onClose(); onRunPlugin(entry.name, false) }}>{tr('dialogs.plugins.run')}</button>
        {:else}
          <button type="button" disabled={busy} on:click={() => { onClose(); onRunPlugin(entry.name, true) }}>{tr('dialogs.plugins.dryRunEllipsis')}</button>
          <button type="button" disabled={busy} on:click={() => { onClose(); onRunPlugin(entry.name, false) }}>{tr('dialogs.plugins.runEllipsis')}</button>
        {/if}
      </div>
    {/each}

    {#if error}<p class="error">{error}</p>{/if}

    <div class="modal-actions">
      <button type="button" on:click={onClose}>{tr('dialogs.common.close')}</button>
    </div>
  </div>
</div>

<style>
  .plugins-dialog {
    width: min(720px, 96vw);
    max-height: 86vh;
    overflow: auto;
    min-width: 0;
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
    table-layout: fixed;
  }

  th, td {
    padding: 6px 8px;
    border-bottom: 1px solid var(--color-divider);
    text-align: left;
    vertical-align: top;
    min-width: 0;
  }

  .source {
    max-width: min(320px, 38vw);
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
    min-width: 0;
  }

  label {
    display: grid;
    gap: 4px;
    font-size: 11px;
    color: var(--color-muted);
  }

  .runners {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 8px;
    margin-bottom: 10px;
    font-size: 12px;
    min-width: 0;
  }

  .empty, .error {
    font-size: 12px;
    margin: 8px 0;
  }

  .error {
    color: var(--color-error);
  }

  button.danger {
    color: var(--color-error);
    background: transparent;
    border-color: var(--color-border);
  }

  @media (max-width: 620px) {
    table,
    thead,
    tbody,
    tr,
    th,
    td {
      display: block;
    }

    thead {
      display: none;
    }

    tr {
      padding: 8px 0;
      border-bottom: 1px solid var(--color-divider);
    }

    td {
      border-bottom: none;
    }

    .source {
      max-width: none;
      white-space: normal;
      word-break: break-all;
    }

    td button {
      width: 100%;
    }

    .runners button {
      flex: 1 1 140px;
    }
  }
</style>
