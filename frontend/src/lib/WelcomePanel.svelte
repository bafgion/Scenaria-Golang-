<script lang="ts">
  import { onMount } from 'svelte'
  import { BRAND_NAME } from './brand'
  import { createTranslator, locale } from './i18n'

  export let startURL = 'https://site.com'
  export let recentProjects: string[] = []
  export let recentFeatures: string[] = []
  export let projectOpen = false
  export let recorded = false
  export let playedSuccess = false
  export let checklistDismissed = false
  export let onOpenProject: () => void
  export let onNewProject: () => void = () => {}
  export let onQuickStart: () => void
  export let onNewScenario: () => void
  export let onOpenFile: () => void
  export let onInsertTemplate: () => void
  export let onOpenExamples: () => void
  export let onOpenRecentProject: (path: string) => void
  export let onOpenRecentFeature: (path: string) => void
  export let onChecklistStep: (step: number) => void
  export let onDismissChecklist: () => void = () => {}
  export let tourElevated = false

  $: tr = createTranslator($locale)

  $: steps = [
    { id: 1, label: tr('welcome.checklist.openProject') },
    { id: 2, label: tr('welcome.checklist.record') },
    { id: 3, label: tr('welcome.checklist.runTest') },
  ]

  $: doneFlags = [projectOpen, recorded, playedSuccess]
  $: currentIndex = doneFlags.findIndex((done) => !done)
  $: activeIndex = currentIndex === -1 ? steps.length - 1 : currentIndex

  let scrollEl: HTMLDivElement
  let bodyEl: HTMLDivElement
  let cardEl: HTMLDivElement

  function syncScrollBody() {
    if (!scrollEl || !bodyEl || !cardEl) return
    const viewportHeight = scrollEl.clientHeight
    const contentH = cardEl.offsetHeight + 48
    bodyEl.style.minHeight = `${Math.max(contentH, viewportHeight)}px`
  }

  onMount(() => {
    const ro = new ResizeObserver(() => syncScrollBody())
    if (scrollEl) ro.observe(scrollEl)
    if (cardEl) ro.observe(cardEl)
    syncScrollBody()
    return () => ro.disconnect()
  })

  $: if (scrollEl && bodyEl && cardEl) {
    recentFeatures
    recentProjects
    checklistDismissed
    projectOpen
    recorded
    playedSuccess
    requestAnimationFrame(syncScrollBody)
  }

  function featureName(path: string): string {
    const parts = path.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || path
  }

  function projectName(path: string): string {
    const parts = path.replace(/\\/g, '/').split('/').filter(Boolean)
    return parts[parts.length - 1] || path
  }
</script>

<div class="welcome" class:onboarding-elevated={tourElevated}>
  <div class="welcome-scroll" bind:this={scrollEl}>
    <div class="welcome-scroll-body" bind:this={bodyEl}>
      <div class="welcome-card" bind:this={cardEl} data-tour="welcome-card">
        <h1>{BRAND_NAME}</h1>

        {#if !checklistDismissed}
          <div class="checklist" data-tour="welcome-checklist">
            {#each steps as step, index}
              {@const done = doneFlags[index]}
              {@const current = index === activeIndex && !done}
              <div class="checklist-row" class:done class:current>
                {#if done}
                  <span class="checklist-icon done">✓</span>
                {:else if current}
                  <span class="checklist-icon current" aria-hidden="true">→</span>
                {:else}
                  <span class="checklist-icon muted">○</span>
                {/if}
                {#if current}
                  <button type="button" class="checklist-link" on:click={() => onChecklistStep(step.id)}>
                    {step.label}
                  </button>
                {:else}
                  <span class:muted={!done && !current}>{step.label}</span>
                {/if}
              </div>
            {/each}
          </div>
          <button type="button" class="checklist-dismiss" on:click={onDismissChecklist}>{tr('welcome.checklist.hide')}</button>
        {/if}

        <div class="quick-start">
          <input bind:value={startURL} placeholder="https://site.com" />
          <button class="primary" on:click={onQuickStart} title={tr('welcome.quickStartTitle')}>
            {tr('welcome.quickStart')}
          </button>
        </div>

        <p class="section-heading">{tr('welcome.gettingStarted')}</p>
        <div class="links">
          <button data-tour="welcome-examples" on:click={onOpenExamples}>{tr('welcome.openExamples')}</button>
          <button on:click={onNewProject}>{tr('welcome.newProject')}</button>
          <button on:click={onOpenProject}>{tr('welcome.openFolder')}</button>
          <button on:click={onNewScenario}>{tr('welcome.newScenario')}</button>
          <button on:click={onOpenFile}>{tr('welcome.openFile')}</button>
          <button on:click={onInsertTemplate}>{tr('welcome.insertTemplate')}</button>
        </div>

        {#if recentFeatures.length > 0}
          <p class="recent-heading">{tr('welcome.recentFiles')}</p>
          <div class="links">
            {#each recentFeatures as path}
              <button on:click={() => onOpenRecentFeature(path)} title={path}>{featureName(path)}</button>
            {/each}
          </div>
        {/if}

        {#if recentProjects.length > 0}
          <p class="recent-heading">{tr('welcome.recentProjects')}</p>
          <div class="links">
            {#each recentProjects as path}
              <button on:click={() => onOpenRecentProject(path)} title={path}>{projectName(path)}</button>
            {/each}
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>

<style>
  .welcome {
    flex: 1;
    min-height: 0;
    background: var(--color-bg);
    overflow: hidden;
  }

  .welcome-scroll {
    height: 100%;
    overflow-x: hidden;
    overflow-y: auto;
    background: var(--color-bg);
    scrollbar-width: thin;
    scrollbar-color: var(--color-input) transparent;
  }

  .welcome-scroll::-webkit-scrollbar {
    width: 10px;
    background: transparent;
  }

  .welcome-scroll::-webkit-scrollbar-track {
    background: transparent;
  }

  .welcome-scroll::-webkit-scrollbar-thumb {
    background: var(--color-input);
    border-radius: 5px;
    min-height: 20px;
    border: 2px solid transparent;
    background-clip: padding-box;
  }

  .welcome-scroll::-webkit-scrollbar-button {
    display: none;
    height: 0;
    width: 0;
  }

  .welcome-scroll-body {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    padding: 24px;
    box-sizing: border-box;
  }

  .welcome-card {
    width: min(520px, 100%);
    padding: 32px;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    background: var(--color-panel);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  }

  h1 {
    margin: 0 0 20px;
    font-size: 22px;
    font-weight: 600;
    text-align: center;
  }

  .checklist {
    margin-bottom: 8px;
    padding: 12px 14px;
    border: 1px solid var(--color-border);
    border-radius: 6px;
    background: var(--color-bg);
  }

  .checklist-row {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 4px 0;
    font-size: 13px;
  }

  .checklist-row.done {
    color: var(--color-muted);
    text-decoration: line-through;
  }

  .checklist-row.current {
    font-weight: 600;
  }

  .checklist-icon {
    width: 16px;
    text-align: center;
    flex-shrink: 0;
  }

  .checklist-icon.muted {
    color: var(--color-muted);
  }

  .checklist-icon.done {
    color: var(--color-success, #2e7d32);
  }

  .checklist-link {
    background: none;
    border: none;
    padding: 0;
    font: inherit;
    color: var(--color-link, #0066cc);
    cursor: pointer;
    text-align: left;
    text-decoration: underline;
  }

  .checklist-dismiss {
    margin: 0 0 16px;
    padding: 0;
    background: none;
    border: none;
    font-size: 12px;
    color: var(--color-muted);
    cursor: pointer;
    text-decoration: underline;
  }

  .quick-start {
    display: flex;
    gap: 8px;
    margin-bottom: 20px;
  }

  .quick-start input {
    flex: 1;
    min-width: 0;
  }

  .section-heading,
  .recent-heading {
    margin: 16px 0 8px;
    font-size: 12px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--color-muted);
  }

  .section-heading {
    margin-top: 0;
  }

  .links {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .links button {
    justify-content: flex-start;
    text-align: left;
    padding: 6px 8px;
    font-size: 13px;
  }

  .muted {
    color: var(--color-muted);
  }
</style>
