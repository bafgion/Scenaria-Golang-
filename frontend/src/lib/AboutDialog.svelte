<script lang="ts">
  import brandMark from '../assets/branding/app-icon-mark.png'
  import { BRAND_DESCRIPTION, BRAND_TAGLINE, BRAND_TITLE } from './brand'

  export let version = ''
  export let onClose: () => void = () => {}

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal about-modal" role="dialog" aria-modal="true" aria-label="О программе" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <img class="about-mark" src={brandMark} alt="" width="72" height="72" />
    <h3>{BRAND_TITLE}</h3>
    <p class="about-desc">{BRAND_DESCRIPTION}</p>
    <p class="about-version">{version}</p>
    <p class="hint">{BRAND_TAGLINE}</p>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={onClose}>OK</button>
    </div>
  </div>
</div>

<style>
  .about-mark {
    display: block;
    margin: 0 auto 12px;
    border-radius: 16px;
  }

  .about-desc {
    font-size: 13px;
    margin: 0 0 8px;
    text-align: center;
  }

  .about-version {
    font-size: 13px;
    margin: 8px 0;
  }

  .hint {
    font-size: 12px;
    color: var(--color-muted);
    margin: 0 0 12px;
  }
</style>
