<script lang="ts">
  export let url = ''
  export let output = ''
  export let testClient = ''
  export let idleSeconds = 30
  export let appendTo = ''
  export let headless = false
  export let filterRecording = false
  export let navOnlyRecording = false
  export let hoverRecord = false
  export let testClients: string[] = []
  export let recording = false
  export let recordPaused = false
  export let onHttpAuth: () => void = () => {}
  export let onStart: () => void = () => {}
  export let onTogglePause: () => void = () => {}
  export let onStop: () => void = () => {}
  export let onClose: () => void = () => {}

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="modal-backdrop" role="presentation" on:click={onClose}>
  <div class="modal record-dialog" role="dialog" aria-label="Live-запись" on:click|stopPropagation>
    <h3>Live-запись</h3>
    <label>URL <input bind:value={url} /></label>
    <label>Файл <input bind:value={output} /></label>
    <label>Дописать в существующий feature
      <input bind:value={appendTo} placeholder="путь к .feature или пусто" />
    </label>
    <label>TestClient
      <select bind:value={testClient}>
        <option value="">(без профиля)</option>
        {#each testClients as client}
          <option value={client}>{client}</option>
        {/each}
      </select>
    </label>
    <label>Idle (сек) <input type="number" bind:value={idleSeconds} min="5" /></label>
    <label class="check-row"><input type="checkbox" bind:checked={headless} /> Headless</label>
    <label class="check-row">
      <input type="checkbox" bind:checked={filterRecording} on:change={() => filterRecording && (navOnlyRecording = false)} />
      Только важные (фильтр записи)
    </label>
    <label class="check-row">
      <input type="checkbox" bind:checked={navOnlyRecording} on:change={() => navOnlyRecording && (filterRecording = false)} />
      Только ссылки
    </label>
    <label class="check-row"><input type="checkbox" bind:checked={hoverRecord} /> Записывать наведение</label>
    <p class="hint">В URL можно указать user:pass@host — пароль сохранится для хоста. Пикер элемента доступен на паузе.</p>
    <div class="modal-actions">
      <button type="button" on:click={onHttpAuth}>HTTP Auth…</button>
      {#if recording}
        <button type="button" on:click={onTogglePause}>{recordPaused ? 'Resume' : 'Pause'}</button>
        <button type="button" on:click={onStop}>Стоп</button>
      {:else}
        <button type="button" class="primary" on:click={onStart}>Начать</button>
      {/if}
      <button type="button" on:click={onClose}>Закрыть</button>
    </div>
  </div>
</div>
