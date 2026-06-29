<script lang="ts">
  import { onMount } from 'svelte'
  import { BrowserInstallStatus, InstallBrowserEngine, ListPlugins } from '../../wailsjs/go/wailsapp/App'
  import type { gui } from '../../wailsjs/go/models'
  import SettingCard from './SettingCard.svelte'
  import NumberInput from './NumberInput.svelte'

  export let browser = 'chromium'
  export let headless = false
  export let workers = 1
  export let loops = 100
  export let filterRecording = false
  export let navOnlyRecording = false
  export let hoverRecord = false
  export let toolbarCompact = false
  export let stepsPanelVisible = true
  export let stepsPanelHeight = 160
  export let checkUpdatesOnStartup = true

  export let onSave: () => void
  export let onCancel: () => void
  export let onOpenPlugins: (() => void) | null = null
  export let onOpenVanessa: (() => void) | null = null
  export let onInstallLog: ((line: string) => void) | null = null

  type TabId = 'record' | 'selectors' | 'plugins' | 'ui'

  let tab: TabId = 'record'
  let search = ''
  let plugins: gui.PluginEntryDTO[] = []
  let browserStatus: gui.BrowserInstallStatusDTO | null = null
  let browserInstallBusy = false
  let browserInstallProgress = ''

  const tabs: { id: TabId; label: string }[] = [
    { id: 'record', label: 'Запись и браузер' },
    { id: 'selectors', label: 'Селекторы' },
    { id: 'plugins', label: 'Плагины' },
    { id: 'ui', label: 'Интерфейс' },
  ]

  const selectorStrategies = [
    'data-testid',
    'id',
    'role / aria-label',
    'name',
    'text',
    'placeholder',
    'CSS (хрупкий)',
  ]

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
    browserInstallProgress = `Установка ${browserStatus?.label || browser}…`
    try {
      const result = await InstallBrowserEngine(browser)
      if (result.output) {
        browserInstallProgress = result.output.trim()
        onInstallLog?.(result.output.trim())
      }
      if (result.error) {
        browserInstallProgress = result.error
        onInstallLog?.(`Ошибка: ${result.error}`)
      }
    } catch (e: any) {
      browserInstallProgress = String(e)
      onInstallLog?.(`Ошибка: ${e}`)
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
    if (/запис|браузер|headless|chromium|playwright/.test(compact)) tab = 'record'
    else if (/селектор|testid|css|стратег/.test(compact)) tab = 'selectors'
    else if (/плагин|vanessa|runner/.test(compact)) tab = 'plugins'
    else if (/интерфейс|панел|toolbar|шаг|обновлен/.test(compact)) tab = 'ui'
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onCancel}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="app-dialog" role="dialog" aria-modal="true" aria-label="Настройки — Scenaria" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <header class="dialog-search-bar">
      <input bind:value={search} placeholder="Поиск настроек" on:input={onSearch} />
    </header>

    <div class="dialog-body">
      <nav class="dialog-sidebar" aria-label="Разделы настроек">
        {#each tabs as t}
          <button type="button" class:active={tab === t.id} on:click={() => pickTab(t.id)}>{t.label}</button>
        {/each}
      </nav>

      <div class="dialog-content">
        {#if tab === 'record'}
          <section class="setting-section">
            <h4 class="setting-section-title">Браузер</h4>
            <p class="setting-section-desc">Поведение окна и сессии при записи и прогоне Playwright.</p>

            <SettingCard title="Без окна браузера" description="Headless — окно не показывается при записи и запуске.">
              <input type="checkbox" bind:checked={headless} />
            </SettingCard>

            <SettingCard title="Движок браузера" description="Playwright: chromium, firefox или webkit.">
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
                    {browserStatus.label}: установлен
                    <span class="detail">{browserStatus.detail}</span>
                  {:else}
                    {browserStatus.label}: не установлен — нужна загрузка перед записью и прогоном.
                  {/if}
                {:else}
                  Проверка движка…
                {/if}
              </p>
              <button
                type="button"
                class="install-btn"
                on:click={installBrowserEngine}
                disabled={browserInstallBusy}
              >
                {browserInstallBusy ? 'Установка…' : browserStatus?.installed ? 'Переустановить' : 'Установить движок'}
              </button>
            </div>
            {#if browserInstallProgress}
              <p class="browser-install-progress">{browserInstallProgress}</p>
            {/if}
          </section>

          <section class="setting-section">
            <h4 class="setting-section-title">Запись шагов</h4>
            <p class="setting-section-desc">Фильтры при записи действий в браузере.</p>

            <SettingCard title="Только важные" description="Пропускать второстепенные события при записи.">
              <input
                type="checkbox"
                bind:checked={filterRecording}
                on:change={() => filterRecording && (navOnlyRecording = false)}
              />
            </SettingCard>

            <SettingCard title="Только ссылки" description="Записывать переходы по ссылкам, без кликов по элементам.">
              <input
                type="checkbox"
                bind:checked={navOnlyRecording}
                on:change={() => navOnlyRecording && (filterRecording = false)}
              />
            </SettingCard>

            <SettingCard title="Записывать наведение" description="Добавлять шаги при наведении курсора на элементы.">
              <input type="checkbox" bind:checked={hoverRecord} />
            </SettingCard>
          </section>

          <section class="setting-section">
            <h4 class="setting-section-title">Запуск</h4>
            <SettingCard title="Параллельные воркеры" description="Число одновременных браузерных сессий при пакетном запуске.">
              <span class="num-with-unit">
                <NumberInput bind:value={workers} min={1} max={16} width="56px" />
                <span>шт.</span>
              </span>
            </SettingCard>
            <SettingCard title="Лимит итераций циклов" description="Максимум повторов для блоков «Повторяю» / «Пока».">
              <NumberInput bind:value={loops} min={1} max={10000} width="72px" />
            </SettingCard>
          </section>
        {:else if tab === 'selectors'}
          <section class="setting-section">
            <h4 class="setting-section-title">Приоритет стратегий</h4>
            <p class="setting-section-desc">
              При записи и подборе селектора Scenaria перебирает стратегии сверху вниз. Более стабильные — выше.
            </p>
            <ul class="selector-list">
              {#each selectorStrategies as name}
                <li>{name}</li>
              {/each}
            </ul>
            <p class="hint">Изменение порядка в этой версии пока недоступно — используется встроенная эвристика.</p>
          </section>
        {:else if tab === 'plugins'}
          <section class="setting-section">
            <h4 class="setting-section-title">Runner'ы и add-on'ы</h4>
            <p class="setting-section-desc">
              Плагины устанавливаются в <code>addons/&lt;имя&gt;/</code> и регистрируются в <code>.scenaria/plugins.json</code>.
            </p>
            <div class="plugin-list">
              <div class="plugin-row">
                <span>Playwright</span>
                <span class="status ok">встроен</span>
              </div>
              <div class="plugin-row">
                <span>Vanessa Automation</span>
                <span class="status" class:ok={vanessaEntry} class:warn={!vanessaEntry}>
                  {vanessaEntry ? 'установлен' : 'недоступен'}
                </span>
              </div>
            </div>
            {#if onOpenVanessa}
              <button type="button" class="dialog-link-btn" on:click={onOpenVanessa}>Настройки Vanessa…</button>
            {/if}
            {#if onOpenPlugins}
              <button type="button" class="dialog-link-btn" on:click={onOpenPlugins}>Управление плагинами…</button>
            {/if}
            <p class="hint">ZIP-плагины: <code>scenaria plugins install …</code> или диалог «Управление плагинами».</p>
          </section>
        {:else}
          <section class="setting-section">
            <h4 class="setting-section-title">Панель инструментов</h4>
            <p class="setting-section-desc">Внешний вид верхней панели действий.</p>
            <SettingCard title="Компактная панель" description="Меньше подписей на кнопках — только иконки.">
              <input type="checkbox" bind:checked={toolbarCompact} />
            </SettingCard>
          </section>

          <section class="setting-section">
            <h4 class="setting-section-title">Панель шагов</h4>
            <p class="setting-section-desc">Список распознанных шагов под редактором сценария.</p>
            <SettingCard title="Показывать панель шагов" description="Отображать разбор шагов Gherkin под редактором.">
              <input type="checkbox" bind:checked={stepsPanelVisible} />
            </SettingCard>
            <SettingCard title="Высота панели" description="Высота области со списком шагов в пикселях.">
              <span class="num-with-unit">
                <NumberInput bind:value={stepsPanelHeight} min={80} max={480} width="56px" />
                <span>px</span>
              </span>
            </SettingCard>
          </section>

          <section class="setting-section">
            <h4 class="setting-section-title">Обновления</h4>
            <SettingCard title="Проверять при запуске" description="Искать новую версию Scenaria при старте IDE.">
              <input type="checkbox" bind:checked={checkUpdatesOnStartup} />
            </SettingCard>
          </section>
        {/if}
      </div>
    </div>

    <footer class="dialog-footer">
      <button type="button" class="primary" on:click={onSave}>OK</button>
      <button type="button" on:click={onCancel}>Отмена</button>
    </footer>
  </div>
</div>

<style>
  code {
    font-family: var(--font-mono);
    font-size: 11px;
  }

  .browser-install-row {
    display: flex;
    gap: 10px;
    align-items: flex-start;
    margin: 0 0 8px;
    padding: 0 2px;
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

  .install-btn {
    flex-shrink: 0;
    padding: 6px 10px;
    font-size: 12px;
  }

  .browser-install-progress {
    margin: 0 0 12px;
    font-size: 11px;
    color: var(--color-muted);
    white-space: pre-wrap;
    font-family: var(--font-mono);
  }
</style>
