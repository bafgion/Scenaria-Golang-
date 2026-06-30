<script lang="ts">
  import { icons } from './icons'

  export let activeKey = ''
  export let welcomeKey = '__welcome__'
  export let tabs: { path: string; content: string; dirty: boolean }[] = []
  export let tabLabel: (path: string) => string = (p) => p
  export let tabUnsaved: (tab: { path: string; content: string; dirty: boolean }) => boolean = (t) => t.dirty

  export let onSelect: (key: string) => void = () => {}
  export let onClose: (path: string) => void = () => {}
  export let onCloseWelcome: () => void = () => {}
  export let welcomeVisible = true
</script>

<div class="editor-tabbar" role="tablist" aria-label="Вкладки редактора">
  {#if welcomeVisible}
    <div
      class="editor-tab welcome"
      class:active={activeKey === welcomeKey}
      role="tab"
      tabindex="0"
      aria-selected={activeKey === welcomeKey}
      title="Стартовая страница"
      on:click={() => onSelect(welcomeKey)}
      on:keydown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault()
          onSelect(welcomeKey)
        }
      }}
    >
      <span class="tab-icon" aria-hidden="true">{@html icons.home}</span>
      <span class="tab-label">Старт</span>
      <button
        type="button"
        class="tab-close"
        aria-label="Закрыть вкладку"
        title="Закрыть"
        on:click|stopPropagation={onCloseWelcome}
      >
        {@html icons.close}
      </button>
    </div>
  {/if}

  {#each tabs as tab (tab.path)}
    <div
      class="editor-tab file"
      class:active={tab.path === activeKey}
      role="tab"
      tabindex="0"
      aria-selected={tab.path === activeKey}
      title={tab.path}
      on:click={() => onSelect(tab.path)}
      on:keydown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault()
          onSelect(tab.path)
        }
      }}
    >
      <span class="tab-label">{tabLabel(tab.path)}{tabUnsaved(tab) ? ' *' : ''}</span>
      <button
        type="button"
        class="tab-close"
        aria-label="Закрыть вкладку"
        title="Закрыть"
        on:click|stopPropagation={() => onClose(tab.path)}
      >
        {@html icons.close}
      </button>
    </div>
  {/each}
</div>

<style>
  .editor-tabbar {
    display: flex;
    height: var(--tab-height);
    max-height: var(--tab-height);
    background: var(--color-toolbar);
    border-bottom: 1px solid var(--color-zone-line);
    flex-shrink: 0;
    min-width: 0;
    overflow: hidden;
  }

  .editor-tab {
    display: inline-flex;
    align-items: center;
    height: var(--tab-height);
    max-height: var(--tab-height);
    border: none;
    border-right: 1px solid var(--color-divider);
    background: var(--color-tab-inactive);
    color: var(--color-muted);
    font-size: 11px;
    line-height: 1;
    cursor: default;
    box-sizing: border-box;
    user-select: none;
    overflow: hidden;
  }

  .editor-tab.welcome {
    position: relative;
    gap: 4px;
    padding: 2px 22px 2px 8px;
    flex-shrink: 0;
    min-width: 56px;
    max-width: none;
    font-family: inherit;
  }

  .editor-tab.file {
    position: relative;
    padding: 2px 22px 2px 8px;
    gap: 0;
    flex: 0 1 auto;
    min-width: 56px;
    max-width: 200px;
  }

  .editor-tab.active {
    background: var(--color-bg);
    color: var(--color-text);
    border-top: 1px solid var(--color-primary);
    border-bottom: 1px solid var(--color-bg);
    margin-bottom: -1px;
    flex-shrink: 0;
    max-width: 280px;
  }

  .editor-tab:focus-visible {
    outline: 1px solid var(--color-primary);
    outline-offset: -1px;
  }

  .tab-icon {
    display: inline-flex;
    width: 16px;
    height: 16px;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .tab-icon :global(svg) {
    width: 16px;
    height: 16px;
  }

  .tab-label {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    min-width: 0;
    flex: 1;
  }

  .tab-close {
    position: absolute;
    right: 4px;
    top: 50%;
    transform: translateY(-50%);
    width: 14px;
    height: 14px;
    padding: 0;
    border: none;
    border-radius: 3px;
    background: transparent;
    color: var(--color-text);
    opacity: 0.72;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    cursor: default;
    z-index: 1;
  }

  .tab-close :global(svg) {
    width: 12px;
    height: 12px;
    display: block;
  }

  .tab-close:hover {
    opacity: 1;
    background: var(--color-input);
  }

  .editor-tab.active .tab-close {
    opacity: 0.85;
  }
</style>
