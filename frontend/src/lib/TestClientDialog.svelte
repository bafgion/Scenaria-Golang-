<script lang="ts">
  import { onMount } from 'svelte'
  import {
    ReadTestClientJSON,
    SaveTestClientJSON,
    DeleteTestClient,
    ListTestClients,
    CaptureBrowserSession,
  } from '../../wailsjs/go/wailsapp/App'

  export let testClients: string[] = []
  export let selectedName = ''
  export let browserOpen = false
  export let suggestName = ''
  export let onUse: (name: string) => void = () => {}
  export let onClose: () => void = () => {}
  export let onClientsChange: (names: string[]) => void = () => {}
  export let onLog: (message: string) => void = () => {}
  export let onAskConfirm: (message: string) => Promise<boolean> = (message) =>
    Promise.resolve(window.confirm(message))

  let editorName = ''
  let jsonText = ''
  let busy = false
  let error = ''
  let isNew = false

  const blankTemplate = (name: string) =>
    JSON.stringify(
      {
        name,
        base_url: 'https://example.com',
        cookies: [],
        local_storage: {},
      },
      null,
      2,
    )

  $: canSave = Boolean(editorName.trim()) && Boolean(jsonText.trim()) && !busy
  $: canUse = Boolean(selectedName) && !isNew
  $: canCapture = Boolean(editorName.trim()) && browserOpen && !busy

  async function selectClient(name: string) {
    isNew = false
    selectedName = name
    editorName = name
    error = ''
    busy = true
    try {
      jsonText = await ReadTestClientJSON(name)
    } catch (e: any) {
      jsonText = ''
      error = String(e)
    } finally {
      busy = false
    }
  }

  function startNew() {
    isNew = true
    selectedName = ''
    editorName = 'new_client'
    jsonText = blankTemplate('new_client')
    error = ''
  }

  async function saveClient() {
    const name = editorName.trim()
    if (!name) return
    busy = true
    error = ''
    try {
      JSON.parse(jsonText)
      await SaveTestClientJSON(name, jsonText)
      testClients = await ListTestClients()
      onClientsChange(testClients)
      isNew = false
      selectedName = name
      onLog(`TestClient сохранён: ${name}`)
    } catch (e: any) {
      error = String(e)
    } finally {
      busy = false
    }
  }

  async function captureFromBrowser() {
    const name = editorName.trim()
    if (!name || !browserOpen) return
    busy = true
    error = ''
    try {
      const summary = await CaptureBrowserSession(name)
      testClients = await ListTestClients()
      onClientsChange(testClients)
      isNew = false
      selectedName = name
      jsonText = await ReadTestClientJSON(name)
      onLog(summary)
      onUse(name)
    } catch (e: any) {
      error = String(e)
    } finally {
      busy = false
    }
  }

  async function deleteClient() {
    if (!selectedName || isNew) return
    if (!(await onAskConfirm(`Удалить TestClient «${selectedName}»?`))) return
    busy = true
    error = ''
    try {
      await DeleteTestClient(selectedName)
      testClients = await ListTestClients()
      onClientsChange(testClients)
      selectedName = testClients[0] || ''
      if (selectedName) await selectClient(selectedName)
      else {
        editorName = ''
        jsonText = ''
        isNew = false
      }
      onLog(`TestClient удалён`)
    } catch (e: any) {
      error = String(e)
    } finally {
      busy = false
    }
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }

  onMount(async () => {
    if (selectedName) await selectClient(selectedName)
    else if (suggestName.trim()) {
      const name = suggestName.trim()
      isNew = true
      editorName = name
      jsonText = blankTemplate(name)
    } else if (testClients.length === 0) startNew()
  })
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal wide tall test-client-dialog" role="dialog" aria-modal="true" aria-label="TestClient" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>TestClient</h3>
    <p class="hint capture-hint">
      Профиль с cookies и localStorage для повторного запуска с той же сессией. Откройте браузер, войдите на сайт, затем нажмите «Захватить из браузера».
    </p>
    <div class="test-client-body">
      <div class="client-list">
        {#each testClients as client}
          <button
            type="button"
            class="client-item"
            class:active={client === selectedName && !isNew}
            on:click={() => selectClient(client)}
          >
            {client}
          </button>
        {:else}
          <p class="hint empty-hint">Нет файлов в .scenaria/test_clients/</p>
        {/each}
        <button type="button" class="client-item new-btn" class:active={isNew} on:click={startNew}>+ Новый…</button>
      </div>
      <div class="client-editor">
        <label>
          Имя файла
          <input bind:value={editorName} disabled={!isNew && Boolean(selectedName)} placeholder="client_name" />
        </label>
        <label class="json-label">
          JSON
          <textarea bind:value={jsonText} spellcheck="false"></textarea>
        </label>
        {#if error}
          <p class="error">{error}</p>
        {/if}
      </div>
    </div>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={captureFromBrowser} disabled={!canCapture} title={browserOpen ? '' : 'Сначала откройте браузер (Ctrl+B)'}>
        Захватить из браузера
      </button>
      <button type="button" class="primary" on:click={saveClient} disabled={!canSave}>Сохранить</button>
      <button type="button" on:click={deleteClient} disabled={!selectedName || isNew || busy}>Удалить</button>
      <button type="button" class="primary" on:click={() => onUse(selectedName)} disabled={!canUse}>Использовать при запуске</button>
      <button type="button" on:click={onClose}>Закрыть</button>
    </div>
  </div>
</div>

<style>
  .test-client-body {
    display: grid;
    grid-template-columns: 180px 1fr;
    gap: 12px;
    min-height: 320px;
  }

  .empty-hint {
    padding: 8px;
    font-size: 12px;
    color: var(--color-muted);
  }

  .new-btn {
    border-style: dashed;
    margin-top: 4px;
  }

  .client-editor {
    display: flex;
    flex-direction: column;
    gap: 8px;
    min-width: 0;
  }

  .json-label {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }

  .json-label textarea {
    flex: 1;
    min-height: 240px;
    font-family: var(--font-mono, monospace);
    font-size: 12px;
    resize: vertical;
  }

  .error {
    color: var(--color-error, #c62828);
    font-size: 12px;
    margin: 0;
  }

  .capture-hint {
    margin: 0 0 10px;
    font-size: 12px;
    line-height: 1.45;
    color: var(--color-muted);
  }
</style>
