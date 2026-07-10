<script lang="ts">
  import { createTranslator, locale } from './i18n'
  import type { gui } from '../../wailsjs/go/models'
  import { flakyLabel, isFlakyPath, type FlakyScenarioStat } from './flakyMetrics'
  import { formatRunResultStatus } from './resultStatus'

  export let entries: gui.RunResultEntry[] = []
  export let flakyByPath: Map<string, FlakyScenarioStat> = new Map()
  export let flakyStepByPath: Map<string, string> = new Map()
  export let onOpenFeature: (path: string) => void = () => {}
  export let onRerunFailed: () => void = () => {}
  export let onClose: () => void = () => {}

  $: tr = createTranslator($locale)

  let filter = 'all'
  let query = ''

  $: filtered = entries.filter((entry) => {
    const status = formatRunResultStatus(entry, { flaky: isFlakyPath(flakyByPath, entry.path) })
    if (filter === 'flaky' && !isFlakyPath(flakyByPath, entry.path)) return false
    if (filter === 'failed' && status.tone !== 'error') return false
    if (filter === 'passed' && !(status.key === 'passed' || status.key === 'passedWithRetries')) return false
    const hay = `${entry.path} ${entry.message} ${entry.runner}`.toLowerCase()
    return !query.trim() || hay.includes(query.trim().toLowerCase())
  })

  function splitPath(path: string): { feature: string; scenario: string } {
    const idx = path.indexOf('::')
    if (idx < 0) return { feature: path, scenario: '' }
    return { feature: path.slice(0, idx), scenario: path.slice(idx + 2) }
  }

  function basename(p: string): string {
    const parts = p.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || p
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

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="palette-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="palette run-history" role="dialog" aria-modal="true" aria-label={tr('dialogs.runHistory.ariaLabel')} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{tr('dialogs.runHistory.title')}</h3>
    <div class="toolbar">
      <input class="search" bind:value={query} placeholder={tr('dialogs.runHistory.searchPlaceholder')} />
      <select bind:value={filter}>
        <option value="all">{tr('dialogs.runHistory.filterAll')}</option>
        <option value="failed">{tr('dialogs.runHistory.filterFailed')}</option>
        <option value="passed">{tr('dialogs.runHistory.filterPassed')}</option>
        <option value="flaky">{tr('dialogs.runHistory.filterFlaky')}</option>
      </select>
      <button type="button" on:click={onRerunFailed}>{tr('dialogs.runHistory.rerunFailed')}</button>
    </div>
    {#if filtered.length === 0}
      <p class="empty">{tr('dialogs.runHistory.empty')}</p>
    {:else}
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>{tr('dialogs.runHistory.colScenario')}</th>
              <th>{tr('dialogs.runHistory.colResult')}</th>
              <th>{tr('dialogs.runHistory.colMessage')}</th>
              <th>{tr('dialogs.runHistory.colTime')}</th>
            </tr>
          </thead>
          <tbody>
            {#each filtered as entry}
              {@const parts = splitPath(entry.path)}
              {@const flakyStat = flakyByPath.get(entry.path)}
              {@const stepHint = flakyStepByPath.get(entry.path)}
              {@const status = formatRunResultStatus(entry, { flaky: flakyStat?.flaky })}
              <tr class:failed={status.tone === 'error'} class:flaky={flakyStat?.flaky} on:dblclick={() => openEntry(entry)} title={tr('dialogs.runHistory.openFeatureTitle')}>
                <td>
                  <div class="scenario">{parts.scenario || basename(parts.feature)}</div>
                  <div class="feature">{basename(parts.feature)}</div>
                  {#if flakyStat?.flaky}
                    <div class="flaky-tag">{flakyLabel(flakyStat)}</div>
                  {/if}
                  {#if stepHint}
                    <div class="step-flaky">{stepHint}</div>
                  {/if}
                </td>
                <td>
                  <span class={`status-badge status-badge--${status.key} status-badge--${status.tone}`}>
                    {tr(`results.status.${status.key}`)}
                  </span>
                </td>
                <td class="msg">{entry.message || tr('dialogs.common.notFound')}</td>
                <td class="at">{formatAt(entry.at)}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
    <div class="actions">
      <button type="button" on:click={onClose}>{tr('dialogs.common.close')}</button>
    </div>
  </div>
</div>

<style>
  .run-history {
    width: min(860px, 96vw);
    max-height: 82vh;
    display: flex;
    flex-direction: column;
  }

  .toolbar {
    display: flex;
    gap: 8px;
    margin-bottom: 10px;
    flex-wrap: wrap;
  }

  .search {
    flex: 1;
    min-width: 180px;
  }

  .table-wrap {
    overflow: auto;
    flex: 1;
    min-height: 200px;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-size: 12px;
  }

  th,
  td {
    padding: 6px 8px;
    border-bottom: 1px solid var(--color-divider);
    text-align: left;
    vertical-align: top;
  }

  th {
    position: sticky;
    top: 0;
    background: var(--color-toolbar);
    color: var(--color-muted);
  }

  tr {
    cursor: default;
  }

  tr.failed td:nth-child(2) {
    color: var(--color-error);
  }

  .status-badge {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 2px 8px;
    border-radius: 999px;
    font-size: 10px;
    font-weight: 600;
    line-height: 1.3;
    border: 1px solid transparent;
    white-space: nowrap;
  }

  .status-badge--neutral {
    color: var(--color-muted);
    background: var(--color-input);
    border-color: var(--color-border);
  }

  .status-badge--success {
    color: var(--color-success);
    background: rgba(34, 197, 94, 0.12);
    border-color: rgba(34, 197, 94, 0.2);
  }

  .status-badge--warning {
    color: var(--color-warning, #d4a017);
    background: rgba(212, 160, 23, 0.14);
    border-color: rgba(212, 160, 23, 0.24);
  }

  .status-badge--error {
    color: var(--color-error);
    background: rgba(239, 68, 68, 0.12);
    border-color: rgba(239, 68, 68, 0.2);
  }

  tr:not(.failed) td:nth-child(2) {
    color: var(--color-success);
  }

  tr.flaky td:first-child {
    border-left: 2px solid var(--color-warning, #d4a017);
  }

  .scenario {
    color: var(--color-text);
  }

  .feature {
    font-size: 10px;
    color: var(--color-muted);
  }

  .flaky-tag {
    font-size: 10px;
    color: var(--color-warning, #d4a017);
    margin-top: 2px;
  }

  .step-flaky {
    font-size: 10px;
    color: var(--color-muted);
  }

  .msg {
    color: var(--color-muted);
    max-width: 320px;
    word-break: break-word;
  }

  .at {
    color: var(--color-muted);
    white-space: nowrap;
    font-size: 10px;
  }

  .empty {
    margin: 12px 0;
    color: var(--color-muted);
    font-size: 12px;
  }

  .actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 12px;
  }
</style>
