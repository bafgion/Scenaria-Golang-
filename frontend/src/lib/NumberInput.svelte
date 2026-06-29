<script lang="ts">
  export let value = 0
  export let min: number | undefined = undefined
  export let max: number | undefined = undefined
  export let step = 1
  export let width = '72px'
  export let disabled = false
  export let inputId: string | undefined = undefined

  function clamp(n: number): number {
    if (!Number.isFinite(n)) return min ?? 0
    let v = n
    if (min !== undefined) v = Math.max(min, v)
    if (max !== undefined) v = Math.min(max, v)
    return v
  }

  function stepUp() {
    value = clamp(value + step)
  }

  function stepDown() {
    value = clamp(value - step)
  }

  function onInput(e: Event) {
    const raw = (e.currentTarget as HTMLInputElement).value
    if (raw === '' || raw === '-') return
    value = clamp(Number(raw))
  }
</script>

<div class="num-input-wrap" class:disabled style="--num-width: {width}">
  <input
    class="num-input-field"
    type="number"
    id={inputId}
    bind:value
    {min}
    {max}
    {step}
    {disabled}
    on:change={onInput}
    on:blur={() => (value = clamp(value))}
  />
  <div class="num-input-spin" aria-hidden="true">
    <button type="button" tabindex="-1" title="Увеличить" {disabled} on:click={stepUp}>
      <svg viewBox="0 0 10 6" width="10" height="6" aria-hidden="true"><path fill="currentColor" d="M1 5 L5 1 L9 5 Z" /></svg>
    </button>
    <button type="button" tabindex="-1" title="Уменьшить" {disabled} on:click={stepDown}>
      <svg viewBox="0 0 10 6" width="10" height="6" aria-hidden="true"><path fill="currentColor" d="M1 1 L5 5 L9 1 Z" /></svg>
    </button>
  </div>
</div>

<style>
  .num-input-wrap {
    display: inline-flex;
    align-items: stretch;
    min-width: var(--num-width, 72px);
    border: 1px solid var(--color-border);
    border-radius: 4px;
    background: var(--color-input);
    overflow: hidden;
  }

  .num-input-field {
    flex: 1;
    min-width: 0;
    width: var(--num-width, 72px);
    padding: 5px 8px;
    border: none;
    background: transparent;
    color: var(--color-text);
    font-size: 13px;
    font-family: inherit;
    appearance: textfield;
    -moz-appearance: textfield;
  }

  .num-input-field:focus {
    outline: none;
  }

  .num-input-field::-webkit-inner-spin-button,
  .num-input-field::-webkit-outer-spin-button {
    -webkit-appearance: none;
    margin: 0;
  }

  .num-input-spin {
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
    border-left: 1px solid var(--color-border);
    background: var(--color-toolbar);
  }

  .num-input-spin button {
    display: flex;
    align-items: center;
    justify-content: center;
    flex: 1;
    min-height: 14px;
    width: 22px;
    padding: 0;
    border: none;
    background: transparent;
    color: var(--color-muted);
    cursor: pointer;
  }

  .num-input-spin button:hover {
    background: var(--color-selected);
    color: #fff;
  }

  .num-input-spin button + button {
    border-top: 1px solid var(--color-border);
  }

  .num-input-wrap.disabled {
    opacity: 0.45;
    pointer-events: none;
  }

  .num-input-spin button:disabled {
    cursor: default;
  }
</style>
