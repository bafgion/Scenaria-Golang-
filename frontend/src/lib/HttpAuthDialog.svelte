<script lang="ts">
  import { onMount } from 'svelte'
  import { createTranslator, locale } from './i18n'
  import { HTTPAuthForHost, ListHTTPAuthHosts, RemoveHTTPAuth, SaveHTTPAuth } from '../../wailsjs/go/wailsapp/App'
  import { gui } from '../../wailsjs/go/models'

  export let initialHost = ''
  export let onCancel: () => void

  $: tr = createTranslator($locale)

  let host = ''
  let username = ''
  let password = ''
  let hosts: string[] = []
  let busy = false
  let message = ''

  onMount(() => {
    host = initialHost
    void refreshHosts().then(() => {
      if (host) void loadHost(host)
    })
  })

  async function refreshHosts() {
    const listed = await ListHTTPAuthHosts().catch(() => [] as string[])
    hosts = listed ?? []
  }

  function onHostSelect(e: Event) {
    void loadHost((e.currentTarget as HTMLSelectElement).value)
  }

  async function loadHost(value: string) {
    host = value
    const creds = await HTTPAuthForHost(value).catch(() => new gui.HTTPAuthCredentials())
    username = creds.username
    password = creds.password
  }

  async function save() {
    if (!host.trim()) return
    busy = true
    message = ''
    try {
      await SaveHTTPAuth(new gui.HTTPAuthRequest({
        host: host.trim(),
        username: username.trim(),
        password,
      }))
      message = tr('dialogs.httpAuth.saved')
      await refreshHosts()
    } catch (e: any) {
      message = String(e)
    } finally {
      busy = false
    }
  }

  async function remove() {
    if (!host.trim()) return
    busy = true
    message = ''
    try {
      await RemoveHTTPAuth(host.trim())
      username = ''
      password = ''
      message = tr('dialogs.httpAuth.removed')
      await refreshHosts()
    } catch (e: any) {
      message = String(e)
    } finally {
      busy = false
    }
  }

  function closeDialog() {
    onCancel()
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      e.stopPropagation()
      closeDialog()
    }
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop modal-layer-top" role="presentation" on:click={closeDialog}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal http-auth" role="dialog" aria-modal="true" aria-label={tr('dialogs.httpAuth.ariaLabel')} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{tr('dialogs.httpAuth.title')}</h3>
    <p class="hint">{tr('dialogs.httpAuth.hint')}</p>
    {#if hosts.length}
      <label>{tr('dialogs.httpAuth.savedHosts')}
        <select on:change={onHostSelect}>
          <option value="">{tr('dialogs.httpAuth.selectHost')}</option>
          {#each hosts as item}
            <option value={item} selected={item === host}>{item}</option>
          {/each}
        </select>
      </label>
    {/if}
    <label>{tr('dialogs.httpAuth.host')} <input bind:value={host} placeholder={tr('dialogs.httpAuth.hostPlaceholder')} /></label>
    <label>{tr('dialogs.httpAuth.username')} <input bind:value={username} autocomplete="username" /></label>
    <label>{tr('dialogs.httpAuth.password')} <input type="password" bind:value={password} autocomplete="current-password" /></label>
    {#if message}<p class="message">{message}</p>{/if}
    <div class="modal-actions">
      <button type="button" class="primary" disabled={busy || !host.trim() || !username.trim()} on:click={save}>
        {tr('dialogs.common.save')}
      </button>
      <button type="button" disabled={busy || !host.trim()} on:click={remove}>{tr('dialogs.common.delete')}</button>
      <button type="button" on:click={closeDialog}>{tr('dialogs.common.close')}</button>
    </div>
  </div>
</div>

<style>
  .http-auth {
    width: min(440px, 92vw);
  }

  .hint,
  .message {
    margin: 0 0 10px;
    font-size: 12px;
    color: var(--color-muted);
  }

  label {
    display: flex;
    flex-direction: column;
    gap: 4px;
    margin-bottom: 10px;
    font-size: 12px;
  }
</style>
