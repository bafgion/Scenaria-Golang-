<script lang="ts">
  export let x = 0
  export let y = 0
  export let onRun: () => void = () => {}
  export let onDryRun: () => void = () => {}
  export let onOpen: () => void = () => {}
  export let onDuplicate: () => void = () => {}
  export let onRename: () => void = () => {}
  export let onMove: () => void = () => {}
  export let onDelete: () => void = () => {}
  export let onReveal: () => void = () => {}
  export let onClose: () => void = () => {}

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="context-menu-backdrop" role="presentation" on:click={onClose} on:contextmenu|preventDefault={onClose}>
  <div class="context-menu" style="left: {x}px; top: {y}px" role="menu" tabindex="0" on:click|stopPropagation on:keydown|stopPropagation>
    <button type="button" on:click={onRun}>Запустить</button>
    <button type="button" on:click={onDryRun}>Dry-run</button>
    <button type="button" on:click={onOpen}>Открыть</button>
    <button type="button" on:click={onDuplicate}>Дублировать</button>
    <button type="button" on:click={onRename}>Переименовать…</button>
    <button type="button" on:click={onMove}>Переместить…</button>
    <button type="button" on:click={onReveal}>Показать в проводнике</button>
    <button type="button" class="danger" on:click={onDelete}>Удалить</button>
  </div>
</div>
