<script lang="ts">
  import { onMount } from 'svelte'
  import { icons } from './icons'

  export let visible = false
  export let recording = false
  export let playing = false
  export let paused = false

  export let onRecord: () => void = () => {}
  export let onPause: () => void = () => {}
  export let onStop: () => void = () => {}
  export let onPicker: () => void = () => {}
  export let onFocusBrowser: () => void = () => {}

  $: pickerEnabled = (recording && paused) || (!recording && !playing)
  $: pauseEnabled = recording && !playing
  $: stopEnabled = recording || playing
  $: recordEnabled = !recording && !playing
  $: focusEnabled = recording || playing

  let overlay: HTMLDivElement
  let pos = { x: 0, y: 0 }
  let dragging = false
  let dragStart = { x: 0, y: 0 }
  let origin = { x: 0, y: 0 }

  $: title = recording
    ? paused
      ? '⏸ Пауза — можно выбрать элемент'
      : '● Идёт запись'
    : playing
      ? '▶ Идёт тест'
      : 'Scenaria — перетащите панель'

  onMount(() => {
    const place = () => {
      if (!overlay) return
      pos = {
        x: Math.max(12, window.innerWidth - overlay.offsetWidth - 24),
        y: 80,
      }
    }
    place()
    window.addEventListener('resize', place)
    return () => window.removeEventListener('resize', place)
  })

  function onPointerDown(e: PointerEvent) {
    if ((e.target as HTMLElement).closest('button')) return
    dragging = true
    dragStart = { x: e.clientX, y: e.clientY }
    origin = { ...pos }
    overlay.setPointerCapture(e.pointerId)
  }

  function releaseDrag(e?: PointerEvent) {
    dragging = false
    if (!overlay) return
    try {
      if (e) overlay.releasePointerCapture(e.pointerId)
    } catch {
      /* ignore */
    }
  }

  function onPointerMove(e: PointerEvent) {
    if (!dragging) return
    pos = {
      x: Math.max(0, origin.x + (e.clientX - dragStart.x)),
      y: Math.max(0, origin.y + (e.clientY - dragStart.y)),
    }
  }

  function onPointerUp(e: PointerEvent) {
    releaseDrag(e)
  }

  function onPointerCancel(e: PointerEvent) {
    releaseDrag(e)
  }

  function onLostPointerCapture() {
    dragging = false
  }
</script>

{#if visible}
  <div
    bind:this={overlay}
    class="browser-overlay"
    style="left: {pos.x}px; top: {pos.y}px"
    on:pointerdown={onPointerDown}
    on:pointermove={onPointerMove}
    on:pointerup={onPointerUp}
    on:pointercancel={onPointerCancel}
    on:lostpointercapture={onLostPointerCapture}
  >
    <p class="title" class:recording class:playing={playing && !recording}>{title}</p>
    <div class="actions">
      <button type="button" disabled={!recordEnabled} on:click={onRecord}>
        {@html icons.record}<span>Запись</span>
      </button>
      <button type="button" disabled={!pauseEnabled} on:click={onPause}>
        <span>{paused ? 'Продолжить' : 'Пауза'}</span>
      </button>
      <button type="button" disabled={!stopEnabled} on:click={onStop}>
        {@html icons.stop}<span>Стоп</span>
      </button>
      <button type="button" disabled={!focusEnabled} on:click={onFocusBrowser} title="Показать окно браузера">
        {@html icons.globe}<span>Браузер</span>
      </button>
      <button
        type="button"
        disabled={!pickerEnabled}
        on:click={onPicker}
        title={recording && !paused ? 'Поставьте запись на паузу' : 'Указать элемент'}
      >
        <span>Указать элемент</span>
      </button>
    </div>
  </div>
{/if}

<style>
  .browser-overlay {
    position: fixed;
    z-index: 10000;
    min-width: 280px;
    background: var(--color-sidebar);
    border: 1px solid var(--color-border);
    border-radius: 6px;
    box-shadow: 0 8px 28px rgba(0, 0, 0, 0.45);
    padding: 8px 10px 10px;
    user-select: none;
  }

  .title {
    margin: 0 0 8px;
    font-size: 11px;
    color: var(--color-muted);
    cursor: move;
  }

  .title.recording {
    color: var(--color-recording);
  }

  .title.playing {
    color: var(--color-primary);
  }

  .actions {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
  }

  .actions button {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 4px 8px;
    font-size: 11px;
    background: var(--color-input);
    border: 1px solid var(--color-border);
    border-radius: 3px;
    color: var(--color-text);
  }

  .actions button:disabled {
    opacity: 0.45;
  }

  .actions button :global(svg) {
    width: 14px;
    height: 14px;
  }
</style>
