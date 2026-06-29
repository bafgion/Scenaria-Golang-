import {
  LoadRecents,
  RememberRecentProject,
  RememberRecentFeature,
} from '../../wailsjs/go/wailsapp/App'

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

function wailsReady(): boolean {
  return typeof window !== 'undefined' && !!(window as unknown as { go?: unknown }).go
}

export async function loadRecents(): Promise<Recents> {
  if (!wailsReady()) {
    return loadLocal()
  }
  try {
    const data = await Promise.race([
      LoadRecents(),
      new Promise<never>((_, reject) => setTimeout(() => reject(new Error('timeout')), 3000)),
    ])
    const recents: Recents = {
      projects: (data.projects || []).slice(0, MAX),
      features: (data.features || []).slice(0, MAX),
    }
    saveLocal(recents)
    return recents
  } catch {
    return loadLocal()
  }
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
