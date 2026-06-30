<script lang="ts">
  import brandMark from '../assets/branding/app-icon-mark.png'

  export let version = 'dev'
  export let message = 'Запуск…'
  export let progress = 0
  export let fading = false
  export let standalone = false
</script>

<div class="splash-overlay" class:fading class:standalone aria-live="polite" aria-busy={!fading}>
  <div class="splash-panel">
    <div class="splash-accent" aria-hidden="true"></div>

    <div class="splash-body">
      <div class="logo-wrap">
        <div class="logo-halo" aria-hidden="true"></div>
        <img class="logo-mark" src={brandMark} alt="" width="72" height="72" />
      </div>

      <h1 class="splash-title">SCENARIA</h1>
      <p class="splash-version">VERSION {version}</p>

      <div class="splash-progress-block">
        <div class="splash-status-row">
          <span class="splash-status">{message}</span>
          <span class="splash-percent">{Math.round(progress)}%</span>
        </div>
        <div class="splash-track" aria-hidden="true">
          <div class="splash-fill" style="width: {Math.max(0, Math.min(100, progress))}%"></div>
        </div>
      </div>

      <p class="splash-tagline">Автотесты сайтов · Gherkin · Playwright</p>
    </div>
  </div>
</div>

<style>
  .splash-overlay {
    position: fixed;
    inset: 0;
    z-index: 3000;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #141618;
    opacity: 1;
    transition: opacity 320ms cubic-bezier(0.33, 1, 0.68, 1);
  }

  .splash-overlay.fading {
    opacity: 0;
    pointer-events: none;
  }

  .splash-overlay.standalone {
    background: linear-gradient(180deg, #2b2f33 0%, #1f2124 42%, #141618 100%);
  }

  .splash-overlay.standalone .splash-panel {
    width: 100%;
    height: 100%;
    min-height: 100%;
    max-width: none;
    border-radius: 0;
    box-shadow: none;
    display: flex;
    flex-direction: column;
  }

  .splash-overlay.standalone .splash-body {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
    padding: 28px 32px 24px;
  }

  .splash-overlay.standalone .splash-tagline {
    margin-top: auto;
  }

  .splash-panel {
    position: relative;
    width: min(540px, calc(100vw - 32px));
    border-radius: 12px;
    overflow: hidden;
    background: linear-gradient(180deg, #2b2f33 0%, #1f2124 42%, #141618 100%);
    box-shadow:
      0 18px 48px rgba(0, 0, 0, 0.45),
      0 6px 18px rgba(0, 0, 0, 0.35),
      0 0 0 1px rgba(255, 255, 255, 0.08);
  }

  .splash-accent {
    height: 4px;
    background: linear-gradient(90deg, #3aa8d8 0%, #5ec8f2 45%, #7fd99a 100%);
  }

  .splash-body {
    padding: 30px 40px 26px;
    text-align: center;
  }

  .logo-wrap {
    position: relative;
    width: 128px;
    height: 128px;
    margin: 0 auto 10px;
    display: grid;
    place-items: center;
  }

  .logo-halo {
    position: absolute;
    inset: 0;
    border-radius: 50%;
    background: radial-gradient(circle, rgba(94, 200, 242, 0.35) 0%, rgba(9, 71, 113, 0.2) 45%, transparent 72%);
    animation: splash-pulse 2.8s ease-in-out infinite;
  }

  .logo-mark {
    position: relative;
    z-index: 1;
    display: block;
    object-fit: contain;
  }

  @keyframes splash-pulse {
    0%,
    100% {
      transform: scale(0.96);
      opacity: 0.85;
    }
    50% {
      transform: scale(1.04);
      opacity: 1;
    }
  }

  .splash-title {
    margin: 0;
    font-size: 30px;
    font-weight: 300;
    letter-spacing: 0.14em;
    color: #f5f5f5;
  }

  .splash-version {
    margin: 6px 0 0;
    font-size: 11px;
    letter-spacing: 0.16em;
    text-transform: uppercase;
    color: #8b8b8b;
  }

  .splash-progress-block {
    margin-top: 34px;
    text-align: left;
  }

  .splash-status-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    margin-bottom: 10px;
    font-size: 12px;
  }

  .splash-status {
    color: #e0e0e0;
  }

  .splash-percent {
    color: #8b8b8b;
    font-variant-numeric: tabular-nums;
  }

  .splash-track {
    height: 5px;
    border-radius: 2px;
    background: #2a2a2a;
    overflow: hidden;
  }

  .splash-fill {
    height: 100%;
    min-width: 6px;
    border-radius: 2px;
    background: linear-gradient(90deg, #3aa8d8 0%, #4db8e8 55%, #7fd99a 100%);
    box-shadow: inset -8px 0 12px rgba(255, 255, 255, 0.15);
    transition: width 280ms ease-out;
  }

  .splash-tagline {
    margin: 18px 0 0;
    padding-top: 18px;
    border-top: 1px solid rgba(255, 255, 255, 0.08);
    font-size: 11px;
    color: #6e6e6e;
  }
</style>
