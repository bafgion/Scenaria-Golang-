import type * as Monaco from 'monaco-editor'

let registered = false

const STEP_KEYWORDS =
  'Допустим|Дано|Когда|Тогда|И|Но|Given|When|Then|And|But'
const GHERKIN_HEADERS =
  'Функционал|Функциональность|Функция|Feature|Сценарий|Scenario|Контекст|Background|Структура сценария|Scenario Outline|Примеры|Examples'
const BLOCK_KEYWORDS = 'Если|Повторяю|Пока|Для каждого|Иначе|Конец если|Конец'

export function registerFeatureLanguage(monaco: typeof Monaco) {
  if (registered) {
    return
  }
  registered = true

  monaco.languages.register({ id: 'scenaria-feature' })

  monaco.languages.setMonarchTokensProvider('scenaria-feature', {
    defaultToken: 'text',
    ignoreCase: true,
    tokenizer: {
      root: [
        [/^\s*#.*$/, 'comment'],
        [/^\s*@[\w-]+(?:\s+@[\w-]+)*\s*$/, 'tag'],
        [
          new RegExp(`^\\s*(${GHERKIN_HEADERS})\\s*:`),
          'keyword.gherkin',
        ],
        [new RegExp(`^\\s*(${STEP_KEYWORDS})\\s+`), 'keyword.step'],
        [new RegExp(`^\\s*(${BLOCK_KEYWORDS})(?:\\s|$)`), 'keyword.block'],
        [/TestClient/, 'type.testclient'],
        [/^\s*\|.*\|\s*$/, 'string.table'],
        [/"(?:[^"\\]|\\.)*"/, 'string'],
      ],
    },
  })

  // VS Code–like dark palette aligned with frontend/src/style.css
  monaco.editor.defineTheme('scenaria-dark', {
    base: 'vs-dark',
    inherit: true,
    rules: [
      { token: 'comment', foreground: '858585', fontStyle: 'italic' },
      { token: 'tag', foreground: '5ec8f2' },
      { token: 'keyword.gherkin', foreground: '569cd6', fontStyle: 'bold' },
      { token: 'keyword.step', foreground: 'c586c0', fontStyle: 'bold' },
      { token: 'keyword.block', foreground: 'dcdcaa', fontStyle: 'bold' },
      { token: 'string', foreground: 'ce9178' },
      { token: 'string.table', foreground: '9cdcfe' },
      { token: 'type.testclient', foreground: '4ec9b0', fontStyle: 'bold' },
    ],
    colors: {
      'editor.background': '#1e1e1e',
      'editor.foreground': '#cccccc',
      'editor.lineHighlightBackground': '#2a2d2e',
      'editorGutter.background': '#252526',
      'editorLineNumber.foreground': '#858585',
      'editorLineNumber.activeForeground': '#cccccc',
      'editor.selectionBackground': '#094771',
      'editor.inactiveSelectionBackground': '#3a3d41',
      'editorCursor.foreground': '#cccccc',
      'editorWidget.background': '#252526',
      'editorWidget.border': '#454545',
    },
  })

  monaco.editor.defineTheme('scenaria-light', {
    base: 'vs',
    inherit: true,
    rules: [
      { token: 'comment', foreground: '6a737d', fontStyle: 'italic' },
      { token: 'tag', foreground: '0550ae' },
      { token: 'keyword.gherkin', foreground: '0550ae', fontStyle: 'bold' },
      { token: 'keyword.step', foreground: '8250df', fontStyle: 'bold' },
      { token: 'keyword.block', foreground: '953800', fontStyle: 'bold' },
      { token: 'string', foreground: '0a3069' },
      { token: 'string.table', foreground: '0550ae' },
      { token: 'type.testclient', foreground: '116329', fontStyle: 'bold' },
    ],
    colors: {
      'editor.background': '#ffffff',
      'editor.foreground': '#24292f',
      'editor.lineHighlightBackground': '#f6f8fa',
      'editorGutter.background': '#f6f8fa',
      'editorLineNumber.foreground': '#8c959f',
      'editorLineNumber.activeForeground': '#24292f',
      'editor.selectionBackground': '#b6e3ff',
      'editor.inactiveSelectionBackground': '#d0d7de',
      'editorCursor.foreground': '#24292f',
      'editorWidget.background': '#ffffff',
      'editorWidget.border': '#d0d7de',
    },
  })
}
