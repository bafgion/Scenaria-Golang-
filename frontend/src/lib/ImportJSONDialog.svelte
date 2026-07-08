<script lang="ts">
  import { onMount } from 'svelte'
  import { createTranslator, locale } from './i18n'
  import { StartImportJSON, PickOpenFile, PickSaveFile, ArtifactExists } from '../../wailsjs/go/wailsapp/App'
  import { startRunResultJob } from './asyncRunResult'

  export let projectPath = ''
  export let currentProjectVersion = 0
  export let onClose: () => void = () => {}
  export let onLog: (message: string) => void = () => {}
  export let onImported: (featurePath: string) => void = () => {}

  $: tr = createTranslator($locale)

  let jsonPath = ''
  let outputPath = ''
  let forceOverwrite = false
  let outputExists = false
  let busy = false
  let error = ''

  $: canImport = Boolean(jsonPath.trim()) && Boolean(outputPath.trim()) && !busy

  onMount(() => {
    void pickJson()
  })

  function basename(path: string): string {
    const parts = path.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || path
  }

  function defaultOutput(path: string): string {
    if (!path || !projectPath) return ''
    const name = basename(path).replace(/\.json$/i, '') + '.feature'
    return `${projectPath.replace(/\\/g, '/')}/${name}`.replace(/\//g, '\\')
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
  }

  async function pickJson() {
    const picked = await PickOpenFile(tr('dialogs.import.json.pickJsonTitle'))
    if (!picked) {
      if (!jsonPath) onClose()
      return
    }
    jsonPath = picked
    outputPath = defaultOutput(picked)
    error = ''
    await checkOutputExists()
  }

  async function pickOutput() {
    const picked = await PickSaveFile(tr('dialogs.import.json.pickSaveTitle'), basename(outputPath || 'imported.feature'))
    if (!picked) return
    outputPath = picked
    await checkOutputExists()
  }

  async function runImport() {
    if (!canImport) return
    if (outputExists && !forceOverwrite) {
      error = tr('dialogs.import.json.fileExistsError')
      return
    }
    busy = true
    error = ''
    onLog(tr('dialogs.import.json.importing', { path: jsonPath }))
    try {
      const result = await startRunResultJob('import-json-finished', () => StartImportJSON({
          jsonPath,
          outputPath: outputPath.trim(),
          force: forceOverwrite,
        }), currentProjectVersion)
      if (result.output) onLog(result.output.trimEnd())
      if (result.error) {
        error = result.error
        onLog(`${tr('dialogs.import.json.errorPrefix')} ${result.error}`)
        return
      }
      onImported(outputPath.trim())
      onClose()
    } finally {
      busy = false
    }
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal wide import-dialog" role="dialog" aria-modal="true" aria-label={tr('dialogs.import.json.ariaLabel')} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{tr('dialogs.import.json.title')}</h3>
    <label>
      {tr('dialogs.import.json.jsonFile')}
      <div class="path-row">
        <input bind:value={jsonPath} readonly placeholder={tr('dialogs.import.json.jsonPlaceholder')} />
        <button type="button" on:click={pickJson} disabled={busy}>{tr('dialogs.common.browse')}</button>
      </div>
    </label>
    <label>
      {tr('dialogs.import.json.feature')}
      <div class="path-row">
        <input bind:value={outputPath} on:change={checkOutputExists} placeholder={tr('dialogs.import.json.featurePlaceholder')} />
        <button type="button" on:click={pickOutput} disabled={busy}>{tr('dialogs.common.browse')}</button>
      </div>
    </label>
    {#if outputExists}
      <label class="check-row warn">
        <input type="checkbox" bind:checked={forceOverwrite} />
        {tr('dialogs.import.json.overwrite')}
      </label>
    {/if}
    {#if error}<p class="error">{error}</p>{/if}
    <div class="modal-actions">
      <button type="button" class="primary" on:click={runImport} disabled={!canImport}>{tr('dialogs.import.json.confirm')}</button>
      <button type="button" on:click={onClose} disabled={busy}>{tr('dialogs.common.cancel')}</button>
    </div>
  </div>
</div>

<style>
  .path-row {
    display: flex;
    gap: 8px;
  }

  .path-row input {
    flex: 1;
  }

  .warn {
    color: var(--color-muted);
    font-size: 12px;
  }

  .error {
    color: var(--color-error, #c62828);
    font-size: 12px;
    margin: 0;
  }
</style>
