<script lang="ts">
  import { createEventDispatcher, onDestroy, tick } from 'svelte'
  import { t } from '../i18n'
  import { ONBOARDING_TOUR_STEPS } from './tourSteps'
  import { canAdvanceTourStep, isTourStepComplete, maxValidTourStepIndex, tourTargetSelector, type TourContext } from './tourState'

  export let active = false
  export let context: TourContext

  const dispatch = createEventDispatcher<{
    skip: void
    complete: void
    stepChange: { index: number; id: string }
  }>()

  const CARD_W = 400
  const CARD_H_EST = 240
  const VIEWPORT_PAD = 16

  let stepIndex = 0
  let targetEl: HTMLElement | null = null
  let hole = { top: 0, left: 0, width: 0, height: 0 }
  let cardStyle = ''
  let cardPlacement: 'top' | 'bottom' | 'left' | 'right' | 'center' = 'center'
  let modalMode = false
  let autoAdvancedFor = -1
  let resizeObserver: ResizeObserver | null = null
  let viewportResizeObserver: ResizeObserver | null = null
  let scrollParent: HTMLElement | null = null
  let onScroll: (() => void) | null = null
  let onWindowResize: (() => void) | null = null
  let onWelcomeScroll: (() => void) | null = null
  let repositionRaf = 0
  let highlightedEl: HTMLElement | null = null
  let layoutToken = 0
  let spotlightReady = false

  function setHighlightedTarget(el: HTMLElement | null) {
    highlightedEl?.classList.remove('onboarding-target')
    highlightedEl = el
    el?.classList.add('onboarding-target')
  }

  $: steps = ONBOARDING_TOUR_STEPS
  $: step = steps[stepIndex]
  $: total = steps.length
  $: canNext = step ? canAdvanceTourStep(step, context) : false
  $: isLast = stepIndex >= total - 1
  $: actionPending = step?.kind === 'action' && !canNext
  $: progressPct = ((stepIndex + 1) / total) * 100
  $: blockers = modalMode || step?.kind === 'action' ? [] : spotlightBlockers(hole)

  $: if (active) {
    const maxStep = maxValidTourStepIndex(context)
    if (stepIndex > maxStep) {
      stepIndex = maxStep
      autoAdvancedFor = -1
      dispatch('stepChange', { index: stepIndex, id: steps[stepIndex].id })
    }
  }

  $: if (active && step) {
    spotlightReady = false
    hole = { top: 0, left: 0, width: 0, height: 0 }
    modalMode = (step.placement ?? 'bottom') === 'center'
    const token = ++layoutToken
    void runLayout(step, token)
  }

  async function runLayout(current: (typeof steps)[number], token: number) {
    await layoutStep(current, token)
  }

  function unionRects(...rects: DOMRect[]): DOMRect {
    const valid = rects.filter((r) => r.width > 0 && r.height > 0)
    if (valid.length === 0) return new DOMRect(0, 0, 0, 0)
    const top = Math.min(...valid.map((r) => r.top))
    const left = Math.min(...valid.map((r) => r.left))
    const right = Math.max(...valid.map((r) => r.right))
    const bottom = Math.max(...valid.map((r) => r.bottom))
    return new DOMRect(left, top, right - left, bottom - top)
  }

  $: if (active && step?.kind === 'action' && isTourStepComplete(step.id, context) && autoAdvancedFor !== stepIndex) {
    autoAdvancedFor = stepIndex
    void advanceSoon()
  }

  async function advanceSoon() {
    await tick()
    await new Promise((r) => requestAnimationFrame(() => requestAnimationFrame(r)))
    if (!active || autoAdvancedFor !== stepIndex) return
    if (isLast) {
      dispatch('complete')
      return
    }
    stepIndex += 1
    dispatch('stepChange', { index: stepIndex, id: steps[stepIndex].id })
  }

  function spotlightBlockers(h: typeof hole) {
    if (!h.width || !h.height) return []
    const vw = window.innerWidth
    const vh = window.innerHeight
    return [
      { top: 0, left: 0, width: vw, height: h.top },
      { top: h.top, left: 0, width: h.left, height: h.height },
      { top: h.top, left: h.left + h.width, width: Math.max(0, vw - h.left - h.width), height: h.height },
      { top: h.top + h.height, left: 0, width: vw, height: Math.max(0, vh - h.top - h.height) },
    ]
  }

  async function resolveTarget(target: string): Promise<HTMLElement | null> {
    const menuTarget =
      target === 'menu-run' ||
      target === 'menu-validate' ||
      target === 'menu-dry-run' ||
      target === 'menu-run-dropdown'
    const attempts = menuTarget ? 24 : 10
    for (let attempt = 0; attempt < attempts; attempt++) {
      await tick()
      if (attempt > 0) {
        await new Promise<void>((r) => requestAnimationFrame(() => r()))
      }
      const dropdown = document.querySelector('[data-tour="menu-run-dropdown"]') as HTMLElement | null
      if (menuTarget && (!dropdown || !isElementVisible(dropdown))) continue
      const el = document.querySelector(tourTargetSelector(target)) as HTMLElement | null
      if (!el || !isElementVisible(el)) continue
      return el
    }
    return document.querySelector(tourTargetSelector(target)) as HTMLElement | null
  }

  async function layoutStep(current: (typeof steps)[number], token = layoutToken) {
    await tick()
    if (token !== layoutToken) return
    clearHighlight()
    spotlightReady = false
    hole = { top: 0, left: 0, width: 0, height: 0 }
    modalMode = (current.placement ?? 'bottom') === 'center'
    const el = await resolveTarget(current.target)
    if (token !== layoutToken) return
    targetEl = el
    let highlightEl = el
    if (el && current.target === 'menu-run') {
      const actionSelector =
        current.id === 'validate'
          ? '[data-tour="menu-validate"]'
          : current.id === 'dry-run'
            ? '[data-tour="menu-dry-run"]'
            : ''
      if (actionSelector) {
        const actionEl = document.querySelector(actionSelector) as HTMLElement | null
        if (actionEl && isElementVisible(actionEl)) highlightEl = actionEl
      }
    }
    setHighlightedTarget(highlightEl)

    if (modalMode) {
      cardPlacement = 'center'
      cardStyle = centerCardStyle()
      spotlightReady = true
      return
    }

    if (!el || !isElementVisible(el)) {
      modalMode = true
      cardPlacement = 'center'
      cardStyle = centerCardStyle()
      spotlightReady = true
      return
    }

    el.scrollIntoView({ block: 'nearest', inline: 'nearest', behavior: 'smooth' })
    await tick()
    if (token !== layoutToken) return
    updateHole()
    positionCard(current.placement ?? 'bottom')
    bindResize(el)
    spotlightReady = true
  }

  function centerCardStyle(): string {
    return 'top: 50%; left: 50%; transform: translate(-50%, -50%); width: min(440px, calc(100vw - 32px));'
  }

  function isElementVisible(el: HTMLElement): boolean {
    const rect = el.getBoundingClientRect()
    if (rect.width < 2 || rect.height < 2) return false
    const style = window.getComputedStyle(el)
    return style.display !== 'none' && style.visibility !== 'hidden' && Number(style.opacity) > 0
  }

  function getTargetRect(el: HTMLElement): DOMRect {
    const tourId = el.getAttribute('data-tour')
    if (tourId === 'menu-run') {
      const trigger = el.querySelector('[data-tour="menu-run-trigger"]') as HTMLElement | null
      const dropdown = el.querySelector('[data-tour="menu-run-dropdown"]') as HTMLElement | null
      return unionRects(
        el.getBoundingClientRect(),
        trigger?.getBoundingClientRect() ?? new DOMRect(),
        dropdown?.getBoundingClientRect() ?? new DOMRect(),
      )
    }
    if (tourId === 'menu-validate' || tourId === 'menu-dry-run') {
      const trigger = document.querySelector('[data-tour="menu-run-trigger"]') as HTMLElement | null
      const dropdown = document.querySelector('[data-tour="menu-run-dropdown"]') as HTMLElement | null
      return unionRects(
        el.getBoundingClientRect(),
        trigger?.getBoundingClientRect() ?? new DOMRect(),
        dropdown?.getBoundingClientRect() ?? new DOMRect(),
      )
    }
    let rect = el.getBoundingClientRect()
    if (tourId === 'menu-run-dropdown') {
      const trigger = document.querySelector('[data-tour="menu-run-trigger"]') as HTMLElement | null
      if (trigger) rect = unionRects(rect, trigger.getBoundingClientRect())
    }
    if (tourId === 'welcome-examples') {
      const card = document.querySelector('[data-tour="welcome-card"]') as HTMLElement | null
      if (card) rect = unionRects(rect, card.getBoundingClientRect())
    }
    return rect
  }

  function scheduleReposition() {
    if (!active || !step) return
    if (modalMode) {
      cardStyle = centerCardStyle()
      return
    }
    if (!targetEl || !isElementVisible(targetEl)) return
    cancelAnimationFrame(repositionRaf)
    repositionRaf = requestAnimationFrame(() => {
      repositionRaf = 0
      if (!targetEl || !step || modalMode) return
      updateHole()
      positionCard(step.placement ?? 'bottom')
    })
  }

  function updateHole() {
    if (!targetEl) return
    const rect = getTargetRect(targetEl)
    const pad = 8
    hole = {
      top: Math.max(0, rect.top - pad),
      left: Math.max(0, rect.left - pad),
      width: rect.width + pad * 2,
      height: rect.height + pad * 2,
    }
  }

  function positionCard(preferred: string) {
    if (!targetEl) {
      cardPlacement = 'center'
      cardStyle = centerCardStyle()
      return
    }
    const rect = getTargetRect(targetEl)
    const margin = 14
    const vw = window.innerWidth
    const vh = window.innerHeight

    let placement = preferred as 'top' | 'bottom' | 'left' | 'right'
    let top = 0
    let left = 0

    const fitsBottom = rect.bottom + margin + CARD_H_EST <= vh - VIEWPORT_PAD
    const fitsTop = rect.top - margin - CARD_H_EST >= VIEWPORT_PAD
    const fitsRight = rect.right + margin + CARD_W <= vw - VIEWPORT_PAD
    const fitsLeft = rect.left - margin - CARD_W >= VIEWPORT_PAD

    if (placement === 'bottom' && !fitsBottom && fitsTop) placement = 'top'
    if (placement === 'top' && !fitsTop && fitsBottom) placement = 'bottom'
    if (placement === 'right' && !fitsRight && fitsLeft) placement = 'left'
    if (placement === 'left' && !fitsLeft && fitsRight) placement = 'right'

    if (placement === 'bottom') {
      top = rect.bottom + margin
      left = rect.left + rect.width / 2 - CARD_W / 2
    } else if (placement === 'top') {
      top = rect.top - margin - CARD_H_EST
      left = rect.left + rect.width / 2 - CARD_W / 2
    } else if (placement === 'right') {
      top = rect.top + rect.height / 2 - CARD_H_EST / 2
      left = rect.right + margin
    } else {
      top = rect.top + rect.height / 2 - CARD_H_EST / 2
      left = rect.left - margin - CARD_W
    }

    left = Math.min(Math.max(VIEWPORT_PAD, left), vw - CARD_W - VIEWPORT_PAD)
    top = Math.min(Math.max(VIEWPORT_PAD, top), vh - CARD_H_EST - VIEWPORT_PAD)

    cardPlacement = placement
    cardStyle = `top: ${top}px; left: ${left}px; width: min(${CARD_W}px, calc(100vw - ${VIEWPORT_PAD * 2}px));`
  }

  function bindResize(el: HTMLElement) {
    resizeObserver?.disconnect()
    resizeObserver = new ResizeObserver(() => scheduleReposition())
    resizeObserver.observe(el)

    viewportResizeObserver?.disconnect()
    viewportResizeObserver = new ResizeObserver(() => scheduleReposition())
    viewportResizeObserver.observe(document.documentElement)
    const ide = document.querySelector('.ide')
    if (ide instanceof HTMLElement) viewportResizeObserver.observe(ide)
    const welcomeScroll = el.closest('.welcome-scroll')
    if (welcomeScroll instanceof HTMLElement) viewportResizeObserver.observe(welcomeScroll)

    if (onWindowResize) window.removeEventListener('resize', onWindowResize)
    onWindowResize = () => scheduleReposition()
    window.addEventListener('resize', onWindowResize)

    if (onWelcomeScroll) {
      const prev = onWelcomeScroll
      el.closest('.welcome-scroll')?.removeEventListener('scroll', prev)
    }
    const scrollHost = el.closest('.welcome-scroll')
    if (scrollHost) {
      onWelcomeScroll = () => scheduleReposition()
      scrollHost.addEventListener('scroll', onWelcomeScroll, { passive: true })
    } else {
      onWelcomeScroll = null
    }

    if (onScroll) scrollParent?.removeEventListener('scroll', onScroll, true)
    onScroll = () => scheduleReposition()
    scrollParent = document.querySelector('.ide')
    scrollParent?.addEventListener('scroll', onScroll, true)
  }

  function clearHighlight() {
    setHighlightedTarget(null)
    resizeObserver?.disconnect()
    resizeObserver = null
    viewportResizeObserver?.disconnect()
    viewportResizeObserver = null
    if (onWindowResize) {
      window.removeEventListener('resize', onWindowResize)
      onWindowResize = null
    }
    if (onWelcomeScroll) {
      document.querySelector('.welcome-scroll')?.removeEventListener('scroll', onWelcomeScroll)
      onWelcomeScroll = null
    }
    if (repositionRaf) {
      cancelAnimationFrame(repositionRaf)
      repositionRaf = 0
    }
    if (onScroll) {
      scrollParent?.removeEventListener('scroll', onScroll, true)
      onScroll = null
    }
  }

  function onBack() {
    if (stepIndex <= 0) return
    autoAdvancedFor = -1
    stepIndex -= 1
    dispatch('stepChange', { index: stepIndex, id: steps[stepIndex].id })
  }

  async function onNext() {
    if (!canNext) return
    if (isLast) {
      dispatch('complete')
      return
    }
    autoAdvancedFor = -1
    stepIndex += 1
    dispatch('stepChange', { index: stepIndex, id: steps[stepIndex].id })
  }

  function onSkip() {
    dispatch('skip')
  }

  export function restart() {
    stepIndex = 0
    autoAdvancedFor = -1
  }

  export async function relayout() {
    if (step) {
      const token = ++layoutToken
      await layoutStep(step, token)
    } else scheduleReposition()
  }

  onDestroy(() => {
    clearHighlight()
    scrollParent = null
  })
</script>

{#if active && step}
  <div class="onboarding-overlay" role="presentation">
    {#if modalMode}
      <div class="onboarding-backdrop"></div>
    {:else if hole.width > 0 && hole.height > 0}
      <div
        class="onboarding-hole"
        style="top: {hole.top}px; left: {hole.left}px; width: {hole.width}px; height: {hole.height}px;"
      ></div>
      {#each blockers as blocker (blocker.top + '-' + blocker.left)}
        <div
          class="onboarding-blocker"
          style="top: {blocker.top}px; left: {blocker.left}px; width: {blocker.width}px; height: {blocker.height}px;"
        ></div>
      {/each}
    {:else if !spotlightReady}
      <div class="onboarding-backdrop"></div>
    {/if}
  </div>
  {#if !modalMode && hole.width > 0 && hole.height > 0}
    <div
      class="onboarding-ring"
      style="top: {hole.top}px; left: {hole.left}px; width: {hole.width}px; height: {hole.height}px;"
    ></div>
  {/if}
  <div
    class="onboarding-card"
    class:centered={cardPlacement === 'center'}
    class:placement-top={cardPlacement === 'top'}
    class:placement-bottom={cardPlacement === 'bottom'}
    class:placement-left={cardPlacement === 'left'}
    class:placement-right={cardPlacement === 'right'}
    style={cardStyle}
    role="dialog"
    aria-modal="true"
    aria-labelledby="onboarding-title"
  >
      <div class="onboarding-progress-row">
        <p class="onboarding-progress">{t('onboarding.stepOf', { current: stepIndex + 1, total })}</p>
        <div class="onboarding-progress-bar" aria-hidden="true">
          <div class="onboarding-progress-fill" style="width: {progressPct}%"></div>
        </div>
      </div>
      <h2 id="onboarding-title">{t(step.titleKey)}</h2>
      <p class="onboarding-body">{t(step.bodyKey)}</p>
      {#if actionPending}
        <p class="onboarding-hint">{t('onboarding.waitAction')}</p>
      {/if}
      <div class="onboarding-actions">
        <button type="button" class="ghost" on:click={onSkip}>{t('onboarding.skip')}</button>
        <div class="onboarding-nav">
          <button type="button" disabled={stepIndex === 0} on:click={onBack}>{t('onboarding.back')}</button>
          <button type="button" class="primary" disabled={!canNext} on:click={onNext}>
            {isLast ? t('onboarding.done') : t('onboarding.next')}
          </button>
        </div>
      </div>
    </div>
{/if}

<style>
  .onboarding-overlay {
    position: fixed;
    inset: 0;
    z-index: var(--z-onboarding-dim);
    pointer-events: none;
  }

  .onboarding-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.55);
    pointer-events: auto;
    z-index: var(--z-onboarding-dim);
  }

  .onboarding-hole {
    position: fixed;
    z-index: var(--z-onboarding-dim);
    border-radius: 10px;
    box-shadow: 0 0 0 9999px rgba(0, 0, 0, 0.55);
    pointer-events: none;
    transition: top 0.2s ease, left 0.2s ease, width 0.2s ease, height 0.2s ease;
  }

  .onboarding-ring {
    position: fixed;
    z-index: var(--z-onboarding-ring);
    border-radius: 10px;
    box-shadow: 0 0 0 2px var(--color-primary);
    pointer-events: none;
    transition: top 0.2s ease, left 0.2s ease, width 0.2s ease, height 0.2s ease;
  }

  .onboarding-blocker {
    position: fixed;
    z-index: var(--z-onboarding-block);
    pointer-events: auto;
  }

  .onboarding-card {
    position: fixed;
    z-index: var(--z-onboarding-card);
    background: var(--color-sidebar);
    border: 1px solid var(--color-border);
    border-radius: 10px;
    padding: 18px 20px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.55);
    color: var(--color-text);
    pointer-events: auto;
    box-sizing: border-box;
  }

  .onboarding-card::before {
    content: '';
    position: absolute;
    width: 12px;
    height: 12px;
    background: var(--color-sidebar);
    border: 1px solid var(--color-border);
    transform: rotate(45deg);
    pointer-events: none;
  }

  .onboarding-card.centered::before,
  .onboarding-card.placement-left::before {
    display: none;
  }

  .onboarding-card.placement-bottom::before {
    top: -7px;
    left: 50%;
    margin-left: -6px;
    border-bottom: none;
    border-right: none;
  }

  .onboarding-card.placement-top::before {
    bottom: -7px;
    left: 50%;
    margin-left: -6px;
    border-top: none;
    border-left: none;
  }

  .onboarding-card.placement-right::before {
    left: -7px;
    top: 50%;
    margin-top: -6px;
    border-bottom: none;
    border-left: none;
  }

  .onboarding-progress-row {
    margin-bottom: 10px;
  }

  .onboarding-progress {
    margin: 0 0 6px;
    font-size: 11px;
    color: var(--color-muted);
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .onboarding-progress-bar {
    height: 3px;
    border-radius: 2px;
    background: var(--color-input);
    overflow: hidden;
  }

  .onboarding-progress-fill {
    height: 100%;
    border-radius: 2px;
    background: var(--color-primary);
    transition: width 0.25s ease;
  }

  h2 {
    margin: 0 0 8px;
    font-size: 17px;
    font-weight: 600;
    line-height: 1.3;
  }

  .onboarding-body {
    margin: 0 0 12px;
    font-size: 13px;
    line-height: 1.5;
    color: var(--color-text);
  }

  .onboarding-hint {
    margin: 0 0 12px;
    font-size: 12px;
    color: var(--color-accent);
  }

  .onboarding-actions {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    flex-wrap: wrap;
    padding-top: 4px;
    border-top: 1px solid var(--color-divider);
  }

  .onboarding-nav {
    display: flex;
    gap: 8px;
    margin-left: auto;
  }

  button {
    font-size: 12px;
    padding: 7px 14px;
    border-radius: 4px;
    border: 1px solid var(--color-border);
    background: var(--color-input);
    color: var(--color-text);
    cursor: pointer;
  }

  button:disabled {
    opacity: 0.45;
    cursor: default;
  }

  button.primary {
    background: var(--color-primary);
    border-color: var(--color-primary);
    color: #fff;
  }

  button.ghost {
    background: transparent;
    border-color: transparent;
    color: var(--color-muted);
    padding-left: 0;
  }

  button.ghost:hover {
    color: var(--color-text);
  }
</style>
