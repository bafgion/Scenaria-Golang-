<script lang="ts">
  import { onMount } from 'svelte'
  import { HTTPAuthForHost, ListHTTPAuthHosts, RemoveHTTPAuth, SaveHTTPAuth } from '../../wailsjs/go/wailsapp/App'
  import { gui } from '../../wailsjs/go/models'

  export let initialHost = ''
  export let onClose: () => void = () => {}

  let host = ''
  let username = ''
  let password = ''
  let hosts: string[] = []
  let busy = false
  let message = ''

  onMount(async () => {
    host = initialHost
    await refreshHosts()
    if (host) await loadHost(host)
  })

  async function refreshHosts() {
    hosts = await ListHTTPAuthHosts().catch(() => [])
  }

  function onHostSelect(e: Event) {
    loadHost((e.currentTarget as HTMLSelectElement).value)
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

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="palette-backdrop" role="presentation" on:click={onClose}>
  <div class="palette http-auth" role="dialog" aria-label="HTTP авторизация" on:click|stopPropagation>
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
    <div class="actions">
      <button type="button" class="primary" disabled={busy || !host.trim() || !username.trim()} on:click={save}>
        Сохранить
      </button>
      <button type="button" disabled={busy || !host.trim()} on:click={remove}>Удалить</button>
      <button type="button" on:click={onClose}>Закрыть</button>
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

  input,
  select {
    padding: 6px 8px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: var(--color-text);
  }

  .actions {
    display: flex;
    gap: 8px;
    justify-content: flex-end;
    flex-wrap: wrap;
  }
</style>
