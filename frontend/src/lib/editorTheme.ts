import type { EditorSettings } from './editorOptions'

export type ResolvedEditorTheme = 'scenaria-dark' | 'scenaria-light'

export function resolveEditorTheme(theme: EditorSettings['theme']): ResolvedEditorTheme {
  if (theme === 'system') {
    if (typeof window !== 'undefined' && window.matchMedia('(prefers-color-scheme: light)').matches) {
      return 'scenaria-light'
    }
    return 'scenaria-dark'
  }
  return theme === 'scenaria-light' ? 'scenaria-light' : 'scenaria-dark'
}

export function subscribeSystemTheme(theme: EditorSettings['theme'], onChange: () => void): () => void {
  if (theme !== 'system' || typeof window === 'undefined') return () => {}
  const mq = window.matchMedia('(prefers-color-scheme: dark)')
  const handler = () => onChange()
  mq.addEventListener('change', handler)
  return () => mq.removeEventListener('change', handler)
}
