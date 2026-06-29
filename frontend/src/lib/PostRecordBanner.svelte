<script lang="ts">
  import type { gui } from '../../wailsjs/go/models'

  export let path = ''
  export let stepCount = 0
  export let hints: gui.ScenarioHintDTO[] = []
  export let onValidate: () => void = () => {}
  export let onSave: () => void = () => {}
  export let onFixHover: () => void = () => {}
  export let onShowStep: (index: number) => void = () => {}
  export let onFixHint: (hint: gui.ScenarioHintDTO) => void = () => {}
  export let onDismissHint: (hint: gui.ScenarioHintDTO) => void = () => {}
  export let onClose: () => void = () => {}

  let selected = 0

  $: menuHoverCount = hints.filter((h) => h.id === 'menu_hover').length
  $: visibleHints = hints
  $: if (selected >= visibleHints.length) selected = Math.max(0, visibleHints.length - 1)
  $: current = visibleHints[selected]

  function basename(p: string): string {
    const parts = p.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || p
  }

  function hintLabel(hint: gui.ScenarioHintDTO): string {
    const prefix = hint.severity === 'warning' ? '⚠ ' : 'ℹ '
    const stepNo = hint.stepIndex + 1
    return `${prefix}${hint.title}${stepNo > 0 ? ` (шаг ${stepNo})` : ''}`
  }
</script>

<div class="post-record-banner" role="status">
  <div class="top">
    <div class="summary">
      {#if visibleHints.length}
        Записано шагов: {stepCount} · подсказок: {visibleHints.length}
        <span class="file"> — {basename(path)}</span>
      {:else}
        Записано шагов: {stepCount} — {basename(path)}
      {/if}
    </div>
    {#if menuHoverCount > 0}
      <span class="legacy-hint">Похоже на hover-меню: {menuHoverCount} клик(ов) без «навожу»</span>
      <button type="button" on:click={onFixHover}>Добавить наведение</button>
    {/if}
    <div class="actions">
      <button type="button" class="primary" on:click={onValidate}>Проверить</button>
      <button type="button" on:click={onSave}>Сохранить</button>
      <button type="button" class="dismiss" on:click={onClose}>Закрыть</button>
    </div>
  </div>

  {#if visibleHints.length > 0}
    <div class="hints">
      <ul>
        {#each visibleHints as hint, i}
          <li>
            <button type="button" class:selected={i === selected} on:click={() => (selected = i)}>
              {hintLabel(hint)}
            </button>
          </li>
        {/each}
      </ul>
      {#if current}
        <div class="hint-actions">
          <button type="button" on:click={() => onShowStep(current.stepIndex + 1)}>Показать шаг</button>
          {#if current.autoFixable}
            <button type="button" on:click={() => onFixHint(current)}>Исправить</button>
          {/if}
          <button type="button" on:click={() => onDismissHint(current)}>Игнорировать</button>
        </div>
      {/if}
    </div>
  {/if}
</div>

<style>
  .post-record-banner {
    background: rgba(137, 209, 133, 0.08);
    border-bottom: 1px solid var(--color-success);
    padding: 6px 10px;
    font-size: 12px;
  }

  .top {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 8px;
  }

  .summary {
    flex: 1;
    min-width: 200px;
    color: var(--color-text);
  }

  .file {
    color: var(--color-muted);
  }

  .legacy-hint {
    color: var(--color-warning, #e5c07b);
    font-size: 11px;
  }

  .actions {
    display: flex;
    gap: 6px;
    flex-wrap: wrap;
  }

  button {
    padding: 4px 10px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    background: var(--color-input);
    color: var(--color-text);
    font-size: 12px;
    cursor: pointer;
  }

  button.primary {
    background: var(--color-accent);
    color: var(--color-on-accent, #fff);
    border-color: var(--color-accent);
  }

  button.dismiss {
    color: var(--color-muted);
  }

  .hints {
    margin-top: 6px;
    border-top: 1px solid var(--color-divider);
    padding-top: 6px;
  }

  .hints ul {
    list-style: none;
    margin: 0 0 6px;
    padding: 0;
    max-height: 100px;
    overflow: auto;
  }

  .hints li button {
    width: 100%;
    text-align: left;
    border: none;
    background: transparent;
    padding: 4px 6px;
  }

  .hints li button:hover,
  .hints li button.selected {
    background: var(--color-input);
  }

  .hint-actions {
    display: flex;
    gap: 6px;
    flex-wrap: wrap;
  }
</style>
