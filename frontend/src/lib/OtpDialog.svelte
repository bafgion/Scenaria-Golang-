<script lang="ts">
  import { onMount, tick } from 'svelte'
  import { createTranslator, locale } from './i18n'

  export let email = ''
  export let onSubmit: (code: string) => void = () => {}
  export let onCancel: () => void = () => {}

  $: tr = createTranslator($locale)

  let code = ''
  let codeInput: HTMLInputElement | null = null

  onMount(() => {
    void tick().then(() => codeInput?.focus())
  })

  function submit() {
    onSubmit(code.trim())
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel()
    if (e.key === 'Enter') submit()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="modal-backdrop modal-layer-top" role="presentation" on:click={onCancel}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal" role="dialog" aria-modal="true" aria-label={tr('dialogs.otp.ariaLabel')} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{tr('dialogs.otp.title')}</h3>
    {#if email}<p class="hint">{email}</p>{/if}
    <input bind:this={codeInput} bind:value={code} placeholder={tr('dialogs.otp.placeholder')} />
    <div class="modal-actions">
      <button type="button" class="primary" on:click={submit}>{tr('dialogs.common.ok')}</button>
      <button type="button" on:click={onCancel}>{tr('dialogs.common.cancel')}</button>
    </div>
  </div>
</div>

<style>
  .hint {
    font-size: 12px;
    color: var(--color-muted);
    margin: 0 0 8px;
  }
</style>
