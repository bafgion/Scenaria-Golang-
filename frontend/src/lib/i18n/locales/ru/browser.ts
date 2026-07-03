export const browser = {
  recording: '● Идёт запись',
  paused: '⏸ Пауза — можно выбрать элемент',
  testing: '▶ Идёт тест',
  resume: 'Продолжить',
  pause: 'Пауза',
  stop: 'Стоп',
  record: 'Запись',
  showBrowserTitle: 'Показать окно браузера',
  pickElement: 'Указать элемент',
  pickElementPauseHint: 'Поставьте запись на паузу',
} as const

export const brand = {
  tagline: 'Автотесты сайтов · Gherkin · Playwright',
  description: 'Запись и запуск автотестов сайтов в браузере.',
  aboutVersion: 'Версия {version}',
  overlayTitle: '{name} — перетащите панель',
} as const
