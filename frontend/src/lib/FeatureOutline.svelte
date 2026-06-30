<script lang="ts">
  import { flattenFeatureSymbols, parseFeatureSymbols } from './gherkinDocumentSymbols'

  export let text = ''
  export let currentLine = 1
  export let onGoto: (line: number) => void = () => {}

  $: rows = flattenFeatureSymbols(parseFeatureSymbols(text))
</script>

{#if rows.length === 0}
  <p class="feature-outline-empty">Нет структуры сценария</p>
{:else}
  <ul class="feature-outline" role="tree" aria-label="Структура сценария">
    {#each rows as row (row.kind + row.line + row.name)}
      <li
        role="treeitem"
        class="feature-outline-item"
        class:active={row.line === currentLine}
        class:kind-feature={row.kind === 'feature'}
        class:kind-scenario={row.kind === 'scenario' || row.kind === 'outline'}
        class:kind-step={row.kind === 'step'}
        aria-selected={row.line === currentLine}
        style={`padding-left: ${8 + row.depth * 12}px`}
      >
        <button type="button" class="feature-outline-btn" on:click={() => onGoto(row.line)}>
          {#if row.detail && row.kind === 'step'}
            <span class="feature-outline-kw">{row.detail}</span>
          {/if}
          <span class="feature-outline-name">{row.name}</span>
          <span class="feature-outline-line">:{row.line}</span>
        </button>
      </li>
    {/each}
  </ul>
{/if}

<style>
  .feature-outline {
    list-style: none;
    margin: 0;
    padding: 4px 0;
  }

  .feature-outline-empty {
    margin: 0;
    padding: 8px 12px;
    font-size: 11px;
    color: var(--color-muted);
  }

  .feature-outline-item {
    margin: 0;
  }

  .feature-outline-btn {
    width: 100%;
    border: none;
    background: transparent;
    color: var(--color-text);
    text-align: left;
    font-size: 11px;
    padding: 2px 8px 2px 0;
    cursor: pointer;
    display: flex;
    gap: 4px;
    align-items: baseline;
  }

  .feature-outline-btn:hover {
    background: var(--color-hover);
  }

  .feature-outline-item.active .feature-outline-btn {
    background: color-mix(in srgb, var(--color-primary) 18%, transparent);
  }

  .feature-outline-kw {
    color: var(--color-accent);
    flex-shrink: 0;
  }

  .feature-outline-name {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
  }

  .feature-outline-line {
    color: var(--color-muted);
    flex-shrink: 0;
    font-size: 10px;
  }

  .kind-feature .feature-outline-name {
    font-weight: 600;
  }

  .kind-scenario .feature-outline-name {
    font-weight: 500;
  }
</style>
