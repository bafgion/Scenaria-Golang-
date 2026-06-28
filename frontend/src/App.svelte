<script lang="ts">
  import { onMount } from 'svelte'
  import {
    Version,
    OpenProject,
    ReadFeature,
    SaveFeature,
    Run,
    Validate,
    SearchSteps,
    ListTestClients,
    InitProject,
  } from '../wailsjs/go/wailsapp/App'

  let version = ''
  let projectPath = ''
  let features: string[] = []
  let tags: string[] = []
  let selectedFeature = ''
  let editorText = ''
  let logText = ''
  let statusText = 'Готово'
  let runTag = ''
  let runTestClient = ''
  let runVars = ''
  let testClients: string[] = []
  let stepQuery = ''
  let stepResults: { category: string; template: string; help: string }[] = []

  onMount(async () => {
    try {
      version = await Version()
    } catch {
      version = 'dev'
    }
  })

  function basename(path: string): string {
    const parts = path.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || path
  }

  async function openProjectDialog() {
    const path = prompt('Путь к папке проекта:', projectPath || '')
    if (!path) return
    try {
      const info = await OpenProject(path)
      projectPath = info.path
      features = info.features || []
      tags = info.tags || []
      logText += `Проект: ${projectPath}\n`
      testClients = await ListTestClients().catch(() => [])
      statusText = `Файлов: ${features.length}`
    } catch (e: any) {
      logText += `Ошибка: ${e}\n`
    }
  }

  async function loadFeature(path: string) {
    selectedFeature = path
    try {
      editorText = await ReadFeature(path)
      statusText = basename(path)
    } catch (e: any) {
      logText += `Ошибка открытия: ${e}\n`
    }
  }

  async function saveFeature() {
    if (!selectedFeature) return
    try {
      await SaveFeature(selectedFeature, editorText)
      logText += `Сохранено: ${selectedFeature}\n`
    } catch (e: any) {
      logText += `Ошибка сохранения: ${e}\n`
    }
  }

  function parseVars(text: string): Record<string, string> {
    const out: Record<string, string> = {}
    for (const line of text.split('\n')) {
      const trimmed = line.trim()
      if (!trimmed || trimmed.startsWith('#')) continue
      const eq = trimmed.indexOf('=')
      if (eq <= 0) continue
      out[trimmed.slice(0, eq).trim()] = trimmed.slice(eq + 1).trim()
    }
    return out
  }

  async function runProject(dryRun: boolean) {
    if (!projectPath) return
    const result = await Run({
      tag: runTag,
      testClient: runTestClient,
      vars: parseVars(runVars),
      dryRun,
      headed: !dryRun,
      engine: dryRun ? '' : 'playwright',
      installPlaywright: !dryRun,
    })
    if (result.output) logText += result.output
    if (result.error) logText += `Ошибка: ${result.error}\n`
  }

  async function validateProject() {
    if (!projectPath) return
    const result = await Validate('chromium', true)
    if (result.output) logText += result.output
    if (result.error) logText += `Ошибка: ${result.error}\n`
  }

  async function initProject() {
    const out = await InitProject()
    if (out) logText += out
    if (projectPath) {
      const info = await OpenProject(projectPath)
      features = info.features || []
    }
  }

  async function searchSteps() {
    stepResults = await SearchSteps(stepQuery)
  }
</script>

<div class="app">
  <header class="toolbar">
    <span class="title">Scenaria {version}</span>
    <button on:click={openProjectDialog}>Открыть…</button>
    <button on:click={initProject} disabled={!projectPath}>Init</button>
    <button on:click={() => runProject(true)} disabled={!projectPath}>Dry-run</button>
    <button on:click={() => runProject(false)} disabled={!projectPath}>Playwright</button>
    <button on:click={validateProject} disabled={!projectPath}>Проверить</button>
    <button on:click={saveFeature} disabled={!selectedFeature}>Сохранить</button>
  </header>

  <aside class="sidebar">
    <p style="font-size:12px;color:#9aa0a6">{projectPath || 'Проект не открыт'}</p>
    {#each features as feature}
      <button
        class="feature-item"
        class:active={feature === selectedFeature}
        on:click={() => loadFeature(feature)}
      >
        {basename(feature)}
      </button>
    {/each}
  </aside>

  <main class="editor">
    <textarea bind:value={editorText} placeholder="Откройте .feature файл…" spellcheck="false"></textarea>
    <div class="status">{statusText}</div>
  </main>

  <aside class="panel">
    <h3>Запуск</h3>
    <div class="run-form">
      <label>Тег <input bind:value={runTag} placeholder="@smoke" /></label>
      <label>TestClient
        <select bind:value={runTestClient}>
          <option value="">(из feature)</option>
          {#each testClients as client}
            <option value={client}>{client}</option>
          {/each}
        </select>
      </label>
      <label>Переменные<textarea bind:value={runVars} rows="3" placeholder="KEY=value"></textarea></label>
      {#if tags.length}
        <small>Теги: {tags.join(', ')}</small>
      {/if}
    </div>
    <h3>Шаги</h3>
    <div class="run-form">
      <input bind:value={stepQuery} placeholder="Поиск шага" on:input={searchSteps} />
      {#each stepResults.slice(0, 8) as step}
        <button
          class="feature-item"
          on:click={() => {
            editorText += (editorText && !editorText.endsWith('\n') ? '\n' : '') + step.template
          }}
        >
          {step.template}
        </button>
      {/each}
    </div>
    <h3>Лог</h3>
    <div class="log">{logText}</div>
  </aside>
</div>
