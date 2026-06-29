export type HotkeyId =
  | 'save'
  | 'save-as'
  | 'run'
  | 'run-current'
  | 'browser'
  | 'record'
  | 'new'
  | 'open'
  | 'find'
  | 'steps-help'
  | 'hotkeys'
  | 'settings'
  | 'palette'
  | 'snippets'
  | 'journal'
  | 'escape'

function hasMod(e: KeyboardEvent): boolean {
  return e.ctrlKey || e.metaKey
}

export function matchHotkey(e: KeyboardEvent): HotkeyId | null {
  if (hasMod(e) && e.shiftKey && e.code === 'KeyS') return 'save-as'
  if (hasMod(e) && !e.shiftKey && e.code === 'KeyS') return 'save'
  if (hasMod(e) && e.shiftKey && e.code === 'Enter') return 'run-current'
  if (hasMod(e) && !e.shiftKey && e.code === 'Enter') return 'run'
  if (hasMod(e) && e.code === 'KeyB') return 'browser'
  if (hasMod(e) && e.code === 'KeyR') return 'record'
  if (hasMod(e) && e.code === 'KeyN') return 'new'
  if (hasMod(e) && e.code === 'KeyO') return 'open'
  if (hasMod(e) && e.code === 'KeyH') return 'find'
  if (e.code === 'F1' && e.shiftKey) return 'hotkeys'
  if (e.code === 'F1' && !e.shiftKey) return 'steps-help'
  if (hasMod(e) && e.code === 'Comma') return 'settings'
  if (hasMod(e) && e.shiftKey && e.code === 'KeyP') return 'palette'
  if (hasMod(e) && e.shiftKey && e.code === 'Space') return 'snippets'
  if (hasMod(e) && e.code === 'Backquote') return 'journal'
  if (e.code === 'Escape') return 'escape'
  return null
}

export function shouldIgnoreAppHotkey(e: KeyboardEvent): boolean {
  const hotkey = matchHotkey(e)
  // F1 обрабатываем глобально (в т.ч. из Monaco) — иначе не доходит до справки.
  if (hotkey === 'steps-help' || hotkey === 'hotkeys') {
    return false
  }

  const target = e.target
  if (!(target instanceof Element)) return false

  if (target.closest('.palette-backdrop, .modal-backdrop')) {
    const tag = (target as HTMLElement).tagName
    if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') {
      const id = matchHotkey(e)
      if (id && id !== 'escape') return true
    }
  }

  if (target.closest('input, textarea, select') && !target.closest('.monaco-wrap')) {
    return true
  }

  // В редакторе сочетания обрабатывает Monaco addCommand — иначе срабатывает дважды.
  if (target.closest('.monaco-wrap')) {
    return true
  }

  return false
}
