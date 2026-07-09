<script lang="ts">
  import { createTranslator, locale } from './i18n'
  import type { gui } from '../../wailsjs/go/models'

  export let selector = ''
  export let choices: gui.PickerStepChoice[] = []
  export let candidates: gui.SelectorCandidate[] = []
  export let warnings: string[] = []
  export let suggestedAction = ''
  export let suggestedChoice = 0
  export let onInsert: (text: string) => void = () => {}
  export let onClose: () => void = () => {}
  export let onSelectorChange: (selector: string) => void = () => {}

  $: tr = createTranslator($locale)

  let selected = 0
  let activeSelector = selector
  let userPickedAction = false

  $: activeSelector = selector
  $: if (selector) userPickedAction = false
  $: if (!userPickedAction && choices.length > 0) {
    const idx = suggestedChoice
    selected = idx >= 0 && idx < choices.length ? idx : 0
  }
  $: preview = choices[selected]?.preview || tr('dialogs.common.notFound')
  $: description = choices[selected]?.description || ''
  $: selectorOnlyLabel = tr('dialogs.pickerStep.selectorOnly')

  function selectCandidate(cand: gui.SelectorCandidate) {
    activeSelector = cand.selector
    onSelectorChange(cand.selector)
  }

  function confirm() {
    const choice = choices[selected]
    if (!choice) return
    if (choice.label === selectorOnlyLabel) {
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
<div class="palette-backdrop modal-layer-top" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="palette picker-step" role="dialog" aria-modal="true" aria-label={tr('dialogs.pickerStep.ariaLabel')} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{tr('dialogs.pickerStep.title')}</h3>
    <p class="selector-preview" title={activeSelector}>{activeSelector.length > 120 ? activeSelector.slice(0, 117) + '…' : activeSelector}</p>
    {#if warnings.length}
      <div class="warnings" role="status">
        {#each warnings as warning}
          <span class="warning-pill">{warning}</span>
        {/each}
      </div>
    {/if}
    {#if suggestedAction}
      <p class="suggested-action">Suggested: {suggestedAction}</p>
    {/if}
    <div class="picker-body">
      <div class="side-panels">
        {#if candidates.length > 1}
          <div class="candidate-panel">
            <div class="caption">Selectors</div>
            <ul class="choice-list">
              {#each candidates as cand, i}
                <li>
                  <button
                    type="button"
                    class:selected={cand.selector === activeSelector}
                    on:click={() => selectCandidate(cand)}
                    title={cand.warnings?.join(', ') || cand.strategy}
                  >
                    <span class="strategy">{cand.strategy}</span>
                    <span class="cand-text">{cand.selector}</span>
                  </button>
                </li>
              {/each}
            </ul>
          </div>
        {/if}
        <div class="action-panel">
          <div class="caption">{tr('dialogs.pickerStep.actions')}</div>
          <ul class="choice-list">
            {#each choices as choice, i}
              <li>
                <button type="button" class:selected={i === selected} on:click={() => { userPickedAction = true; selected = i }}>
                  {choice.label}
                </button>
              </li>
            {/each}
          </ul>
        </div>
      </div>
      <div class="preview-pane">
        <div class="caption">{tr('dialogs.pickerStep.scenarioExample')}</div>
        <pre>{preview}</pre>
        <p class="hint">{description}</p>
      </div>
    </div>
    <div class="actions">
      <button type="button" class="primary" on:click={confirm}>{tr('dialogs.pickerStep.insert')}</button>
      <button type="button" on:click={onClose}>{tr('dialogs.common.cancel')}</button>
    </div>
  </div>
</div>

<style>
  .picker-step {
    width: min(720px, 94vw);
  }

  .selector-preview {
    margin: 0 0 12px;
    font-family: var(--font-mono, monospace);
    font-size: 11px;
    color: var(--color-muted);
    word-break: break-all;
  }

  .warnings {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin-bottom: 8px;
  }

  .warning-pill {
    font-size: 10px;
    padding: 2px 8px;
    border-radius: 999px;
    background: rgba(255, 193, 7, 0.18);
    color: #c9a227;
  }

  .suggested-action {
    margin: 0 0 8px;
    font-size: 11px;
    color: var(--color-muted);
  }

  .picker-body {
    display: grid;
    grid-template-columns: 220px 1fr;
    gap: 12px;
    min-height: 220px;
  }

  .side-panels {
    display: flex;
    flex-direction: column;
    gap: 10px;
    min-width: 0;
  }

  .caption {
    font-size: 11px;
    color: var(--color-muted);
    margin-bottom: 4px;
  }

  .choice-list {
    list-style: none;
    margin: 0;
    padding: 0;
    border: 1px solid var(--color-border);
    border-radius: 4px;
    overflow: auto;
    max-height: 180px;
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
    background: var(--color-selected);
  }

  .candidate-panel .strategy {
    display: block;
    font-size: 10px;
    color: var(--color-muted);
    text-transform: uppercase;
  }

  .candidate-panel .cand-text {
    display: block;
    font-family: var(--font-mono, monospace);
    font-size: 10px;
    word-break: break-all;
  }

  .preview-pane {
    display: flex;
    flex-direction: column;
    gap: 8px;
    min-width: 0;
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
