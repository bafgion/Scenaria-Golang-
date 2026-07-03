import { brand, browser } from './browser'
import { catalog } from './catalog'
import { common } from './common'
import { confirm } from './confirm'
import { dialogs } from './dialogs'
import { editor } from './editor'
import { filePicker } from './filePicker'
import { journal } from './journal'
import { menus, panels } from './menus'
import { onboarding } from './onboarding'
import { palette } from './palette'
import { results } from './results'
import { settings, splash } from './settings'
import { statusBar } from './statusBar'
import { toolbar } from './toolbar'
import { welcome } from './welcome'

export const ru = {
  brand,
  browser,
  catalog,
  common,
  confirm,
  dialogs,
  editor,
  filePicker,
  journal,
  menus,
  onboarding,
  panels,
  palette,
  results,
  settings,
  splash,
  statusBar,
  toolbar,
  welcome,
} as const

export type Messages = typeof ru
