<script lang="ts">
  import type { gui } from '../../wailsjs/go/models'

  export let selector = ''
  export let choices: gui.PickerStepChoice[] = []
  export let onInsert: (text: string) => void = () => {}
  export let onClose: () => void = () => {}

  let selected = 0

  $: preview = choices[selected]?.preview || '—'
  $: description = choices[selected]?.description || ''

  function confirm() {
    const choice = choices[selected]
    if (!choice) return
    if (choice.label === 'Только селектор') {
      onInsert(choice.preview)
    } else {
      onInsert((choice.preview.endsWith('\n') ? choice.preview : choice.preview + '\n'))
    }
    onClose()
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      e.preventDefault()
      onClose()
    }
    if (e.key === 'Enter') {
      e.preventDefault()
      confirm()
    }
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="palette-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="palette picker-step" role="dialog" aria-modal="true" aria-label="Шаг для элемента" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>Элемент выбран — укажите шаг</h3>
    <p class="selector-preview" title={selector}>{selector.length > 120 ? selector.slice(0, 117) + '…' : selector}</p>
    <div class="picker-body">
      <ul class="choice-list">
        {#each choices as choice, i}
          <li>
            <button type="button" class:selected={i === selected} on:click={() => (selected = i)}>
              {choice.label}
            </button>
          </li>
        {/each}
      </ul>
      <div class="preview-pane">
        <div class="caption">Пример в сценарии</div>
        <pre>{preview}</pre>
        <p class="hint">{description}</p>
      </div>
    </div>
    <div class="actions">
      <button type="button" class="primary" on:click={confirm}>Вставить</button>
      <button type="button" on:click={onClose}>Отмена</button>
    </div>
  </div>
</div>

<style>
  .picker-step {
    width: min(640px, 94vw);
  }

  .selector-preview {
    margin: 0 0 12px;
    font-family: var(--font-mono, monospace);
    font-size: 11px;
    color: var(--color-muted);
    word-break: break-all;
  }

  .picker-body {
    display: grid;
    grid-template-columns: 180px 1fr;
    gap: 12px;
    min-height: 220px;
  }

  .choice-list {
    list-style: none;
    margin: 0;
    padding: 0;
    border: 1px solid var(--color-border);
    border-radius: 4px;
    overflow: auto;
  }

  .choice-list button {
    display: block;
    width: 100%;
    text-align: left;
    padding: 8px 10px;
    border: none;
    background: transparent;
    color: var(--color-text);
    font-size: 12px;
  }

  .choice-list button.selected,
  .choice-list button:hover {
    background: var(--color-input);
  }

  .preview-pane {
    display: flex;
    flex-direction: column;
    gap: 8px;
    min-width: 0;
  }

  .caption {
    font-size: 11px;
    color: var(--color-muted);
  }

  pre {
    margin: 0;
    padding: 10px;
    background: var(--color-input);
    border: 1px solid var(--color-border);
    border-radius: 4px;
    font-size: 12px;
    white-space: pre-wrap;
    word-break: break-word;
    flex: 1;
  }

  .hint {
    margin: 0;
    font-size: 11px;
    color: var(--color-muted);
  }

  .actions {
    display: flex;
    gap: 8px;
    justify-content: flex-end;
    margin-top: 12px;
  }
</style>
