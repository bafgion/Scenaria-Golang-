import type { editor as MonacoEditor } from 'monaco-editor'

export type TabViewState = {
  scrollTop: number
  position: { lineNumber: number; column: number }
}

/** Per-tab cursor and scroll for Monaco editor tab switches. */
export class MonacoTabViewStateStore {
  private states = new Map<string, TabViewState>()

  private key(path: string | null): string {
    return path ?? '__welcome__'
  }

  capture(editor: MonacoEditor.IStandaloneCodeEditor, path: string | null): void {
    const position = editor.getPosition()
    if (!position) return
    this.states.set(this.key(path), {
      scrollTop: editor.getScrollTop(),
      position: { ...position },
    })
  }

  restore(editor: MonacoEditor.IStandaloneCodeEditor, path: string | null): void {
    const state = this.states.get(this.key(path))
    if (!state) return
    editor.setPosition(state.position)
    editor.setScrollTop(state.scrollTop)
    editor.revealPositionInCenterIfOutsideViewport(state.position)
  }

  drop(path: string): void {
    this.states.delete(this.key(path))
  }
}
