<script lang="ts">
  export let browser = 'chromium'
  export let syntaxOnly = false
  export let scope: 'project' | 'current' = 'project'
  export let canValidateCurrent = false
  export let currentFileName = ''
  export let onConfirm: (payload: { browser: string; syntaxOnly: boolean; scope: 'project' | 'current' }) => void = () => {}
  export let onCancel: () => void = () => {}

  function confirm() {
    onConfirm({ browser, syntaxOnly, scope })
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onCancel}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal validate-dialog" role="dialog" aria-modal="true" aria-label="Проверка сценария" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>Проверка сценария</h3>
    <fieldset class="scope">
      <legend>Область</legend>
      <label class="check-row">
        <input type="radio" bind:group={scope} value="project" />
        Весь проект
      </label>
      <label class="check-row" class:disabled={!canValidateCurrent}>
        <input type="radio" bind:group={scope} value="current" disabled={!canValidateCurrent} />
        Текущий файл{#if currentFileName} — {currentFileName}{/if}
      </label>
    </fieldset>
    <label class="check-row">
      <input type="checkbox" bind:checked={syntaxOnly} />
      Только синтаксис (без браузера)
    </label>
    <label>
      Браузер для проверки селекторов
      <select bind:value={browser} disabled={syntaxOnly}>
        <option value="chromium">chromium</option>
        <option value="firefox">firefox</option>
        <option value="webkit">webkit</option>
      </select>
    </label>
    <p class="hint">Проверяет шаги DSL и при необходимости селекторы в браузере.</p>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={confirm}>Проверить</button>
      <button type="button" on:click={onCancel}>Отмена</button>
    </div>
  </div>
</div>

<style>
  h3 {
    margin: 0 0 12px;
    font-size: 14px;
  }

  fieldset.scope {
    margin: 0 0 12px;
    padding: 8px 10px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
  }

  legend {
    font-size: 11px;
    color: var(--color-muted);
    padding: 0 4px;
  }

  label {
    display: grid;
    gap: 4px;
    margin-bottom: 10px;
    font-size: 11px;
    color: var(--color-muted);
  }

  .check-row {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 12px;
    color: var(--color-text);
    margin-bottom: 6px;
  }

  .check-row.disabled {
    opacity: 0.5;
  }

  .hint {
    margin: 0 0 12px;
    font-size: 11px;
    color: var(--color-muted);
    line-height: 1.4;
  }
</style>
