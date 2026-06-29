<script lang="ts">
  import type { RunForm } from './runTypes'

  export let title = 'Запуск сценария'
  export let form: RunForm
  export let testClients: string[] = []
  export let tags: string[] = []
  export let onConfirm: () => void = () => {}
  export let onCancel: () => void = () => {}

  function pickTag(tag: string) {
    form = { ...form, tag }
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="modal-backdrop" role="presentation" on:click={onCancel}>
  <div class="modal wide run-dialog" role="dialog" aria-label={title} on:click|stopPropagation>
    <h3>{title}</h3>
    <label>Тег <input bind:value={form.tag} placeholder="@smoke" /></label>
    {#if tags.length > 0}
      <div class="tag-chips">
        {#each tags as tag}
          <button type="button" class="chip" class:active={form.tag === tag} on:click={() => pickTag(tag)}>{tag}</button>
        {/each}
      </div>
    {/if}
    <label>TestClient
      <select bind:value={form.testClient}>
        <option value="">(из feature / не задан)</option>
        {#each testClients as client}
          <option value={client}>{client}</option>
        {/each}
      </select>
    </label>
    <div class="row-2">
      <label>Движок
        <select bind:value={form.engine} disabled={form.dryRun}>
          <option value="playwright">playwright</option>
          <option value="stub">stub</option>
        </select>
      </label>
      <label>Браузер
        <select bind:value={form.browser} disabled={form.dryRun}>
          <option value="chromium">chromium</option>
          <option value="firefox">firefox</option>
          <option value="webkit">webkit</option>
        </select>
      </label>
    </div>
    <div class="row-2">
      <label>Параллельные воркеры
        <input type="number" bind:value={form.workers} min="1" max="16" disabled={form.dryRun} />
      </label>
      <label>Slow-mo (мс)
        <input type="number" bind:value={form.slowMo} min="0" step="50" disabled={form.dryRun} />
      </label>
    </div>
    <label>Переменные (NAME=VALUE)
      <textarea bind:value={form.vars} placeholder="BASE_URL=https://example.com"></textarea>
    </label>
    <label>Base URL (переопределение)
      <input bind:value={form.baseUrl} placeholder="https://example.com" disabled={form.dryRun} />
    </label>
    <label class="check-row"><input type="checkbox" bind:checked={form.dryRun} /> Dry-run (без браузера)</label>
    <label class="check-row"><input type="checkbox" bind:checked={form.headed} disabled={form.dryRun} /> Headed (видимый браузер)</label>
    <label class="check-row"><input type="checkbox" bind:checked={form.installPW} disabled={form.dryRun} /> Установить Playwright при необходимости</label>
    <label class="check-row"><input type="checkbox" bind:checked={form.allure} disabled={form.dryRun} /> Allure (.scenaria/allure-results)</label>
    <label class="check-row"><input type="checkbox" bind:checked={form.trace} disabled={form.dryRun} /> Trace</label>
    <label class="check-row"><input type="checkbox" bind:checked={form.video} disabled={form.dryRun} /> Video</label>
    <label class="check-row"><input type="checkbox" bind:checked={form.html} disabled={form.dryRun} /> HTML-отчёт (.scenaria/report.html)</label>
    <label class="check-row"><input type="checkbox" bind:checked={form.junit} disabled={form.dryRun} /> JUnit (.scenaria/junit.xml)</label>
    <label class="check-row"><input type="checkbox" bind:checked={form.summaryJson} disabled={form.dryRun} /> Summary JSON (.scenaria/summary.json)</label>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={onConfirm}>Запустить</button>
      <button type="button" on:click={onCancel}>Отмена</button>
    </div>
  </div>
</div>

<style>
  .run-dialog {
    max-height: 90vh;
    overflow: auto;
  }

  .row-2 {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 10px;
  }

  .tag-chips {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin: -4px 0 8px;
  }

  .chip {
    padding: 2px 8px;
    font-size: 11px;
    border: 1px solid var(--color-border);
    border-radius: 10px;
    background: var(--color-input);
    color: var(--color-muted);
  }

  .chip.active {
    border-color: var(--color-primary);
    color: var(--color-text);
    background: var(--color-selected);
  }
</style>
