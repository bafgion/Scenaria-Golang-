<script lang="ts">
  export let email = ''
  export let onSubmit: (code: string) => void = () => {}
  export let onCancel: () => void = () => {}

  let code = ''

  function submit() {
    onSubmit(code.trim())
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel()
    if (e.key === 'Enter') submit()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
<div class="modal-backdrop" role="presentation">
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal" role="dialog" aria-modal="true" aria-label="Код из почты" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>Код из почты</h3>
    {#if email}<p class="hint">{email}</p>{/if}
    <input bind:value={code} placeholder="123456" autofocus />
    <div class="modal-actions">
      <button type="button" class="primary" on:click={submit}>OK</button>
      <button type="button" on:click={onCancel}>Отмена</button>
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
