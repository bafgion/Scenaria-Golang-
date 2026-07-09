import { buildSessionTabsSnapshot, sessionTabPathsFromSettings, untitledContentMap } from '../lib/sessionTabs'
import type { TabBody } from '../lib/tabMemory'
import { tabEditorText } from '../lib/tabMemory'
import { syncUntitledCounterFromPaths } from '../lib/untitled'
import type { DialogBindController } from './dialogBindController'
import type { PaletteTr } from './paletteCommandsController'
import type { createRecentsStore } from '../stores/recentsStore'
import type { createRecorderPrefsStore } from '../stores/recorderPrefsStore'
import type { createSettingsStore } from '../stores/settingsStore'
import type { createTabsStore } from '../stores/tabsStore'
import type { createTestClientStore } from '../stores/testClientStore'
import type { createUiPrefsStore } from '../stores/uiPrefsStore'
import { gui } from '../../wailsjs/go/models'

export type WorkspaceSessionStores = {
  settingsStore: ReturnType<typeof createSettingsStore>
  recorderPrefsStore: ReturnType<typeof createRecorderPrefsStore>
  uiPrefsStore: ReturnType<typeof createUiPrefsStore>
  recentsStore: ReturnType<typeof createRecentsStore>
  dialogBinds: DialogBindController
  tabsStore: ReturnType<typeof createTabsStore>
  testClientStore: ReturnType<typeof createTestClientStore>
}

export type ApplyEditorTextOptions = {
  switchTab?: boolean
  tabPath?: string | null
  skipValidate?: boolean
}

export type WorkspaceSessionContext = {
  stores: WorkspaceSessionStores
  welcomeKey: string
  getProjectPath: () => string
  getTabs: () => TabBody[]
  getActiveTab: () => string
  getEditorText: () => string
  syncActiveTabContent: () => void
  isUntitled: (path: string) => boolean
  saveSettings: (dto: gui.AppSettingsDTO) => Promise<void>
  saveFeatureDraft: (path: string, text: string) => Promise<void>
  openProject: (path: string) => Promise<gui.ProjectInfo>
  resolveProjectPath: (path: string) => Promise<string>
  applyProjectScan: (info: gui.ProjectInfo, fallbackPath?: string) => void
  listTestClients: () => Promise<string[]>
  loadFeature: (path: string) => Promise<void>
  applyEditorText: (text: string, opts: ApplyEditorTextOptions) => Promise<void>
  trimTabsMemory: () => void
  appendLog: (line: string) => void
  setStatus: (msg: string, tone?: 'normal' | 'error' | 'success' | 'busy') => void
  tr: PaletteTr
}

export function hasRestorableWorkspaceSession(s: Partial<gui.AppSettingsDTO> | null | undefined): boolean {
  if (!s) return false
  if ((s.sessionProject || '').trim()) return true
  if ((s.activeTab || '').trim()) return true
  if ((s.openTabs || []).some((p) => (p || '').trim())) return true
  if ((s.untitledTabs || []).some((t) => (t?.path || '').trim())) return true
  return false
}

export function createWorkspaceSessionController(ctx: WorkspaceSessionContext) {
  function buildSettingsDTO(): gui.AppSettingsDTO {
    ctx.syncActiveTabContent()
    const tabs = ctx.getTabs()
    const activeTab = ctx.getActiveTab()
    const editorText = ctx.getEditorText()
    const sessionTabs = buildSessionTabsSnapshot(tabs, activeTab, editorText, ctx.welcomeKey)
    const recents = ctx.stores.recentsStore.snapshot()
    return gui.AppSettingsDTO.createFrom({
      sessionProject: ctx.getProjectPath(),
      openTabs: sessionTabs.openTabs,
      untitledTabs: sessionTabs.untitledTabs,
      activeTab: sessionTabs.activeTab,
      ...ctx.stores.settingsStore.dtoFields(ctx.stores.settingsStore.snapshot()),
      ...ctx.stores.recorderPrefsStore.dtoFields(ctx.stores.recorderPrefsStore.snapshot()),
      ...ctx.stores.uiPrefsStore.dtoFields(ctx.stores.uiPrefsStore.snapshot()),
      recentProjects: recents.projects,
      recentFeatures: recents.features,
    })
  }

  async function persistSettings() {
    const { dialogBinds } = ctx.stores
    dialogBinds.flushUiPrefsBindLocals()
    dialogBinds.flushRecorderPrefsBindLocals()
    dialogBinds.flushRecordFormBindLocals()
    dialogBinds.flushSettingsBindLocals()
    await ctx.saveSettings(buildSettingsDTO())
  }

  async function autosaveDirtyDrafts() {
    if (!ctx.getProjectPath()) return
    ctx.syncActiveTabContent()
    const activeTab = ctx.getActiveTab()
    const editorText = ctx.getEditorText()
    for (const tab of ctx.getTabs()) {
      if (!tab.dirty || ctx.isUntitled(tab.path)) continue
      const text = tab.path === activeTab ? editorText : tabEditorText(tab)
      try {
        await ctx.saveFeatureDraft(tab.path, text)
      } catch {
        /* offline */
      }
    }
  }

  async function restoreWorkspaceSession(s: gui.AppSettingsDTO) {
    const proj = (s.sessionProject || '').trim()
    if (!proj) return
    const resolvedProj = (await ctx.resolveProjectPath(proj)).trim()
    if (!resolvedProj) return
    try {
      const info = await ctx.openProject(resolvedProj)
      ctx.applyProjectScan(info, resolvedProj)
      ctx.stores.testClientStore.setClients(await ctx.listTestClients().catch((): string[] => []))
    } catch {
      ctx.appendLog(ctx.tr('journal.session.projectNotFound', { path: resolvedProj }))
      ctx.setStatus(ctx.tr('journal.status.sessionProjectNotFound'), 'error')
      return
    }
    try {
      const untitledBodies = untitledContentMap(s.untitledTabs)
      syncUntitledCounterFromPaths([...(s.openTabs || []), ...untitledBodies.keys()])
      const tabPaths = sessionTabPathsFromSettings(s.openTabs, s.untitledTabs)
      const tabs = ctx.getTabs()
      for (const p of tabPaths) {
        if (ctx.isUntitled(p)) {
          const content = untitledBodies.get(p)
          if (content === undefined || tabs.some((t) => t.path === p)) continue
          ctx.stores.tabsStore.appendTab({ path: p, content, dirty: true })
          continue
        }
        try {
          await ctx.loadFeature(p)
        } catch {
          /* skip missing files */
        }
      }
      const active = (s.activeTab || '').trim()
      if (active) {
        ctx.stores.tabsStore.setWelcomeVisible(false)
        if (ctx.isUntitled(active)) {
          const tab = ctx.getTabs().find((t) => t.path === active)
          if (tab) {
            await ctx.applyEditorText(tabEditorText(tab), {
              switchTab: true,
              tabPath: active,
              skipValidate: true,
            })
            ctx.stores.tabsStore.setActiveTab(active)
            ctx.trimTabsMemory()
          }
        } else {
          await ctx.loadFeature(active)
        }
      } else if (tabPaths.length > 0) {
        ctx.stores.tabsStore.setWelcomeVisible(false)
      }
    } catch {
      /* ignore broken session */
    }
  }

  return {
    buildSettingsDTO,
    persistSettings,
    autosaveDirtyDrafts,
    restoreWorkspaceSession,
  }
}

export type WorkspaceSessionController = ReturnType<typeof createWorkspaceSessionController>
