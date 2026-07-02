<script lang="ts">
  import { onMount, tick } from 'svelte'
  import { brandOverlayTitle } from './brand'
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
  function defaultPos() {
    if (typeof window === 'undefined') return { x: 12, y: 12 }
    const w = 280
    const h = 80
    return {
      x: Math.max(12, window.innerWidth - w - 24),
      y: Math.max(12, window.innerHeight - h - 24),
    }
  }

  let pos = defaultPos()
  let userMoved = false

  function place() {
    if (!overlay || userMoved) return
    const w = overlay.offsetWidth || 280
    const h = overlay.offsetHeight || 80
    pos = {
      x: Math.max(12, window.innerWidth - w - 24),
      y: Math.max(12, window.innerHeight - h - 24),
    }
  }

  $: if (!visible) userMoved = false
  $: if (visible) void tick().then(() => requestAnimationFrame(place))

  onMount(() => {
    void tick().then(() => requestAnimationFrame(place))
    const onResize = () => place()
    window.addEventListener('resize', onResize)
    return () => window.removeEventListener('resize', onResize)
  })

  let dragging = false
  let dragStart = { x: 0, y: 0 }
  let origin = { x: 0, y: 0 }

  $: title = recording
    ? paused
      ? '⏸ Пауза — можно выбрать элемент'
      : '● Идёт запись'
    : playing
      ? '▶ Идёт тест'
      : brandOverlayTitle()

  function onPointerDown(e: PointerEvent) {
    if ((e.target as HTMLElement).closest('button')) return
    userMoved = true
    dragging = true
    dragStart = { x: e.clientX, y: e.clientY }
    origin = { ...pos }
    ;(e.currentTarget as HTMLElement).setPointerCapture(e.pointerId)
  }

  function releaseDrag(e?: PointerEvent) {
    dragging = false
    if (!e) return
    try {
      ;(e.currentTarget as HTMLElement).releasePointerCapture(e.pointerId)
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
  >
    <p
      class="title"
      class:recording
      class:playing={playing && !recording}
      on:pointerdown={onPointerDown}
      on:pointermove={onPointerMove}
      on:pointerup={onPointerUp}
      on:pointercancel={onPointerCancel}
      on:lostpointercapture={onLostPointerCapture}
    >
      {title}
    </p>
    <div class="overlay-actions">
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
    pointer-events: none;
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
    pointer-events: auto;
  }

  .overlay-actions {
    pointer-events: auto;
  }

  .overlay-actions :global(button) {
    pointer-events: auto;
  }

  .title.recording {
    color: var(--color-recording);
  }

  .title.playing {
    color: var(--color-primary);
  }
</style>
