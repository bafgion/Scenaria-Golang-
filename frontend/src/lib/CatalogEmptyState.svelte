<script lang="ts">
  import { onMount } from 'svelte'
  import { catalogEmptyIcon, type CatalogEmptyKind } from './icons'

  export let title = ''
  export let hint = ''
  export let kind: CatalogEmptyKind = 'no_project'

  let rootEl: HTMLDivElement
  let compact = false

  $: iconHtml = catalogEmptyIcon(kind)
  $: displayHint = compact && hint ? hint.split('\n')[0] : hint

  onMount(() => {
    const ro = new ResizeObserver(() => {
      compact = (rootEl?.clientWidth ?? 0) < 210
    })
    if (rootEl) ro.observe(rootEl)
    return () => ro.disconnect()
  })
</script>

<div class="catalog-empty-state" bind:this={rootEl}>
  <div class="catalog-empty-card">
    <div class="catalog-empty-icon" aria-hidden="true">
      {@html iconHtml}
    </div>
    <p class="catalog-empty-title">{title}</p>
    <p class="catalog-empty-hint">{displayHint}</p>
  </div>
</div>

<style>
  .catalog-empty-state {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 12px 8px;
    min-height: 0;
  }

  .catalog-empty-card {
    max-width: calc(100% - 12px);
    width: 100%;
    padding: 16px 14px;
    text-align: center;
    background: #2a2a2b;
    border: 1px solid var(--color-border);
    border-radius: 6px;
    box-sizing: border-box;
  }

  .catalog-empty-icon {
    width: 52px;
    height: 52px;
    margin: 0 auto 8px;
    display: grid;
    place-items: center;
    color: var(--color-text);
    background: var(--color-input);
    border: 1px solid var(--color-border);
    border-radius: 26px;
    box-sizing: border-box;
  }

  .catalog-empty-icon :global(svg) {
    width: 28px;
    height: 28px;
  }

  .catalog-empty-title {
    margin: 0;
    font-size: 12px;
    font-weight: 600;
    color: var(--color-text);
    line-height: 1.35;
  }

  .catalog-empty-hint {
    margin: 8px 0 0;
    font-size: 11px;
    color: var(--color-muted);
    line-height: 1.3;
    white-space: pre-line;
  }
</style>
