<script lang="ts">
  import { createTranslator, locale } from './i18n'
  import { gui } from '../../wailsjs/go/models'
  import { flakyLabel, type FlakyScenarioStat } from './flakyMetrics'
  import { isUntitled, untitledLabel } from './untitled'

  export let entries: gui.RunResultEntry[] = []
  export let flakyByPath: Map<string, FlakyScenarioStat> = new Map()
  export let flakyStepByPath: Map<string, string> = new Map()
  export let artifacts: gui.ProjectArtifacts = new gui.ProjectArtifacts()
  export let onOpenFeature: (path: string) => void = () => {}
  export let onRerun: () => void = () => {}
  export let onRunFlaky: (entry: gui.RunResultEntry) => void = () => {}
  export let onOpenFolder: (path: string) => void = () => {}
  export let onServeAllure: (path: string) => void = () => {}
  export let onOpenHtmlReport: (path: string) => void = () => {}
  export let onOpenTrace: (path: string) => void = () => {}
  export let onGotoFailedStep: (entry: gui.RunResultEntry) => void = () => {}
  export let allureInstalled = true
  export let allureRunning = false
  export let onOpenAllureInstall: () => void = () => {}

  $: tr = createTranslator($locale)

  function splitPath(path: string): { feature: string; scenario: string } {
    const idx = path.indexOf('::')
    if (idx < 0) return { feature: path, scenario: '' }
    return { feature: path.slice(0, idx), scenario: path.slice(idx + 2) }
  }

  function basename(p: string): string {
    const parts = p.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || p
  }

  function featureDisplayName(path: string): string {
    if (isUntitled(path)) return untitledLabel(path)
    return basename(path)
  }

  function formatAt(at: string): string {
    if (!at) return ''
    try {
      return new Date(at).toLocaleString()
    } catch {
      return at
    }
  }
  function openEntry(entry: gui.RunResultEntry) {
    const feature = splitPath(entry.path).feature
    if (feature) onOpenFeature(feature)
  }
</script>

<div class="results-panel">
  <div class="results-toolbar">
    <div class="artifact-btns">
      {#if artifacts.allureDir}
        {#if allureInstalled}
          <button type="button" on:click={() => onServeAllure(artifacts.allureDir)}>{allureRunning ? tr('results.allureAgain') : tr('results.allureServe')}</button>
        {:else}
          <button type="button" class="warn-btn" on:click={onOpenAllureInstall}>{tr('results.allureNotFound')}</button>
        {/if}
        <button type="button" on:click={() => onOpenFolder(artifacts.allureDir)}>{tr('results.allureFolder')}</button>
      {/if}
      {#if artifacts.htmlReport}
        <button type="button" on:click={() => onOpenHtmlReport(artifacts.htmlReport)}>{tr('results.htmlReport')}</button>
      {/if}
      {#if artifacts.junitReport}
        <button type="button" on:click={() => onOpenFolder(artifacts.junitReport)}>JUnit</button>
      {/if}
      {#if artifacts.summaryJson}
        <button type="button" on:click={() => onOpenFolder(artifacts.summaryJson)}>Summary JSON</button>
      {/if}
      {#if artifacts.tracesDir}
        <button type="button" on:click={() => onOpenTrace(artifacts.tracesDir)}>{tr('results.traceViewer')}</button>
        <button type="button" on:click={() => onOpenFolder(artifacts.tracesDir)}>{tr('results.traceFolder')}</button>
      {/if}
      {#if artifacts.videosDir}
        <button type="button" on:click={() => onOpenFolder(artifacts.videosDir)}>Video</button>
      {/if}
    </div>
    <button type="button" class="rerun" on:click={onRerun}>{tr('results.rerunFailed')}</button>
  </div>
  {#if entries.length === 0}
    <p class="empty">{tr('results.empty')}</p>
  {:else}
    <table class="results-table">
      <thead>
        <tr>
          <th>{tr('results.col.scenario')}</th>
          <th>{tr('results.col.result')}</th>
          <th>{tr('results.col.message')}</th>
          <th>{tr('results.col.time')}</th>
        </tr>
      </thead>
      <tbody>
        {#each entries as entry}
          {@const parts = splitPath(entry.path)}
          {@const flakyStat = flakyByPath.get(entry.path)}
          {@const stepHint = flakyStepByPath.get(entry.path)}
          <tr
            class:failed={!entry.success}
            class:flaky={flakyStat?.flaky}
            class:clickable={!!parts.feature}
            on:dblclick={() => openEntry(entry)}
            title={tr('results.dblClickHint')}
          >
            <td>
              <div class="scenario-name">{parts.scenario || basename(parts.feature)}</div>
              <div class="feature-name">{featureDisplayName(parts.feature)}</div>
              {#if flakyStat?.flaky}
                <div class="flaky-row">
                  <div class="flaky-tag">{flakyLabel(flakyStat)}</div>
                  <button type="button" class="flaky-rerun" on:click|stopPropagation={() => onRunFlaky(entry)}>
                    {tr('results.run3x')}
                  </button>
                </div>
              {/if}
              {#if stepHint}
                <div class="step-flaky">{stepHint}</div>
              {/if}
            </td>
            <td class="status">
              {entry.success ? '✓ OK' : '✗ FAIL'}
              {#if !entry.success && entry.failed_step != null && entry.failed_step >= 0}
                <button type="button" class="goto-step" on:click|stopPropagation={() => onGotoFailedStep(entry)}>
                  {tr('results.gotoStep', { n: entry.failed_step + 1 })}
                </button>
              {/if}
            </td>
            <td class="msg">{entry.message || tr('dialogs.common.notFound')}</td>
            <td class="at">{formatAt(entry.at)}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>

<style>
  .results-panel {
    height: 100%;
    overflow: auto;
    padding: 8px 12px;
  }

  .results-toolbar {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
    flex-wrap: wrap;
  }

  .artifact-btns {
    display: flex;
    gap: 4px;
    flex-wrap: wrap;
  }

  .rerun {
    margin-left: auto;
  }

  .warn-btn {
    color: var(--color-warning, #dcdcaa);
  }

  .empty {
    font-size: 12px;
    color: var(--color-muted);
    margin: 8px 0;
  }

  .results-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 12px;
  }

  .results-table th {
    text-align: left;
    padding: 6px 8px;
    background: var(--color-toolbar);
    color: var(--color-muted);
    border-bottom: 1px solid var(--color-border);
    position: sticky;
    top: 0;
  }

  .results-table td {
    padding: 6px 8px;
    border-bottom: 1px solid var(--color-divider);
    vertical-align: top;
  }

  tr.clickable {
    cursor: pointer;
  }

  tr.clickable:hover td {
    background: var(--color-list-hover, rgba(255, 255, 255, 0.04));
  }

  tr.failed td.status {
    color: var(--color-error);
  }

  tr.flaky td:first-child {
    border-left: 2px solid var(--color-warning, #d4a017);
  }

  .flaky-tag {
    font-size: 10px;
    color: var(--color-warning, #d4a017);
  }

  .flaky-row {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: 2px;
    flex-wrap: wrap;
  }

  .flaky-rerun {
    font-size: 10px;
    padding: 1px 6px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: var(--color-primary);
    cursor: pointer;
  }

  .flaky-rerun:hover {
    background: var(--color-selected);
  }

  .step-flaky {
    font-size: 10px;
    color: var(--color-muted);
  }

  tr:not(.failed) td.status {
    color: var(--color-success);
  }

  .scenario-name {
    color: var(--color-text);
  }

  .feature-name {
    font-size: 10px;
    color: var(--color-muted);
  }

  .msg {
    color: var(--color-muted);
    max-width: 280px;
    word-break: break-word;
  }

  .at {
    color: var(--color-muted);
    white-space: nowrap;
    font-size: 10px;
  }
</style>
