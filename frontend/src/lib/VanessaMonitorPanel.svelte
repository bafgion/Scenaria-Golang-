<script lang="ts">
  import { gui } from '../../wailsjs/go/models'

  export let snapshot: gui.VanessaRunSnapshotDTO = new gui.VanessaRunSnapshotDTO()
  export let running = false
  export let onClose: () => void = () => {}

  function basename(path: string): string {
    const parts = (path || '').replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || path
  }

  $: progress =
    snapshot.totalPlanned > 0
      ? Math.min(100, Math.round((snapshot.completedCases / snapshot.totalPlanned) * 100))
      : 0
</script>

<div class="vanessa-monitor" role="status">
  <div class="head">
    <strong>Vanessa {running ? '— выполняется' : '— завершено'}</strong>
    <button type="button" class="close" on:click={onClose}>×</button>
  </div>
  {#if snapshot.runDir}
    <p class="meta">Каталог: {basename(snapshot.runDir)}</p>
  {/if}
  {#if snapshot.currentScenario}
    <p class="current">Сейчас: {snapshot.currentScenario}</p>
  {/if}
  <div class="progress-wrap">
    <div class="progress-bar" style="width: {progress}%"></div>
  </div>
  <p class="counts">{snapshot.completedCases} / {snapshot.totalPlanned} сценариев</p>
  {#if snapshot.cases.length > 0}
    <table>
      <thead>
        <tr><th>Сценарий</th><th>Результат</th></tr>
      </thead>
      <tbody>
        {#each snapshot.cases as item}
          <tr class:fail={!item.success}>
            <td>{item.name}</td>
            <td>{item.success ? '✓' : '✗'} {item.message}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {:else if running}
    <p class="empty">Ожидание JUnit…</p>
  {:else}
    <p class="empty">Нет данных JUnit</p>
  {/if}
</div>

<style>
  .vanessa-monitor {
    margin: 8px 12px;
    padding: 10px 12px;
    border: 1px solid var(--color-border);
    border-radius: 4px;
    background: var(--color-panel, var(--color-input));
    max-height: 280px;
    overflow: auto;
  }

  .head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 6px;
    font-size: 12px;
  }

  .close {
    border: none;
    background: transparent;
    font-size: 16px;
    color: var(--color-muted);
    cursor: pointer;
  }

  .meta, .current, .counts, .empty {
    margin: 0 0 6px;
    font-size: 11px;
    color: var(--color-muted);
  }

  .current {
    color: var(--color-text);
  }

  .progress-wrap {
    height: 6px;
    background: var(--color-border);
    border-radius: 3px;
    overflow: hidden;
    margin-bottom: 6px;
  }

  .progress-bar {
    height: 100%;
    background: var(--color-accent);
    transition: width 0.3s ease;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-size: 11px;
  }

  th, td {
    padding: 4px 6px;
    border-bottom: 1px solid var(--color-divider);
    text-align: left;
  }

  tr.fail td {
    color: var(--color-error);
  }
</style>
