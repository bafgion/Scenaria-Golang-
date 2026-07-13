export type RunForm = {
  tag: string
  scenario: string
  testClient: string
  vars: string
  engine: string
  dryRun: boolean
  headed: boolean
  installPW: boolean
  allure: boolean
  trace: boolean
  video: boolean
  html: boolean
  junit: boolean
  summaryJson: boolean
  htmlTimestamp: boolean
  htmlLightMode: boolean
  reuseLiveBrowser: boolean
  continueOnFail: boolean
  workers: number
  slowMo: number
  browser: string
  baseUrl: string
  startStep: number
  endStep: number
}

export type RunFormMode = 'batch' | 'single' | 'tag' | 'step-range'

export function defaultRunForm(partial?: Partial<RunForm>): RunForm {
  return {
    tag: '',
    scenario: '',
    testClient: '',
    vars: '',
    engine: 'playwright',
    dryRun: false,
    headed: false,
    installPW: false,
    allure: false,
    trace: false,
    video: false,
    html: false,
    junit: false,
    summaryJson: false,
    htmlTimestamp: false,
    htmlLightMode: true,
    reuseLiveBrowser: false,
    continueOnFail: false,
    workers: 1,
    slowMo: 0,
    browser: 'chromium',
    baseUrl: '',
    startStep: -1,
    endStep: -1,
    ...partial,
  }
}

export type RunFormModeDefaults = Partial<RunForm>

export function runFormFromMode(lastRun: RunForm, mode: RunFormMode, defaults: RunFormModeDefaults = {}): RunForm {
  const form: RunForm = {
    ...lastRun,
    ...defaults,
  }

  switch (mode) {
    case 'batch':
      form.tag = ''
      form.scenario = ''
      form.startStep = -1
      form.endStep = -1
      break
    case 'single':
      form.tag = ''
      form.startStep = -1
      form.endStep = -1
      break
    case 'tag':
      form.scenario = ''
      form.startStep = -1
      form.endStep = -1
      break
    case 'step-range':
      form.tag = ''
      break
  }

  return form
}

export function currentScenarioRunFormFrom(
  lastRun: RunForm,
  defaults: RunFormModeDefaults = {},
  liveBrowserOpen = false,
): RunForm {
  const form = runFormFromMode(lastRun, 'step-range', defaults)
  if (!form.dryRun && liveBrowserOpen) {
    form.reuseLiveBrowser = true
    form.workers = 1
  }
  return form
}

export function batchRunFormFrom(lastRun: RunForm, dryRun: boolean): RunForm {
  return runFormFromMode(lastRun, 'batch', { dryRun })
}
