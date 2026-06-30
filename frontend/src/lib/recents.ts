import {
  LoadRecents,
  RememberRecentProject,
  RememberRecentFeature,
} from '../../wailsjs/go/wailsapp/App'
import { callWailsWithTimeout, wailsReady } from './wailsTimeout'

const KEY = 'scenaria.recents'
const MAX = 6

export type Recents = {
  projects: string[]
  features: string[]
}

const empty: Recents = { projects: [], features: [] }

function loadLocal(): Recents {
  try {
    const raw = localStorage.getItem(KEY)
    if (!raw) return { ...empty }
    const data = JSON.parse(raw) as Recents
    return {
      projects: (data.projects || []).slice(0, MAX),
      features: (data.features || []).slice(0, MAX),
    }
  } catch {
    return { ...empty }
  }
}

function saveLocal(recents: Recents) {
  try {
    localStorage.setItem(KEY, JSON.stringify(recents))
  } catch {
    /* ignore */
  }
}

export async function loadRecents(): Promise<Recents> {
  if (!wailsReady()) {
    return loadLocal()
  }
  const data = await callWailsWithTimeout('LoadRecents', LoadRecents(), 3000)
  if (!data) {
    return loadLocal()
  }
  const recents: Recents = {
    projects: (data.projects || []).slice(0, MAX),
    features: (data.features || []).slice(0, MAX),
  }
  saveLocal(recents)
  return recents
}

export async function rememberProject(path: string) {
  const norm = path.trim()
  if (!norm) return
  const local = loadLocal()
  local.projects = [norm, ...local.projects.filter((x) => x !== norm)].slice(0, MAX)
  saveLocal(local)
  try {
    await RememberRecentProject(norm)
  } catch {
    /* dev without wails */
  }
}

export async function rememberFeature(path: string) {
  const norm = path.trim()
  if (!norm) return
  const local = loadLocal()
  local.features = [norm, ...local.features.filter((x) => x !== norm)].slice(0, MAX)
  saveLocal(local)
  try {
    await RememberRecentFeature(norm)
  } catch {
    /* dev without wails */
  }
}
