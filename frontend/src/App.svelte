<script lang="ts">
  import { onDestroy, onMount } from 'svelte'
  import MonacoEditor from './lib/MonacoEditor.svelte'
  import { EventsOn } from '../wailsjs/runtime/runtime'
  import {
    Version,
    OpenProject,
    ReadFeature,
    SaveFeature,
    Run,
    Validate,
    ValidateFeature,
    SearchSteps,
    ListTestClients,
    InitProject,
    PickProjectFolder,
    PickSaveFile,
    PickOpenFile,
    Export,
    ImportJSON,
    RunVanessa,
    StartRecord,
    PauseRecording,
    ResumeRecording,
    CancelRecording,
    LoadSettings,
    SaveSettings,
    SubmitOTPCode,
    CancelOTP,
  } from '../wailsjs/go/wailsapp/App'
  import { gui } from '../wailsjs/go/models'

  type EditorTab = { path: string; content: string; dirty: boolean }

  let version = ''
  let projectPath = ''
  let features: string[] = []
  let tags: string[] = []
  let tabs: EditorTab[] = []
  let activeTab = ''
  let editorText = ''
  let logText = ''
  let statusText = 'Готово'
  let runTag = ''
  let runTestClient = ''
  let runVars = ''
  let testClients: string[] = []
  let stepQuery = ''
  let stepResults: { category: string; template: string; help: string }[] = []
  let monaco: MonacoEditor | undefined

  let showSettings = false
  let showRecord = false
  let showOtp = false
  let recording = false
  let recordPaused = false

  let settingsBrowser = 'chromium'
  let settingsHeadless = false
  let settingsWorkers = 1
  let settingsLoops = 100

  let recordURL = 'https://example.com'
  let recordOutput = 'recorded.feature'
  let recordIdle = 30

  let otpEmail = ''
  let otpCode = ''

  const unsubscribers: (() => void)[] = []

  onMount(async () => {
    try {
      version = await Version()
    } catch {
      version = 'dev'
    }

    unsubscribers.push(
      EventsOn('otp-prompt', (email: string) => {
        otpEmail = email || ''
        otpCode = ''
        showOtp = true
      }),
    )
    unsubscribers.push(
      EventsOn('record-started', () => {
        recording = true
        recordPaused = false
        logText += 'Запись начата…\n'
      }),
    )
    unsubscribers.push(
      EventsOn('record-finished', async (result: gui.RunResult) => {
        recording = false
        recordPaused = false
        showRecord = false
        if (result.output) logText += result.output
        if (result.error) logText += `Ошибка записи: ${result.error}\n`
        await refreshProject()
      }),
    )

    const onKey = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key === 's') {
        e.preventDefault()
        saveFeature()
      }
    }
    window.addEventListener('keydown', onKey)
    unsubscribers.push(() => window.removeEventListener('keydown', onKey))
  })

  onDestroy(() => {
    for (const off of unsubscribers) off()
  })

  function basename(path: string): string {
    const parts = path.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || path
  }

  async function refreshProject() {
    if (!projectPath) return
    const info = await OpenProject(projectPath)
    features = info.features || []
    tags = info.tags || []
    testClients = await ListTestClients().catch(() => [])
    statusText = `Файлов: ${features.length}`
  }

  async function openProjectDialog() {
    let path = ''
    try {
      path = await PickProjectFolder()
    } catch {
      path = ''
    }
    if (!path) {
      path = prompt('Путь к папке проекта:', projectPath || '') || ''
    }
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

  function syncActiveTabContent() {
    tabs = tabs.map((t) => (t.path === activeTab ? { ...t, content: editorText, dirty: true } : t))
  }

  async function loadFeature(path: string) {
    syncActiveTabContent()
    const existing = tabs.find((t) => t.path === path)
    if (existing) {
      activeTab = path
      editorText = existing.content
      statusText = basename(path)
      await validateEditor()
      return
    }
    try {
      const content = await ReadFeature(path)
      tabs = [...tabs, { path, content, dirty: false }]
      activeTab = path
      editorText = content
      statusText = basename(path)
      await validateEditor()
    } catch (e: any) {
      logText += `Ошибка открытия: ${e}\n`
    }
  }

  function selectTab(path: string) {
    if (path === activeTab) return
    loadFeature(path)
  }

  function closeTab(path: string, event?: Event) {
    event?.stopPropagation()
    tabs = tabs.filter((t) => t.path !== path)
    if (activeTab === path) {
      const next = tabs[tabs.length - 1]
      if (next) {
        activeTab = next.path
        editorText = next.content
        statusText = basename(next.path)
      } else {
        activeTab = ''
        editorText = ''
        statusText = 'Готово'
      }
    }
  }

  async function saveFeature() {
    if (!activeTab) return
    syncActiveTabContent()
    try {
      await SaveFeature(activeTab, editorText)
      tabs = tabs.map((t) => (t.path === activeTab ? { ...t, content: editorText, dirty: false } : t))
      logText += `Сохранено: ${activeTab}\n`
      statusText = `${basename(activeTab)} — сохранено`
    } catch (e: any) {
      logText += `Ошибка сохранения: ${e}\n`
    }
  }

  async function validateEditor() {
    try {
      const issues = await ValidateFeature(editorText)
      monaco?.setMarkers(issues || [])
      statusText =
        issues && issues.length > 0
          ? `Ошибки шагов: ${issues.length}`
          : activeTab
            ? `${basename(activeTab)} — OK`
            : 'Готово'
    } catch {
      /* dev without wails */
    }
  }

  async function onEditorChange(text: string) {
    editorText = text
    syncActiveTabContent()
    await validateEditor()
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

  async function validateProject(browser: boolean) {
    if (!projectPath) return
    const result = await Validate('chromium', !browser)
    if (result.output) logText += result.output
    if (result.error) logText += `Ошибка: ${result.error}\n`
  }

  async function initProject() {
    const out = await InitProject()
    if (out) logText += out
    await refreshProject()
  }

  async function searchSteps() {
    stepResults = await SearchSteps(stepQuery)
  }

  function insertStep(template: string) {
    const line = template.endsWith('\n') ? template : template + '\n'
    monaco?.insertAtCursor(line)
  }

  async function openSettings() {
    const s = await LoadSettings()
    settingsBrowser = s.browser || 'chromium'
    settingsHeadless = s.headless
    settingsWorkers = s.parallelWorkers || 1
    settingsLoops = s.maxLoopIterations || 100
    showSettings = true
  }

  async function applySettings() {
    await SaveSettings({
      browser: settingsBrowser,
      headless: settingsHeadless,
      parallelWorkers: settingsWorkers,
      maxLoopIterations: settingsLoops,
    })
    showSettings = false
    logText += 'Настройки сохранены.\n'
  }

  async function exportFeature() {
    if (!activeTab) return
    const out = await PickSaveFile('Экспорт', basename(activeTab).replace('.feature', '.json'))
    if (!out) return
    const result = await Export({
      inputPath: activeTab,
      output: out,
      format: out.endsWith('.json') ? 'json' : 'feature',
      baseURL: '',
    })
    if (result.output) logText += result.output
    if (result.error) logText += `Ошибка: ${result.error}\n`
  }

  async function importJSON() {
    const jsonPath = await PickOpenFile('Импорт JSON')
    if (!jsonPath) return
    const out = await PickSaveFile('Сохранить feature', 'imported.feature')
    if (!out) return
    const result = await ImportJSON({ jsonPath, outputPath: out })
    if (result.output) logText += result.output
    if (result.error) logText += `Ошибка: ${result.error}\n`
    await refreshProject()
    await loadFeature(out)
  }

  async function runVanessa(dry: boolean) {
    const result = await RunVanessa(dry)
    if (result.output) logText += result.output
    if (result.error) logText += `Ошибка: ${result.error}\n`
  }

  function beginRecord() {
    showRecord = true
  }

  async function startRecord() {
    await StartRecord({
      url: recordURL,
      output: recordOutput,
      idleSeconds: recordIdle,
      headless: false,
    })
  }

  async function toggleRecordPause() {
    if (recordPaused) {
      await ResumeRecording()
      recordPaused = false
    } else {
      await PauseRecording()
      recordPaused = true
    }
  }

  async function stopRecord() {
    await CancelRecording()
  }

  async function submitOtp() {
    await SubmitOTPCode(otpCode)
    showOtp = false
  }

  async function cancelOtp() {
    await CancelOTP()
    showOtp = false
  }
</script>

<div class="app">
  <header class="toolbar">
    <span class="title">Scenaria {version}</span>
    <button on:click={openProjectDialog}>Открыть…</button>
    <button on:click={initProject} disabled={!projectPath}>Init</button>
    <button on:click={() => runProject(true)} disabled={!projectPath}>Dry-run</button>
    <button on:click={() => runProject(false)} disabled={!projectPath}>Playwright</button>
    <button on:click={() => validateProject(false)} disabled={!projectPath}>Проверить</button>
    <button on:click={() => validateProject(true)} disabled={!projectPath}>В браузере</button>
    <button on:click={beginRecord} disabled={!projectPath || recording}>Запись…</button>
    <button on:click={() => runVanessa(true)} disabled={!projectPath}>VA dry</button>
    <button on:click={() => runVanessa(false)} disabled={!projectPath}>VA run</button>
    <button on:click={exportFeature} disabled={!activeTab}>Экспорт</button>
    <button on:click={importJSON} disabled={!projectPath}>Импорт</button>
    <button on:click={openSettings}>Настройки</button>
    <button on:click={saveFeature} disabled={!activeTab}>Сохранить</button>
  </header>

  <aside class="sidebar">
    <p class="hint">{projectPath || 'Проект не открыт'}</p>
    {#each features as feature}
      <button
        class="feature-item"
        class:active={feature === activeTab}
        on:click={() => loadFeature(feature)}
      >
        {basename(feature)}
      </button>
    {/each}
  </aside>

  <main class="editor">
    {#if tabs.length > 0}
      <div class="tabbar">
        {#each tabs as tab}
          <button class="tab" class:active={tab.path === activeTab} on:click={() => selectTab(tab.path)}>
            {basename(tab.path)}{tab.dirty ? ' *' : ''}
            <button type="button" class="close" on:click={(e) => closeTab(tab.path, e)} aria-label="Закрыть">×</button>
          </button>
        {/each}
      </div>
    {/if}
    <MonacoEditor bind:this={monaco} bind:value={editorText} on:change={(e) => onEditorChange(e.detail)} />
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
      {#if tags.length}<small>Теги: {tags.join(', ')}</small>{/if}
    </div>
    <h3>Шаги</h3>
    <div class="run-form">
      <input bind:value={stepQuery} placeholder="Поиск шага" on:input={searchSteps} />
      {#each stepResults.slice(0, 8) as step}
        <button class="feature-item" on:click={() => insertStep(step.template)}>{step.template}</button>
      {/each}
    </div>
    <h3>Лог</h3>
    <div class="log">{logText}</div>
  </aside>
</div>

{#if showSettings}
  <div class="modal-backdrop">
    <div class="modal">
      <h3>Настройки</h3>
      <label>Браузер
        <select bind:value={settingsBrowser}>
          <option value="chromium">chromium</option>
          <option value="firefox">firefox</option>
          <option value="webkit">webkit</option>
        </select>
      </label>
      <label><input type="checkbox" bind:checked={settingsHeadless} /> Headless</label>
      <label>Воркеры <input type="number" bind:value={settingsWorkers} min="1" /></label>
      <label>Лимит циклов <input type="number" bind:value={settingsLoops} min="1" /></label>
      <div class="modal-actions">
        <button on:click={applySettings}>Сохранить</button>
        <button on:click={() => (showSettings = false)}>Отмена</button>
      </div>
    </div>
  </div>
{/if}

{#if showRecord}
  <div class="modal-backdrop">
    <div class="modal">
      <h3>Live-запись</h3>
      <label>URL <input bind:value={recordURL} /></label>
      <label>Файл <input bind:value={recordOutput} /></label>
      <label>Idle (сек) <input type="number" bind:value={recordIdle} min="5" /></label>
      <div class="modal-actions">
        {#if recording}
          <button on:click={toggleRecordPause}>{recordPaused ? 'Resume' : 'Pause'}</button>
          <button on:click={stopRecord}>Стоп</button>
        {:else}
          <button on:click={startRecord}>Начать</button>
        {/if}
        <button on:click={() => (showRecord = false)}>Закрыть</button>
      </div>
    </div>
  </div>
{/if}

{#if showOtp}
  <div class="modal-backdrop">
    <div class="modal">
      <h3>Код из почты</h3>
      {#if otpEmail}<p>{otpEmail}</p>{/if}
      <input bind:value={otpCode} placeholder="123456" />
      <div class="modal-actions">
        <button on:click={submitOtp}>OK</button>
        <button on:click={cancelOtp}>Отмена</button>
      </div>
    </div>
  </div>
{/if}
