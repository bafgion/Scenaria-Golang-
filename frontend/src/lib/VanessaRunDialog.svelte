<script lang="ts">
  import { onMount } from 'svelte'
  import { ListVanessaRunDirs, ReadVanessaSettingsJSON } from '../../wailsjs/go/wailsapp/App'

  export let dryRun = false
  export let preferRerun = false
  export let tag = ''
  export let excludeTags = ''
  export let scenario = ''
  export let rerunFailedRunDir = ''
  export let installEpf = false
  export let epfUrl = ''
  export let epfDest = ''
  export let platformExe = ''
  export let epfPath = ''
  export let ibConnection = ''
  export let reportAllure = false
  export let vaDir = ''
  export let vaFiles = ''
  export let tags: string[] = []
  export let scenarios: string[] = []
  export let onConfirm: () => void = () => {}
  export let onCancel: () => void = () => {}

  let runDirs: string[] = []
  let loadingDirs = false
  let showAdvanced = false

  function pickTag(value: string) {
    tag = value
  }

  function shortDir(path: string): string {
    const parts = path.replace(/\\/g, '/').split('/')
    return parts.slice(-2).join('/')
  }

  onMount(async () => {
    loadingDirs = true
    try {
      runDirs = await ListVanessaRunDirs(15)
      if (!rerunFailedRunDir && runDirs.length > 0 && preferRerun) {
        rerunFailedRunDir = runDirs[0]
      }
    } catch {
      runDirs = []
    } finally {
      loadingDirs = false
    }
    try {
      const raw = await ReadVanessaSettingsJSON()
      const cfg = JSON.parse(raw)
      if (!platformExe && cfg.platform_executable) platformExe = cfg.platform_executable
      if (!epfPath && cfg.epf_path) epfPath = cfg.epf_path
      if (!ibConnection && cfg.ib_connection_string) ibConnection = cfg.ib_connection_string
      if (!reportAllure && cfg.report_allure) reportAllure = Boolean(cfg.report_allure)
    } catch {
      /* defaults stay empty */
    }
  })

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="modal-backdrop" role="presentation" on:click={onCancel}>
  <div class="modal wide tall" role="dialog" aria-label="Запуск Vanessa" on:click|stopPropagation>
    <h3>Vanessa {dryRun ? '(dry-run)' : ''}</h3>
    <label>Тег <input bind:value={tag} placeholder="@smoke" /></label>
    {#if tags.length > 0}
      <div class="tag-chips">
        {#each tags as t}
          <button type="button" class="chip" class:active={tag === t} on:click={() => pickTag(t)}>{t}</button>
        {/each}
      </div>
    {/if}
    <label>Исключить теги (через запятую) <input bind:value={excludeTags} placeholder="@wip, @draft" /></label>
    <label>Имя сценария <input bind:value={scenario} placeholder="Необязательно" list="vanessa-scenario-list" /></label>
    {#if scenarios.length > 0}
      <datalist id="vanessa-scenario-list">
        {#each scenarios as name}
          <option value={name}></option>
        {/each}
      </datalist>
      <div class="tag-chips">
        {#each scenarios as name}
          <button type="button" class="chip" class:active={scenario === name} on:click={() => (scenario = name)}>{name}</button>
        {/each}
      </div>
    {/if}
    <label>
      Перезапуск упавших (run-dir)
      <select bind:value={rerunFailedRunDir} disabled={loadingDirs || runDirs.length === 0}>
        <option value="">(не использовать)</option>
        {#each runDirs as dir}
          <option value={dir}>{shortDir(dir)}</option>
        {/each}
      </select>
    </label>
    {#if runDirs.length === 0 && !loadingDirs}
      <p class="hint">Нет каталогов run-* — сначала выполните обычный запуск Vanessa.</p>
    {/if}
    <label class="check-row">
      <input type="checkbox" bind:checked={installEpf} />
      Установить/обновить EPF (--epf-install)
    </label>
    {#if installEpf}
      <label>URL EPF (необязательно) <input bind:value={epfUrl} placeholder="https://…/vanessa-automation.epf" /></label>
      <label>Путь назначения EPF <input bind:value={epfDest} placeholder="C:\vanessa\vanessa-automation.epf" /></label>
    {/if}
    <button type="button" class="advanced-toggle" on:click={() => (showAdvanced = !showAdvanced)}>
      {showAdvanced ? '▼' : '▶'} Параметры 1C и пути
    </button>
    {#if showAdvanced}
      <label>Платформа 1C (--platform-exe) <input bind:value={platformExe} placeholder="C:\Program Files\1cv8\bin\1cv8.exe" /></label>
      <label>EPF (--epf) <input bind:value={epfPath} placeholder="C:\vanessa\vanessa-automation.epf" /></label>
      <label>Строка подключения IB (--ib) <input bind:value={ibConnection} placeholder='Srvr="localhost";Ref="base";' /></label>
      <label class="check-row"><input type="checkbox" bind:checked={reportAllure} /> Allure-отчёт (--allure)</label>
      <label>Папка сценариев (--dir) <input bind:value={vaDir} placeholder="features\smoke" /></label>
      <label>Файлы сценариев (--files) <input bind:value={vaFiles} placeholder="a.feature,b.feature" /></label>
      <p class="hint">Пустые поля берутся из <code>.scenaria/vanessa.json</code>. Без --dir/--files запускается весь проект.</p>
    {:else}
      <p class="hint">По умолчанию — настройки из <code>.scenaria/vanessa.json</code>.</p>
    {/if}
    <div class="modal-actions">
      <button type="button" class="primary" on:click={onConfirm}>{dryRun ? 'Dry-run' : 'Запустить'}</button>
      <button type="button" on:click={onCancel}>Отмена</button>
    </div>
  </div>
</div>

<style>
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

  .hint {
    font-size: 11px;
    color: var(--color-muted);
    margin: 0;
  }

  .advanced-toggle {
    margin: 8px 0;
    padding: 4px 0;
    border: none;
    background: none;
    color: var(--color-muted);
    font-size: 12px;
    cursor: pointer;
    text-align: left;
  }

  .advanced-toggle:hover {
    color: var(--color-text);
  }
</style>
