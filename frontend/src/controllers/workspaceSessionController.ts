import { buildSessionTabsSnapshot, resolveRestoredActiveTab, sessionTabPathsFromSettings, untitledContentMap } from '../lib/sessionTabs'
import type { TabBody } from '../lib/tabMemory'
import { tabEditorText } from '../lib/tabMemory'
import { syncUntitledCounterFromPaths } from '../lib/untitled'
import {
  clearUntitledRecoveryAll,
  clearUntitledRecoveryPath,
  journalUntitledTabs,
  loadUntitledRecoveryJournal,
  mergeUntitledSessionWithRecovery,
  type UntitledRecoveryStorage,
} from '../lib/untitledRecoveryJournal'
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
  hydrate?: boolean
}

export type WorkspaceSessionContext = {
  stores: WorkspaceSessionStores
  welcomeKey: string
  getProjectPath: () => string
  getTabs: () => TabBody[]
  getActiveTab: () => string
  getEditorText: () => string | null
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
  recoveryStorage?: UntitledRecoveryStorage | null
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
    const sessionTabs = buildSessionTabsSnapshot(tabs, activeTab, ctx.getEditorText, ctx.welcomeKey)
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
    journalUntitledTabs(ctx.getTabs(), ctx.getActiveTab(), ctx.getEditorText, ctx.recoveryStorage)
    await ctx.saveSettings(buildSettingsDTO())
  }

  function journalCurrentUntitledTabs() {
    journalUntitledTabs(ctx.getTabs(), ctx.getActiveTab(), ctx.getEditorText, ctx.recoveryStorage)
  }

  function clearUntitledJournalPath(path: string) {
    clearUntitledRecoveryPath(path, ctx.recoveryStorage)
  }

  function clearUntitledJournal() {
    clearUntitledRecoveryAll(ctx.recoveryStorage)
  }

  async function autosaveDirtyDrafts() {
    if (!ctx.getProjectPath()) return
    ctx.syncActiveTabContent()
    const activeTab = ctx.getActiveTab()
    for (const tab of ctx.getTabs()) {
      if (!tab.dirty || ctx.isUntitled(tab.path)) continue
      const liveText = tab.path === activeTab ? ctx.getEditorText() : null
      const text = liveText !== null ? liveText : tabEditorText(tab)
      try {
        await ctx.saveFeatureDraft(tab.path, text)
      } catch {
        /* offline */
      }
    }
  }

  async function restoreWorkspaceSession(s: gui.AppSettingsDTO) {
    const proj = (s.sessionProject || '').trim()
    const hasOpenTabs = (s.openTabs || []).some((p) => (p || '').trim())
    const recoveryTabs = loadUntitledRecoveryJournal(ctx.recoveryStorage)
    const mergedUntitledTabs = mergeUntitledSessionWithRecovery(s.untitledTabs, recoveryTabs)
    const hasUntitled = mergedUntitledTabs.some((t) => (t?.path || '').trim())
    if (!proj && !hasOpenTabs && !hasUntitled) return

    const untitledBodies = untitledContentMap(mergedUntitledTabs)
    syncUntitledCounterFromPaths([...(s.openTabs || []), ...untitledBodies.keys()])
    const tabPaths = sessionTabPathsFromSettings(s.openTabs, mergedUntitledTabs)
    const restoreUntitledTabs = () => {
      for (const p of tabPaths) {
        if (!ctx.isUntitled(p)) continue
        const content = untitledBodies.get(p)
        if (content === undefined || ctx.getTabs().some((t) => t.path === p)) continue
        ctx.stores.tabsStore.appendTab({ path: p, content, draft: content, dirty: true })
      }
    }
    const restoreTabOrder = () => {
      const tabs = ctx.getTabs()
      const byPath = new Map(tabs.map((tab) => [tab.path, tab]))
      const ordered = tabPaths.map((path) => byPath.get(path)).filter((tab): tab is TabBody => Boolean(tab))
      const orderedPaths = new Set(ordered.map((tab) => tab.path))
      const remaining = tabs.filter((tab) => !orderedPaths.has(tab.path))
      if (ordered.length > 0) {
        ctx.stores.tabsStore.setTabs([...ordered, ...remaining])
      }
    }
    const activateRestoredTab = async () => {
      restoreTabOrder()
      const active = (s.activeTab || '').trim()
      const focusPath = resolveRestoredActiveTab(active, tabPaths, ctx.getTabs(), ctx.welcomeKey)
      if (!focusPath || focusPath === ctx.welcomeKey) return
      ctx.stores.tabsStore.patch({ welcomeTabVisible: false, activeTab: focusPath })
      if (ctx.isUntitled(focusPath)) {
        const tab = ctx.getTabs().find((t) => t.path === focusPath)
        if (tab) {
          await ctx.applyEditorText(tabEditorText(tab), {
            switchTab: true,
            tabPath: focusPath,
            skipValidate: true,
            hydrate: true,
          })
          ctx.trimTabsMemory()
        }
      }
    }

    try {
      restoreUntitledTabs()
    } catch {
      /* ignore broken untitled session */
    }

    if (proj) {
      const resolvedProj = (await ctx.resolveProjectPath(proj)).trim()
      if (!resolvedProj) {
        await activateRestoredTab()
        return
      }
      try {
        const info = await ctx.openProject(resolvedProj)
        ctx.applyProjectScan(info, resolvedProj)
        const clients = await Promise.resolve(ctx.listTestClients()).catch((): string[] => [])
        ctx.stores.testClientStore.setClients(clients ?? [])
      } catch {
        ctx.appendLog(ctx.tr('journal.session.projectNotFound', { path: resolvedProj }))
        ctx.setStatus(ctx.tr('journal.status.sessionProjectNotFound'), 'error')
        await activateRestoredTab()
        return
      }
    }
    try {
      for (const p of tabPaths) {
        if (ctx.isUntitled(p)) {
          continue
        }
        try {
          await ctx.loadFeature(p)
        } catch {
          /* skip missing files */
        }
      }
      await activateRestoredTab()
    } catch {
      /* ignore broken session */
    }
  }

  return {
    buildSettingsDTO,
    persistSettings,
    journalCurrentUntitledTabs,
    clearUntitledJournalPath,
    clearUntitledJournal,
    autosaveDirtyDrafts,
    restoreWorkspaceSession,
  }
}

export type WorkspaceSessionController = ReturnType<typeof createWorkspaceSessionController>
