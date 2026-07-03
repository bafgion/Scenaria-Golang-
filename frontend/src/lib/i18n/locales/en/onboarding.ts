export const onboarding = {
  skip: 'Skip tour',
  back: 'Back',
  next: 'Next',
  done: 'Done',
  stepOf: 'Step {current} of {total}',
  waitAction: 'Complete the action on screen to continue',
  steps: {
    welcome: {
      title: 'Welcome to Scenaria',
      body: 'This short tour covers examples, the editor, validation, and dry-run. The Start checklist stays visible as you progress.',
    },
    openExamples: {
      title: 'Open examples',
      body: 'Click «Open example scenarios». The built-in tutorial project loads without picking a folder.',
    },
    pickFeature: {
      title: 'Pick a scenario',
      body: 'In the catalog on the left, open any .feature file, e.g. «01-pervaya-proverka».',
    },
    editor: {
      title: 'Gherkin editor',
      body: 'Scenario steps are in Russian. Error highlighting, Ctrl+Space for completions, F1 for the step catalog.',
    },
    validate: {
      title: 'Syntax validation',
      body: 'The «Run & record» menu is open — choose «Validate…» (below run items). Keep «Syntax only» in the dialog and confirm.',
    },
    dryRun: {
      title: 'Dry-run',
      body: 'In the same menu choose «Dry-run» — a trial run without opening the browser.',
    },
    journal: {
      title: 'Journal',
      body: 'Open the «Journal» tab at the bottom (or Ctrl+`). Validation and run output appear here.',
    },
    finish: {
      title: 'All set!',
      body: 'Tour complete. Restart anytime via «Help → Training…». Happy testing!',
    },
  },
} as const
