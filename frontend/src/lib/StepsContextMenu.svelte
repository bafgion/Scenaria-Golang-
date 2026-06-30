<script lang="ts">
  export let x = 0
  export let y = 0
  export let onRunFrom: () => void = () => {}
  export let onRunTo: () => void = () => {}
  export let onDryRunFrom: () => void = () => {}
  export let onDryRunTo: () => void = () => {}
  export let onGoto: () => void = () => {}
  export let onHelp: () => void = () => {}
  export let onClose: () => void = () => {}

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="context-menu-backdrop" role="presentation" on:click={onClose} on:contextmenu|preventDefault={onClose}>
  <div class="context-menu" style="left: {x}px; top: {y}px" role="menu" tabindex="0" on:click|stopPropagation on:keydown|stopPropagation>
    <button type="button" on:click={onRunFrom}>Запустить с шага</button>
    <button type="button" on:click={onRunTo}>Запустить до шага</button>
    <button type="button" on:click={onDryRunFrom}>Dry-run с шага</button>
    <button type="button" on:click={onDryRunTo}>Dry-run до шага</button>
    <button type="button" on:click={onGoto}>Перейти к строке</button>
    <button type="button" on:click={onHelp}>Справка по шагу</button>
  </div>
</div>
