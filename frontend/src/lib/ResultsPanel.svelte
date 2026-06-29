<script lang="ts">
  import { gui } from '../../wailsjs/go/models'

  export let entries: gui.RunResultEntry[] = []
  export let artifacts: gui.ProjectArtifacts = new gui.ProjectArtifacts()
  export let onRerun: () => void = () => {}
  export let onOpenFolder: (path: string) => void = () => {}

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
</script>

<div class="results-panel">
  <div class="results-toolbar">
    <div class="artifact-btns">
      {#if artifacts.allureDir}
        <button type="button" on:click={() => onOpenFolder(artifacts.allureDir)}>Allure</button>
      {/if}
      {#if artifacts.htmlReport}
        <button type="button" on:click={() => onOpenFolder(artifacts.htmlReport)}>HTML-отчёт</button>
      {/if}
      {#if artifacts.tracesDir}
        <button type="button" on:click={() => onOpenFolder(artifacts.tracesDir)}>Trace</button>
      {/if}
      {#if artifacts.videosDir}
        <button type="button" on:click={() => onOpenFolder(artifacts.videosDir)}>Video</button>
      {/if}
    </div>
    <button type="button" class="rerun" on:click={onRerun}>Перезапустить упавшие</button>
  </div>
  {#if entries.length === 0}
    <p class="empty">Результаты прогона появятся здесь после запуска тестов</p>
  {:else}
    <table class="results-table">
      <thead>
        <tr>
          <th>Сценарий</th>
          <th>Результат</th>
          <th>Сообщение</th>
          <th>Время</th>
        </tr>
      </thead>
      <tbody>
        {#each entries as entry}
          {@const parts = splitPath(entry.path)}
          <tr class:failed={!entry.success}>
            <td>
              <div class="scenario-name">{parts.scenario || basename(parts.feature)}</div>
              <div class="feature-name">{basename(parts.feature)}</div>
            </td>
            <td class="status">{entry.success ? '✓ OK' : '✗ FAIL'}</td>
            <td class="msg">{entry.message || '—'}</td>
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

  .results-toolbar button {
    padding: 4px 10px;
    font-size: 11px;
    background: var(--color-input);
    border: 1px solid var(--color-border);
    border-radius: 2px;
    color: var(--color-text);
  }

  .rerun {
    margin-left: auto;
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

  tr.failed td.status {
    color: var(--color-error);
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
