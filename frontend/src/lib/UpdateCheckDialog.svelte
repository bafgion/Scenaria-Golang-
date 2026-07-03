<script lang="ts">
  import { createTranslator, locale } from './i18n'
  import type { gui } from '../../wailsjs/go/models'
  import { BRAND_NAME } from './brand'

  export let currentVersion = ''
  export let info: gui.UpdateInfoDTO | null = null
  export let message = ''
  export let hasUpdate = false
  export let downloading = false
  export let progress: gui.UpdateProgressDTO | null = null
  export let canAutoApply = false
  export let onClose: () => void = () => {}
  export let onOpenRelease: () => void = () => {}
  export let onDownload: () => void = () => {}
  export let onApply: () => void = () => {}

  $: tr = createTranslator($locale)

  $: progressPercent = Math.max(0, Math.min(100, Number(progress?.percent ?? 0)))
  $: progressLabel = progress?.message || (downloading ? tr('dialogs.update.updating') : '')
  $: progressIndeterminate = downloading && progressPercent <= 0 && (!progress?.stage || progress.stage === 'download')

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape' && !downloading) onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={() => !downloading && onClose()}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal update-dialog" role="dialog" aria-modal="true" aria-label={tr('dialogs.update.ariaLabel')} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{tr('dialogs.update.title', { brand: BRAND_NAME })}</h3>
    <p class="version">{tr('dialogs.update.currentVersion', { version: info?.currentVersion || currentVersion || tr('dialogs.common.notFound') })}</p>
    {#if hasUpdate}
      <p class="update-available">{tr('dialogs.update.available', { version: info?.latestVersion || tr('dialogs.update.availableFallback') })}</p>
    {:else}
      <p class="up-to-date">{tr('dialogs.update.upToDate')}</p>
    {/if}
    {#if downloading && progressLabel}
      <div class="progress-block" aria-live="polite">
        <div class="progress-label">{progressLabel}</div>
        <div class="progress-track" role="progressbar" aria-valuemin="0" aria-valuemax="100" aria-valuenow={progressPercent} aria-busy={progressIndeterminate}>
          <div class="progress-fill" class:indeterminate={progressIndeterminate} style:width="{progressPercent}%"></div>
        </div>
      </div>
    {/if}
    {#if message}
      <pre class="details">{message.trim()}</pre>
    {/if}
    <div class="modal-actions">
      {#if hasUpdate && info?.htmlUrl}
        <button type="button" disabled={downloading} on:click={onOpenRelease}>{tr('dialogs.update.releasePage')}</button>
      {/if}
      {#if hasUpdate && canAutoApply && info?.downloadUrl}
        <button type="button" class="primary" disabled={downloading} on:click={onApply}>
          {downloading ? tr('dialogs.update.updating') : tr('dialogs.update.install')}
        </button>
      {/if}
      {#if hasUpdate && info?.downloadUrl}
        <button type="button" disabled={downloading} on:click={onDownload}>
          {downloading ? tr('dialogs.update.downloading') : tr('dialogs.update.download')}
        </button>
      {/if}
      <button type="button" disabled={downloading} on:click={onClose}>{tr('dialogs.common.close')}</button>
    </div>
  </div>
</div>

<style>
  h3 {
    margin: 0 0 8px;
    font-size: 14px;
  }

  .version {
    margin: 0 0 8px;
    font-size: 12px;
    color: var(--color-muted);
  }

  .update-available {
    margin: 0 0 8px;
    font-size: 13px;
    color: var(--color-brand);
  }

  .up-to-date {
    margin: 0 0 8px;
    font-size: 13px;
    color: var(--color-success);
  }

  .details {
    margin: 0 0 12px;
    padding: 8px;
    font-size: 11px;
    line-height: 1.4;
    white-space: pre-wrap;
    background: var(--color-input);
    border: 1px solid var(--color-border);
    border-radius: 4px;
    max-height: 160px;
    overflow: auto;
    font-family: var(--font-mono);
  }

  .progress-block {
    margin: 0 0 12px;
  }

  .progress-label {
    margin: 0 0 6px;
    font-size: 12px;
    color: var(--color-text);
  }

  .progress-track {
    height: 8px;
    border-radius: 4px;
    background: var(--color-input);
    border: 1px solid var(--color-border);
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    min-width: 0;
    background: var(--color-primary);
    transition: width 0.2s ease;
  }

  .progress-fill.indeterminate {
    width: 40% !important;
    animation: update-progress-slide 1.1s ease-in-out infinite;
  }

  @keyframes update-progress-slide {
    0% {
      transform: translateX(-120%);
    }
    100% {
      transform: translateX(320%);
    }
  }
</style>
