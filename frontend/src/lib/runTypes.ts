export type RunForm = {
  tag: string
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
  workers: number
  slowMo: number
  browser: string
}

export function defaultRunForm(partial?: Partial<RunForm>): RunForm {
  return {
    tag: '',
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
    workers: 1,
    slowMo: 0,
    browser: 'chromium',
    ...partial,
  }
}
