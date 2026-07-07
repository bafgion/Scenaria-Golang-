<script lang="ts">
  import { onMount } from 'svelte'
  import { BrowserInstallStatus, StartInstallBrowserEngine, ListPlugins } from '../../wailsjs/go/wailsapp/App'
  import type { gui } from '../../wailsjs/go/models'
  import { startRunResultJob } from './asyncRunResult'
  import { BRAND_NAME } from './brand'
  import SettingCard from './SettingCard.svelte'
  import { DEFAULT_EDITOR_SETTINGS, type EditorSettings } from './editorOptions'
  import { createTranslator, locale, setLocale, type Locale } from './i18n'

  export let browser = 'chromium'
  export let headless = false
  export let workers = 1
  export let slowMo = 0
  export let loops = 100
  export let filterRecording = false
  export let navOnlyRecording = false
  export let hoverRecord = false
  export let scrollBeforeClick = false
  export let hoverRecordMinMs = 600
  export let pickerDuringRecording = false
  export let toolbarCompact = false
  export let stepsPanelVisible = true
  export let stepsPanelHeight = 160
  export let checkUpdatesOnStartup = true
  export let uiLocale: Locale = 'ru'
  export let selectorClickStrategies: string[] = ['text', 'contextual', 'aria', 'title', 'testid', 'id']
  export let selectorInputStrategies: string[] = ['testid', 'id', 'label', 'placeholder', 'aria', 'name']
  export let navWaitUntil = 'domcontentloaded'
  export let htmlReportOpenMode: 'full' | 'light' = 'full'
  export let projectOpen = false
  export let editorSettings: EditorSettings = { ...DEFAULT_EDITOR_SETTINGS }

  export let onSave: () => void
  export let onApply: (() => void | Promise<void | string>) | null = null
  export let onCancel: () => void
  export let onOpenPlugins: (() => void) | null = null
  export let onOpenVanessa: (() => void) | null = null
  export let onInstallLog: ((line: string) => void) | null = null

  type TabId = 'record' | 'selectors' | 'plugins' | 'ui' | 'editor'

  let tab: TabId = 'record'
  let search = ''
  let plugins: gui.PluginEntryDTO[] = []
  let browserStatus: gui.BrowserInstallStatusDTO | null = null
  let browserInstallBusy = false
  let browserInstallProgress = ''
  let applyBusy = false
  let applyNotice = ''
  let applyNoticeTimer: ReturnType<typeof setTimeout> | null = null

  $: tr = createTranslator($locale)

  $: tabs = [
    { id: 'record' as TabId, label: tr('settings.tabs.record') },
    { id: 'selectors' as TabId, label: tr('settings.tabs.selectors') },
    { id: 'plugins' as TabId, label: tr('settings.tabs.plugins') },
    { id: 'editor' as TabId, label: tr('settings.tabs.editor') },
    { id: 'ui' as TabId, label: tr('settings.tabs.ui') },
  ]

  function onUiLocaleChange() {
    setLocale(uiLocale)
  }

  async function handleApply() {
    if (!onApply || applyBusy) return
    applyBusy = true
    applyNotice = ''
    try {
      const result = await onApply()
      applyNotice = typeof result === 'string' && result.trim() ? result.trim() : tr('common.settingsApplied')
      if (applyNoticeTimer) clearTimeout(applyNoticeTimer)
      applyNoticeTimer = setTimeout(() => {
        applyNotice = ''
        applyNoticeTimer = null
      }, 5000)
    } catch (err) {
      applyNotice = `${tr('common.error')}: ${err instanceof Error ? err.message : String(err)}`
    } finally {
      applyBusy = false
    }
  }

  $: slowMoPresets = [
    [0, tr('settings.slowMo.fast')],
    [100, tr('settings.slowMo.normal')],
    [250, tr('settings.slowMo.slow')],
    [500, tr('settings.slowMo.tutorial')],
  ] as [number, string][]
  const defaultClickStrategies = ['text', 'contextual', 'aria', 'title', 'testid', 'id']
  const defaultInputStrategies = ['label', 'placeholder', 'aria', 'name', 'testid', 'id']

  const strategyTechnicalLabels: Record<string, string> = {
    testid: 'data-testid',
    id: 'ID (#)',
    aria: 'aria-label',
    label: 'label:has-text',
    placeholder: 'placeholder',
    name: 'name',
  }

  function labelForStrategy(key: string): string {
    if (key === 'contextual') return tr('settings.strategies.contextual')
    if (key === 'text') return tr('settings.strategies.text')
    return strategyTechnicalLabels[key] || key
  }

  function moveStrategy(list: string[], index: number, delta: number): string[] {
    const next = index + delta
    if (next < 0 || next >= list.length) return list
    const copy = [...list]
    const [item] = copy.splice(index, 1)
    copy.splice(next, 0, item)
    return copy
  }

  function moveClickStrategy(index: number, delta: number) {
    selectorClickStrategies = moveStrategy(selectorClickStrategies, index, delta)
  }

  function moveInputStrategy(index: number, delta: number) {
    selectorInputStrategies = moveStrategy(selectorInputStrategies, index, delta)
  }

  onMount(async () => {
    try {
      plugins = await ListPlugins()
    } catch {
      plugins = []
    }
    await refreshBrowserEngineStatus()
  })

  $: vanessaEntry = plugins.find((p) => p.vanessa || p.name.toLowerCase() === 'vanessa')

  async function refreshBrowserEngineStatus() {
    if (browserInstallBusy) return
    try {
      browserStatus = await BrowserInstallStatus(browser)
    } catch {
      browserStatus = null
    }
  }

  async function onBrowserChange() {
    browserInstallProgress = ''
    await refreshBrowserEngineStatus()
  }

  async function installBrowserEngine() {
    if (browserInstallBusy) return
    browserInstallBusy = true
    browserInstallProgress = tr('settings.browser.installing', { engine: browserStatus?.label || browser })
    try {
      const result = await startRunResultJob('browser-install-finished', () => StartInstallBrowserEngine(browser))
      if (result.output) {
        browserInstallProgress = result.output.trim()
        onInstallLog?.(result.output.trim())
      }
      if (result.error) {
        browserInstallProgress = result.error
        onInstallLog?.(`${tr('settings.browser.error')}: ${result.error}`)
      }
    } catch (e: any) {
      browserInstallProgress = String(e)
      onInstallLog?.(`${tr('settings.browser.error')}: ${e}`)
    } finally {
      browserInstallBusy = false
      await refreshBrowserEngineStatus()
    }
  }

  function pickTab(id: TabId) {
    tab = id
  }

  function onSearch() {
    const q = search.trim().toLowerCase()
    if (!q) return
    const compact = q.replace(/\s+/g, '')
    if (/запис|record|браузер|browser|headless|chromium|playwright/.test(compact)) tab = 'record'
    else if (/селектор|selector|testid|css|стратег|strateg/.test(compact)) tab = 'selectors'
    else if (/плагин|plugin|vanessa|runner/.test(compact)) tab = 'plugins'
    else if (/редактор|editor|monaco|шрифт|font|миникарт|minimap|перенос|wrap|tab|fold|sticky|подсказк|hint|сценари|scenario/.test(compact)) tab = 'editor'
    else if (/отч[её]т|report|html/.test(compact)) tab = 'record'
    else if (/интерфейс|interface|ui|панел|panel|toolbar|шаг|step|обновлен|update/.test(compact)) tab = 'ui'
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel()
  }

  function resetToDefaults() {
    browser = 'chromium'
    headless = false
    workers = 1
    slowMo = 0
    loops = 100
    navWaitUntil = 'domcontentloaded'
    filterRecording = false
    navOnlyRecording = false
    hoverRecord = false
    scrollBeforeClick = false
    hoverRecordMinMs = 600
    pickerDuringRecording = false
    toolbarCompact = false
    uiLocale = 'ru'
    setLocale('ru')
    stepsPanelVisible = true
    stepsPanelHeight = 160
    checkUpdatesOnStartup = true
    selectorClickStrategies = [...defaultClickStrategies]
    selectorInputStrategies = [...defaultInputStrategies]
    editorSettings = { ...DEFAULT_EDITOR_SETTINGS }
  }

  $: extremeRunWarning =
    workers >= 8 && slowMo >= 1000
      ? tr('settings.warnings.manyWorkers')
      : workers >= 12
        ? tr('settings.warnings.highWorkers')
        : slowMo >= 2000
          ? tr('settings.warnings.highSlowMo')
          : ''
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onCancel}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="app-dialog" role="dialog" aria-modal="true" aria-label="{tr('settings.title')} — {BRAND_NAME}" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <header class="dialog-search-bar">
      <input bind:value={search} placeholder={tr('settings.searchPlaceholder')} on:input={onSearch} />
    </header>

    <div class="dialog-body">
      <nav class="dialog-sidebar" aria-label={tr('settings.sectionsNav')}>
        {#each tabs as t}
          <button type="button" class:active={tab === t.id} on:click={() => pickTab(t.id)}>{t.label}</button>
        {/each}
      </nav>

      <div class="dialog-content">
        {#if tab === 'record'}
          <section class="setting-section">
            <h4 class="setting-section-title">{tr('settings.sections.browser.title')}</h4>
            <p class="setting-section-desc">{tr('settings.sections.browser.desc')}</p>

            <SettingCard title={tr('settings.cards.headless.title')} description={tr('settings.cards.headless.description')}>
              <input type="checkbox" bind:checked={headless} />
            </SettingCard>

            <SettingCard title={tr('settings.cards.browserEngine.title')} description={tr('settings.cards.browserEngine.description')}>
              <select bind:value={browser} on:change={onBrowserChange} disabled={browserInstallBusy}>
                <option value="chromium">Chromium</option>
                <option value="firefox">Firefox</option>
                <option value="webkit">WebKit</option>
              </select>
            </SettingCard>

            <div class="browser-install-row">
              <p
                class="browser-engine-status"
                class:ok={browserStatus?.installed}
                class:warn={browserStatus && !browserStatus.installed}
              >
                {#if browserStatus}
                  {#if browserStatus.installed}
                    {tr('settings.browser.installed', { label: browserStatus.label })}
                    <span class="detail">{browserStatus.detail}</span>
                  {:else}
                    {tr('settings.browser.notInstalled', { label: browserStatus.label })}
                  {/if}
                {:else}
                  {tr('settings.browser.checking')}
                {/if}
              </p>
              <button
                type="button"
                class="install-btn"
                on:click={installBrowserEngine}
                disabled={browserInstallBusy}
              >
                {browserInstallBusy
                  ? tr('settings.browser.installingBtn')
                  : browserStatus?.installed
                    ? tr('settings.browser.reinstall')
                    : tr('settings.browser.installEngine')}
              </button>
            </div>
            {#if browserInstallProgress}
              <p class="browser-install-progress">{browserInstallProgress}</p>
            {/if}
          </section>

          <section class="setting-section">
            <h4 class="setting-section-title">{tr('settings.sections.recording.title')}</h4>
            <p class="setting-section-desc">{tr('settings.sections.recording.desc')}</p>

            <SettingCard title={tr('settings.cards.importantOnly.title')} description={tr('settings.cards.importantOnly.description')}>
              <input
                type="checkbox"
                bind:checked={filterRecording}
                on:change={() => filterRecording && (navOnlyRecording = false)}
              />
            </SettingCard>

            <SettingCard title={tr('settings.cards.linksOnly.title')} description={tr('settings.cards.linksOnly.description')}>
              <input
                type="checkbox"
                bind:checked={navOnlyRecording}
                on:change={() => navOnlyRecording && (filterRecording = false)}
              />
            </SettingCard>

            <SettingCard title={tr('settings.cards.hoverRecord.title')} description={tr('settings.cards.hoverRecord.description')}>
              <input type="checkbox" bind:checked={hoverRecord} />
            </SettingCard>

            <SettingCard title={tr('settings.cards.hoverMin.title')} description={tr('settings.cards.hoverMin.description')}>
              <span class="num-with-unit">
                <input type="number" class="setting-number" bind:value={hoverRecordMinMs} min={100} max={5000} step={50} />
                <span>{tr('settings.units.ms')}</span>
              </span>
            </SettingCard>

            <SettingCard title={tr('settings.cards.scrollBeforeClick.title')} description={tr('settings.cards.scrollBeforeClick.description')}>
              <input type="checkbox" bind:checked={scrollBeforeClick} />
            </SettingCard>

            <SettingCard title={tr('settings.cards.pickerDuringRecording.title')} description={tr('settings.cards.pickerDuringRecording.description')}>
              <input type="checkbox" bind:checked={pickerDuringRecording} />
            </SettingCard>
          </section>

          <section class="setting-section">
            <h4 class="setting-section-title">{tr('settings.sections.run.title')}</h4>
            <SettingCard title={tr('settings.cards.workers.title')} description={tr('settings.cards.workers.description')}>
              <span class="num-with-unit">
                <input type="number" class="setting-number" bind:value={workers} min={1} max={16} />
                <span>{tr('settings.units.pcs')}</span>
              </span>
            </SettingCard>
            <SettingCard title={tr('settings.cards.slowMo.title')} description={tr('settings.cards.slowMo.description')}>
              <span class="num-with-unit">
                <input type="number" class="setting-number" bind:value={slowMo} min={0} max={5000} step={50} />
                <span>{tr('settings.units.ms')}</span>
              </span>
              <div class="slowmo-presets">
                {#each slowMoPresets as [ms, label]}
                  <button type="button" class="preset-chip" class:active={slowMo === ms} on:click={() => (slowMo = ms)}>{label}</button>
                {/each}
              </div>
            </SettingCard>
            <SettingCard title={tr('settings.cards.navWait.title')} description={tr('settings.cards.navWait.description')}>
              <select bind:value={navWaitUntil}>
                <option value="load">{tr('settings.navWait.load')}</option>
                <option value="domcontentloaded">{tr('settings.navWait.domcontentloaded')}</option>
                <option value="networkidle">{tr('settings.navWait.networkidle')}</option>
                <option value="commit">{tr('settings.navWait.commit')}</option>
              </select>
            </SettingCard>
            <SettingCard title={tr('settings.cards.loops.title')} description={tr('settings.cards.loops.description')}>
              <input type="number" class="setting-number" bind:value={loops} min={1} max={10000} />
            </SettingCard>
            {#if extremeRunWarning}
              <p class="setting-run-warning">{extremeRunWarning}</p>
            {/if}
          </section>

          {#if projectOpen}
            <section class="setting-section">
              <h4 class="setting-section-title">{tr('settings.sections.reports.title')}</h4>
              <p class="setting-section-desc">{tr('settings.sections.reports.desc')}</p>
              <SettingCard title={tr('settings.cards.htmlReportOpen.title')} description={tr('settings.cards.htmlReportOpen.description')}>
                <select bind:value={htmlReportOpenMode}>
                  <option value="full">{tr('settings.cards.htmlReportOpen.full')}</option>
                  <option value="light">{tr('settings.cards.htmlReportOpen.light')}</option>
                </select>
              </SettingCard>
            </section>
          {/if}
        {:else if tab === 'selectors'}
          <section class="setting-section">
            <h4 class="setting-section-title">{tr('settings.sections.selectors.title')}</h4>
            <p class="setting-section-desc">
              {tr('settings.sections.selectors.desc', { brand: BRAND_NAME })}
            </p>
            <h5 class="strategy-group-title">{tr('settings.strategies.clicksTitle')}</h5>
            <ul class="selector-list editable">
              {#each selectorClickStrategies as key, i}
                <li>
                  <span class="strategy-name">{labelForStrategy(key)}</span>
                  <span class="strategy-actions">
                    <button type="button" class="btn-compact" title={tr('settings.strategies.moveUp')} disabled={i === 0} on:click={() => moveClickStrategy(i, -1)}>↑</button>
                    <button type="button" class="btn-compact" title={tr('settings.strategies.moveDown')} disabled={i === selectorClickStrategies.length - 1} on:click={() => moveClickStrategy(i, 1)}>↓</button>
                  </span>
                </li>
              {/each}
            </ul>
            <button type="button" class="dialog-link-btn" on:click={() => (selectorClickStrategies = [...defaultClickStrategies])}>{tr('settings.strategies.resetClicks')}</button>
            <h5 class="strategy-group-title">{tr('settings.strategies.inputsTitle')}</h5>
            <ul class="selector-list editable">
              {#each selectorInputStrategies as key, i}
                <li>
                  <span class="strategy-name">{labelForStrategy(key)}</span>
                  <span class="strategy-actions">
                    <button type="button" class="btn-compact" title={tr('settings.strategies.moveUp')} disabled={i === 0} on:click={() => moveInputStrategy(i, -1)}>↑</button>
                    <button type="button" class="btn-compact" title={tr('settings.strategies.moveDown')} disabled={i === selectorInputStrategies.length - 1} on:click={() => moveInputStrategy(i, 1)}>↓</button>
                  </span>
                </li>
              {/each}
            </ul>
            <button type="button" class="dialog-link-btn" on:click={() => (selectorInputStrategies = [...defaultInputStrategies])}>{tr('settings.strategies.resetInputs')}</button>
          </section>
        {:else if tab === 'plugins'}
          <section class="setting-section">
            <h4 class="setting-section-title">{tr('settings.sections.plugins.title')}</h4>
            <p class="setting-section-desc">
              {tr('settings.sections.plugins.desc')}
            </p>
            <div class="plugin-list">
              <div class="plugin-row">
                <span>{tr('settings.plugins.playwright')}</span>
                <span class="status ok">{tr('settings.plugins.builtIn')}</span>
              </div>
              <div class="plugin-row">
                <span>{tr('settings.plugins.vanessa')}</span>
                <span class="status" class:ok={vanessaEntry} class:warn={!vanessaEntry}>
                  {vanessaEntry ? tr('settings.plugins.installed') : tr('settings.plugins.unavailable')}
                </span>
              </div>
            </div>
            {#if onOpenVanessa}
              <button type="button" class="dialog-link-btn" on:click={onOpenVanessa}>{tr('settings.vanessaSettings')}</button>
            {/if}
            {#if onOpenPlugins}
              <button type="button" class="dialog-link-btn" on:click={onOpenPlugins}>{tr('settings.managePlugins')}</button>
            {/if}
            <p class="hint">{tr('settings.plugins.zipHint')}</p>
          </section>
        {:else if tab === 'editor'}
          <section class="setting-section">
            <h4 class="setting-section-title">{tr('settings.sections.editorFont.title')}</h4>
            <p class="setting-section-desc">{tr('settings.sections.editorFont.desc')}</p>
            <SettingCard title={tr('settings.cards.fontSize.title')} description={tr('settings.cards.fontSize.description')}>
              <input type="number" class="setting-number" bind:value={editorSettings.fontSize} min={8} max={32} />
            </SettingCard>
            <SettingCard title={tr('settings.cards.fontFamily.title')} description={tr('settings.cards.fontFamily.description')}>
              <input type="text" class="setting-text setting-text-mono" bind:value={editorSettings.fontFamily} />
            </SettingCard>
            <SettingCard title={tr('settings.cards.theme.title')} description={tr('settings.cards.theme.description')}>
              <select bind:value={editorSettings.theme}>
                <option value="scenaria-dark">{tr('settings.editor.themeDark')}</option>
                <option value="scenaria-light">{tr('settings.editor.themeLight')}</option>
                <option value="system">{tr('settings.editor.themeSystem')}</option>
              </select>
            </SettingCard>
            <SettingCard title={tr('settings.cards.wordWrap.title')} description={tr('settings.cards.wordWrap.description')}>
              <select bind:value={editorSettings.wordWrap}>
                <option value="on">{tr('settings.editor.wordWrapOn')}</option>
                <option value="off">{tr('settings.editor.wordWrapOff')}</option>
              </select>
            </SettingCard>
            <SettingCard title={tr('settings.cards.minimap.title')} description={tr('settings.cards.minimap.description')}>
              <input type="checkbox" bind:checked={editorSettings.minimap} />
            </SettingCard>
            <SettingCard title={tr('settings.cards.lineNumbers.title')} description={tr('settings.cards.lineNumbers.description')}>
              <select bind:value={editorSettings.lineNumbers}>
                <option value="on">{tr('settings.editor.lineNumbersOn')}</option>
                <option value="relative">{tr('settings.editor.lineNumbersRelative')}</option>
                <option value="off">{tr('settings.editor.lineNumbersOff')}</option>
              </select>
            </SettingCard>
            <SettingCard title={tr('settings.cards.renderWhitespace.title')} description={tr('settings.cards.renderWhitespace.description')}>
              <select bind:value={editorSettings.renderWhitespace}>
                <option value="none">{tr('settings.editor.whitespaceNone')}</option>
                <option value="boundary">{tr('settings.editor.whitespaceBoundary')}</option>
                <option value="selection">{tr('settings.editor.whitespaceSelection')}</option>
                <option value="trailing">{tr('settings.editor.whitespaceTrailing')}</option>
                <option value="all">{tr('settings.editor.whitespaceAll')}</option>
              </select>
            </SettingCard>
          </section>

          <section class="setting-section">
            <h4 class="setting-section-title">{tr('settings.sections.editorInput.title')}</h4>
            <SettingCard title={tr('settings.cards.tabSize.title')} description={tr('settings.cards.tabSize.description')}>
              <input type="number" class="setting-number" bind:value={editorSettings.tabSize} min={1} max={8} />
            </SettingCard>
            <SettingCard title={tr('settings.cards.insertSpaces.title')} description={tr('settings.cards.insertSpaces.description')}>
              <input type="checkbox" bind:checked={editorSettings.insertSpaces} />
            </SettingCard>
            <SettingCard title={tr('settings.cards.folding.title')} description={tr('settings.cards.folding.description')}>
              <input type="checkbox" bind:checked={editorSettings.folding} />
            </SettingCard>
            <SettingCard title={tr('settings.cards.stickyScroll.title')} description={tr('settings.cards.stickyScroll.description')}>
              <input type="checkbox" bind:checked={editorSettings.stickyScroll} />
            </SettingCard>
            <SettingCard title={tr('settings.cards.autoClosingQuotes.title')} description={tr('settings.cards.autoClosingQuotes.description')}>
              <select bind:value={editorSettings.autoClosingQuotes}>
                <option value="languageDefined">{tr('settings.editor.autoQuotesLanguage')}</option>
                <option value="always">{tr('settings.editor.autoQuotesAlways')}</option>
                <option value="beforeWhitespace">{tr('settings.editor.autoQuotesBeforeWhitespace')}</option>
                <option value="never">{tr('settings.editor.autoQuotesNever')}</option>
              </select>
            </SettingCard>
            <SettingCard title={tr('settings.cards.formatOnSave.title')} description={tr('settings.cards.formatOnSave.description')}>
              <input type="checkbox" bind:checked={editorSettings.formatOnSave} />
            </SettingCard>
            <SettingCard title={tr('settings.cards.stepHover.title')} description={tr('settings.cards.stepHover.description')}>
              <input type="checkbox" bind:checked={editorSettings.stepHover} />
            </SettingCard>
            <SettingCard title={tr('settings.cards.validateOnType.title')} description={tr('settings.cards.validateOnType.description')}>
              <input type="checkbox" bind:checked={editorSettings.validateOnType} />
            </SettingCard>
          </section>

          <section class="setting-section">
            <h4 class="setting-section-title">{tr('settings.sections.editorNavigation.title')}</h4>
            <p class="setting-section-desc">{tr('settings.sections.editorNavigation.desc')}</p>
            <SettingCard title={tr('settings.cards.breadcrumbs.title')} description={tr('settings.cards.breadcrumbs.description')}>
              <input type="checkbox" bind:checked={editorSettings.breadcrumbs} />
            </SettingCard>
            <SettingCard title={tr('settings.cards.symbolOutline.title')} description={tr('settings.cards.symbolOutline.description')}>
              <input type="checkbox" bind:checked={editorSettings.symbolOutline} />
            </SettingCard>
            <SettingCard title={tr('settings.cards.stepsPanelView.title')} description={tr('settings.cards.stepsPanelView.description')}>
              <select bind:value={editorSettings.stepsPanelView}>
                <option value="outline">{tr('settings.editor.stepsPanelOutline')}</option>
                <option value="steps">{tr('settings.editor.stepsPanelSteps')}</option>
              </select>
            </SettingCard>
            <SettingCard title={tr('settings.cards.codeLens.title')} description={tr('settings.cards.codeLens.description')}>
              <input type="checkbox" bind:checked={editorSettings.codeLens} />
            </SettingCard>
            <SettingCard title={tr('settings.cards.inlayHints.title')} description={tr('settings.cards.inlayHints.description')}>
              <input type="checkbox" bind:checked={editorSettings.inlayHints} />
            </SettingCard>
          </section>

          <section class="setting-section">
            <h4 class="setting-section-title">{tr('settings.sections.editorHints.title')}</h4>
            <p class="setting-section-desc">{tr('settings.sections.editorHints.desc')}</p>
            <SettingCard title={tr('settings.cards.scenarioHints.title')} description={tr('settings.cards.scenarioHints.description')}>
              <input type="checkbox" bind:checked={editorSettings.scenarioHints} />
            </SettingCard>
            <SettingCard title={tr('settings.cards.scenarioHintsAfterRecord.title')} description={tr('settings.cards.scenarioHintsAfterRecord.description')}>
              <input type="checkbox" bind:checked={editorSettings.scenarioHintsAfterRecord} disabled={!editorSettings.scenarioHints} />
            </SettingCard>
            <SettingCard title={tr('settings.cards.scenarioHintsShowWarning.title')} description={tr('settings.cards.scenarioHintsShowWarning.description')}>
              <input type="checkbox" bind:checked={editorSettings.scenarioHintsShowWarning} disabled={!editorSettings.scenarioHints} />
            </SettingCard>
            <SettingCard title={tr('settings.cards.scenarioHintsShowInfo.title')} description={tr('settings.cards.scenarioHintsShowInfo.description')}>
              <input type="checkbox" bind:checked={editorSettings.scenarioHintsShowInfo} disabled={!editorSettings.scenarioHints} />
            </SettingCard>
            <SettingCard title={tr('settings.cards.scenarioHintsAutoFixOnSave.title')} description={tr('settings.cards.scenarioHintsAutoFixOnSave.description')}>
              <input type="checkbox" bind:checked={editorSettings.scenarioHintsAutoFixOnSave} disabled={!editorSettings.scenarioHints} />
            </SettingCard>
          </section>
        {:else}
          <section class="setting-section">
            <h4 class="setting-section-title">{tr('settings.sections.uiLanguage.title')}</h4>
            <p class="setting-section-desc">{tr('settings.sections.uiLanguage.desc')}</p>
            <SettingCard title={tr('settings.uiLocale')} description="">
              <select bind:value={uiLocale} on:change={onUiLocaleChange}>
                <option value="ru">{tr('settings.uiLocaleRu')}</option>
                <option value="en">{tr('settings.uiLocaleEn')}</option>
              </select>
            </SettingCard>
          </section>

          <section class="setting-section">
            <h4 class="setting-section-title">{tr('settings.sections.toolbar.title')}</h4>
            <p class="setting-section-desc">{tr('settings.sections.toolbar.desc')}</p>
            <SettingCard title={tr('settings.cards.toolbarCompact.title')} description={tr('settings.cards.toolbarCompact.description')}>
              <input type="checkbox" bind:checked={toolbarCompact} />
            </SettingCard>
          </section>

          <section class="setting-section">
            <h4 class="setting-section-title">{tr('settings.sections.stepsPanel.title')}</h4>
            <p class="setting-section-desc">{tr('settings.sections.stepsPanel.desc')}</p>
            <SettingCard title={tr('settings.cards.stepsPanelVisible.title')} description={tr('settings.cards.stepsPanelVisible.description')}>
              <input type="checkbox" bind:checked={stepsPanelVisible} />
            </SettingCard>
            <SettingCard title={tr('settings.cards.stepsPanelHeight.title')} description={tr('settings.cards.stepsPanelHeight.description')}>
              <span class="num-with-unit">
                <input type="number" class="setting-number" bind:value={stepsPanelHeight} min={80} max={480} />
                <span>{tr('settings.units.px')}</span>
              </span>
            </SettingCard>
          </section>

          <section class="setting-section">
            <h4 class="setting-section-title">{tr('settings.sections.updates.title')}</h4>
            <SettingCard title={tr('settings.cards.checkUpdates.title')} description={tr('settings.cards.checkUpdates.description', { brand: BRAND_NAME })}>
              <input type="checkbox" bind:checked={checkUpdatesOnStartup} />
            </SettingCard>
          </section>
        {/if}
      </div>
    </div>

    <footer class="dialog-footer">
      <button type="button" class="reset-btn" on:click={resetToDefaults}>{tr('settings.resetDefaults')}</button>
      {#if applyNotice}
        <span class="apply-notice" role="status">{applyNotice}</span>
      {/if}
      <span class="dialog-footer-spacer"></span>
      {#if onApply}
        <button type="button" class:applied={!!applyNotice && !applyBusy} disabled={applyBusy} on:click={handleApply}>
          {applyBusy ? tr('settings.applyBusy') : applyNotice ? tr('settings.applyDone') : tr('settings.apply')}
        </button>
      {/if}
      <button type="button" class="primary" on:click={onSave}>{tr('settings.ok')}</button>
      <button type="button" on:click={onCancel}>{tr('settings.cancel')}</button>
    </footer>
  </div>
</div>

<style>
  .browser-install-row {
    display: flex;
    gap: 10px;
    align-items: flex-start;
    margin: 0 0 8px;
    padding: 0 2px;
    min-height: 44px;
  }

  .browser-engine-status {
    flex: 1;
    margin: 0;
    font-size: 12px;
    line-height: 1.45;
    color: var(--color-muted);
  }

  .browser-engine-status.ok {
    color: var(--color-success, #4ec9b0);
  }

  .browser-engine-status.warn {
    color: var(--color-warning, #dcdcaa);
  }

  .browser-engine-status .detail {
    display: block;
    margin-top: 4px;
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--color-muted);
    word-break: break-all;
  }

  .browser-install-progress {
    margin: 0 0 12px;
    min-height: 1.2em;
    font-size: 11px;
    color: var(--color-muted);
    white-space: pre-wrap;
    font-family: var(--font-mono);
  }

  .slowmo-presets {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin-top: 8px;
  }

  .setting-run-warning {
    margin: 0 2px 12px;
    font-size: 12px;
    color: var(--color-warning, #dcdcaa);
  }

  .dialog-footer {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .dialog-footer-spacer {
    flex: 1;
  }

  .apply-notice {
    max-width: min(420px, 42vw);
    font-size: 12px;
    line-height: 1.35;
    color: var(--color-success, #4ec9b0);
  }

  .dialog-footer button.applied:not(:disabled) {
    color: var(--color-success, #4ec9b0);
    border-color: var(--color-success, #4ec9b0);
  }

  .reset-btn {
    font-size: 12px;
  }

  .strategy-group-title {
    margin: 12px 0 6px;
    font-size: 12px;
    font-weight: 600;
    color: var(--color-text);
  }

  .selector-list.editable {
    list-style: none;
    margin: 0 0 8px;
    padding: 0;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    overflow: hidden;
  }

  .selector-list.editable li {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    padding: 6px 8px;
    border-bottom: 1px solid var(--color-divider);
    font-size: 12px;
  }

  .selector-list.editable li:last-child {
    border-bottom: none;
  }

  .strategy-actions {
    display: flex;
    gap: 4px;
  }

  .strategy-actions :global(.btn-compact) {
    min-width: 24px;
  }
</style>
