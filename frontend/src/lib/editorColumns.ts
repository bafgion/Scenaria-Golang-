/**
 * Monaco column positions are UTF-16 code units; Go completions use rune indices.
 */

export function monacoColumnToRuneIndex(line: string, col0Based: number): number {
  if (col0Based <= 0) return 0
  let utf16 = 0
  let runeIdx = 0
  for (const ch of line) {
    if (utf16 >= col0Based) return runeIdx
    const cp = ch.codePointAt(0) ?? 0
    utf16 += cp > 0xffff ? 2 : 1
    runeIdx++
  }
  return runeIdx
}

export function runeIndexToMonacoColumn(line: string, runeIndex: number): number {
  if (runeIndex <= 0) return 0
  let utf16 = 0
  let runeIdx = 0
  for (const ch of line) {
    if (runeIdx >= runeIndex) return utf16
    const cp = ch.codePointAt(0) ?? 0
    utf16 += cp > 0xffff ? 2 : 1
    runeIdx++
  }
  return utf16
}

export function runeSlice(line: string, start: number, end: number): string {
  return [...line].slice(start, end).join('')
}
