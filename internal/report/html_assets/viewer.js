(function () {
  const dataEl = document.getElementById('scenaria-report-data');
  if (!dataEl) return;
  const DATA = JSON.parse(dataEl.textContent || '{}');
  const STRINGS = {
    ru: {
      scenarios: 'Сценарии', timeline: 'Timeline', actionLog: 'Action log', inspector: 'Inspector',
      noFilter: 'Нет сценариев по фильтру', pickScenario: 'Выберите сценарий слева', pickStep: 'Кликните шаг',
      noStep: 'Нет шага', light: 'Лёгкий режим', full: 'Полный режим', trimmed: 'Скриншоты урезаны (лимит размера)',
      regression: 'Регрессия', history: 'История прогонов', openIde: 'Open in Scenaria IDE', rerunIde: 'Re-run in IDE',
      traceIde: 'Trace in IDE', exportNoFailed: 'Нет упавших сценариев',
      ciCompare: 'CI сравнение', newFailures: 'Новые падения', fixed: 'Исправлено', durationDelta: 'Δ длительность',
      pageContext: 'Страница', traceOffset: 'Смещение в trace',
      runDiff: 'Diff прогонов', traceAtStep: 'Trace на шаге', copyTraceSeek: 'Copy seek hint',
      domSnapshot: 'DOM', a11ySnapshot: 'Accessibility', traceTimeline: 'Trace timeline',
      stepTitle: 'Шаг #{n}', error: 'Ошибка', tips: 'Советы', screenshot: 'Скриншот',
      screenshotExpand: 'Открыть скриншот', screenshotHint: 'Нажмите для увеличения', screenshotClose: 'Закрыть',
      lastRun: 'Прошлый прогон', changed: '(изменился)', actions: 'Действия',
      flakyHistory: 'Шаг падал {n} раз в недавней истории', flakyInHistory: 'flaky (история)',
      historicallyUnstable: 'исторически нестабильно', failedTimesHistory: 'неудач в истории: {n}',
      stepRetry: 'повтор шага ×{n}', canceled: 'отменено', notStarted: 'не запущено',
      passed: 'успешно', failed: 'ошибка', skipped: 'пропущено', running: 'выполняется', notRun: 'не запущено',
      stepCol: 'Шаг',
      when: 'Когда', status: 'Статус', message: 'Сообщение', stepInTrace: 'шаг #{n}',
      copySelector: 'Copy selector', copyGherkin: 'Copy Gherkin', copyRerun: 'Copy re-run',
      traceFile: 'Trace file', traceDrop: 'Перетащите .zip trace сюда — offline-просмотр во вкладке Trace',
      traceZipRequired: 'Нужен Playwright trace .zip',
      traceDropOpen: 'Откройте trace.playwright.dev и перетащите файл: {name}',
      filterFailedOnly: 'Только упавшие', searchPlaceholder: 'Поиск…',
      tagPlaceholder: 'Фильтр по тегу (@smoke)', durationPlaceholder: 'Мин. длительность (ms)',
      exportFailed: 'Export failed .feature', network: 'Network', gherkin: 'Gherkin', selector: 'Selector',
      traceView: 'Trace', traceViewEmpty: 'Trace доступен только для упавших сценариев с архивом',
      traceViewTrimmed: 'События trace урезаны — откройте файл через trace.playwright.dev или IDE',
      traceZipParseFailed: 'Не удалось прочитать trace.zip',
      traceImported: 'Trace импортирован: {n} событий',
      openTraceExternal: 'Открыть в trace.playwright.dev',
      validation: 'Проверка селекторов',
      validateAction: 'Действие',
      validateMatches: 'Совпадения',
      validateFound: 'found',
      validateMissing: 'missing',
      validateWarning: 'warning',
      validateModeStatic: 'static',
      validateModeFlow: 'flow-aware',
    },
    en: {
      scenarios: 'Scenarios', timeline: 'Timeline', actionLog: 'Action log', inspector: 'Inspector',
      noFilter: 'No scenarios match filter', pickScenario: 'Select a scenario', pickStep: 'Click a step',
      noStep: 'No step', light: 'Light mode', full: 'Full mode', trimmed: 'Screenshots trimmed (size limit)',
      regression: 'Regression', history: 'Run history', openIde: 'Open in Scenaria IDE', rerunIde: 'Re-run in IDE',
      traceIde: 'Trace in IDE', exportNoFailed: 'No failed scenarios',
      ciCompare: 'CI compare', newFailures: 'New failures', fixed: 'Fixed', durationDelta: 'Duration Δ',
      pageContext: 'Page', traceOffset: 'Trace offset',
      runDiff: 'Run diff', traceAtStep: 'Trace at step', copyTraceSeek: 'Copy seek hint',
      domSnapshot: 'DOM', a11ySnapshot: 'Accessibility', traceTimeline: 'Trace timeline',
      stepTitle: 'Step #{n}', error: 'Error', tips: 'Tips', screenshot: 'Screenshot',
      screenshotExpand: 'Open screenshot', screenshotHint: 'Click to enlarge', screenshotClose: 'Close',
      lastRun: 'Previous run', changed: '(changed)', actions: 'Actions',
      flakyHistory: 'Step failed {n} times recently', flakyInHistory: 'flaky (history)',
      historicallyUnstable: 'historically unstable', failedTimesHistory: 'failed {n} times in history',
      stepRetry: 'step retry ×{n}', canceled: 'canceled', notStarted: 'not started',
      passed: 'passed', failed: 'failed', skipped: 'skipped', running: 'running', notRun: 'not started',
      stepCol: 'Step',
      when: 'When', status: 'Status', message: 'Message', stepInTrace: 'step #{n}',
      copySelector: 'Copy selector', copyGherkin: 'Copy Gherkin', copyRerun: 'Copy re-run',
      traceFile: 'Trace file', traceDrop: 'Drop .zip trace here — offline viewing in Trace tab',
      traceZipRequired: 'Playwright trace .zip required',
      traceDropOpen: 'Open trace.playwright.dev and drop file: {name}',
      filterFailedOnly: 'Failed only', searchPlaceholder: 'Search…',
      tagPlaceholder: 'Filter by tag (@smoke)', durationPlaceholder: 'Min duration (ms)',
      exportFailed: 'Export failed .feature', network: 'Network', gherkin: 'Gherkin', selector: 'Selector',
      traceView: 'Trace', traceViewEmpty: 'Trace is only available for failed scenarios with an archive',
      traceViewTrimmed: 'Trace events were trimmed — open the file in trace.playwright.dev or the IDE',
      traceZipParseFailed: 'Failed to parse trace.zip',
      traceImported: 'Trace imported: {n} events',
      openTraceExternal: 'Open in trace.playwright.dev',
      validation: 'Selector validation',
      validateAction: 'Action',
      validateMatches: 'Matches',
      validateFound: 'found',
      validateMissing: 'missing',
      validateWarning: 'warning',
      validateModeStatic: 'static',
      validateModeFlow: 'flow-aware',
    },
  };
  Object.assign(STRINGS.ru, {
    scenarios: '\u0421\u0446\u0435\u043d\u0430\u0440\u0438\u0438',
    timeline: '\u0428\u0430\u0433\u0438',
    actionLog: '\u0416\u0443\u0440\u043d\u0430\u043b \u0434\u0435\u0439\u0441\u0442\u0432\u0438\u0439',
    inspector: '\u0418\u043d\u0441\u043f\u0435\u043a\u0442\u043e\u0440',
    noFilter: '\u041d\u0435\u0442 \u0441\u0446\u0435\u043d\u0430\u0440\u0438\u0435\u0432 \u043f\u043e \u0444\u0438\u043b\u044c\u0442\u0440\u0443',
    pickScenario: '\u0412\u044b\u0431\u0435\u0440\u0438\u0442\u0435 \u0441\u0446\u0435\u043d\u0430\u0440\u0438\u0439 \u0441\u043b\u0435\u0432\u0430',
    pickStep: '\u041a\u043b\u0438\u043a\u043d\u0438\u0442\u0435 \u0448\u0430\u0433',
    noStep: '\u041d\u0435\u0442 \u0448\u0430\u0433\u0430',
    light: '\u041b\u0435\u0433\u043a\u0438\u0439 \u0440\u0435\u0436\u0438\u043c',
    full: '\u041f\u043e\u043b\u043d\u044b\u0439 \u0440\u0435\u0436\u0438\u043c',
    trimmed: '\u0421\u043a\u0440\u0438\u043d\u0448\u043e\u0442\u044b \u0443\u0440\u0435\u0437\u0430\u043d\u044b (\u043b\u0438\u043c\u0438\u0442 \u0440\u0430\u0437\u043c\u0435\u0440\u0430)',
    regression: '\u0420\u0435\u0433\u0440\u0435\u0441\u0441\u0438\u044f',
    history: '\u0418\u0441\u0442\u043e\u0440\u0438\u044f \u043f\u0440\u043e\u0433\u043e\u043d\u043e\u0432',
    openIde: '\u041e\u0442\u043a\u0440\u044b\u0442\u044c \u0432 Scenaria IDE',
    rerunIde: '\u041f\u043e\u0432\u0442\u043e\u0440\u0438\u0442\u044c \u0432 IDE',
    traceIde: '\u041e\u0442\u043a\u0440\u044b\u0442\u044c trace \u0432 IDE',
    exportNoFailed: '\u041d\u0435\u0442 \u0443\u043f\u0430\u0432\u0448\u0438\u0445 \u0441\u0446\u0435\u043d\u0430\u0440\u0438\u0435\u0432',
    ciCompare: 'CI-\u0441\u0440\u0430\u0432\u043d\u0435\u043d\u0438\u0435',
    newFailures: '\u041d\u043e\u0432\u044b\u0435 \u043f\u0430\u0434\u0435\u043d\u0438\u044f',
    fixed: '\u0418\u0441\u043f\u0440\u0430\u0432\u043b\u0435\u043d\u043e',
    durationDelta: '\u0414\u0435\u043b\u044c\u0442\u0430 \u0434\u043b\u0438\u0442\u0435\u043b\u044c\u043d\u043e\u0441\u0442\u0438',
    pageContext: '\u0421\u0442\u0440\u0430\u043d\u0438\u0446\u0430',
    traceOffset: '\u0421\u043c\u0435\u0449\u0435\u043d\u0438\u0435 \u0432 trace',
    runDiff: '\u0421\u0440\u0430\u0432\u043d\u0435\u043d\u0438\u0435 \u043f\u0440\u043e\u0433\u043e\u043d\u043e\u0432',
    traceAtStep: 'Trace \u043d\u0430 \u0448\u0430\u0433\u0435',
    copyTraceSeek: '\u041a\u043e\u043f\u0438\u0440\u043e\u0432\u0430\u0442\u044c seek-\u043f\u043e\u0434\u0441\u043a\u0430\u0437\u043a\u0443',
    domSnapshot: 'DOM-\u0441\u043d\u0438\u043c\u043e\u043a',
    a11ySnapshot: '\u0421\u043d\u0438\u043c\u043e\u043a \u0434\u043e\u0441\u0442\u0443\u043f\u043d\u043e\u0441\u0442\u0438',
    traceTimeline: '\u0422\u0430\u0439\u043c\u043b\u0430\u0439\u043d trace',
    stepTitle: '\u0428\u0430\u0433 #{n}',
    error: '\u041e\u0448\u0438\u0431\u043a\u0430',
    tips: '\u0421\u043e\u0432\u0435\u0442\u044b',
    screenshot: '\u0421\u043a\u0440\u0438\u043d\u0448\u043e\u0442',
    screenshotExpand: '\u041e\u0442\u043a\u0440\u044b\u0442\u044c \u0441\u043a\u0440\u0438\u043d\u0448\u043e\u0442',
    screenshotHint: '\u041d\u0430\u0436\u043c\u0438\u0442\u0435 \u0434\u043b\u044f \u0443\u0432\u0435\u043b\u0438\u0447\u0435\u043d\u0438\u044f',
    screenshotClose: '\u0417\u0430\u043a\u0440\u044b\u0442\u044c',
    lastRun: '\u041f\u0440\u043e\u0448\u043b\u044b\u0439 \u043f\u0440\u043e\u0433\u043e\u043d',
    changed: '(\u0438\u0437\u043c\u0435\u043d\u0438\u043b\u0441\u044f)',
    actions: '\u0414\u0435\u0439\u0441\u0442\u0432\u0438\u044f',
    flakyHistory: '\u0428\u0430\u0433 \u043f\u0430\u0434\u0430\u043b {n} \u0440\u0430\u0437 \u0432 \u043d\u0435\u0434\u0430\u0432\u043d\u0435\u0439 \u0438\u0441\u0442\u043e\u0440\u0438\u0438',
    flakyInHistory: 'flaky (\u0438\u0441\u0442\u043e\u0440\u0438\u044f)',
    canceled: '\u043e\u0442\u043c\u0435\u043d\u0435\u043d\u043e',
    notStarted: '\u043d\u0435 \u0437\u0430\u043f\u0443\u0449\u0435\u043d\u043e',
    stepCol: '\u0428\u0430\u0433',
    when: '\u041a\u043e\u0433\u0434\u0430',
    status: '\u0421\u0442\u0430\u0442\u0443\u0441',
    message: '\u0421\u043e\u043e\u0431\u0449\u0435\u043d\u0438\u0435',
    stepInTrace: '\u0448\u0430\u0433 #{n}',
    copySelector: '\u041a\u043e\u043f\u0438\u0440\u043e\u0432\u0430\u0442\u044c \u0441\u0435\u043b\u0435\u043a\u0442\u043e\u0440',
    copyGherkin: '\u041a\u043e\u043f\u0438\u0440\u043e\u0432\u0430\u0442\u044c Gherkin',
    copyRerun: '\u041a\u043e\u043f\u0438\u0440\u043e\u0432\u0430\u0442\u044c \u043a\u043e\u043c\u0430\u043d\u0434\u0443 \u043f\u043e\u0432\u0442\u043e\u0440\u0430',
    traceFile: '\u0424\u0430\u0439\u043b trace',
    traceDrop: '\u041f\u0435\u0440\u0435\u0442\u0430\u0449\u0438\u0442\u0435 .zip trace \u0441\u044e\u0434\u0430 - offline-\u043f\u0440\u043e\u0441\u043c\u043e\u0442\u0440 \u0432\u043e \u0432\u043a\u043b\u0430\u0434\u043a\u0435 Trace',
    traceZipRequired: '\u041d\u0443\u0436\u0435\u043d Playwright trace .zip',
    traceDropOpen: '\u041e\u0442\u043a\u0440\u043e\u0439\u0442\u0435 trace.playwright.dev \u0438 \u043f\u0435\u0440\u0435\u0442\u0430\u0449\u0438\u0442\u0435 \u0444\u0430\u0439\u043b: {name}',
    filterFailedOnly: '\u0422\u043e\u043b\u044c\u043a\u043e \u0443\u043f\u0430\u0432\u0448\u0438\u0435',
    searchPlaceholder: '\u041f\u043e\u0438\u0441\u043a...',
    tagPlaceholder: '\u0424\u0438\u043b\u044c\u0442\u0440 \u043f\u043e \u0442\u0435\u0433\u0443 (@smoke)',
    durationPlaceholder: '\u041c\u0438\u043d. \u0434\u043b\u0438\u0442\u0435\u043b\u044c\u043d\u043e\u0441\u0442\u044c (ms)',
    exportFailed: '\u042d\u043a\u0441\u043f\u043e\u0440\u0442 \u0443\u043f\u0430\u0432\u0448\u0438\u0445 .feature',
    network: '\u0421\u0435\u0442\u044c',
    gherkin: '\u0428\u0430\u0433 Gherkin',
    selector: '\u0421\u0435\u043b\u0435\u043a\u0442\u043e\u0440',
    traceView: 'Trace',
    traceViewEmpty: 'Trace \u0434\u043e\u0441\u0442\u0443\u043f\u0435\u043d \u0442\u043e\u043b\u044c\u043a\u043e \u0434\u043b\u044f \u0443\u043f\u0430\u0432\u0448\u0438\u0445 \u0441\u0446\u0435\u043d\u0430\u0440\u0438\u0435\u0432 \u0441 \u0430\u0440\u0445\u0438\u0432\u043e\u043c',
    traceViewTrimmed: '\u0421\u043e\u0431\u044b\u0442\u0438\u044f trace \u0443\u0440\u0435\u0437\u0430\u043d\u044b - \u043e\u0442\u043a\u0440\u043e\u0439\u0442\u0435 \u0444\u0430\u0439\u043b \u0447\u0435\u0440\u0435\u0437 trace.playwright.dev \u0438\u043b\u0438 IDE',
    traceZipParseFailed: '\u041d\u0435 \u0443\u0434\u0430\u043b\u043e\u0441\u044c \u043f\u0440\u043e\u0447\u0438\u0442\u0430\u0442\u044c trace.zip',
    traceImported: 'Trace \u0438\u043c\u043f\u043e\u0440\u0442\u0438\u0440\u043e\u0432\u0430\u043d: {n} \u0441\u043e\u0431\u044b\u0442\u0438\u0439',
    openTraceExternal: '\u041e\u0442\u043a\u0440\u044b\u0442\u044c \u0432 trace.playwright.dev',
    modeFull: '\u041f\u043e\u043b\u043d\u044b\u0439',
    modeLight: '\u041b\u0435\u0433\u043a\u0438\u0439',
    modeUnavailable: '\u0421\u043d\u0430\u0447\u0430\u043b\u0430 \u0441\u0433\u0435\u043d\u0435\u0440\u0438\u0440\u0443\u0439\u0442\u0435 paired HTML-\u043e\u0442\u0447\u0435\u0442 \u0434\u043b\u044f \u044d\u0442\u043e\u0433\u043e \u0440\u0435\u0436\u0438\u043c\u0430',
    searchLabel: '\u041f\u043e\u0438\u0441\u043a \u0441\u0446\u0435\u043d\u0430\u0440\u0438\u0435\u0432',
    failedOnlyLabel: '\u0422\u043e\u043b\u044c\u043a\u043e \u0443\u043f\u0430\u0432\u0448\u0438\u0435',
    tagFilterLabel: '\u0424\u0438\u043b\u044c\u0442\u0440 \u043f\u043e \u0442\u0435\u0433\u0443',
    durationFilterLabel: '\u041c\u0438\u043d\u0438\u043c\u0430\u043b\u044c\u043d\u0430\u044f \u0434\u043b\u0438\u0442\u0435\u043b\u044c\u043d\u043e\u0441\u0442\u044c \u0432 \u043c\u0438\u043b\u043b\u0438\u0441\u0435\u043a\u0443\u043d\u0434\u0430\u0445',
  });
  Object.assign(STRINGS.en, {
    durationDelta: 'Duration delta',
    screenshotExpand: 'Open screenshot',
    screenshotHint: 'Click to enlarge',
    screenshotClose: 'Close',
    traceDrop: 'Drop .zip trace here - offline viewing in Trace tab',
    searchPlaceholder: 'Search...',
    traceViewTrimmed: 'Trace events were trimmed - open the file in trace.playwright.dev or the IDE',
    modeFull: 'Full',
    modeLight: 'Light',
    modeUnavailable: 'Generate the paired HTML report for this mode first',
    searchLabel: 'Search scenarios',
    failedOnlyLabel: 'Failed only',
    tagFilterLabel: 'Filter by tag',
    durationFilterLabel: 'Minimum duration in milliseconds',
  });
  function t(key) {
    const loc = DATA.locale === 'en' ? 'en' : 'ru';
    return (STRINGS[loc] || STRINGS.ru)[key] || key;
  }
  function tf(key, vars) {
    let s = t(key);
    if (vars) {
      Object.keys(vars).forEach((k) => {
        s = s.replace(new RegExp('\\{' + k + '\\}', 'g'), vars[k]);
      });
    }
    return s;
  }
  const state = {
    scenarioId: null,
    stepIndex: null,
    centerTab: 'timeline',
    filterFailed: false,
    filterTag: '',
    filterDurationMin: 0,
    search: '',
    importedTrace: null,
  };

  const $ = (sel) => document.querySelector(sel);
  const esc = (s) => String(s ?? '').replace(/[&<>"]/g, (c) => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;'}[c]));

  function screenshotLightboxEl() {
    return $('#screenshot-lightbox');
  }

  function openScreenshotLightbox(src, caption) {
    const box = screenshotLightboxEl();
    if (!box || !src) return;
    const img = $('#screenshot-lightbox-img');
    const cap = $('#screenshot-lightbox-caption');
    const backdrop = box.querySelector('.screenshot-lightbox-backdrop');
    const closeBtn = box.querySelector('.screenshot-lightbox-close');
    if (img) {
      img.src = src;
      img.alt = caption || t('screenshot');
    }
    if (cap) cap.textContent = caption || '';
    if (backdrop) backdrop.setAttribute('aria-label', t('screenshotClose'));
    if (closeBtn) closeBtn.setAttribute('aria-label', t('screenshotClose'));
    box.hidden = false;
    box.setAttribute('aria-hidden', 'false');
    document.body.classList.add('screenshot-lightbox-open');
    closeBtn?.focus();
  }

  function closeScreenshotLightbox() {
    const box = screenshotLightboxEl();
    if (!box || box.hidden) return;
    box.hidden = true;
    box.setAttribute('aria-hidden', 'true');
    document.body.classList.remove('screenshot-lightbox-open');
    const img = $('#screenshot-lightbox-img');
    if (img) img.removeAttribute('src');
  }

  function bindScreenshotLightbox() {
    const box = screenshotLightboxEl();
    if (!box || box.dataset.bound) return;
    box.dataset.bound = '1';
    box.querySelector('.screenshot-lightbox-backdrop')?.addEventListener('click', closeScreenshotLightbox);
    box.querySelector('.screenshot-lightbox-close')?.addEventListener('click', closeScreenshotLightbox);
  }

  function activateOnEnterSpace(el, fn) {
    el.addEventListener('click', fn);
    el.addEventListener('keydown', (e) => {
      if (e.key !== 'Enter' && e.key !== ' ') return;
      e.preventDefault();
      fn();
    });
  }

  function normalizeStatus(s) {
    const status = String(s || '').trim().toLowerCase();
    if (!status) return 'skipped';
    if (['passed', 'pass', 'ok', 'success'].includes(status)) return 'passed';
    if (['failed', 'fail', 'broken', 'error'].includes(status)) return 'failed';
    if (['canceled', 'cancelled', 'aborted'].includes(status)) return 'canceled';
    if (['not-started', 'notstarted', 'not-start', 'notstart', 'not started', 'not run', 'not-run'].includes(status)) return 'not-started';
    if (['running', 'in-progress', 'inprogress'].includes(status)) return 'running';
    return 'skipped';
  }

  function statusClass(s) {
    return normalizeStatus(s);
  }

  function statusLabel(s) {
    const status = normalizeStatus(s);
    if (status === 'passed') return t('passed');
    if (status === 'failed') return t('failed');
    if (status === 'canceled') return t('canceled');
    if (status === 'not-started') return t('notStarted');
    if (status === 'running') return t('running');
    return t('skipped');
  }

  function renderStatusBadge(status, extraClass) {
    const cls = ['badge', statusClass(status), extraClass || ''].filter(Boolean).join(' ');
    return '<span class="' + esc(cls) + '">' + esc(statusLabel(status)) + '</span>';
  }

  function historyMetaText(sc) {
    const flaky = scenarioFlakyStat(sc);
    if (!flaky) return '';
    if (flaky.Flaky) {
      return t('historicallyUnstable');
    }
    if ((flaky.Failures || 0) > 0) {
      return tf('failedTimesHistory', { n: flaky.Failures });
    }
    return '';
  }

  function normScenarioKey(s) {
    return String(s || '').replace(/\\/g, '/').toLowerCase();
  }

  function scenarioFlakyStat(sc) {
    const flakyList = (DATA.flaky && DATA.flaky.scenarios) || [];
    const keys = [sc.id, (sc.feature_path || '') + '::' + (sc.scenario || '')].map(normScenarioKey);
    return flakyList.find((x) => keys.includes(normScenarioKey(x.path)));
  }

  function flakyScenariosInRun() {
    return (DATA.scenarios || []).filter((sc) => historyMetaText(sc));
  }

  function filteredScenarios() {
    return (DATA.scenarios || []).filter((sc) => {
      if (state.filterFailed && sc.status !== 'failed') return false;
      if (state.filterTag && !(sc.tags || []).some((t) => t.includes(state.filterTag))) return false;
      if (state.filterDurationMin > 0 && (sc.duration_ms || 0) < state.filterDurationMin) return false;
      if (state.search) {
        const hay = (sc.scenario + ' ' + sc.feature_path).toLowerCase();
        if (!hay.includes(state.search.toLowerCase())) return false;
      }
      return true;
    });
  }

  function postBridge(path, body) {
    const base = (DATA.bridge_url || '').replace(/\/$/, '');
    if (!base) return Promise.reject(new Error('no bridge'));
    return fetch(base + path, {
      method: 'POST',
      headers: Object.assign(
        { 'Content-Type': 'application/json' },
        DATA.bridge_token ? { 'X-Scenaria-Bridge-Token': DATA.bridge_token } : {},
      ),
      body: JSON.stringify(body),
    });
  }

  function renderSparkline(values) {
    if (!values || values.length < 2) return '';
    const max = Math.max.apply(null, values.concat([1]));
    const bars = values.map((v, i) => {
      const h = Math.max(4, Math.round((v / max) * 24));
      const cls = 'spark-bar' + (i === values.length - 1 ? ' current' : '');
      return '<span class="' + cls + '" style="height:' + h + 'px" title="' + v + ' ms"></span>';
    }).join('');
    return '<div class="sparkline">' + bars + '</div>';
  }

  function renderCICompare() {
    const box = $('#ci-compare');
    if (!box) return;
    const ci = DATA.ci_compare;
    if (!ci) {
      box.innerHTML = '';
      box.hidden = true;
      return;
    }
    box.hidden = false;
    let html = '<strong>' + t('ciCompare') + '</strong>';
    if (ci.previous_at) html += ' <span class="sub">@ ' + esc(ci.previous_at) + '</span>';
    if ((ci.new_failures || []).length) {
      html += '<div class="ci-block fail">' + t('newFailures') + ': ' +
        ci.new_failures.map((x) => esc(x.scenario)).join(', ') + '</div>';
    }
    if ((ci.fixed || []).length) {
      html += '<div class="ci-block pass">' + t('fixed') + ': ' +
        ci.fixed.map((x) => esc(x.scenario)).join(', ') + '</div>';
    }
    if ((ci.duration_deltas || []).length) {
      html += '<div class="ci-block">' + t('durationDelta') + ': ' +
        ci.duration_deltas.slice(0, 3).map((d) => esc(d.scenario) + ' ' + (d.delta_ms > 0 ? '+' : '') + d.delta_ms + 'ms').join('; ') + '</div>';
    }
    box.innerHTML = html;
  }

  function renderRunDiff(sc) {
    const diff = sc.run_diff;
    if (!diff || !(diff.rows || []).length) return '';
    const cols = diff.columns || [];
    let html = '<div class="section"><label>' + t('runDiff') + '</label><table class="run-diff"><thead><tr><th>' + t('stepCol') + '</th>';
    cols.forEach((c) => {
      html += '<th title="' + esc(c.at || '') + '">' + esc(c.label) + (c.current ? ' *' : '') + '</th>';
    });
    html += '</tr></thead><tbody>';
    diff.rows.forEach((row) => {
      html += '<tr><td>#' + (row.index + 1) + ' ' + esc(row.text || '') + '</td>';
      (row.cells || []).forEach((cell) => {
        html += '<td class="diff-' + esc(cell) + '">' + esc(cell) + '</td>';
      });
      html += '</tr>';
    });
    return html + '</tbody></table></div>';
  }

  function formatTraceSeek(ms) {
    if (ms == null || ms < 0) return '';
    if (ms < 1000) return ms + 'ms';
    return (ms / 1000).toFixed(2) + 's';
  }

  function updateHeaderHeightVar() {
    const header = document.querySelector('header');
    if (!header) return;
    document.documentElement.style.setProperty('--report-header-h', Math.ceil(header.getBoundingClientRect().height) + 'px');
  }

  function applyStaticLocale() {
    const search = $('#filter-search');
    if (search) {
      search.placeholder = t('searchPlaceholder');
      search.setAttribute('aria-label', t('searchLabel'));
    }
    const tag = $('#filter-tag');
    if (tag) {
      tag.placeholder = t('tagPlaceholder');
      tag.setAttribute('aria-label', t('tagFilterLabel'));
    }
    const dur = $('#filter-duration');
    if (dur) {
      dur.placeholder = t('durationPlaceholder');
      dur.setAttribute('aria-label', t('durationFilterLabel'));
    }
    const cb = $('#filter-failed');
    if (cb) cb.setAttribute('aria-label', t('failedOnlyLabel'));
    const failedLbl = document.querySelector('label[for="filter-failed"]') || $('#filter-failed')?.parentElement;
    if (failedLbl && failedLbl.tagName === 'LABEL') {
      failedLbl.lastChild.textContent = ' ' + t('filterFailedOnly');
    } else {
      const cb = $('#filter-failed');
      if (cb && cb.parentElement) cb.parentElement.childNodes[cb.parentElement.childNodes.length - 1].textContent = ' ' + t('filterFailedOnly');
    }
    const exp = $('#export-failed');
    if (exp) exp.textContent = t('exportFailed');
    const toggle = $('#toggle-sidebar');
    if (toggle) toggle.textContent = t('scenarios');
    const full = $('#mode-full');
    if (full) full.textContent = t('modeFull');
    const light = $('#mode-light');
    if (light) light.textContent = t('modeLight');
  }

  function reportStateHash() {
    const params = new URLSearchParams();
    if (state.scenarioId) params.set('scenario', state.scenarioId);
    if (state.stepIndex != null) params.set('step', String(state.stepIndex));
    if (state.centerTab) params.set('tab', state.centerTab);
    if (state.search) params.set('q', state.search);
    if (state.filterFailed) params.set('failed', '1');
    if (state.filterTag) params.set('tag', state.filterTag);
    if (state.filterDurationMin) params.set('duration', String(state.filterDurationMin));
    const raw = params.toString();
    return raw ? '#state=' + raw : '';
  }

  function restoreReportStateFromHash() {
    const hash = window.location.hash || '';
    if (!hash.startsWith('#state=')) return false;
    const params = new URLSearchParams(hash.slice('#state='.length));
    state.scenarioId = params.get('scenario') || null;
    const step = params.get('step');
    state.stepIndex = step == null ? null : Number(step);
    state.centerTab = params.get('tab') || 'timeline';
    state.search = params.get('q') || '';
    state.filterFailed = params.get('failed') === '1';
    state.filterTag = params.get('tag') || '';
    state.filterDurationMin = Number(params.get('duration')) || 0;
    const search = $('#filter-search');
    if (search) search.value = state.search;
    const failed = $('#filter-failed');
    if (failed) failed.checked = state.filterFailed;
    const tag = $('#filter-tag');
    if (tag) tag.value = state.filterTag;
    const dur = $('#filter-duration');
    if (dur && state.filterDurationMin) dur.value = String(state.filterDurationMin);
    const sc = (DATA.scenarios || []).find((x) => x.id === state.scenarioId);
    if (!sc) return false;
    if (!(sc.steps || []).some((st) => st.index === state.stepIndex)) {
      state.stepIndex = (sc.steps || [])[0]?.index ?? null;
    }
    return true;
  }

  function modeHref(mode) {
    const links = DATA.mode_links || {};
    if (mode === 'full') return links.full_href || '';
    return links.light_href || '';
  }

  function modeAvailable(mode) {
    const links = DATA.mode_links || {};
    if (mode === 'full') return !!links.full_available;
    return !!links.light_available;
  }

  function applyModeSwitch() {
    const links = DATA.mode_links || {};
    const current = links.current || (DATA.light_mode ? 'light' : 'full');
    document.querySelectorAll('.mode-option').forEach((btn) => {
      const mode = btn.dataset.mode;
      const active = mode === current;
      const available = active || modeAvailable(mode);
      btn.classList.toggle('active', active);
      btn.disabled = !available;
      btn.title = available ? '' : t('modeUnavailable');
      btn.setAttribute('aria-pressed', active ? 'true' : 'false');
    });
  }

  function renderHistoryRuns(sc) {
    const runs = sc.history_runs || [];
    if (!runs.length) return '';
    let html = '<div class="section"><label>' + t('history') + '</label><table class="history-table"><thead><tr><th>' + t('when') + '</th><th>' + t('status') + '</th><th>ms</th><th>' + t('stepCol') + '</th><th>' + t('message') + '</th></tr></thead><tbody>';
    runs.forEach((run, i) => {
      const changed = i === 0 && sc.history && sc.history.changed ? ' class="changed"' : '';
      html += '<tr' + changed + '><td>' + esc(run.at) + '</td><td>' + esc(statusLabel(run.status)) + '</td><td>' +
        (run.duration_ms || '—') + '</td><td>' +
        (run.failed_step != null ? '#' + (run.failed_step + 1) : '—') + '</td><td>' + esc(run.message || '') + '</td></tr>';
    });
    return html + '</tbody></table></div>';
  }

  function renderHeader() {
    const s = DATA.summary || {};
    $('#brand-title').textContent = (DATA.brand || 'Scenaria') + ' Report';
    $('#meta-line').textContent = [
      DATA.generated_at || '',
      'Mode: ' + (DATA.mode || ''),
      ((DATA.scenarios || []).length || s.scenarios || 0) + ' scenarios',
    ].join(' · ');
    $('#stat-pass').textContent = (s.passed || 0) + ' passed';
    $('#stat-fail').textContent = (s.failed || 0) + ' failed';
    $('#stat-skip').textContent = (s.skipped || 0) + ' skipped';
    $('#stat-cancel').textContent = (s.canceled || 0) + ' ' + t('canceled');
    $('#stat-not-started').textContent = (s.not_started || 0) + ' ' + t('notStarted');
    const flaky = flakyScenariosInRun().length;
    $('#stat-flaky').textContent = flaky ? flaky + ' ' + t('historicallyUnstable') : '';
    const slow = (DATA.slow_steps || [])[0];
    $('#stat-slow').textContent = slow ? 'Slowest: ' + slow.duration_ms + 'ms' : '';
    $('#light-badge').textContent = DATA.light_mode ? t('light') : t('full');
    if (DATA.artifacts_trimmed) {
      $('#stat-flaky').textContent = ($('#stat-flaky').textContent + ' · ' + t('trimmed')).trim();
    }
    document.documentElement.lang = DATA.locale === 'en' ? 'en' : 'ru';
    renderCICompare();
    applyStaticLocale();
    const sb = document.querySelector('.panel.sidebar .panel-head');
    if (sb) sb.textContent = t('scenarios');
    const insp = $('#inspector-title');
    if (insp) insp.textContent = t('inspector');
    $('#tab-timeline').textContent = t('timeline');
    $('#tab-actionlog').textContent = t('actionLog');
    applyModeSwitch();
    updateCenterTabs();
    updateHeaderHeightVar();
  }

  function scenarioTraceEvents(sc) {
    if (!sc) return [];
    if (state.importedTrace && state.importedTrace.scenarioId === sc.id && state.importedTrace.events) {
      return state.importedTrace.events;
    }
    return sc.trace_events || [];
  }

  function scenarioHasTraceView(sc) {
    if (!sc || sc.status !== 'failed') return false;
    if (state.importedTrace && state.importedTrace.scenarioId === sc.id && (state.importedTrace.events || []).length) {
      return true;
    }
    return !!(sc.trace_path || (sc.trace_events && sc.trace_events.length));
  }

  function updateCenterTabs() {
    const sc = (DATA.scenarios || []).find((x) => x.id === state.scenarioId);
    const traceTab = $('#tab-trace');
    const hasTrace = scenarioHasTraceView(sc);
    if (traceTab) {
      traceTab.hidden = !hasTrace;
      traceTab.textContent = t('traceView');
    }
    if (!hasTrace && state.centerTab === 'trace') {
      state.centerTab = 'timeline';
    }
    document.querySelectorAll('.tab').forEach((x) => {
      const active = x.dataset.tab === state.centerTab;
      x.classList.toggle('active', active);
      x.setAttribute('aria-selected', active ? 'true' : 'false');
      x.tabIndex = active ? 0 : -1;
    });
    const tl = $('#timeline');
    const al = $('#actionlog');
    const tv = $('#traceview');
    if (tl) tl.hidden = state.centerTab !== 'timeline';
    if (al) al.hidden = state.centerTab !== 'actionlog';
    if (tv) tv.hidden = state.centerTab !== 'trace';
  }

  function renderAllViews() {
    renderTimeline();
    renderActionLog();
    renderTraceView();
    renderInspector();
    updateCenterTabs();
  }

  function renderTraceView() {
    const sc = (DATA.scenarios || []).find((x) => x.id === state.scenarioId);
    const wrap = $('#traceview');
    if (!wrap) return;
    if (!sc || !scenarioHasTraceView(sc)) {
      wrap.innerHTML = '<div class="empty">' + t('traceViewEmpty') + '</div>';
      return;
    }
    const events = scenarioTraceEvents(sc);
    const steps = sc.steps || [];
    const maxMs = Math.max(
      sc.duration_ms || 0,
      events.reduce((m, e) => Math.max(m, e.offset_ms || 0), 0),
      steps.reduce((m, st) => Math.max(m, st.trace_offset_ms || 0), 0),
      1,
    );
    let html = '<div class="trace-mini">';
    if (!events.length && sc.trace_path) {
      html += '<div class="trace-hint">' + t('traceViewTrimmed') + '</div>';
    }
    html += '<div class="trace-rail">';
    steps.forEach((st) => {
      if (st.trace_offset_ms == null) return;
      const pct = Math.min(100, (st.trace_offset_ms / maxMs) * 100);
      html += '<div class="trace-step-marker' + (st.index === state.stepIndex ? ' active' : '') +
        '" style="left:' + pct + '%" data-idx="' + st.index + '"></div>';
    });
    events.forEach((ev, i) => {
      const pct = Math.min(100, ((ev.offset_ms || 0) / maxMs) * 100);
      const active = steps.find((st) => st.index === state.stepIndex);
      const isActive = active && Math.abs((active.trace_offset_ms || 0) - (ev.offset_ms || 0)) < 80;
      html += '<div class="trace-tick' + (isActive ? ' active' : '') + '" style="left:' + pct +
        '%" data-ms="' + ev.offset_ms + '" data-i="' + i + '"></div>';
    });
    html += '</div><div class="trace-split">';
    html += '<div class="trace-events-panel"><h4>' + t('traceTimeline') + '</h4>';
    if (!events.length) {
      html += '<div class="mono sub">' + esc(sc.trace_command || sc.trace_path || '') + '</div>';
    } else {
      events.forEach((ev) => {
        const active = steps.find((st) => st.index === state.stepIndex);
        const cls = active && Math.abs((active.trace_offset_ms || 0) - (ev.offset_ms || 0)) < 80 ? ' active' : '';
        html += '<div class="trace-event-row' + cls + '" data-ms="' + ev.offset_ms + '"><span class="ms">+' +
          ev.offset_ms + 'ms</span> ' + esc(ev.title || ev.kind || '') + '</div>';
      });
    }
    html += '</div><div class="trace-steps-panel"><h4>' + t('timeline') + '</h4>';
    steps.forEach((st) => {
      const cls = [statusClass(st.status), st.index === state.stepIndex ? 'active' : ''].filter(Boolean).join(' ');
      html += '<div class="trace-step-row ' + cls + '" data-idx="' + st.index + '"><span>#' + (st.index + 1) +
        '</span> ' + esc(st.gherkin || st.text) +
        (st.trace_offset_ms != null ? ' <span class="sub">@ ' + st.trace_offset_ms + 'ms</span>' : '') + '</div>';
    });
    html += '</div></div></div>';
    wrap.innerHTML = html;
    wrap.querySelectorAll('.trace-event-row, .trace-tick').forEach((el) => {
      el.addEventListener('click', () => {
        const ms = Number(el.dataset.ms);
        state.stepIndex = stepForTraceOffset(sc, ms);
        $('#inspector-panel').classList.add('open');
        renderAllViews();
      });
    });
    wrap.querySelectorAll('.trace-step-row, .trace-step-marker').forEach((el) => {
      el.addEventListener('click', () => {
        state.stepIndex = Number(el.dataset.idx);
        $('#inspector-panel').classList.add('open');
        renderAllViews();
      });
    });
  }

  function renderScenarioList() {
    const list = $('#scenario-list');
    const items = filteredScenarios();
    if (!items.length) {
      list.innerHTML = '<div class="empty">' + t('noFilter') + '</div>';
      return;
    }
    list.innerHTML = items.map((sc) => {
      const tags = (sc.tags || []).slice(0, 3).map((t) => '<span class="badge">' + esc(t) + '</span>').join('');
      const spark = renderSparkline(sc.duration_sparkline);
      const active = sc.id === state.scenarioId;
      const historyMeta = historyMetaText(sc);
      const historyLine = historyMeta ? '<div class="scenario-meta-history">' + esc(historyMeta) + '</div>' : '';
      return '<div class="scenario-item' + (active ? ' active' : '') + '" role="button" tabindex="0" aria-selected="' + (active ? 'true' : 'false') + '" aria-label="' + esc(sc.scenario + ' — ' + statusLabel(sc.status)) + '" data-id="' + esc(sc.id) + '">' +
        renderStatusBadge(sc.status) +
        '<div class="title">' + esc(sc.scenario) + spark + '</div>' +
        '<div class="sub">' + esc(sc.feature_path) + '</div>' +
        historyLine +
        (tags ? '<div class="sub">' + tags + '</div>' : '') +
      '</div>';
    }).join('');
    list.querySelectorAll('.scenario-item').forEach((el) => {
      activateOnEnterSpace(el, () => selectScenario(el.dataset.id));
    });
  }

  function selectScenario(id, stepIndex) {
    state.scenarioId = id;
    const sc = (DATA.scenarios || []).find((x) => x.id === id);
    if (!sc) return;
    if (stepIndex == null) {
      const failed = sc.steps.findIndex((st) => st.status === 'failed');
      stepIndex = failed >= 0 ? failed : 0;
    }
    state.stepIndex = stepIndex;
    $('#sidebar').classList.remove('mobile-open');
    renderScenarioList();
    renderAllViews();
    $('#inspector-panel').classList.add('open');
  }

  function renderTimeline() {
    const sc = (DATA.scenarios || []).find((x) => x.id === state.scenarioId);
    const wrap = $('#timeline');
    if (!sc) {
      wrap.innerHTML = '<div class="empty">' + t('pickScenario') + '</div>';
      return;
    }
    $('#timeline-title').textContent = sc.scenario + (sc.example_index ? ' (example #' + sc.example_index + ')' : '') +
      (sc.duration_ms ? ' · ' + sc.duration_ms + ' ms' : '');
    const spark = renderSparkline(sc.duration_sparkline);
    const steps = sc.steps || [];
    const body = [];
    const renderStepNode = (st) => {
      const flaky = st.flaky_failures >= 2 ? '<span class="flaky-pill">' + esc(tf('failedTimesHistory', { n: st.flaky_failures })) + '</span> ' : '';
      const iteration = st.iteration_path ? '<span class="retry-pill">' + esc(st.iteration_path) + '</span> ' : '';
      const terminal = st.terminal_action ? '<span class="retry-pill">' + esc(st.terminal_action) + '</span> ' : '';
      const retries = st.retry_attempts > 0 ? '<span class="retry-pill">' + esc(tf('stepRetry', { n: st.retry_attempts })) + '</span> ' : '';
      const stepSpark = renderSparkline(st.duration_sparkline);
      const active = st.index === state.stepIndex;
      return '<div class="step-node ' + statusClass(st.status) + (st.flaky_failures >= 2 ? ' flaky' : '') + (active ? ' active' : '') + '" role="button" tabindex="0" aria-current="' + (active ? 'step' : 'false') + '" aria-label="' + esc('#' + (st.index + 1) + ' ' + (st.gherkin || st.text) + ' — ' + statusLabel(st.status)) + '" data-idx="' + st.index + '">' +
        '<div class="line">#' + (st.index + 1) + (st.line ? ' · line ' + st.line : '') + ' ' + flaky + iteration + terminal + retries + stepSpark + '</div>' +
        '<div class="text">' + esc(st.gherkin || st.text) + '</div>' +
        (st.duration_ms ? '<div class="dur">' + st.duration_ms + ' ms</div>' : '') +
      '</div>';
    };
    for (let i = 0; i < steps.length;) {
      const st = steps[i];
      const path = st.iteration_path || '';
      if (!path) {
        body.push(renderStepNode(st));
        i += 1;
        continue;
      }
      const group = [st];
      let j = i + 1;
      while (j < steps.length && steps[j].iteration_path === path) {
        group.push(steps[j]);
        j += 1;
      }
      const failed = group.some((step) => normalizeStatus(step.status) === 'failed');
      const active = group.some((step) => step.index === state.stepIndex);
      const open = failed || active;
      const groupBody = group.map((step) => renderStepNode(step)).join('');
      body.push('<details class="iteration-group' + (failed ? ' failed' : '') + '" ' + (open ? 'open' : '') + '>' +
        '<summary class="iteration-summary">' +
        '<span class="iteration-label">' + esc(path) + '</span>' +
        '<span class="iteration-meta">' + esc(group.length + ' steps') + '</span>' +
        '</summary>' +
        '<div class="iteration-steps">' + groupBody + '</div>' +
      '</details>');
      i = j;
    }
    wrap.innerHTML = (spark ? '<div class="timeline-spark">' + spark + '</div>' : '') + '<div class="timeline">' + body.join('') + '</div>';
    wrap.querySelectorAll('.step-node').forEach((el) => {
      activateOnEnterSpace(el, () => {
        state.stepIndex = Number(el.dataset.idx);
        $('#inspector-panel').classList.add('open');
        renderAllViews();
      });
    });
  }

  function renderActionLog() {
    const sc = (DATA.scenarios || []).find((x) => x.id === state.scenarioId);
    const wrap = $('#actionlog');
    if (!sc) {
      wrap.innerHTML = '<div class="empty">' + t('pickScenario') + '</div>';
      return;
    }
    const rows = (sc.steps || []).map((st) => {
      const cls = [statusClass(st.status), st.index === state.stepIndex ? 'active' : ''].filter(Boolean).join(' ');
      const spark = renderSparkline(st.duration_sparkline);
      return '<tr class="' + cls + '" role="button" tabindex="0" aria-label="' + esc('#' + (st.index + 1) + ' ' + (st.gherkin || st.text) + ' — ' + statusLabel(st.status)) + '" data-idx="' + st.index + '"><td>#' + (st.index + 1) + '</td><td>' + esc(statusLabel(st.status)) +
        '</td><td>' + esc(st.gherkin || st.text) + spark + '</td><td>' + (st.duration_ms || '—') + '</td></tr>';
    }).join('');
    wrap.innerHTML = '<table class="action-log"><thead><tr><th>#</th><th>Status</th><th>Step</th><th>ms</th></tr></thead><tbody>' +
      rows + '</tbody></table>';
    wrap.querySelectorAll('tr[data-idx]').forEach((el) => {
      activateOnEnterSpace(el, () => {
        state.stepIndex = Number(el.dataset.idx);
        $('#inspector-panel').classList.add('open');
        renderAllViews();
      });
    });
  }

  function stepForTraceOffset(sc, ms) {
    const steps = sc.steps || [];
    if (!steps.length) return 0;
    const hasOffsets = steps.some((st) => st.trace_offset_ms != null);
    if (hasOffsets) {
      let best = steps[0].index;
      let bestDelta = Infinity;
      steps.forEach((st) => {
        const off = st.trace_offset_ms || 0;
        const d = Math.abs(off - ms);
        if (d < bestDelta) {
          bestDelta = d;
          best = st.index;
        }
      });
      return best;
    }
    const events = scenarioTraceEvents(sc);
    if (!events.length) return steps[0].index;
    let evIdx = 0;
    let bestDelta = Infinity;
    events.forEach((ev, i) => {
      const d = Math.abs((ev.offset_ms || 0) - ms);
      if (d < bestDelta) {
        bestDelta = d;
        evIdx = i;
      }
    });
    return steps[Math.min(evIdx, steps.length - 1)].index;
  }

  function renderTraceTimeline(sc, activeStep) {
    const events = scenarioTraceEvents(sc);
    if (!events.length) return '';
    let html = '<div class="section"><label>' + t('traceTimeline') + '</label><div class="trace-timeline">';
    events.forEach((ev) => {
      const active = activeStep && Math.abs((activeStep.trace_offset_ms || 0) - ev.offset_ms) < 80;
      html += '<div class="trace-event' + (active ? ' active' : '') + '" data-ms="' + ev.offset_ms + '">' +
        '<span class="ms">+' + ev.offset_ms + 'ms</span> ' + esc(ev.title || ev.kind || '') + '</div>';
    });
    return html + '</div></div>';
  }

  function renderRegressions(sc) {
    const items = sc.regressions || [];
    if (!items.length) return '';
    return '<div class="section"><label>' + t('regression') + '</label>' +
      items.map((r) => '<div class="regression"><strong>' + esc(r.kind) + '</strong> #' + (r.step + 1) +
        ' ' + esc(r.text || '') + '<br><span class="sub">' + esc(r.detail || '') + '</span></div>').join('') +
      '</div>';
  }

  function validationStatusLabel(status) {
    if (status === 'found') return t('validateFound');
    if (status === 'missing') return t('validateMissing');
    if (status === 'warning') return t('validateWarning');
    return status || '';
  }

  function validationModeLabel(mode) {
    if (mode === 'flow') return t('validateModeFlow');
    if (mode === 'static') return t('validateModeStatic');
    return mode || '';
  }

  function renderValidationBanner(sc) {
    if (!sc.validation_limitation && !sc.validation_mode) return '';
    let html = '<div class="validation-banner">';
    if (sc.validation_mode) {
      html += '<span class="validation-mode">' + esc(validationModeLabel(sc.validation_mode)) + '</span>';
    }
    if (sc.validation_limitation) {
      html += '<span class="validation-limitation">' + esc(sc.validation_limitation) + '</span>';
    }
    return html + '</div>';
  }

  function renderScenarioValidationTable(sc) {
    const rows = (sc.steps || []).filter((step) => step.validation);
    if (!rows.length) return '';
    const showAction = rows.some((step) => step.validation && step.validation.action_kind);
    const showMatches = rows.some((step) => step.validation && (step.validation.match_count || 0) > 0);
    let html = '<div class="section"><label>' + t('validation') + '</label><table class="validation-table"><thead><tr>' +
      '<th>' + t('stepCol') + '</th><th>' + t('status') + '</th>';
    if (showAction) html += '<th>' + t('validateAction') + '</th>';
    html += '<th>' + t('selector') + '</th>';
    if (showMatches) html += '<th>' + t('validateMatches') + '</th>';
    html += '<th>' + t('message') + '</th></tr></thead><tbody>';
    rows.forEach((step) => {
      const v = step.validation;
      const status = v.status || '';
      html += '<tr class="validation-row status-' + esc(status) + '">' +
        '<td>#' + (step.index + 1) + '</td>' +
        '<td><span class="validation-badge status-' + esc(status) + '">' + esc(validationStatusLabel(status)) + '</span></td>';
      if (showAction) html += '<td>' + esc(v.action_kind || '—') + '</td>';
      html += '<td class="mono">' + esc(step.selector || '—') + '</td>';
      if (showMatches) {
        html += '<td>' + ((v.match_count || 0) > 0 ? esc(String(v.match_count)) : '—') + '</td>';
      }
      html += '<td>' + esc(v.message || '') + '</td></tr>';
    });
    return html + '</tbody></table></div>';
  }

  function renderStepValidation(st) {
    const v = st.validation;
    if (!v) return '';
    let html = '<div class="section"><label>' + t('validation') + '</label><table class="validation-table validation-detail">';
    html += '<tr><th>' + t('status') + '</th><td><span class="validation-badge status-' + esc(v.status || '') + '">' +
      esc(validationStatusLabel(v.status)) + '</span></td></tr>';
    if (v.action_kind) {
      html += '<tr><th>' + t('validateAction') + '</th><td>' + esc(v.action_kind) + '</td></tr>';
    }
    if (st.selector) {
      html += '<tr><th>' + t('selector') + '</th><td class="mono">' + esc(st.selector) + '</td></tr>';
    }
    if ((v.match_count || 0) > 0) {
      html += '<tr><th>' + t('validateMatches') + '</th><td>' + esc(String(v.match_count)) + '</td></tr>';
    }
    if (v.message) {
      html += '<tr><th>' + t('message') + '</th><td>' + esc(v.message) + '</td></tr>';
    }
    return html + '</table></div>';
  }

  function renderInspector() {
    const sc = (DATA.scenarios || []).find((x) => x.id === state.scenarioId);
    const box = $('#inspector');
    if (!sc) {
      box.innerHTML = '<div class="empty">' + t('pickStep') + '</div>';
      return;
    }
    const st = (sc.steps || []).find((x) => x.index === state.stepIndex);
    if (!st) {
      box.innerHTML = '<div class="empty">' + t('noStep') + '</div>';
      return;
    }
    let html = '<h3>' + tf('stepTitle', { n: st.index + 1 }) + '</h3>';
    html += renderValidationBanner(sc);
    html += renderScenarioValidationTable(sc);
    html += '<div class="section"><label>' + t('gherkin') + '</label><div class="mono">' + esc(st.gherkin || st.text) + '</div></div>';
    if (st.iteration_path) {
      html += '<div class="section"><label>Iteration path</label><div class="mono">' + esc(st.iteration_path) + '</div></div>';
    }
    if (st.terminal_action) {
      html += '<div class="section"><label>Terminal action</label><div class="mono">' + esc(st.terminal_action) + '</div></div>';
    }
    if (st.selector) {
      html += '<div class="section"><label>' + t('selector') + '</label><div class="mono" id="sel-text">' + esc(st.selector) + '</div>' +
        '<button type="button" id="copy-sel" style="margin-top:6px">' + t('copySelector') + '</button></div>';
    }
    html += renderStepValidation(st);
    if (st.error) {
      html += '<div class="section"><label>' + t('error') + '</label><div class="mono">' + esc(st.error) + '</div></div>';
    }
    if (st.network) {
      html += '<div class="section"><label>' + t('network') + '</label><div class="mono">' + esc(st.network) + '</div></div>';
    } else if (
      state.importedTrace && state.importedTrace.scenarioId === sc.id &&
      (state.importedTrace.network || []).length
    ) {
      html += '<div class="section"><label>' + t('network') + ' (trace)</label>' +
        state.importedTrace.network.map((n) =>
          '<div class="mono">+' + n.offset_ms + 'ms ' + esc(n.snippet) + '</div>',
        ).join('') + '</div>';
    }
    if (st.page_context) {
      html += '<div class="section"><label>' + t('pageContext') + '</label><div class="mono">' + esc(st.page_context) + '</div></div>';
    }
    if (st.dom_snapshot) {
      html += '<details class="section snapshot"><summary>' + t('domSnapshot') + '</summary><pre class="mono snapshot-pre">' + esc(st.dom_snapshot) + '</pre></details>';
    }
    if (st.a11y_snapshot) {
      html += '<details class="section snapshot"><summary>' + t('a11ySnapshot') + '</summary><pre class="mono snapshot-pre">' + esc(st.a11y_snapshot) + '</pre></details>';
    }
    if (st.trace_offset_ms != null && st.trace_offset_ms >= 0 && scenarioHasTraceView(sc) && sc.trace_path) {
      const seek = formatTraceSeek(st.trace_offset_ms);
      html += '<div class="section"><label>' + t('traceOffset') + '</label><div class="mono">~' + seek +
        ' <span class="sub">(' + tf('stepInTrace', { n: st.index + 1 }) + ')</span></div></div>';
    }
    if (st.flaky_failures >= 2) {
      html += '<div class="section"><label>History metadata</label><div class="history-box">' + tf('failedTimesHistory', { n: st.flaky_failures }) + '<br><span class="sub">' + t('historicallyUnstable') + '</span></div></div>';
    }
    if ((st.tips || []).length) {
      html += '<div class="section"><label>' + t('tips') + '</label>' + st.tips.map((tip) => '<div class="tip">' + esc(tip) + '</div>').join('') + '</div>';
    }
    const shot = st.screenshot || (st.status === 'failed' ? sc.screenshot : '');
    if (shot) {
      const caption = tf('stepTitle', { n: st.index + 1 });
      html += '<div class="section screenshot-section"><label>' + t('screenshot') + '</label>' +
        '<button type="button" class="screenshot-open" aria-label="' + esc(t('screenshotExpand')) + '">' +
        '<img class="screenshot screenshot-thumb" src="' + esc(shot) + '" alt="' + esc(caption) + '" loading="lazy">' +
        '</button>' +
        '<p class="screenshot-hint">' + esc(t('screenshotHint')) + '</p></div>';
    }
    if (sc.history) {
      html += '<div class="section"><label>' + t('lastRun') + '</label><div class="history-box">' +
        esc(statusLabel(sc.history.last_status)) + ' @ ' + esc(sc.history.last_at) +
        (sc.history.changed ? ' <strong>' + t('changed') + '</strong>' : '') +
        (sc.history.last_message ? '<br>' + esc(sc.history.last_message) : '') +
      '</div></div>';
    }
    html += renderRegressions(sc);
    html += renderRunDiff(sc);
    if (scenarioHasTraceView(sc)) {
      html += renderTraceTimeline(sc, st);
    }
    html += renderHistoryRuns(sc);
    html += '<div class="section"><label>' + t('actions') + '</label>' +
      '<button type="button" id="copy-gherkin">' + t('copyGherkin') + '</button> ' +
      '<button type="button" id="copy-rerun">' + t('copyRerun') + '</button>';
    if (DATA.bridge_url) {
      html += ' <button type="button" id="open-ide" class="primary">' + t('openIde') + '</button>';
      html += ' <button type="button" id="rerun-ide">' + t('rerunIde') + '</button>';
    }
    if (sc.trace_command && sc.status === 'failed' && sc.trace_path) {
      html += ' <button type="button" id="copy-trace">Copy trace cmd</button>';
      html += ' <a class="btn" href="https://trace.playwright.dev/" target="_blank" rel="noopener">Open Playwright Trace</a>';
      if (DATA.bridge_url && sc.trace_path) {
        html += ' <button type="button" id="trace-ide">' + t('traceIde') + '</button>';
        if (st.trace_offset_ms != null) {
          html += ' <button type="button" id="trace-step">' + t('traceAtStep') + '</button>';
          html += ' <button type="button" id="copy-trace-seek">' + t('copyTraceSeek') + '</button>';
        }
      }
    }
    html += '</div>';
    if (sc.status === 'failed' && (sc.trace_path || scenarioHasTraceView(sc))) {
      html += '<div class="section"><label>' + t('traceFile') + '</label>';
      if (sc.trace_path) {
        html += '<div class="mono">' + esc(sc.trace_path) + '</div>';
      }
      html += '<div class="trace-drop" id="trace-drop">' + t('traceDrop') + '</div></div>';
    }
    box.innerHTML = html;

    const copy = (text) => navigator.clipboard && navigator.clipboard.writeText(text);
    const gherkinBtn = $('#copy-gherkin');
    if (gherkinBtn) gherkinBtn.onclick = () => copy(st.gherkin || st.text);
    const selBtn = $('#copy-sel');
    if (selBtn) selBtn.onclick = () => copy(st.selector);
    const rerunBtn = $('#copy-rerun');
    if (rerunBtn) rerunBtn.onclick = () => copy(sc.rerun_command || '');
    const traceBtn = $('#copy-trace');
    if (traceBtn) traceBtn.onclick = () => copy(sc.trace_command || '');
    const openIde = $('#open-ide');
    if (openIde) {
      openIde.onclick = () => {
        postBridge('/goto', {
          feature_path: sc.feature_path,
          scenario: sc.scenario,
          leaf_index: st.index,
          line: st.line || 0,
        }).then(() => { openIde.textContent = 'Sent to IDE ✓'; }).catch(() => alert('Scenaria IDE не запущена или мост недоступен'));
      };
    }
    const rerunIde = $('#rerun-ide');
    if (rerunIde) {
      rerunIde.onclick = () => {
        postBridge('/rerun', { feature_path: sc.feature_path, scenario: sc.scenario })
          .then(() => { rerunIde.textContent = 'OK ✓'; })
          .catch(() => alert('Scenaria IDE'));
      };
    }
    const traceIde = $('#trace-ide');
    if (traceIde) {
      traceIde.onclick = () => {
        postBridge('/trace', { trace_path: sc.trace_path, report_dir: DATA.report_dir || '' })
          .then(() => { traceIde.textContent = 'OK ✓'; })
          .catch(() => alert('Scenaria IDE'));
      };
    }
    const traceStep = $('#trace-step');
    if (traceStep) {
      traceStep.onclick = () => {
        postBridge('/trace', {
          trace_path: sc.trace_path,
          report_dir: DATA.report_dir || '',
          trace_offset_ms: st.trace_offset_ms || 0,
          step_index: st.index,
        }).then(() => { traceStep.textContent = 'OK ✓'; }).catch(() => alert('Scenaria IDE'));
      };
    }
    const copySeek = $('#copy-trace-seek');
    if (copySeek) {
      copySeek.onclick = () => {
        const hint = 'Seek trace to ~' + formatTraceSeek(st.trace_offset_ms) + ' (step #' + (st.index + 1) + ')';
        copy(hint);
      };
    }
    const shotOpen = box.querySelector('.screenshot-open');
    if (shotOpen) {
      activateOnEnterSpace(shotOpen, () => {
        const img = shotOpen.querySelector('img');
        if (img && img.src) openScreenshotLightbox(img.getAttribute('src') || img.src, tf('stepTitle', { n: st.index + 1 }));
      });
    }
    bindTraceDrop();
    box.querySelectorAll('.trace-event').forEach((el) => {
      el.addEventListener('click', () => {
        const ms = Number(el.dataset.ms);
        state.stepIndex = stepForTraceOffset(sc, ms);
        renderAllViews();
      });
    });
  }

  function bindTraceDrop() {
    const zone = $('#trace-drop');
    if (!zone) return;
    zone.addEventListener('dragover', (e) => { e.preventDefault(); zone.classList.add('hover'); });
    zone.addEventListener('dragleave', () => zone.classList.remove('hover'));
    zone.addEventListener('drop', (e) => {
      e.preventDefault();
      zone.classList.remove('hover');
      const file = e.dataTransfer && e.dataTransfer.files && e.dataTransfer.files[0];
      if (!file || !file.name.endsWith('.zip')) {
        alert(t('traceZipRequired'));
        return;
      }
      const sc = (DATA.scenarios || []).find((x) => x.id === state.scenarioId);
      if (!sc || sc.status !== 'failed') return;
      file.arrayBuffer().then((buf) => {
        if (!window.ScenariaTraceZip || !window.ScenariaTraceZip.parsePlaywrightTraceZip) {
          throw new Error('zip parser missing');
        }
        return window.ScenariaTraceZip.parsePlaywrightTraceZip(buf);
      }).then((parsed) => {
        state.importedTrace = {
          scenarioId: sc.id,
          events: parsed.events || [],
          network: parsed.network || [],
        };
        if (parsed.events && parsed.events.length) state.centerTab = 'trace';
        renderAllViews();
        alert(tf('traceImported', { n: (parsed.events || []).length }));
      }).catch(() => {
        if (window.confirm(t('traceZipParseFailed') + '\n\n' + t('openTraceExternal') + '?')) {
          window.open('https://trace.playwright.dev/', '_blank');
        }
      });
    });
  }

  function bindFilters() {
    $('#filter-failed').addEventListener('change', (e) => {
      state.filterFailed = e.target.checked;
      renderScenarioList();
    });
    $('#filter-tag').addEventListener('input', (e) => {
      state.filterTag = e.target.value.trim();
      renderScenarioList();
    });
    $('#filter-duration').addEventListener('input', (e) => {
      state.filterDurationMin = Number(e.target.value) || 0;
      renderScenarioList();
    });
    $('#filter-search').addEventListener('input', (e) => {
      state.search = e.target.value.trim();
      renderScenarioList();
    });
    $('#toggle-sidebar').addEventListener('click', () => {
      $('#sidebar').classList.toggle('mobile-open');
    });
    const closeInspector = $('#close-inspector');
    if (closeInspector) {
      closeInspector.addEventListener('click', () => {
        $('#inspector-panel').classList.remove('open');
      });
    }
    document.querySelectorAll('.mode-option').forEach((btn) => {
      btn.addEventListener('click', () => {
        const mode = btn.dataset.mode;
        const links = DATA.mode_links || {};
        const current = links.current || (DATA.light_mode ? 'light' : 'full');
        if (mode === current || !modeAvailable(mode)) return;
        const href = modeHref(mode);
        if (!href) return;
        window.location.href = href + reportStateHash();
      });
    });
    $('#export-failed').addEventListener('click', exportFailedFeature);
    document.querySelectorAll('.tab').forEach((tab) => {
      tab.addEventListener('click', () => {
        state.centerTab = tab.dataset.tab;
        updateCenterTabs();
        if (state.centerTab === 'actionlog') renderActionLog();
        if (state.centerTab === 'trace') renderTraceView();
      });
    });
    bindKeyboardNav();
  }

  function bindKeyboardNav() {
    document.addEventListener('keydown', (e) => {
      const lightbox = screenshotLightboxEl();
      if (lightbox && !lightbox.hidden) {
        if (e.key === 'Escape') {
          e.preventDefault();
          closeScreenshotLightbox();
        }
        return;
      }
      if (e.target && (e.target.matches('input, textarea, select') || e.target.isContentEditable)) return;
      if (e.target && e.target.matches('.scenario-item') && (e.key === 'ArrowDown' || e.key === 'ArrowUp')) {
        e.preventDefault();
        const items = Array.from(document.querySelectorAll('.scenario-item'));
        const idx = items.indexOf(e.target);
        const next = e.key === 'ArrowDown' ? items[Math.min(idx + 1, items.length - 1)] : items[Math.max(idx - 1, 0)];
        if (next) next.focus();
        return;
      }
      if (e.target && e.target.matches('.tab') && (e.key === 'ArrowRight' || e.key === 'ArrowLeft')) {
        e.preventDefault();
        const tabs = Array.from(document.querySelectorAll('.tab:not([hidden])'));
        const idx = tabs.indexOf(e.target);
        const next = e.key === 'ArrowRight' ? tabs[(idx + 1) % tabs.length] : tabs[(idx - 1 + tabs.length) % tabs.length];
        if (next) {
          next.focus();
          next.click();
        }
        return;
      }
      if (e.target && e.target.matches('button, .mode-option, .scenario-item')) return;
      const sc = (DATA.scenarios || []).find((x) => x.id === state.scenarioId);
      if (!sc || !(sc.steps || []).length) return;
      const idx = (sc.steps || []).findIndex((x) => x.index === state.stepIndex);
      if (e.key === 'ArrowDown' || e.key === 'j') {
        e.preventDefault();
        const next = sc.steps[Math.min(idx + 1, sc.steps.length - 1)];
        if (next) {
          state.stepIndex = next.index;
          renderAllViews();
        }
      } else if (e.key === 'ArrowUp' || e.key === 'k') {
        e.preventDefault();
        const prev = sc.steps[Math.max(idx - 1, 0)];
        if (prev) {
          state.stepIndex = prev.index;
          renderAllViews();
        }
      }
    });
  }

  function exportFailedFeature() {
    const failed = (DATA.scenarios || []).filter((sc) => sc.status === 'failed');
    if (!failed.length) {
      alert(t('exportNoFailed'));
      return;
    }
  const lines = ['# language: ru', 'Функционал: Failed steps export', ''];
    failed.forEach((sc) => {
      lines.push('  Сценарий: ' + sc.scenario);
      const step = (sc.steps || []).find((st) => st.status === 'failed');
      if (step) lines.push('    ' + (step.gherkin || step.text));
      lines.push('');
    });
    const blob = new Blob([lines.join('\n')], { type: 'text/plain' });
    const a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    a.download = 'failed-steps.feature';
    a.click();
  }

  function autoSelectFirstFailure() {
    const failed = (DATA.scenarios || []).find((sc) => sc.status === 'failed');
    if (failed) {
      selectScenario(failed.id);
      return;
    }
    const first = (DATA.scenarios || [])[0];
    if (first) selectScenario(first.id);
  }

  renderHeader();
  window.addEventListener('resize', updateHeaderHeightVar);
  bindScreenshotLightbox();
  bindFilters();
  renderScenarioList();
  if (restoreReportStateFromHash()) {
    renderScenarioList();
    renderAllViews();
  } else {
    autoSelectFirstFailure();
  }
})();
