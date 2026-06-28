import type * as Monaco from 'monaco-editor'

let registered = false

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
        [/@[\w-]+/, 'tag'],
        [
          /^(Функционал|Функциональность|Функция|Feature|Сценарий|Scenario|Контекст|Background|Структура сценария|Scenario Outline|Примеры|Examples)\s*:/,
          'keyword.gherkin',
        ],
        [/^(Допустим|Когда|Тогда|И|Но|Given|When|Then|And|But)\s+/, 'keyword.step'],
        [/^(Если|Повторяю|Пока|Для каждого|Иначе|Конец если|Конец)\b/, 'keyword.block'],
        [/TestClient/, 'type.testclient'],
        [/"(?:[^"\\]|\\.)*"/, 'string'],
      ],
    },
  })

  monaco.editor.defineTheme('scenaria-dark', {
    base: 'vs-dark',
    inherit: true,
    rules: [
      { token: 'comment', foreground: '6a737d', fontStyle: 'italic' },
      { token: 'tag', foreground: '79c0ff' },
      { token: 'keyword.gherkin', foreground: 'd2a8ff', fontStyle: 'bold' },
      { token: 'keyword.step', foreground: 'ff7b72', fontStyle: 'bold' },
      { token: 'keyword.block', foreground: 'ffa657', fontStyle: 'bold' },
      { token: 'string', foreground: 'a5d6ff' },
      { token: 'type.testclient', foreground: 'ffa657', fontStyle: 'bold' },
    ],
    colors: {
      'editor.background': '#12141a',
      'editor.lineHighlightBackground': '#1a1d26',
      'editorGutter.background': '#12141a',
    },
  })
}
