<script lang="ts">
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

  export let onSave: () => void
  export let onCancel: () => void

  let tab: 'ui' | 'record' | 'plugins' = 'ui'
  let search = ''

  const tabs = [
    { id: 'ui' as const, label: 'Интерфейс' },
    { id: 'record' as const, label: 'Запись и браузер' },
    { id: 'plugins' as const, label: 'Плагины' },
  ]

  function pickTab(id: typeof tab) {
    tab = id
    search = ''
  }

  function onSearch() {
    const q = search.trim().toLowerCase()
    if (!q) return
    if ('записьбраузерheadless'.includes(q.replace(/\s/g, ''))) tab = 'record'
    else if ('плагинvanessa'.includes(q.replace(/\s/g, ''))) tab = 'plugins'
    else tab = 'ui'
  }
</script>

<div class="modal-backdrop" role="presentation">
  <div class="settings-dialog" role="dialog" aria-label="Настройки">
    <header class="settings-head">
      <h3>Настройки</h3>
      <input bind:value={search} placeholder="Поиск настроек…" on:input={onSearch} />
    </header>
    <div class="settings-body">
      <nav class="settings-nav">
        {#each tabs as t}
          <button class:active={tab === t.id} on:click={() => pickTab(t.id)}>{t.label}</button>
        {/each}
      </nav>
      <div class="settings-content">
        {#if tab === 'ui'}
          <label class="check-row">
            <input type="checkbox" bind:checked={toolbarCompact} />
            Компактная панель инструментов
          </label>
          <label class="check-row">
            <input type="checkbox" bind:checked={stepsPanelVisible} />
            Показывать панель шагов
          </label>
          <label>
            Высота панели шагов (px)
            <input type="number" bind:value={stepsPanelHeight} min="80" max="480" />
          </label>
          <label>
            Параллельные воркеры
            <input type="number" bind:value={workers} min="1" />
          </label>
          <label>
            Лимит итераций циклов
            <input type="number" bind:value={loops} min="1" />
          </label>
        {:else if tab === 'record'}
          <label>
            Браузер по умолчанию
            <select bind:value={browser}>
              <option value="chromium">chromium</option>
              <option value="firefox">firefox</option>
              <option value="webkit">webkit</option>
            </select>
          </label>
          <label class="check-row">
            <input type="checkbox" bind:checked={headless} />
            Без окна браузера (headless)
          </label>
          <label class="check-row">
            <input type="checkbox" bind:checked={filterRecording} on:change={() => filterRecording && (navOnlyRecording = false)} />
            Только важные (фильтр записи)
          </label>
          <label class="check-row">
            <input type="checkbox" bind:checked={navOnlyRecording} on:change={() => navOnlyRecording && (filterRecording = false)} />
            Только ссылки
          </label>
          <label class="check-row">
            <input type="checkbox" bind:checked={hoverRecord} />
            Записывать наведение
          </label>
        {:else}
          <p class="hint">Vanessa и другие плагины настраиваются в <code>.scenaria/vanessa.json</code> и через меню «Запись и тест».</p>
          <p class="hint">Установка ZIP-плагинов — в CLI: <code>scenaria plugins install …</code></p>
        {/if}
      </div>
    </div>
    <footer class="settings-foot">
      <button class="primary" on:click={onSave}>ОК</button>
      <button on:click={onCancel}>Отмена</button>
    </footer>
  </div>
</div>

<style>
  .settings-dialog {
    width: min(640px, 92vw);
    max-height: 88vh;
    background: var(--color-sidebar);
    border: 1px solid var(--color-border);
    border-radius: 4px;
    display: flex;
    flex-direction: column;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.45);
  }

  .settings-head {
    padding: 14px 16px 10px;
    border-bottom: 1px solid var(--color-border);
  }

  .settings-head h3 {
    margin: 0 0 10px;
    font-size: 14px;
  }

  .settings-head input {
    width: 100%;
    padding: 6px 8px;
    background: var(--color-input);
    border: 1px solid var(--color-border);
    border-radius: 3px;
    color: var(--color-text);
  }

  .settings-body {
    display: grid;
    grid-template-columns: 180px 1fr;
    min-height: 280px;
    flex: 1;
    overflow: hidden;
  }

  .settings-nav {
    border-right: 1px solid var(--color-border);
    padding: 8px 0;
    display: flex;
    flex-direction: column;
  }

  .settings-nav button {
    text-align: left;
    padding: 8px 14px;
    border: none;
    background: transparent;
    color: var(--color-text);
    font-size: 13px;
  }

  .settings-nav button.active {
    background: var(--color-selected);
  }

  .settings-content {
    padding: 16px;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .settings-content label {
    display: flex;
    flex-direction: column;
    gap: 4px;
    font-size: 12px;
    color: var(--color-text);
  }

  .settings-content select,
  .settings-content input[type='number'] {
    padding: 6px 8px;
    background: var(--color-input);
    border: 1px solid var(--color-border);
    border-radius: 3px;
    color: var(--color-text);
  }

  .check-row {
    flex-direction: row !important;
    align-items: center;
    gap: 8px !important;
  }

  .check-row input {
    width: auto;
  }

  .settings-foot {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    padding: 10px 16px;
    border-top: 1px solid var(--color-border);
  }

  .settings-foot button {
    padding: 6px 14px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: #fff;
  }

  .settings-foot .primary {
    background: var(--color-primary);
    border-color: var(--color-primary);
  }

  .hint {
    font-size: 12px;
    color: var(--color-muted);
    line-height: 1.5;
    margin: 0;
  }
</style>
