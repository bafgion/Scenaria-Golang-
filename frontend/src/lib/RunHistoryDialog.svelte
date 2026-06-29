<script lang="ts">
  import type { gui } from '../../wailsjs/go/models'

  export let entries: gui.RunResultEntry[] = []
  export let onOpenFeature: (path: string) => void = () => {}
  export let onRerunFailed: () => void = () => {}
  export let onClose: () => void = () => {}

  let filter = 'all'
  let query = ''

  $: filtered = entries.filter((entry) => {
    if (filter === 'failed' && entry.success) return false
    if (filter === 'passed' && !entry.success) return false
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

<div class="palette-backdrop" role="presentation" on:click={onClose}>
  <div class="palette run-history" role="dialog" aria-label="История запусков" on:click|stopPropagation>
    <h3>История запусков</h3>
    <div class="toolbar">
      <input class="search" bind:value={query} placeholder="Поиск по сценарию или сообщению…" />
      <select bind:value={filter}>
        <option value="all">Все</option>
        <option value="failed">Только упавшие</option>
        <option value="passed">Только успешные</option>
      </select>
      <button type="button" on:click={onRerunFailed}>Перезапустить упавшие</button>
    </div>
    {#if filtered.length === 0}
      <p class="empty">Нет записей для отображения</p>
    {:else}
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Сценарий</th>
              <th>Результат</th>
              <th>Сообщение</th>
              <th>Время</th>
            </tr>
          </thead>
          <tbody>
            {#each filtered as entry}
              {@const parts = splitPath(entry.path)}
              <tr class:failed={!entry.success} on:dblclick={() => openEntry(entry)} title="Двойной клик — открыть feature">
                <td>
                  <div class="scenario">{parts.scenario || basename(parts.feature)}</div>
                  <div class="feature">{basename(parts.feature)}</div>
                </td>
                <td>{entry.success ? 'OK' : 'FAIL'}</td>
                <td class="msg">{entry.message || '—'}</td>
                <td class="at">{formatAt(entry.at)}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
    <div class="actions">
      <button type="button" on:click={onClose}>Закрыть</button>
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
    padding: 6px 8px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: var(--color-text);
  }

  select {
    padding: 6px 8px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: var(--color-text);
  }

  .toolbar button {
    padding: 6px 10px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: var(--color-text);
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

  tr:not(.failed) td:nth-child(2) {
    color: var(--color-success);
  }

  .scenario {
    color: var(--color-text);
  }

  .feature {
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
