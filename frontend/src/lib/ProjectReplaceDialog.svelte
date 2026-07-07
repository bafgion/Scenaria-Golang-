<script lang="ts">
  import { onMount, tick } from 'svelte'
  import { createTranslator, locale } from './i18n'

  export let findText = ''
  export let replaceText = ''
  export let caseSensitive = false
  export let busy = false
  export let onConfirm: () => void | Promise<void> = () => {}
  export let onClose: () => void = () => {}

  $: tr = createTranslator($locale)

  let findInput: HTMLInputElement | null = null

  onMount(() => {
    void tick().then(() => findInput?.focus())
  })

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      e.preventDefault()
      onClose()
    }
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="palette-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="palette find-replace" role="dialog" aria-modal="true" aria-label={tr('dialogs.project.replace.ariaLabel')} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{tr('dialogs.project.replace.title')}</h3>
    <p class="hint">{tr('dialogs.project.replace.hint')}</p>
    <label>{tr('dialogs.project.replace.find')} <input bind:this={findInput} bind:value={findText} /></label>
    <label>{tr('dialogs.project.replace.replace')} <input bind:value={replaceText} /></label>
    <label class="check-row">
      <input type="checkbox" bind:checked={caseSensitive} /> {tr('dialogs.project.replace.caseSensitive')}
    </label>
    <div class="actions">
      <button type="button" class="primary" disabled={busy || !findText} on:click={() => onConfirm()}>
        {busy ? tr('dialogs.project.replace.replacing') : tr('dialogs.project.replace.replaceAll')}
      </button>
      <button type="button" on:click={onClose}>{tr('dialogs.common.cancel')}</button>
    </div>
  </div>
</div>

<style>
  .find-replace .actions {
    margin-top: 4px;
  }
</style>
