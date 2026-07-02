<script lang="ts">
  import { gui } from '../../wailsjs/go/models'

  export let entry: gui.RunResultEntry | null = null
  export let onGotoFailedStep: (entry: gui.RunResultEntry) => void = () => {}

  function splitPath(path: string): { feature: string; scenario: string } {
    const idx = path.indexOf('::')
    if (idx < 0) return { feature: path, scenario: '' }
    return { feature: path.slice(0, idx), scenario: path.slice(idx + 2) }
  }
</script>

<div class="error-panel">
  {#if !entry}
    <p class="empty">Нет ошибок последнего прогона</p>
  {:else}
    <h4>Ошибка теста</h4>
    {@const parts = splitPath(entry.path)}
    <dl>
      <dt>Сценарий</dt>
      <dd>{parts.scenario || '—'}</dd>
      <dt>Файл</dt>
      <dd>{parts.feature}</dd>
      {#if entry.failed_step != null && entry.failed_step >= 0}
        <dt>Упавший шаг</dt>
        <dd>
          <button type="button" class="goto-step" on:click={() => onGotoFailedStep(entry)}>
            Перейти к шагу {entry.failed_step + 1}
          </button>
        </dd>
      {/if}
      <dt>Сообщение</dt>
      <dd class="msg">{entry.message || '—'}</dd>
    </dl>
  {/if}
</div>

<style>
  .error-panel {
    padding: 12px 16px;
    height: 100%;
    overflow: auto;
  }

  h4 {
    margin: 0 0 12px;
    font-size: 13px;
    color: var(--color-error);
  }

  dl {
    margin: 0;
    font-size: 12px;
  }

  dt {
    color: var(--color-muted);
    margin-top: 8px;
  }

  dd {
    margin: 2px 0 0;
    color: var(--color-text);
  }

  .goto-step {
    font-size: 11px;
    padding: 2px 8px;
    border-radius: 4px;
    border: 1px solid var(--color-border);
    background: var(--color-input);
    color: var(--color-text);
    cursor: pointer;
  }

  .goto-step:hover {
    border-color: var(--color-primary);
  }

  .msg {
    font-family: var(--font-mono);
    white-space: pre-wrap;
    color: var(--color-error);
  }

  .empty {
    color: var(--color-muted);
    font-size: 12px;
  }
</style>
