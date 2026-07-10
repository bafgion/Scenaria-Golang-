export const confirm = {
  deleteFeature: {
    title: 'Удалить сценарий',
    message: 'Удалить «{name}» без возможности восстановления? Файл не попадёт в Корзину или Recycle Bin.',
    confirmLabel: 'Удалить',
  },
  openOtherProject: {
    title: 'Открыть другой проект',
    messageDirty: 'Есть {count} несохранённых файл(ов). Закрыть текущий проект и открыть новый?',
    messageClean: 'Закрыть текущий проект и открыть новый?',
    confirmLabel: 'Открыть',
  },
  closeProject: {
    title: 'Закрыть проект',
    message: 'Есть {count} несохранённых файл(ов). Закрыть проект без сохранения?',
    confirmLabel: 'Закрыть',
  },
  diskChanged: {
    title: 'Файл изменён на диске',
    message:
      '«{name}» был изменён вне редактора. Перезагрузить с диска? Несохранённые правки в редакторе будут потеряны.',
    confirmLabel: 'Перезагрузить',
  },
  recordingActive: {
    title: 'Запись активна',
    message:
      'Шаги записываются в «{from}». Переключить вкладку? Дальнейшие шаги будут записываться в «{to}».',
    confirmLabel: 'Переключить',
    dontAskAgainLabel: 'Больше не спрашивать',
  },
  recordingTargetTab: {
    title: 'Целевая вкладка записи',
    messagePaused: '«{name}» — цель записи (на паузе). Закрыть вкладку?',
    messageRecording: 'Запись идёт в «{name}». Поставьте на паузу или закройте вкладку осознанно.',
    confirmLabelClose: 'Закрыть',
    confirmLabelForceClose: 'Закрыть всё равно',
  },
  headless: {
    title: 'Режим headless',
    message:
      'Переключение между headed и headless режимом применится к текущей сессии и может перезапустить браузер. Cookies, данные сессии и состояние страницы могут быть очищены. Продолжить?',
    confirmLabel: 'Применить',
  },
  generic: {
    title: 'Подтверждение',
    confirmLabelDelete: 'Удалить',
  },
  closeBrowser: {
    message: 'Есть несохранённые изменения. Закрыть браузер?',
  },
} as const
