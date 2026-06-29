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
  workers: number
  slowMo: number
  browser: string
  baseUrl: string
}

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
    workers: 1,
    slowMo: 0,
    browser: 'chromium',
    baseUrl: '',
    ...partial,
  }
}
