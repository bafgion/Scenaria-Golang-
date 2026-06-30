<script lang="ts">
  import { onMount } from 'svelte'
  import { BrowserInstallStatus, InstallBrowserEngine, ListPlugins } from '../../wailsjs/go/wailsapp/App'
  import type { gui } from '../../wailsjs/go/models'
  import { BRAND_NAME } from './brand'
  import SettingCard from './SettingCard.svelte'
  import { DEFAULT_EDITOR_SETTINGS, type EditorSettings } from './editorOptions'

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
  export let toolbarCompact = false
  export let stepsPanelVisible = true
  export let stepsPanelHeight = 160
  export let checkUpdatesOnStartup = true
  export let selectorClickStrategies: string[] = ['testid', 'id', 'aria', 'contextual', 'text']
  export let selectorInputStrategies: string[] = ['testid', 'id', 'label', 'placeholder', 'aria', 'name']

  export let onSave: () => void
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

  const tabs: { id: TabId; label: string }[] = [
    { id: 'record', label: 'Запись и браузер' },
    { id: 'selectors', label: 'Селекторы' },
    { id: 'plugins', label: 'Плагины' },
    { id: 'editor', label: 'Редактор' },
    { id: 'ui', label: 'Интерфейс' },
  ]

  const slowMoPresets: [number, string][] = [
    [0, 'Быстро'],
    [100, 'Норма'],
    [250, 'Медленно'],
    [500, 'Учебный'],
  ]
  const defaultClickStrategies = ['testid', 'id', 'aria', 'contextual', 'text']
  const defaultInputStrategies = ['testid', 'id', 'label', 'placeholder', 'aria', 'name']

  const strategyLabels: Record<string, string> = {
    testid: 'data-testid',
    id: 'ID (#)',
    aria: 'aria-label',
    contextual: 'Контекстный has-text',
    text: 'has-text по тексту',
    label: 'label:has-text',
    placeholder: 'placeholder',
    name: 'name',
  }

  function labelForStrategy(key: string): string {
    return strategyLabels[key] || key
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

  export let editorSettings: EditorSettings = { ...DEFAULT_EDITOR_SETTINGS }

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
    else if (/редактор|monaco|шрифт|миникарт|перенос|tab|fold|sticky|подсказк|hint|сценари/.test(compact)) tab = 'editor'
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
  <div class="app-dialog" role="dialog" aria-modal="true" aria-label="Настройки — {BRAND_NAME}" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
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

            <SettingCard title="Минимальное наведение" description="Сколько миллисекунд курсор должен оставаться на элементе перед записью hover.">
              <span class="num-with-unit">
                <input type="number" class="setting-number" bind:value={hoverRecordMinMs} min={100} max={5000} step={50} />
                <span>мс</span>
              </span>
            </SettingCard>

            <SettingCard title="Прокрутка перед кликом" description="Перед записью клика прокручивать элемент в видимую область (как при воспроизведении).">
              <input type="checkbox" bind:checked={scrollBeforeClick} />
            </SettingCard>
          </section>

          <section class="setting-section">
            <h4 class="setting-section-title">Запуск</h4>
            <SettingCard title="Параллельные воркеры" description="Число одновременных браузерных сессий при пакетном запуске.">
              <span class="num-with-unit">
                <input type="number" class="setting-number" bind:value={workers} min={1} max={16} />
                <span>шт.</span>
              </span>
            </SettingCard>
            <SettingCard title="Скорость выполнения (slow-mo)" description="Пауза между действиями Playwright в миллисекундах. 0 — максимально быстро; 100–300 — удобно наблюдать шаги в браузере.">
              <span class="num-with-unit">
                <input type="number" class="setting-number" bind:value={slowMo} min={0} max={5000} step={50} />
                <span>мс</span>
              </span>
              <div class="slowmo-presets">
                {#each slowMoPresets as [ms, label]}
                  <button type="button" class="preset-chip" class:active={slowMo === ms} on:click={() => (slowMo = ms)}>{label}</button>
                {/each}
              </div>
            </SettingCard>
            <SettingCard title="Лимит итераций циклов" description="Максимум повторов для блоков «Повторяю» / «Пока».">
              <input type="number" class="setting-number" bind:value={loops} min={1} max={10000} />
            </SettingCard>
          </section>
        {:else if tab === 'selectors'}
          <section class="setting-section">
            <h4 class="setting-section-title">Приоритет стратегий</h4>
            <p class="setting-section-desc">
              При записи и подборе селектора {BRAND_NAME} перебирает стратегии сверху вниз. Более стабильные — выше.
            </p>
            <h5 class="strategy-group-title">Клики и кнопки</h5>
            <ul class="selector-list editable">
              {#each selectorClickStrategies as key, i}
                <li>
                  <span class="strategy-name">{labelForStrategy(key)}</span>
                  <span class="strategy-actions">
                    <button type="button" title="Выше" disabled={i === 0} on:click={() => moveClickStrategy(i, -1)}>↑</button>
                    <button type="button" title="Ниже" disabled={i === selectorClickStrategies.length - 1} on:click={() => moveClickStrategy(i, 1)}>↓</button>
                  </span>
                </li>
              {/each}
            </ul>
            <button type="button" class="dialog-link-btn" on:click={() => (selectorClickStrategies = [...defaultClickStrategies])}>Сбросить клики</button>
            <h5 class="strategy-group-title">Поля ввода</h5>
            <ul class="selector-list editable">
              {#each selectorInputStrategies as key, i}
                <li>
                  <span class="strategy-name">{labelForStrategy(key)}</span>
                  <span class="strategy-actions">
                    <button type="button" title="Выше" disabled={i === 0} on:click={() => moveInputStrategy(i, -1)}>↑</button>
                    <button type="button" title="Ниже" disabled={i === selectorInputStrategies.length - 1} on:click={() => moveInputStrategy(i, 1)}>↓</button>
                  </span>
                </li>
              {/each}
            </ul>
            <button type="button" class="dialog-link-btn" on:click={() => (selectorInputStrategies = [...defaultInputStrategies])}>Сбросить поля</button>
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
        {:else if tab === 'editor'}
          <section class="setting-section">
            <h4 class="setting-section-title">Шрифт и отображение</h4>
            <p class="setting-section-desc">Параметры Monaco-редактора сценариев.</p>
            <SettingCard title="Размер шрифта" description="От 8 до 32 px.">
              <input type="number" class="setting-number" bind:value={editorSettings.fontSize} min={8} max={32} />
            </SettingCard>
            <SettingCard title="Шрифт" description="Моноширинный шрифт редактора.">
              <input type="text" class="setting-text setting-text-mono" bind:value={editorSettings.fontFamily} />
            </SettingCard>
            <SettingCard title="Тема" description="Тёмная или светлая подсветка Gherkin.">
              <select bind:value={editorSettings.theme}>
                <option value="scenaria-dark">Тёмная</option>
                <option value="scenaria-light">Светлая</option>
              </select>
            </SettingCard>
            <SettingCard title="Перенос строк" description="Переносить длинные шаги по ширине редактора.">
              <select bind:value={editorSettings.wordWrap}>
                <option value="on">Включён</option>
                <option value="off">Выключен</option>
              </select>
            </SettingCard>
            <SettingCard title="Миникарта" description="Обзорная карта кода справа.">
              <input type="checkbox" bind:checked={editorSettings.minimap} />
            </SettingCard>
            <SettingCard title="Номера строк" description="Отображение номеров строк в gutter.">
              <select bind:value={editorSettings.lineNumbers}>
                <option value="on">Обычные</option>
                <option value="relative">Относительные</option>
                <option value="off">Скрыть</option>
              </select>
            </SettingCard>
            <SettingCard title="Пробелы" description="Когда показывать невидимые символы.">
              <select bind:value={editorSettings.renderWhitespace}>
                <option value="none">Не показывать</option>
                <option value="boundary">На границах слов</option>
                <option value="selection">В выделении</option>
                <option value="trailing">В конце строк</option>
                <option value="all">Всегда</option>
              </select>
            </SettingCard>
          </section>

          <section class="setting-section">
            <h4 class="setting-section-title">Ввод и поведение</h4>
            <SettingCard title="Размер табуляции" description="Ширина отступа Tab в пробелах.">
              <input type="number" class="setting-number" bind:value={editorSettings.tabSize} min={1} max={8} />
            </SettingCard>
            <SettingCard title="Пробелы вместо Tab" description="Вставлять пробелы при нажатии Tab.">
              <input type="checkbox" bind:checked={editorSettings.insertSpaces} />
            </SettingCard>
            <SettingCard title="Складывание блоков" description="Сворачивать блоки «Если» / «Повторяю».">
              <input type="checkbox" bind:checked={editorSettings.folding} />
            </SettingCard>
            <SettingCard title="Sticky scroll" description="Закреплять заголовки сценариев при прокрутке.">
              <input type="checkbox" bind:checked={editorSettings.stickyScroll} />
            </SettingCard>
            <SettingCard title="Авто-закрытие кавычек" description="Поведение при вводе кавычек.">
              <select bind:value={editorSettings.autoClosingQuotes}>
                <option value="languageDefined">По языку</option>
                <option value="always">Всегда</option>
                <option value="beforeWhitespace">Перед пробелом</option>
                <option value="never">Никогда</option>
              </select>
            </SettingCard>
            <SettingCard title="Форматировать при сохранении" description="Нормализовать отступы и убрать лишние пустые строки между шагами при Ctrl+S.">
              <input type="checkbox" bind:checked={editorSettings.formatOnSave} />
            </SettingCard>
            <SettingCard title="Подсказки при наведении" description="Показывать справку по шагу при hover в редакторе.">
              <input type="checkbox" bind:checked={editorSettings.stepHover} />
            </SettingCard>
            <SettingCard title="Проверка при вводе" description="Валидировать сценарий с задержкой при редактировании.">
              <input type="checkbox" bind:checked={editorSettings.validateOnType} />
            </SettingCard>
          </section>

          <section class="setting-section">
            <h4 class="setting-section-title">Навигация</h4>
            <p class="setting-section-desc">Структура feature-файла в breadcrumbs и панели шагов.</p>
            <SettingCard title="Breadcrumbs" description="Цепочка заголовков над редактором (Функционал → Сценарий → шаг).">
              <input type="checkbox" bind:checked={editorSettings.breadcrumbs} />
            </SettingCard>
            <SettingCard title="Структура в панели шагов" description="Вкладка «Структура» с деревом сценария и переходом по клику.">
              <input type="checkbox" bind:checked={editorSettings.symbolOutline} />
            </SettingCard>
            <SettingCard title="Вкладка панели по умолчанию" description="Что показывать под редактором при открытии сценария.">
              <select bind:value={editorSettings.stepsPanelView}>
                <option value="outline">Структура</option>
                <option value="steps">Таблица шагов</option>
              </select>
            </SettingCard>
            <SettingCard title="Code Lens для запуска" description="Кнопки «▶ Запустить» над сценариями и шагами в редакторе.">
              <input type="checkbox" bind:checked={editorSettings.codeLens} />
            </SettingCard>
            <SettingCard title="Inlay hints" description="Серые подсказки справа от шага: click → selector, fill → значение.">
              <input type="checkbox" bind:checked={editorSettings.inlayHints} />
            </SettingCard>
          </section>

          <section class="setting-section">
            <h4 class="setting-section-title">Подсказки сценария</h4>
            <p class="setting-section-desc">Эвристики качества шагов в редакторе (маркеры и quick fix).</p>
            <SettingCard title="Показывать подсказки" description="Маркеры warning/info в редакторе и панели «Проверка».">
              <input type="checkbox" bind:checked={editorSettings.scenarioHints} />
            </SettingCard>
            <SettingCard title="После записи" description="Анализировать сценарий сразу после остановки записи.">
              <input type="checkbox" bind:checked={editorSettings.scenarioHintsAfterRecord} disabled={!editorSettings.scenarioHints} />
            </SettingCard>
            <SettingCard title="Предупреждения" description="Подсказки уровня warning (хрупкие селекторы, дубли).">
              <input type="checkbox" bind:checked={editorSettings.scenarioHintsShowWarning} disabled={!editorSettings.scenarioHints} />
            </SettingCard>
            <SettingCard title="Информация" description="Подсказки уровня info (улучшения без критичных рисков).">
              <input type="checkbox" bind:checked={editorSettings.scenarioHintsShowInfo} disabled={!editorSettings.scenarioHints} />
            </SettingCard>
            <SettingCard title="Авто-исправление при сохранении" description="Применять autoFixable подсказки при Ctrl+S.">
              <input type="checkbox" bind:checked={editorSettings.scenarioHintsAutoFixOnSave} disabled={!editorSettings.scenarioHints} />
            </SettingCard>
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
                <input type="number" class="setting-number" bind:value={stepsPanelHeight} min={80} max={480} />
                <span>px</span>
              </span>
            </SettingCard>
          </section>

          <section class="setting-section">
            <h4 class="setting-section-title">Обновления</h4>
            <SettingCard title="Проверять при запуске" description="Искать новую версию {BRAND_NAME} при старте IDE.">
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

  .install-btn {
    flex-shrink: 0;
    padding: 6px 10px;
    font-size: 12px;
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

  .preset-chip {
    padding: 3px 8px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: transparent;
    color: var(--color-muted);
    font-size: 11px;
    cursor: pointer;
  }

  .preset-chip:hover {
    color: var(--color-text);
    border-color: var(--color-primary);
  }

  .preset-chip.active {
    color: var(--color-text);
    border-color: var(--color-primary);
    background: rgba(0, 122, 204, 0.15);
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

  .strategy-actions button {
    min-width: 24px;
    padding: 2px 6px;
    font-size: 11px;
  }
</style>

