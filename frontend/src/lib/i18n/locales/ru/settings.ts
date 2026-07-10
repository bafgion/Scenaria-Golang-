export const settings = {
  title: 'Настройки',
  searchPlaceholder: 'Поиск настроек',
  sectionsNav: 'Разделы настроек',
  tabs: {
    record: 'Запись и браузер',
    selectors: 'Селекторы',
    plugins: 'Плагины',
    editor: 'Редактор',
    ui: 'Интерфейс',
  },
  sections: {
    browser: {
      title: 'Браузер',
      desc: 'Поведение окна и сессии при записи и прогоне Playwright.',
    },
    recording: {
      title: 'Запись шагов',
      desc: 'Фильтры при записи действий в браузере.',
    },
    run: {
      title: 'Запуск',
    },
    reports: {
      title: 'Отчёты',
      desc: 'Сохраняется в .scenaria/project.json текущего проекта.',
    },
    selectors: {
      title: 'Приоритет стратегий',
      desc: 'При записи и подборе селектора {brand} перебирает стратегии сверху вниз. Более стабильные — выше.',
    },
    plugins: {
      title: "Runner'ы и add-on'ы",
      desc: 'Плагины устанавливаются в addons/<name>/ и регистрируются в .scenaria/plugins.json.',
    },
    editorFont: {
      title: 'Шрифт и отображение',
      desc: 'Параметры Monaco-редактора сценариев.',
    },
    editorInput: {
      title: 'Ввод и поведение',
    },
    editorNavigation: {
      title: 'Навигация',
      desc: 'Структура feature-файла в breadcrumbs и панели шагов.',
    },
    editorHints: {
      title: 'Подсказки сценария',
      desc: 'Эвристики качества шагов в редакторе (маркеры и quick fix).',
    },
    uiLanguage: {
      title: 'Язык интерфейса',
      desc: 'Язык меню, диалогов и подписей.',
    },
    toolbar: {
      title: 'Панель инструментов',
      desc: 'Внешний вид верхней панели действий.',
    },
    stepsPanel: {
      title: 'Панель шагов',
      desc: 'Список распознанных шагов под редактором сценария.',
    },
    updates: {
      title: 'Обновления',
    },
  },
  cards: {
    headless: {
      title: 'Без окна браузера',
      description: 'Без окна браузера — окно скрыто при записи и запуске.',
    },
    browserEngine: {
      title: 'Движок браузера',
      description: 'Playwright: chromium, firefox или webkit.',
    },
    importantOnly: {
      title: 'Только важные',
      description: 'Пропускать второстепенные события при записи.',
    },
    linksOnly: {
      title: 'Только ссылки',
      description: 'Записывать переходы по ссылкам, без кликов по элементам.',
    },
    hoverRecord: {
      title: 'Записывать наведение',
      description: 'Добавлять шаги при наведении курсора на элементы.',
    },
    recordUrlWait: {
      title: 'Ожидать адрес после клика',
      description: 'После клика с переходом записывать `ожидаю адрес`, а не `открыт` перед кликом.',
    },
    hoverMin: {
      title: 'Минимальное наведение',
      description: 'Сколько миллисекунд курсор должен оставаться на элементе перед записью hover.',
    },
    scrollBeforeClick: {
      title: 'Прокрутка перед кликом',
      description: 'Перед записью клика прокручивать элемент в видимую область (как при воспроизведении).',
    },
    pickerDuringRecording: {
      title: 'Picker во время записи',
      description: 'Разрешить «Указать элемент» без паузы записи (по умолчанию нужна пауза).',
    },
    workers: {
      title: 'Параллельные воркеры',
      description: 'Число одновременных браузерных сессий при пакетном запуске.',
    },
    slowMo: {
      title: 'Скорость выполнения (slow-mo)',
      description: 'Пауза между действиями Playwright в миллисекундах. 0 — максимально быстро; 100–300 — удобно наблюдать шаги в браузере.',
    },
    navWait: {
      title: 'Ожидание навигации (nav-wait-until)',
      description: 'Когда считать переход по URL завершённым при шагах «Перейти» и записи.',
    },
    loops: {
      title: 'Лимит итераций циклов',
      description: 'Максимум повторов для блоков «Повторяю» / «Пока».',
    },
    htmlReportOpen: {
      title: 'HTML-отчёт при открытии',
      description: 'Какой вариант открывать из панели результатов: полный или лёгкий. Если файл ещё не сгенерирован, откроется доступный.',
      full: 'Полный (скриншоты, trace, DOM)',
      light: 'Лёгкий (быстрее, меньше размер)',
    },
    fontSize: {
      title: 'Размер шрифта',
      description: 'От 8 до 32 px.',
    },
    fontFamily: {
      title: 'Шрифт',
      description: 'Моноширинный шрифт редактора.',
    },
    theme: {
      title: 'Тема',
      description: 'Тёмная, светлая или как в системе.',
    },
    wordWrap: {
      title: 'Перенос строк',
      description: 'Переносить длинные шаги по ширине редактора.',
    },
    minimap: {
      title: 'Миникарта',
      description: 'Обзорная карта кода справа.',
    },
    lineNumbers: {
      title: 'Номера строк',
      description: 'Отображение номеров строк в gutter.',
    },
    renderWhitespace: {
      title: 'Пробелы',
      description: 'Когда показывать невидимые символы.',
    },
    tabSize: {
      title: 'Размер табуляции',
      description: 'Ширина отступа Tab в пробелах.',
    },
    insertSpaces: {
      title: 'Пробелы вместо Tab',
      description: 'Вставлять пробелы при нажатии Tab.',
    },
    folding: {
      title: 'Складывание блоков',
      description: 'Сворачивать блоки «Если» / «Повторяю».',
    },
    stickyScroll: {
      title: 'Sticky scroll',
      description: 'Закреплять заголовки сценариев при прокрутке.',
    },
    autoClosingQuotes: {
      title: 'Авто-закрытие кавычек',
      description: 'Поведение при вводе кавычек.',
    },
    formatOnSave: {
      title: 'Форматировать при сохранении',
      description: 'Нормализовать отступы и убрать лишние пустые строки между шагами при Ctrl+S.',
    },
    stepHover: {
      title: 'Подсказки при наведении',
      description: 'Показывать справку по шагу при hover в редакторе.',
    },
    validateOnType: {
      title: 'Проверка при вводе',
      description: 'Валидировать сценарий с задержкой при редактировании.',
    },
    breadcrumbs: {
      title: 'Breadcrumbs',
      description: 'Цепочка заголовков над редактором (Функционал → Сценарий → шаг).',
    },
    symbolOutline: {
      title: 'Структура в панели шагов',
      description: 'Вкладка «Структура» с деревом сценария и переходом по клику.',
    },
    stepsPanelView: {
      title: 'Вкладка панели по умолчанию',
      description: 'Что показывать под редактором при открытии сценария.',
    },
    codeLens: {
      title: 'Code Lens для запуска',
      description: 'Кнопки «▶ Запустить» над сценариями и шагами в редакторе.',
    },
    inlayHints: {
      title: 'Inlay hints',
      description: 'Серые подсказки справа от шага: click → selector, fill → значение.',
    },
    scenarioHints: {
      title: 'Показывать подсказки',
      description: 'Маркеры warning/info в редакторе и панели «Проверка».',
    },
    scenarioHintsAfterRecord: {
      title: 'После записи',
      description: 'Анализировать сценарий сразу после остановки записи.',
    },
    scenarioHintsShowWarning: {
      title: 'Предупреждения',
      description: 'Подсказки уровня warning (хрупкие селекторы, дубли).',
    },
    scenarioHintsShowInfo: {
      title: 'Информация',
      description: 'Подсказки уровня info (улучшения без критичных рисков).',
    },
    scenarioHintsAutoFixOnSave: {
      title: 'Авто-исправление при сохранении',
      description: 'Применять autoFixable подсказки при Ctrl+S.',
    },
    toolbarCompact: {
      title: 'Компактная панель',
      description: 'Меньше подписей на кнопках — только иконки.',
    },
    stepsPanelVisible: {
      title: 'Показывать панель шагов',
      description: 'Отображать разбор шагов Gherkin под редактором.',
    },
    stepsPanelHeight: {
      title: 'Высота панели',
      description: 'Высота области со списком шагов в пикселях.',
    },
    checkUpdates: {
      title: 'Проверять при запуске',
      description: 'Искать новую версию {brand} при старте IDE.',
    },
  },
  browser: {
    installed: '{label}: установлен',
    notInstalled: '{label}: не установлен — нужна загрузка перед записью и прогоном.',
    checking: 'Проверка движка…',
    installing: 'Установка {engine}…',
    installingBtn: 'Установка…',
    reinstall: 'Переустановить',
    installEngine: 'Установить движок',
    error: 'Ошибка',
  },
  strategies: {
    clicksTitle: 'Клики и кнопки',
    inputsTitle: 'Поля ввода',
    moveUp: 'Выше',
    moveDown: 'Ниже',
    resetClicks: 'Сбросить клики',
    resetInputs: 'Сбросить поля',
    contextual: 'Контекстный has-text',
    text: 'has-text по тексту',
    libraryTitle: 'UI-библиотеки',
    libraryDesc: 'Дополнительные подсказки для Material UI и Ant Design. На обычных сайтах без разметки библиотек поведение не меняется.',
    libraryMui: 'Material UI (MUI)',
    libraryAnt: 'Ant Design',
  },
  plugins: {
    playwright: 'Playwright',
    builtIn: 'встроен',
    vanessa: 'Vanessa Automation',
    installed: 'установлен',
    unavailable: 'недоступен',
    zipHint: 'ZIP-плагины: scenaria plugins install … или диалог «Управление плагинами».',
  },
  navWait: {
    load: 'load — полная загрузка страницы',
    domcontentloaded: 'domcontentloaded — DOM готов (по умолчанию)',
    networkidle: 'networkidle — сеть простаивает',
    commit: 'commit — первый ответ навигации',
  },
  editor: {
    themeDark: 'Тёмная',
    themeLight: 'Светлая',
    themeSystem: 'Как в системе',
    wordWrapOn: 'Включён',
    wordWrapOff: 'Выключен',
    lineNumbersOn: 'Обычные',
    lineNumbersRelative: 'Относительные',
    lineNumbersOff: 'Скрыть',
    whitespaceNone: 'Не показывать',
    whitespaceBoundary: 'На границах слов',
    whitespaceSelection: 'В выделении',
    whitespaceTrailing: 'В конце строк',
    whitespaceAll: 'Всегда',
    autoQuotesLanguage: 'По языку',
    autoQuotesAlways: 'Всегда',
    autoQuotesBeforeWhitespace: 'Перед пробелом',
    autoQuotesNever: 'Никогда',
    stepsPanelOutline: 'Структура',
    stepsPanelSteps: 'Таблица шагов',
  },
  units: {
    ms: 'мс',
    pcs: 'шт.',
    px: 'px',
  },
  warnings: {
    manyWorkers: 'Много воркеров и высокий slow-mo — прогон может быть очень долгим.',
    highWorkers: 'Большое число воркеров нагружает систему.',
    highSlowMo: 'Очень высокий slow-mo — шаги выполняются с большой паузой.',
  },
  uiLocale: 'Язык интерфейса',
  uiLocaleRu: 'Русский',
  uiLocaleEn: 'English',
  slowMo: {
    fast: 'Быстро',
    normal: 'Норма',
    slow: 'Медленно',
    tutorial: 'Учебный',
  },
  resetDefaults: 'Сбросить по умолчанию',
  managePlugins: 'Управление плагинами…',
  vanessaSettings: 'Настройки Vanessa…',
  apply: 'Применить',
  applyBusy: '…',
  applyDone: '✓',
  ok: 'OK',
  cancel: 'Отмена',
  hotkeysEditorNote: 'В редакторе: Ctrl+F — найти, Ctrl+H — найти и заменить.',
} as const

export const splash = {
  envSetup: 'Настройка окружения…',
  connecting: 'Подключение к приложению…',
  loadingSettings: 'Загрузка настроек…',
  initializing: 'Инициализация…',
  ready: 'Готово',
  starting: 'Запуск…',
} as const
