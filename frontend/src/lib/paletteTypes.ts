export type PaletteCommand = {
  id: string
  label: string
  group: string
  shortcut?: string
  run: () => void
}
