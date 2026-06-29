<script lang="ts">
  import { gui } from '../../wailsjs/go/models'

  export let issues: gui.ValidationIssue[] = []
  export let hint = ''
  export let onGotoLine: (line: number) => void = () => {}
</script>

<div class="validate-panel">
  {#if issues.length === 0}
    <p class="empty">{hint || 'Ошибок в шагах сценария нет'}</p>
  {:else}
    <table class="validate-table">
      <thead>
        <tr>
          <th>Строка</th>
          <th>Сообщение</th>
        </tr>
      </thead>
      <tbody>
        {#each issues as issue}
          <tr>
            <td>
              <button type="button" class="line-btn" on:click={() => onGotoLine(issue.line)}>{issue.line}</button>
            </td>
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

  .msg {
    color: var(--color-error);
    font-family: var(--font-mono);
    word-break: break-word;
  }
</style>
