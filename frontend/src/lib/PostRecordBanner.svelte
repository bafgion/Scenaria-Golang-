<script lang="ts">
  import { createTranslator, locale } from './i18n'

  export let path = ''
  export let stepCount = 0
  export let onValidate: () => void = () => {}
  export let onSave: () => void = () => {}
  export let onShowDiff: () => void = () => {}
  export let onClose: () => void = () => {}

  $: tr = createTranslator($locale)

  function basename(p: string): string {
    const parts = p.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || p
  }
</script>

<div class="post-record-banner" role="status">
  <div class="top">
    <div class="summary">
      {tr('dialogs.postRecord.banner.summary', { stepCount, fileName: basename(path) })}
      <span class="hint-note">{tr('dialogs.postRecord.banner.hint')}</span>
    </div>
    <div class="banner-actions">
      <button type="button" on:click={onShowDiff}>{tr('dialogs.postRecord.banner.compare')}</button>
      <button type="button" class="primary" on:click={onValidate}>{tr('dialogs.postRecord.banner.validate')}</button>
      <button type="button" on:click={onSave}>{tr('dialogs.common.save')}</button>
      <button type="button" class="dismiss" on:click={onClose}>{tr('dialogs.common.close')}</button>
    </div>
  </div>
</div>

<style>
  .post-record-banner {
    background: rgba(137, 209, 133, 0.08);
    border-bottom: 1px solid var(--color-success);
    padding: 6px 10px;
    font-size: 12px;
  }

  .top {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 8px;
  }

  .summary {
    flex: 1;
    min-width: 200px;
    color: var(--color-text);
  }

  .hint-note {
    display: block;
    margin-top: 2px;
    color: var(--color-muted);
    font-size: 11px;
  }
</style>
