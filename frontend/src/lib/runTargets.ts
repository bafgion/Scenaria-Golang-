import { isUntitled, untitledLabel } from './untitled'
import type { TabBody } from './tabMemory'
import { tabEditorText, tabNeedsDiskReload } from './tabMemory'

function basename(path: string): string {
  const normalized = path.replace(/\\/g, '/')
  return normalized.split('/').pop() || normalized
}

/** Match a run/history path to an open editor tab (untitled label or internal id). */
export function findTabForRunTarget(
  path: string,
  tabs: TabBody[],
  activeTab: string,
): TabBody | undefined {
  const direct = tabs.find((t) => t.path === path)
  if (direct) return direct

  const base = basename(path)
  if (activeTab && isUntitled(activeTab) && untitledLabel(activeTab) === base) {
    return tabs.find((t) => t.path === activeTab)
  }
  return tabs.find((t) => isUntitled(t.path) && untitledLabel(t.path) === base)
}

/** Map display/history paths back to the tab id used in the editor. */
export function resolveLogicalRunTarget(
  path: string,
  tabs: TabBody[],
  activeTab: string,
): string {
  return findTabForRunTarget(path, tabs, activeTab)?.path ?? path
}

export function editorTextForRunTarget(
  path: string,
  tabs: TabBody[],
  activeTab: string,
  editorText: string,
  liveActiveText?: string,
): string {
  const resolved = resolveLogicalRunTarget(path, tabs, activeTab)
  if (resolved === activeTab) {
    return liveActiveText ?? editorText
  }
  const tab = tabs.find((t) => t.path === resolved)
  return tab ? tabEditorText(tab) : ''
}

export type MaterializeRunTargetDeps = {
  writeTempFeature: (content: string) => Promise<string>
  readFeature: (path: string) => Promise<string>
}

/** Prepare on-disk paths for run/validate, including unsaved and untitled tabs. */
export async function materializeRunTargetPaths(
  paths: string[],
  tabs: TabBody[],
  activeTab: string,
  editorText: string,
  liveActiveText: string | undefined,
  deps: MaterializeRunTargetDeps,
): Promise<string[]> {
  const diskTargets: string[] = []
  for (const rawPath of paths) {
    const path = resolveLogicalRunTarget(rawPath, tabs, activeTab)
    const tab = tabs.find((t) => t.path === path)

    if (!tab) {
      if (isUntitled(path)) continue
      diskTargets.push(path)
      continue
    }

    if (isUntitled(tab.path) || tab.dirty) {
      let content = editorTextForRunTarget(tab.path, tabs, activeTab, editorText, liveActiveText)
      if (!content && tabNeedsDiskReload(tab)) {
        content = await deps.readFeature(tab.path)
      }
      diskTargets.push(await deps.writeTempFeature(content))
    } else {
      diskTargets.push(tab.path)
    }
  }
  return diskTargets
}
