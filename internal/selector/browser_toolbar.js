(() => {
  if (window.__scenariaToolbar) return;

  const state = { recording: false, paused: false, browserOnly: true, stepCount: 0 };
  let pending = null;
  let drag = null;

  const ICONS = {
    record: '<svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor"><circle cx="12" cy="12" r="7"/></svg>',
    pause: '<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2"><rect x="6" y="5" width="3" height="14" rx="1"/><rect x="15" y="5" width="3" height="14" rx="1"/></svg>',
    play: '<svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor"><path d="M8 5v14l11-7z"/></svg>',
    stop: '<svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor"><rect x="6" y="6" width="12" height="12" rx="1.5"/></svg>',
    picker: '<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2"><path d="m12 3 7 7-4 1-1 4-7-7 4-1 1-4z"/><path d="M3 21l4-4"/></svg>',
  };

  const host = document.createElement('div');
  host.id = 'scenaria-browser-toolbar';
  host.setAttribute('data-scenaria-ui', '1');
  host.innerHTML = `
    <style>
      #scenaria-browser-toolbar {
        all: initial;
        position: fixed;
        top: 10px;
        left: 50%;
        transform: translateX(-50%);
        z-index: 2147483646;
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 6px 8px 6px 10px;
        font-family: "Segoe UI", system-ui, -apple-system, sans-serif;
        font-size: 12px;
        color: #ccc;
        background: rgba(30, 30, 30, 0.94);
        border: 1px solid #454545;
        border-radius: 10px;
        box-shadow: 0 8px 28px rgba(0, 0, 0, 0.45);
        backdrop-filter: blur(8px);
        user-select: none;
        max-width: min(96vw, 720px);
        box-sizing: border-box;
      }
      #scenaria-browser-toolbar * { box-sizing: border-box; }
      #scenaria-browser-toolbar .sc-brand {
        display: flex;
        align-items: center;
        gap: 6px;
        padding-right: 8px;
        border-right: 1px solid #3c3c3c;
        cursor: move;
        white-space: nowrap;
      }
      #scenaria-browser-toolbar .sc-dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: #5ec8f2;
        flex-shrink: 0;
      }
      #scenaria-browser-toolbar .sc-dot.recording { background: #f14c4c; box-shadow: 0 0 6px rgba(241,76,76,.55); }
      #scenaria-browser-toolbar .sc-dot.paused { background: #cca700; }
      #scenaria-browser-toolbar .sc-title { font-weight: 600; letter-spacing: 0.04em; font-size: 11px; color: #e8e8e8; }
      #scenaria-browser-toolbar .sc-status {
        flex: 1 1 120px;
        min-width: 80px;
        margin: 0;
        font-size: 11px;
        color: #858585;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
      #scenaria-browser-toolbar .sc-status.recording { color: #f48771; }
      #scenaria-browser-toolbar .sc-status.paused { color: #cca700; }
      #scenaria-browser-toolbar .sc-actions { display: flex; gap: 4px; flex-shrink: 0; }
      #scenaria-browser-toolbar button {
        all: unset;
        cursor: pointer;
        display: inline-flex;
        align-items: center;
        gap: 5px;
        padding: 5px 9px;
        border-radius: 6px;
        border: 1px solid transparent;
        background: #3c3c3c;
        color: #e8e8e8;
        font: inherit;
        line-height: 1;
        transition: background 0.12s, border-color 0.12s;
      }
      #scenaria-browser-toolbar button:hover:not(:disabled) { background: #4a4a4a; }
      #scenaria-browser-toolbar button:disabled { opacity: 0.38; cursor: default; }
      #scenaria-browser-toolbar button.sc-primary:not(:disabled) {
        background: #094771;
        border-color: #007acc;
        color: #fff;
      }
      #scenaria-browser-toolbar button.sc-danger:not(:disabled):hover {
        background: #5a1f1f;
        border-color: #f14c4c;
      }
      #scenaria-browser-toolbar button svg { display: block; flex-shrink: 0; }
      #scenaria-browser-toolbar .sc-label { font-size: 11px; }
      @media (max-width: 640px) {
        #scenaria-browser-toolbar .sc-label { display: none; }
        #scenaria-browser-toolbar button { padding: 6px 8px; }
      }
    </style>
    <div class="sc-brand" data-drag="1">
      <span class="sc-dot"></span>
      <span class="sc-title">SCENARIA</span>
    </div>
    <p class="sc-status"></p>
    <div class="sc-actions">
      <button type="button" data-action="record">${ICONS.record}<span class="sc-label">Запись</span></button>
      <button type="button" data-action="pause">${ICONS.pause}<span class="sc-label">Пауза</span></button>
      <button type="button" data-action="stop" class="sc-danger">${ICONS.stop}<span class="sc-label">Стоп</span></button>
      <button type="button" data-action="picker">${ICONS.picker}<span class="sc-label">Элемент</span></button>
    </div>
  `;

  const dot = host.querySelector('.sc-dot');
  const status = host.querySelector('.sc-status');
  const brand = host.querySelector('.sc-brand');
  const btnPause = host.querySelector('button[data-action="pause"]');

  function render() {
    if (!status || !dot) return;
    const rec = state.recording;
    const paused = state.paused;
    dot.classList.toggle('recording', rec && !paused);
    dot.classList.toggle('paused', rec && paused);
    status.classList.toggle('recording', rec && !paused);
    status.classList.toggle('paused', rec && paused);

    if (rec && paused) {
      status.textContent = 'Пауза — можно выбрать элемент';
    } else if (rec) {
      const n = state.stepCount || 0;
      status.textContent = n > 0 ? `● Запись · ${n} шаг(ов)` : '● Идёт запись';
    } else {
      status.textContent = 'Браузер открыт — запись по кнопке «Запись»';
    }

    host.querySelectorAll('button[data-action]').forEach((btn) => {
      const action = btn.getAttribute('data-action');
      if (action === 'record') {
        btn.disabled = rec;
        btn.classList.remove('sc-primary');
      } else if (action === 'pause') btn.disabled = !rec;
      else if (action === 'stop') btn.disabled = false;
      else if (action === 'picker') btn.disabled = rec && !paused;
    });

    if (btnPause) {
      const label = btnPause.querySelector('.sc-label');
      if (rec) {
        btnPause.innerHTML = (paused ? ICONS.play : ICONS.pause) +
          `<span class="sc-label">${paused ? 'Продолжить' : 'Пауза'}</span>`;
      } else if (label) {
        label.textContent = 'Пауза';
      }
    }
  }

  host.addEventListener('click', (e) => {
    const btn = e.target.closest('button[data-action]');
    if (!btn || btn.disabled) return;
    e.preventDefault();
    e.stopPropagation();
    const action = btn.getAttribute('data-action');
    if (action === 'pause') {
      pending = state.paused ? 'resume' : 'pause';
    } else {
      pending = action;
    }
  }, true);

  brand?.addEventListener('pointerdown', (e) => {
    if (e.button !== 0) return;
    const rect = host.getBoundingClientRect();
    drag = { x: e.clientX, y: e.clientY, left: rect.left, top: rect.top };
    host.style.transform = 'none';
    host.style.left = `${rect.left}px`;
    host.style.top = `${rect.top}px`;
    brand.setPointerCapture(e.pointerId);
  });
  brand?.addEventListener('pointermove', (e) => {
    if (!drag) return;
    host.style.left = `${Math.max(4, drag.left + (e.clientX - drag.x))}px`;
    host.style.top = `${Math.max(4, drag.top + (e.clientY - drag.y))}px`;
  });
  const endDrag = (e) => {
    if (!drag) return;
    drag = null;
    try { brand?.releasePointerCapture(e.pointerId); } catch (_) {}
  };
  brand?.addEventListener('pointerup', endDrag);
  brand?.addEventListener('pointercancel', endDrag);

  const mount = () => {
    if (document.getElementById('scenaria-browser-toolbar')) return;
    (document.body || document.documentElement).appendChild(host);
    render();
  };

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', mount, { once: true });
  } else {
    mount();
  }

  window.__scenariaToolbar = {
    setState(patch) {
      Object.assign(state, patch || {});
      render();
    },
    takeAction() {
      const action = pending;
      pending = null;
      return action;
    },
  };
})();
