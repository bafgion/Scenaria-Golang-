export type StepCatalogEntry = {
  category: string
  template: string
  help: string
}

export type SnippetEntry = StepCatalogEntry

export type StepHelpEntry = {
  label: string
  action: string
  category: string
  description: string
  template: string
  example: string
  parameters: string[]
  help: string
}
