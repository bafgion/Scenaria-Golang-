<script lang="ts">
  import { onMount } from 'svelte'
  import { createTranslator, locale } from './i18n'
  import {
    StartExport,
    PreviewExport,
    PickSaveFile,
    ArtifactExists,
  } from '../../wailsjs/go/wailsapp/App'
  import { startRunResultJob } from './asyncRunResult'
  import type { gui } from '../../wailsjs/go/models'

  export let inputPath = ''
  export let featureText = ''
  export let currentProjectVersion = 0
  export let onClose: () => void = () => {}
  export let onLog: (message: string) => void = () => {}

  $: tr = createTranslator($locale)

  const exportExtensions: Record<string, string> = {
    json: '.json',
    feature: '.feature',
    ts: '.ts',
    python: '.py',
  }

  let format = 'json'
  let outputPath = ''
  let baseURL = ''
  let forceOverwrite = false
  let outputExists = false
  let busy = false
  let preview: gui.ExportPreview | null = null
  let previewError = ''

  $: needsBaseURL = format === 'ts' || format === 'python'
  $: baseURLWarning = needsBaseURL && !baseURL.trim()
  $: canExport = Boolean(outputPath.trim()) && !busy && preview !== null

  onMount(() => {
    syncOutputPath()
    void refreshPreview()
  })

  function basename(path: string): string {
    const parts = path.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || path
  }

  function syncOutputPath() {
    if (!inputPath) return
    const base = inputPath.replace(/\.[^./\\]+$/, '')
    outputPath = base + (exportExtensions[format] || '.json')
    void checkOutputExists()
  }

  async function checkOutputExists() {
    const path = outputPath.trim()
    if (!path) {
      outputExists = false
      return
    }
    try {
      outputExists = await ArtifactExists(path)
    } catch {
      outputExists = false
    }
    if (!outputExists) forceOverwrite = false
  }

  async function refreshPreview() {
    previewError = ''
    try {
      preview = await PreviewExport(featureText)
    } catch (e: any) {
      preview = null
      previewError = String(e)
    }
  }

  async function onFormatChange() {
    syncOutputPath()
  }

  async function browseOutput() {
    const ext = exportExtensions[format] || '.json'
    const picked = await PickSaveFile(tr('dialogs.export.pickSaveTitle'), `export${ext}`)
    if (picked) {
      outputPath = picked
      await checkOutputExists()
    }
  }

  async function onOutputPathInput() {
    await checkOutputExists()
  }

  async function confirmExport() {
    if (!inputPath || !outputPath.trim() || busy) return
    if (outputExists && !forceOverwrite) {
      previewError = tr('dialogs.export.fileExistsError')
      return
    }
    busy = true
    previewError = ''
    try {
      const result = await startRunResultJob('export-finished', () => StartExport({
          inputPath,
          output: outputPath.trim(),
          format,
          baseURL: baseURL.trim(),
          force: forceOverwrite,
        }), currentProjectVersion)
      if (result.output) onLog(result.output.trimEnd())
      if (result.error) {
        onLog(`${tr('dialogs.export.errorPrefix')} ${result.error}`)
        previewError = result.error
        return
      }
      onLog(tr('dialogs.export.exported', { path: outputPath }))
      onClose()
    } catch (e: any) {
      previewError = String(e)
      onLog(`${tr('dialogs.export.errorPrefix')} ${e}`)
    } finally {
      busy = false
    }
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }

  function hintSeverityClass(severity: string): string {
    if (severity === 'warning') return 'warn'
    if (severity === 'error') return 'error'
    return 'info'
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="palette-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="palette export-dialog" role="dialog" aria-modal="true" aria-label={tr('dialogs.export.ariaLabel')} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{tr('dialogs.export.title')}</h3>
    <p class="hint">{tr('dialogs.export.source')} <code>{basename(inputPath)}</code></p>

    {#if preview}
      <div class="preview-summary">
        {#if preview.scenarioTitle}
          <span>{tr('dialogs.export.scenarioTitle', { title: preview.scenarioTitle })}</span>
        {/if}
        <span>{tr('dialogs.export.stepCount', { count: preview.stepCount })}</span>
      </div>
    {/if}

    <label>{tr('dialogs.export.format')}
      <select bind:value={format} on:change={onFormatChange} disabled={busy}>
        <option value="json">{tr('dialogs.export.formatJson')}</option>
        <option value="feature">{tr('dialogs.export.formatFeature')}</option>
        <option value="ts">{tr('dialogs.export.formatTs')}</option>
        <option value="python">{tr('dialogs.export.formatPython')}</option>
      </select>
    </label>

    <label>{tr('dialogs.export.file')}
      <div class="input-row">
        <input bind:value={outputPath} on:blur={onOutputPathInput} disabled={busy} />
        <button type="button" on:click={browseOutput} disabled={busy}>{tr('dialogs.common.browse')}</button>
      </div>
    </label>

    {#if outputExists}
      <label class="overwrite">
        <input type="checkbox" bind:checked={forceOverwrite} disabled={busy} />
        {tr('dialogs.export.overwrite')}
      </label>
    {/if}

    <label>{tr('dialogs.export.baseUrl')}
      <input bind:value={baseURL} placeholder={tr('dialogs.export.baseUrlPlaceholder')} disabled={busy} />
    </label>
    {#if baseURLWarning}
      <p class="warning">{tr('dialogs.export.baseUrlWarning')}</p>
    {/if}

    {#if preview && preview.issues.length > 0}
      <div class="issues">
        <strong>{tr('dialogs.export.issues', { count: preview.issues.length })}</strong>
        <ul>
          {#each preview.issues.slice(0, 8) as issue}
            <li>{tr('dialogs.export.issueLine', { line: issue.line, message: issue.message })}</li>
          {/each}
          {#if preview.issues.length > 8}
            <li>{tr('dialogs.export.andMore', { count: preview.issues.length - 8 })}</li>
          {/if}
        </ul>
      </div>
    {/if}

    {#if preview && preview.hints.length > 0}
      <div class="hints">
        <strong>{tr('dialogs.export.hints', { count: preview.hints.length })}</strong>
        <ul>
          {#each preview.hints.slice(0, 5) as hint}
            <li class={hintSeverityClass(hint.severity)}>{hint.title}</li>
          {/each}
          {#if preview.hints.length > 5}
            <li>{tr('dialogs.export.andMore', { count: preview.hints.length - 5 })}</li>
          {/if}
        </ul>
      </div>
    {/if}

    {#if previewError}<p class="error">{previewError}</p>{/if}

    <div class="actions">
      <button type="button" class="primary" disabled={!canExport} on:click={confirmExport}>
        {busy ? tr('dialogs.export.exporting') : tr('dialogs.export.confirm')}
      </button>
      <button type="button" on:click={onClose} disabled={busy}>{tr('dialogs.common.cancel')}</button>
    </div>
  </div>
</div>

<style>
  .export-dialog {
    width: min(560px, 96vw);
    max-height: 86vh;
    overflow: auto;
  }

  h3 {
    margin: 0 0 8px;
    font-size: 14px;
  }

  .hint {
    margin: 0 0 12px;
    font-size: 11px;
    color: var(--color-muted);
  }

  .preview-summary {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
    font-size: 12px;
    margin-bottom: 12px;
    color: var(--color-muted);
  }

  label {
    display: grid;
    gap: 4px;
    font-size: 11px;
    color: var(--color-muted);
    margin-bottom: 10px;
  }

  .overwrite {
    display: flex;
    align-items: center;
    gap: 8px;
    color: var(--color-warning);
  }

  .input-row {
    display: flex;
    gap: 8px;
  }

  .input-row input {
    flex: 1;
    min-width: 0;
  }

  .warning {
    margin: -4px 0 10px;
    font-size: 11px;
    color: var(--color-warning);
  }

  .issues, .hints {
    margin-bottom: 12px;
    font-size: 11px;
  }

  .issues ul, .hints ul {
    margin: 6px 0 0;
    padding-left: 18px;
  }

  .issues li {
    color: var(--color-error);
  }

  .hints li.warn {
    color: var(--color-warning);
  }

  .hints li.error {
    color: var(--color-error);
  }

  .error {
    color: var(--color-error);
    font-size: 11px;
    margin: 0 0 8px;
  }
</style>
