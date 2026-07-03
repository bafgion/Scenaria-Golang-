<script lang="ts">
  import { onMount } from 'svelte'
  import { HTTPAuthForHost, ListHTTPAuthHosts, RemoveHTTPAuth, SaveHTTPAuth } from '../../wailsjs/go/wailsapp/App'
  import { gui } from '../../wailsjs/go/models'

  export let initialHost = ''
  export let onCancel: () => void

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
      message = 'Сохранено'
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
      message = 'Удалено'
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
  <div class="modal http-auth" role="dialog" aria-modal="true" aria-label="HTTP авторизация" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>HTTP Basic Auth</h3>
    <p class="hint">Логин и пароль сохраняются для хоста и применяются при записи и запуске.</p>
    {#if hosts.length}
      <label>Сохранённые хосты
        <select on:change={onHostSelect}>
          <option value="">— выбрать —</option>
          {#each hosts as item}
            <option value={item} selected={item === host}>{item}</option>
          {/each}
        </select>
      </label>
    {/if}
    <label>Хост <input bind:value={host} placeholder="example.com" /></label>
    <label>Логин <input bind:value={username} autocomplete="username" /></label>
    <label>Пароль <input type="password" bind:value={password} autocomplete="current-password" /></label>
    {#if message}<p class="message">{message}</p>{/if}
    <div class="modal-actions">
      <button type="button" class="primary" disabled={busy || !host.trim() || !username.trim()} on:click={save}>
        Сохранить
      </button>
      <button type="button" disabled={busy || !host.trim()} on:click={remove}>Удалить</button>
      <button type="button" on:click={closeDialog}>Закрыть</button>
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
