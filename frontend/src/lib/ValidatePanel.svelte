<script lang="ts">
  import { createTranslator, locale } from './i18n'
  import { gui } from '../../wailsjs/go/models'

  export let issues: gui.ValidationIssue[] = []
  export let hint = ''
  export let cliLog = ''
  export let activeLine = 0
  export let onGotoLine: (line: number) => void = () => {}

  $: tr = createTranslator($locale)

  function statusLabel(status: string | undefined): string {
    switch (status) {
      case 'found':
        return tr('results.validate.found')
      case 'missing':
        return tr('results.validate.missing')
      case 'warning':
        return tr('results.validate.warning')
      default:
        return status || tr('dialogs.common.notFound')
    }
  }

  $: panelLimitation = issues.find((i) => i.limitation)?.limitation || ''
  $: panelMode = issues.find((i) => i.mode)?.mode || ''
  $: showMatches = issues.some((i) => (i.matchCount ?? 0) > 0)
  $: showAction = issues.some((i) => !!i.actionKind)

  function modeLabel(mode: string | undefined): string {
    if (mode === 'flow') return tr('results.validate.modeFlow')
    if (mode === 'static') return tr('results.validate.modeStatic')
    return mode || ''
  }
</script>

<div class="validate-panel">
  {#if cliLog}
    <pre class="cli-log">{cliLog}</pre>
  {/if}
  {#if panelLimitation}
    <p class="limitation-banner">
      {#if panelMode}
        <span class="mode-tag">{modeLabel(panelMode)}</span>
      {/if}
      {panelLimitation}
    </p>
  {/if}
  {#if issues.length === 0}
    <p class="empty">{hint || tr('results.validate.empty')}</p>
  {:else}
    <table class="validate-table">
      <thead>
        <tr>
          <th>{tr('results.validate.line')}</th>
          <th>{tr('results.validate.status')}</th>
          {#if showAction}<th>{tr('results.validate.action')}</th>{/if}
          <th>{tr('results.validate.selector')}</th>
          {#if showMatches}<th>{tr('results.validate.matches')}</th>{/if}
          <th>{tr('results.validate.message')}</th>
        </tr>
      </thead>
      <tbody>
        {#each issues as issue}
          <tr class:active={activeLine === issue.line} class:found={issue.status === 'found'} class:warning={issue.status === 'warning'}>
            <td>
              <button type="button" class="line-btn" on:click={() => onGotoLine(issue.line)}>{issue.line}</button>
            </td>
            <td><span class="status-badge" class:found={issue.status === 'found'} class:warning={issue.status === 'warning'} class:missing={issue.status === 'missing' || !issue.status}>{statusLabel(issue.status)}</span></td>
            {#if showAction}<td class="action">{issue.actionKind || '—'}</td>{/if}
            <td class="selector">{issue.selector || tr('dialogs.common.notFound')}</td>
            {#if showMatches}<td class="matches">{(issue.matchCount ?? 0) > 0 ? issue.matchCount : '—'}</td>{/if}
            <td class="msg">{issue.message}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>

<style>
  .validate-panel {
    height: 100%;
    overflow: auto;
    padding: 8px 12px;
  }

  .cli-log {
    margin: 0 0 10px;
    padding: 8px;
    font-family: Consolas, 'Courier New', monospace;
    font-size: 11px;
    line-height: 1.4;
    white-space: pre-wrap;
    background: var(--color-input);
    border: 1px solid var(--color-border);
    border-radius: 3px;
    max-height: 200px;
    overflow: auto;
  }

  .limitation-banner {
    margin: 0 0 10px;
    padding: 8px 10px;
    font-size: 11px;
    line-height: 1.45;
    color: var(--color-muted);
    background: color-mix(in srgb, var(--color-warning, #dcdcaa) 10%, var(--color-input));
    border: 1px solid var(--color-border);
    border-radius: 3px;
  }

  .mode-tag {
    display: inline-block;
    margin-right: 8px;
    padding: 1px 6px;
    border-radius: 10px;
    font-size: 10px;
    text-transform: uppercase;
    background: rgba(255, 255, 255, 0.06);
    color: var(--color-text);
  }

  .empty {
    font-size: 12px;
    color: var(--color-muted);
    margin: 8px 0;
  }

  .validate-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 12px;
  }

  .validate-table th {
    text-align: left;
    padding: 6px 8px;
    background: var(--color-toolbar);
    color: var(--color-muted);
    border-bottom: 1px solid var(--color-border);
    position: sticky;
    top: 0;
  }

  .validate-table td {
    padding: 6px 8px;
    border-bottom: 1px solid var(--color-divider);
    vertical-align: top;
  }

  .validate-table tr.active td {
    background: color-mix(in srgb, var(--color-primary) 14%, transparent);
  }

  .validate-table tr.found td.msg,
  .validate-table tr.found .selector {
    color: var(--color-success, #4ec9b0);
  }

  .validate-table tr.warning td.msg {
    color: var(--color-warning, #dcdcaa);
  }

  .line-btn {
    padding: 0 6px;
    min-width: 28px;
    background: var(--color-input);
    border: 1px solid var(--color-border);
    border-radius: 2px;
    color: var(--color-accent);
    font-family: var(--font-mono);
    cursor: pointer;
  }

  .selector {
    font-family: var(--font-mono);
    word-break: break-word;
    color: var(--color-muted);
    max-width: 220px;
  }

  .action,
  .matches {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--color-muted);
    white-space: nowrap;
  }

  .msg {
    color: var(--color-error);
    font-family: var(--font-mono);
    word-break: break-word;
  }

  .status-badge {
    display: inline-block;
    padding: 1px 6px;
    border-radius: 10px;
    font-size: 10px;
    text-transform: uppercase;
    background: rgba(255, 255, 255, 0.06);
    color: var(--color-muted);
  }

  .status-badge.found {
    color: var(--color-success, #4ec9b0);
  }

  .status-badge.warning {
    color: var(--color-warning, #dcdcaa);
  }

  .status-badge.missing {
    color: var(--color-error);
  }
</style>
