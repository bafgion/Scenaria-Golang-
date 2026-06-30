<script lang="ts">
  import { onMount } from 'svelte'
  import { BRAND_NAME } from './brand'

  export let startURL = 'https://site.com'
  export let recentProjects: string[] = []
  export let recentFeatures: string[] = []
  export let projectOpen = false
  export let recorded = false
  export let playedSuccess = false
  export let checklistDismissed = false
  export let onOpenProject: () => void
  export let onQuickStart: () => void
  export let onNewScenario: () => void
  export let onOpenFile: () => void
  export let onInsertTemplate: () => void
  export let onOpenExamples: () => void
  export let onOpenRecentProject: (path: string) => void
  export let onOpenRecentFeature: (path: string) => void
  export let onChecklistStep: (step: number) => void

  const steps = [
    { id: 1, label: 'Открыть проект' },
    { id: 2, label: 'Записать сценарий' },
    { id: 3, label: 'Запустить тест' },
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

<div class="welcome">
  <div class="welcome-scroll" bind:this={scrollEl}>
    <div class="welcome-scroll-body" bind:this={bodyEl}>
      <div class="welcome-card" bind:this={cardEl}>
        <h1>{BRAND_NAME}</h1>

        {#if !checklistDismissed}
          <div class="checklist">
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
        {/if}

        <div class="quick-start">
          <input bind:value={startURL} placeholder="https://site.com" />
          <button class="primary" on:click={onQuickStart} title="Открыть браузер и начать запись">
            Быстрый старт
          </button>
        </div>

        <p class="section-heading">Начало работы</p>
        <div class="links">
          <button on:click={onOpenExamples}>Открыть примеры сценариев</button>
          <button on:click={onOpenProject}>Открыть папку…</button>
          <button on:click={onNewScenario}>Новый сценарий</button>
          <button on:click={onOpenFile}>Открыть файл…</button>
          <button on:click={onInsertTemplate}>Вставить шаблон сценария</button>
        </div>

        {#if recentFeatures.length > 0}
          <p class="recent-heading">Недавние файлы</p>
          <div class="links">
            {#each recentFeatures as path}
              <button on:click={() => onOpenRecentFeature(path)} title={path}>{featureName(path)}</button>
            {/each}
          </div>
        {/if}

        {#if recentProjects.length > 0}
          <p class="recent-heading">Недавние проекты</p>
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
    min-width: 280px;
    padding: 24px 28px;
    background: #252526;
    border: 1px solid var(--color-divider);
    border-radius: 8px;
    box-sizing: border-box;
  }

  h1 {
    margin: 0;
    font-size: 24px;
    font-weight: 300;
    letter-spacing: 0.01em;
    color: var(--color-text);
  }

  .checklist {
    margin: 4px 0 8px;
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .checklist-row {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    line-height: 1.6;
    color: var(--color-text);
  }

  .checklist-row.done {
    color: var(--color-success);
  }

  .checklist-row .muted {
    color: var(--color-muted);
  }

  .checklist-icon {
    width: 14px;
    flex-shrink: 0;
    text-align: center;
  }

  .checklist-icon.done {
    color: var(--color-success);
  }

  .checklist-icon.current {
    color: var(--color-text);
  }

  .checklist-icon.muted {
    color: var(--color-muted);
  }

  .checklist-link {
    padding: 0;
    border: none;
    background: transparent;
    color: var(--color-text);
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    text-align: left;
  }

  .checklist-link:hover {
    color: var(--color-primary);
    text-decoration: underline;
  }

  .quick-start {
    display: flex;
    gap: 8px;
    margin: 6px 0 12px;
  }

  .quick-start input {
    flex: 1;
    min-width: 0;
    padding: 6px 10px;
    background: var(--color-input);
    border: 1px solid var(--color-border);
    border-radius: 3px;
    color: var(--color-text);
    font-size: 13px;
  }

  .quick-start .primary {
    padding: 6px 14px;
    background: var(--color-primary);
    border: none;
    border-radius: 3px;
    color: #ffffff;
    white-space: nowrap;
    cursor: pointer;
    font-size: 13px;
    font-weight: 400;
  }

  .section-heading {
    margin: 8px 0 0;
    font-size: 13px;
    font-weight: 600;
    color: var(--color-text);
  }

  .recent-heading {
    margin: 12px 0 0;
    font-size: 13px;
    font-weight: 600;
    color: var(--color-muted);
  }

  .links {
    display: flex;
    flex-direction: column;
    gap: 0;
    margin-top: 2px;
  }

  .links button {
    text-align: left;
    padding: 4px 0;
    border: none;
    background: transparent;
    color: var(--color-primary);
    font-size: 13px;
    cursor: pointer;
  }

  .links button:hover {
    text-decoration: underline;
  }
</style>
