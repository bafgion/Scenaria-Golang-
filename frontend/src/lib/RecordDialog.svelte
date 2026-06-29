<script lang="ts">
  import StepsInsertDialog from './StepsInsertDialog.svelte'
  import NumberInput from './NumberInput.svelte'

  export let mode: 'live' | 'baseline' = 'live'
  export let url = ''
  export let output = ''
  export let testClient = ''
  export let idleSeconds = 30
  export let appendTo = ''
  export let headless = false
  export let filterRecording = false
  export let navOnlyRecording = false
  export let hoverRecord = false
  export let featureName = 'Записанный сценарий'
  export let scenarioName = 'Запись'
  export let testClients: string[] = []
  export let recording = false
  export let recordPaused = false
  export let baselineBusy = false
  export let onHttpAuth: () => void = () => {}
  export let onStart: () => void = () => {}
  export let onTogglePause: () => void = () => {}
  export let onStop: () => void = () => {}
  export let onSaveBaseline: (payload: {
    output: string
    featureName: string
    scenarioName: string
    steps: string[]
  }) => void = () => {}
  export let onClose: () => void = () => {}

  let steps: string[] = []
  let newStep = ''
  let showStepPicker = false
  let baselineInitUrl = ''

  $: if (mode === 'baseline' && url !== baselineInitUrl) {
    baselineInitUrl = url
    const start = url.trim() || 'https://example.com'
    steps = [`открываю "${start}"`]
  }

  function addStep() {
    const text = newStep.trim()
    if (!text) return
    steps = [...steps, text]
    newStep = ''
  }

  function removeStep(index: number) {
    steps = steps.filter((_, i) => i !== index)
  }

  function moveStep(index: number, delta: number) {
    const next = index + delta
    if (next < 0 || next >= steps.length) return
    const copy = [...steps]
    const tmp = copy[index]
    copy[index] = copy[next]
    copy[next] = tmp
    steps = copy
  }

  function saveBaseline() {
    onSaveBaseline({
      output: output.trim(),
      featureName: featureName.trim(),
      scenarioName: scenarioName.trim(),
      steps: steps.map((s) => s.trim()).filter(Boolean),
    })
  }

  $: previewText = buildPreview(featureName, scenarioName, steps)

  function buildPreview(feature: string, scenario: string, stepList: string[]): string {
    const title = feature.trim() || 'Записанный сценарий'
    const scen = scenario.trim() || 'Базовый сценарий'
    const lines = stepList.map((s) => `    Когда ${s.trim()}`).join('\n')
    return `# language: ru\nФункционал: ${title}\n  Сценарий: ${scen}\n${lines || '    Когда выполняю действие'}`
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape' && !showStepPicker) onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

{#if showStepPicker}
  <StepsInsertDialog
    onInsert={(template) => {
      steps = [...steps, template]
      showStepPicker = false
    }}
    onClose={() => (showStepPicker = false)}
  />
{:else}
  <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
  <div class="modal-backdrop" role="presentation" on:click={onClose}>
    <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
    <div class="modal wide record-dialog" role="dialog" aria-modal="true" aria-label="Запись сценария" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="tabs" role="tablist">
        <button
          type="button"
          role="tab"
          class:active={mode === 'live'}
          disabled={recording}
          on:click={() => (mode = 'live')}
        >
          Live-запись
        </button>
        <button
          type="button"
          role="tab"
          class:active={mode === 'baseline'}
          disabled={recording}
          on:click={() => (mode = 'baseline')}
        >
          Из шагов
        </button>
      </div>

      {#if mode === 'live'}
        <h3>Live-запись</h3>
        <label>URL <input bind:value={url} disabled={recording} /></label>
        <label>Функционал <input bind:value={featureName} placeholder="Записанный сценарий" disabled={recording} /></label>
        <label>Сценарий <input bind:value={scenarioName} placeholder="Запись" disabled={recording} /></label>
        <label>Файл <input bind:value={output} disabled={recording} /></label>
        <label>Дописать в существующий feature
          <input bind:value={appendTo} placeholder="путь к .feature или пусто" disabled={recording} />
        </label>
        <label>TestClient
          <select bind:value={testClient} disabled={recording}>
            <option value="">(без профиля)</option>
            {#each testClients as client}
              <option value={client}>{client}</option>
            {/each}
          </select>
        </label>
        <label for="record-idle">Idle (сек)</label>
        <NumberInput inputId="record-idle" bind:value={idleSeconds} min={5} disabled={recording} width="72px" />
        <label class="check-row"><input type="checkbox" bind:checked={headless} disabled={recording} /> Headless</label>
        <label class="check-row">
          <input type="checkbox" bind:checked={filterRecording} disabled={recording} on:change={() => filterRecording && (navOnlyRecording = false)} />
          Только важные (фильтр записи)
        </label>
        <label class="check-row">
          <input type="checkbox" bind:checked={navOnlyRecording} disabled={recording} on:change={() => navOnlyRecording && (filterRecording = false)} />
          Только ссылки
        </label>
        <label class="check-row"><input type="checkbox" bind:checked={hoverRecord} disabled={recording} /> Записывать наведение</label>
        <p class="hint">В URL можно указать user:pass@host — пароль сохранится для хоста. Пикер элемента доступен на паузе.</p>
        <div class="modal-actions">
          <button type="button" on:click={onHttpAuth} disabled={recording}>HTTP Auth…</button>
          {#if recording}
            <button type="button" on:click={onTogglePause}>{recordPaused ? 'Resume' : 'Pause'}</button>
            <button type="button" on:click={onStop}>Стоп</button>
          {:else}
            <button type="button" class="primary" on:click={onStart}>Начать</button>
          {/if}
          <button type="button" on:click={onClose}>Закрыть</button>
        </div>
      {:else}
        <h3>Запись из шагов</h3>
        <p class="hint">Создаёт .feature без браузера — аналог <code>scenaria record --step …</code>.</p>
        <label>Файл <input bind:value={output} disabled={baselineBusy} /></label>
        <label>Функционал <input bind:value={featureName} disabled={baselineBusy} /></label>
        <label>Сценарий <input bind:value={scenarioName} disabled={baselineBusy} /></label>
        <label>Стартовый URL (для первого шага) <input bind:value={url} disabled={baselineBusy} /></label>

        <div class="steps-header">
          <span>Шаги ({steps.length})</span>
          <button type="button" disabled={baselineBusy} on:click={() => (showStepPicker = true)}>Из каталога…</button>
        </div>
        <ol class="step-list">
          {#each steps as step, index}
            <li>
              <span class="step-text">{step}</span>
              <span class="step-actions">
                <button type="button" title="Выше" disabled={baselineBusy || index === 0} on:click={() => moveStep(index, -1)}>↑</button>
                <button type="button" title="Ниже" disabled={baselineBusy || index === steps.length - 1} on:click={() => moveStep(index, 1)}>↓</button>
                <button type="button" title="Удалить" disabled={baselineBusy} on:click={() => removeStep(index)}>×</button>
              </span>
            </li>
          {/each}
        </ol>

        <div class="add-step">
          <input bind:value={newStep} placeholder="Текст шага…" disabled={baselineBusy} on:keydown={(e) => e.key === 'Enter' && addStep()} />
          <button type="button" disabled={baselineBusy} on:click={addStep}>Добавить</button>
        </div>

        <div class="preview-label">Предпросмотр Gherkin
          <pre class="preview">{previewText}</pre>
        </div>

        <div class="modal-actions">
          <button type="button" class="primary" disabled={baselineBusy} on:click={saveBaseline}>Сохранить feature</button>
          <button type="button" disabled={baselineBusy} on:click={onClose}>Отмена</button>
        </div>
      {/if}
    </div>
  </div>
{/if}

<style>
  .record-dialog {
    max-height: 90vh;
    overflow: auto;
  }

  .tabs {
    display: flex;
    gap: 4px;
    margin-bottom: 12px;
    border-bottom: 1px solid var(--color-border);
    padding-bottom: 8px;
  }

  .tabs button {
    padding: 4px 10px;
    font-size: 12px;
    border: 1px solid transparent;
    border-radius: 3px 3px 0 0;
    background: transparent;
    color: var(--color-muted);
    cursor: pointer;
  }

  .tabs button.active {
    color: var(--color-text);
    border-color: var(--color-border);
    border-bottom-color: var(--color-panel, var(--color-input));
    background: var(--color-panel, var(--color-input));
  }

  h3 {
    margin: 0 0 12px;
    font-size: 14px;
  }

  label {
    display: grid;
    gap: 4px;
    margin-bottom: 8px;
    font-size: 11px;
    color: var(--color-muted);
  }

  input,
  select {
    padding: 6px 8px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: var(--color-text);
    font-size: 12px;
  }

  .check-row {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 12px;
    color: var(--color-text);
    margin-bottom: 6px;
  }

  .hint {
    margin: 0 0 12px;
    font-size: 11px;
    color: var(--color-muted);
    line-height: 1.4;
  }

  .steps-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin: 12px 0 6px;
    font-size: 12px;
  }

  .step-list {
    margin: 0 0 10px;
    padding-left: 20px;
    max-height: 200px;
    overflow: auto;
    font-size: 12px;
  }

  li {
    display: flex;
    align-items: flex-start;
    gap: 8px;
    margin-bottom: 6px;
  }

  .step-text {
    flex: 1;
    word-break: break-word;
  }

  .step-actions {
    display: flex;
    gap: 2px;
    flex-shrink: 0;
  }

  .step-actions button {
    padding: 2px 6px;
    font-size: 11px;
  }

  .add-step {
    display: grid;
    grid-template-columns: 1fr auto;
    gap: 8px;
    margin-bottom: 12px;
  }

  .preview-label {
    display: grid;
    gap: 4px;
    margin-bottom: 12px;
    font-size: 11px;
    color: var(--color-muted);
  }

  .preview {
    margin: 0;
    padding: 8px;
    font-family: Consolas, 'Courier New', monospace;
    font-size: 11px;
    line-height: 1.45;
    white-space: pre-wrap;
    max-height: 120px;
    overflow: auto;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: var(--color-text);
  }

  button {
    padding: 6px 12px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: var(--color-text);
    font-size: 12px;
  }

  button.primary {
    background: var(--color-accent);
    color: var(--color-on-accent, #fff);
    border-color: var(--color-accent);
  }
</style>
