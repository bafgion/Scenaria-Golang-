<script lang="ts">
  import { onMount } from 'svelte'
  import { createTranslator, locale } from './i18n'
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

  $: tr = createTranslator($locale)

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

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onCancel}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal wide tall" role="dialog" aria-modal="true" aria-label={tr('dialogs.vanessa.run.ariaLabel')} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{dryRun ? tr('dialogs.vanessa.run.titleDryRun') : tr('dialogs.vanessa.run.title')}</h3>
    <label>{tr('dialogs.vanessa.run.tag')} <input bind:value={tag} placeholder={tr('dialogs.vanessa.run.tagPlaceholder')} /></label>
    {#if tags.length > 0}
      <div class="tag-chips">
        {#each tags as t}
          <button type="button" class="chip" class:active={tag === t} on:click={() => pickTag(t)}>{t}</button>
        {/each}
      </div>
    {/if}
    <label>{tr('dialogs.vanessa.run.excludeTags')} <input bind:value={excludeTags} placeholder={tr('dialogs.vanessa.run.excludeTagsPlaceholder')} /></label>
    <label>{tr('dialogs.vanessa.run.scenario')} <input bind:value={scenario} placeholder={tr('dialogs.vanessa.run.scenarioPlaceholder')} list="vanessa-scenario-list" /></label>
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
      {tr('dialogs.vanessa.run.rerunFailed')}
      <select bind:value={rerunFailedRunDir} disabled={loadingDirs || runDirs.length === 0}>
        <option value="">{tr('dialogs.vanessa.run.rerunNone')}</option>
        {#each runDirs as dir}
          <option value={dir}>{shortDir(dir)}</option>
        {/each}
      </select>
    </label>
    {#if runDirs.length === 0 && !loadingDirs}
      <p class="hint">{tr('dialogs.vanessa.run.noRunDirs')}</p>
    {/if}
    <label class="check-row">
      <input type="checkbox" bind:checked={installEpf} />
      {tr('dialogs.vanessa.run.installEpf')}
    </label>
    {#if installEpf}
      <label>{tr('dialogs.vanessa.run.epfUrl')} <input bind:value={epfUrl} placeholder={tr('dialogs.vanessa.run.epfUrlPlaceholder')} /></label>
      <label>{tr('dialogs.vanessa.run.epfDest')} <input bind:value={epfDest} placeholder={tr('dialogs.vanessa.run.epfDestPlaceholder')} /></label>
    {/if}
    <button type="button" class="advanced-toggle" on:click={() => (showAdvanced = !showAdvanced)}>
      {showAdvanced ? tr('dialogs.vanessa.run.advancedExpanded') : tr('dialogs.vanessa.run.advancedCollapsed')}
    </button>
    {#if showAdvanced}
      <label>{tr('dialogs.vanessa.run.platformExe')} <input bind:value={platformExe} placeholder={tr('dialogs.vanessa.run.platformExePlaceholder')} /></label>
      <label>{tr('dialogs.vanessa.run.epfPath')} <input bind:value={epfPath} placeholder={tr('dialogs.vanessa.run.epfPathPlaceholder')} /></label>
      <label>{tr('dialogs.vanessa.run.ibConnection')} <input bind:value={ibConnection} placeholder={tr('dialogs.vanessa.run.ibConnectionPlaceholder')} /></label>
      <label class="check-row"><input type="checkbox" bind:checked={reportAllure} /> {tr('dialogs.vanessa.run.reportAllure')}</label>
      <label>{tr('dialogs.vanessa.run.vaDir')} <input bind:value={vaDir} placeholder={tr('dialogs.vanessa.run.vaDirPlaceholder')} /></label>
      <label>{tr('dialogs.vanessa.run.vaFiles')} <input bind:value={vaFiles} placeholder={tr('dialogs.vanessa.run.vaFilesPlaceholder')} /></label>
      <p class="hint">{tr('dialogs.vanessa.run.advancedHint')}</p>
    {:else}
      <p class="hint">{tr('dialogs.vanessa.run.defaultHint')}</p>
    {/if}
    <div class="modal-actions">
      <button type="button" class="primary" on:click={onConfirm}>{dryRun ? tr('dialogs.vanessa.run.dryRun') : tr('dialogs.vanessa.run.start')}</button>
      <button type="button" on:click={onCancel}>{tr('dialogs.common.cancel')}</button>
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
