import { describe, expect, it, vi } from 'vitest'
import { createWorkspaceSessionController, hasRestorableWorkspaceSession } from './workspaceSessionController'
import { createSettingsStore } from '../stores/settingsStore'
import { createRecorderPrefsStore } from '../stores/recorderPrefsStore'
import { createUiPrefsStore } from '../stores/uiPrefsStore'
import { createRecentsStore } from '../stores/recentsStore'
import { createDialogBindController } from './dialogBindController'
import { createRecordFormStore } from '../stores/recordFormStore'
import { createValidateDialogStore } from '../stores/validateDialogStore'
import { createFeatureDialogStore } from '../stores/featureDialogStore'
import { createProjectReplaceStore } from '../stores/projectReplaceStore'
import { createTestClientStore } from '../stores/testClientStore'
import { createVanessaRunStore } from '../stores/vanessaRunStore'
import { createPluginRunStore } from '../stores/pluginRunStore'
import { createRunFormStore } from '../stores/runFormStore'
import { createSettingsDialogStore } from '../stores/settingsDialogStore'
import { createTabsStore } from '../stores/tabsStore'

describe('workspaceSessionController', () => {
  it('detects saved workspace sessions independently from onboarding state', () => {
    expect(hasRestorableWorkspaceSession(null)).toBe(false)
    expect(hasRestorableWorkspaceSession({})).toBe(false)
    expect(hasRestorableWorkspaceSession({ sessionProject: ' C:/proj ' })).toBe(true)
    expect(hasRestorableWorkspaceSession({ openTabs: [''] })).toBe(false)
    expect(hasRestorableWorkspaceSession({ openTabs: ['C:/proj/a.feature'] })).toBe(true)
    expect(hasRestorableWorkspaceSession({ activeTab: 'C:/proj/a.feature' })).toBe(true)
    expect(hasRestorableWorkspaceSession({ untitledTabs: [{ path: '__untitled_1__', content: 'draft' }] })).toBe(true)
  })

  it('builds settings DTO from store snapshots', () => {
    const settingsStore = createSettingsStore()
    const dialogBinds = createDialogBindController({
      settingsStore,
      uiPrefsStore: createUiPrefsStore(),
      recorderPrefsStore: createRecorderPrefsStore(),
      recordFormStore: createRecordFormStore(),
      validateDialogStore: createValidateDialogStore(),
      featureDialogStore: createFeatureDialogStore(),
      projectReplaceStore: createProjectReplaceStore(),
      testClientStore: createTestClientStore(),
      vanessaRunStore: createVanessaRunStore(),
      pluginRunStore: createPluginRunStore(),
      runFormStore: createRunFormStore(),
      settingsDialogStore: createSettingsDialogStore(),
    })
    const recentsStore = createRecentsStore()
    recentsStore.patch({ projects: ['a'], features: ['b.feature'] })
    const applyProjectScan = vi.fn()
    const controller = createWorkspaceSessionController({
      stores: {
        settingsStore,
        recorderPrefsStore: createRecorderPrefsStore(),
        uiPrefsStore: createUiPrefsStore(),
        recentsStore,
        dialogBinds,
        tabsStore: createTabsStore('__welcome__'),
        testClientStore: createTestClientStore(),
      },
      welcomeKey: '__welcome__',
      getProjectPath: () => '/proj',
      getTabs: () => [{ path: '/proj/a.feature', content: 'x', dirty: false }],
      getActiveTab: () => '/proj/a.feature',
      getEditorText: () => 'Feature:',
      syncActiveTabContent: vi.fn(),
      isUntitled: () => false,
      saveSettings: vi.fn(),
      saveFeatureDraft: vi.fn(),
      openProject: vi.fn(),
      resolveProjectPath: vi.fn(async (path: string) => path),
      applyProjectScan,
      listTestClients: vi.fn(),
      loadFeature: vi.fn(),
      applyEditorText: vi.fn(),
      trimTabsMemory: vi.fn(),
      appendLog: vi.fn(),
      setStatus: vi.fn(),
      tr: (key) => key,
    })

    const dto = controller.buildSettingsDTO()
    expect(dto.sessionProject).toBe('/proj')
    expect(dto.openTabs).toEqual(['/proj/a.feature'])
    expect(dto.activeTab).toBe('/proj/a.feature')
    expect(dto.recentProjects).toEqual(['a'])
    expect(dto.recentFeatures).toEqual(['b.feature'])
  })

  it('restores session through resolved project path', async () => {
    const settingsStore = createSettingsStore()
    const dialogBinds = createDialogBindController({
      settingsStore,
      uiPrefsStore: createUiPrefsStore(),
      recorderPrefsStore: createRecorderPrefsStore(),
      recordFormStore: createRecordFormStore(),
      validateDialogStore: createValidateDialogStore(),
      featureDialogStore: createFeatureDialogStore(),
      projectReplaceStore: createProjectReplaceStore(),
      testClientStore: createTestClientStore(),
      vanessaRunStore: createVanessaRunStore(),
      pluginRunStore: createPluginRunStore(),
      runFormStore: createRunFormStore(),
      settingsDialogStore: createSettingsDialogStore(),
    })
    const openProject = vi.fn().mockResolvedValue({ path: '/repo/examples', features: [], tags: [], featureTags: {}, version: 1 })
    const resolveProjectPath = vi.fn(async (path: string) => (path === 'examples' ? '/repo/examples' : path))
    const applyProjectScan = vi.fn()
    const controller = createWorkspaceSessionController({
      stores: {
        settingsStore,
        recorderPrefsStore: createRecorderPrefsStore(),
        uiPrefsStore: createUiPrefsStore(),
        recentsStore: createRecentsStore(),
        dialogBinds,
        tabsStore: createTabsStore('__welcome__'),
        testClientStore: createTestClientStore(),
      },
      welcomeKey: '__welcome__',
      getProjectPath: () => '',
      getTabs: () => [],
      getActiveTab: () => '',
      getEditorText: () => '',
      syncActiveTabContent: vi.fn(),
      isUntitled: () => false,
      saveSettings: vi.fn(),
      saveFeatureDraft: vi.fn(),
      openProject,
      resolveProjectPath,
      applyProjectScan,
      listTestClients: vi.fn(),
      loadFeature: vi.fn(),
      applyEditorText: vi.fn(),
      trimTabsMemory: vi.fn(),
      appendLog: vi.fn(),
      setStatus: vi.fn(),
      tr: (key) => key,
    })

    await controller.restoreWorkspaceSession({ sessionProject: 'examples' } as never)
    expect(resolveProjectPath).toHaveBeenCalledWith('examples')
    expect(openProject).toHaveBeenCalledWith('/repo/examples')
    expect(applyProjectScan).toHaveBeenCalledWith(expect.objectContaining({ path: '/repo/examples' }), '/repo/examples')
  })

  it('does not rewrite absolute examples project paths during session restore', async () => {
    const settingsStore = createSettingsStore()
    const dialogBinds = createDialogBindController({
      settingsStore,
      uiPrefsStore: createUiPrefsStore(),
      recorderPrefsStore: createRecorderPrefsStore(),
      recordFormStore: createRecordFormStore(),
      validateDialogStore: createValidateDialogStore(),
      featureDialogStore: createFeatureDialogStore(),
      projectReplaceStore: createProjectReplaceStore(),
      testClientStore: createTestClientStore(),
      vanessaRunStore: createVanessaRunStore(),
      pluginRunStore: createPluginRunStore(),
      runFormStore: createRunFormStore(),
      settingsDialogStore: createSettingsDialogStore(),
    })
    const absoluteExamples = 'C:\\Users\\bafgion\\Documents\\Projects\\scenaria_go\\examples'
    const openProject = vi.fn().mockResolvedValue({ path: absoluteExamples, features: [], tags: [], featureTags: {}, version: 1 })
    const resolveProjectPath = vi.fn(async (path: string) => path)
    const controller = createWorkspaceSessionController({
      stores: {
        settingsStore,
        recorderPrefsStore: createRecorderPrefsStore(),
        uiPrefsStore: createUiPrefsStore(),
        recentsStore: createRecentsStore(),
        dialogBinds,
        tabsStore: createTabsStore('__welcome__'),
        testClientStore: createTestClientStore(),
      },
      welcomeKey: '__welcome__',
      getProjectPath: () => '',
      getTabs: () => [],
      getActiveTab: () => '',
      getEditorText: () => '',
      syncActiveTabContent: vi.fn(),
      isUntitled: () => false,
      saveSettings: vi.fn(),
      saveFeatureDraft: vi.fn(),
      openProject,
      resolveProjectPath,
      applyProjectScan: vi.fn(),
      listTestClients: vi.fn(),
      loadFeature: vi.fn(),
      applyEditorText: vi.fn(),
      trimTabsMemory: vi.fn(),
      appendLog: vi.fn(),
      setStatus: vi.fn(),
      tr: (key) => key,
    })

    await controller.restoreWorkspaceSession({ sessionProject: absoluteExamples } as never)
    expect(resolveProjectPath).toHaveBeenCalledWith(absoluteExamples)
    expect(openProject).toHaveBeenCalledWith(absoluteExamples)
  })

  it('activates restored feature tab when saved active tab is welcome', async () => {
    const settingsStore = createSettingsStore()
    const dialogBinds = createDialogBindController({
      settingsStore,
      uiPrefsStore: createUiPrefsStore(),
      recorderPrefsStore: createRecorderPrefsStore(),
      recordFormStore: createRecordFormStore(),
      validateDialogStore: createValidateDialogStore(),
      featureDialogStore: createFeatureDialogStore(),
      projectReplaceStore: createProjectReplaceStore(),
      testClientStore: createTestClientStore(),
      vanessaRunStore: createVanessaRunStore(),
      pluginRunStore: createPluginRunStore(),
      runFormStore: createRunFormStore(),
      settingsDialogStore: createSettingsDialogStore(),
    })
    const tabsStore = createTabsStore('__welcome__')
    const loadFeature = vi.fn(async (path: string) => {
      tabsStore.appendTab({ path, content: `Feature: ${path}`, dirty: false })
    })
    const controller = createWorkspaceSessionController({
      stores: {
        settingsStore,
        recorderPrefsStore: createRecorderPrefsStore(),
        uiPrefsStore: createUiPrefsStore(),
        recentsStore: createRecentsStore(),
        dialogBinds,
        tabsStore,
        testClientStore: createTestClientStore(),
      },
      welcomeKey: '__welcome__',
      getProjectPath: () => '/proj',
      getTabs: () => tabsStore.snapshot().tabs,
      getActiveTab: () => tabsStore.snapshot().activeTab,
      getEditorText: () => '',
      syncActiveTabContent: vi.fn(),
      isUntitled: () => false,
      saveSettings: vi.fn(),
      saveFeatureDraft: vi.fn(),
      openProject: vi.fn().mockResolvedValue({ path: '/proj', features: [], tags: [], featureTags: {}, version: 1 }),
      resolveProjectPath: vi.fn(async (path: string) => path),
      applyProjectScan: vi.fn(),
      listTestClients: vi.fn().mockResolvedValue([]),
      loadFeature,
      applyEditorText: vi.fn(),
      trimTabsMemory: vi.fn(),
      appendLog: vi.fn(),
      setStatus: vi.fn(),
      tr: (key) => key,
    })

    await controller.restoreWorkspaceSession({
      sessionProject: '/proj',
      openTabs: ['/proj/old.feature'],
      activeTab: '__welcome__',
    } as never)

    const snap = tabsStore.snapshot()
    expect(snap.welcomeTabVisible).toBe(false)
    expect(snap.activeTab).toBe('/proj/old.feature')
    expect(loadFeature).toHaveBeenCalledWith('/proj/old.feature')
  })

  it('restores active untitled tab alongside saved feature tabs', async () => {
    const settingsStore = createSettingsStore()
    const dialogBinds = createDialogBindController({
      settingsStore,
      uiPrefsStore: createUiPrefsStore(),
      recorderPrefsStore: createRecorderPrefsStore(),
      recordFormStore: createRecordFormStore(),
      validateDialogStore: createValidateDialogStore(),
      featureDialogStore: createFeatureDialogStore(),
      projectReplaceStore: createProjectReplaceStore(),
      testClientStore: createTestClientStore(),
      vanessaRunStore: createVanessaRunStore(),
      pluginRunStore: createPluginRunStore(),
      runFormStore: createRunFormStore(),
      settingsDialogStore: createSettingsDialogStore(),
    })
    const tabsStore = createTabsStore('__welcome__')
    const untitled = '__untitled__:2/demo.feature'
    const loadFeature = vi.fn(async (path: string) => {
      tabsStore.appendTab({ path, content: `Feature: ${path}`, dirty: false })
    })
    const applyEditorText = vi.fn()
    const controller = createWorkspaceSessionController({
      stores: {
        settingsStore,
        recorderPrefsStore: createRecorderPrefsStore(),
        uiPrefsStore: createUiPrefsStore(),
        recentsStore: createRecentsStore(),
        dialogBinds,
        tabsStore,
        testClientStore: createTestClientStore(),
      },
      welcomeKey: '__welcome__',
      getProjectPath: () => '/proj',
      getTabs: () => tabsStore.snapshot().tabs,
      getActiveTab: () => tabsStore.snapshot().activeTab,
      getEditorText: () => '',
      syncActiveTabContent: vi.fn(),
      isUntitled: (path) => path.startsWith('__untitled__:'),
      saveSettings: vi.fn(),
      saveFeatureDraft: vi.fn(),
      openProject: vi.fn().mockResolvedValue({ path: '/proj', features: [], tags: [], featureTags: {}, version: 1 }),
      resolveProjectPath: vi.fn(async (path: string) => path),
      applyProjectScan: vi.fn(),
      listTestClients: vi.fn().mockResolvedValue([]),
      loadFeature,
      applyEditorText,
      trimTabsMemory: vi.fn(),
      appendLog: vi.fn(),
      setStatus: vi.fn(),
      tr: (key) => key,
    })

    await controller.restoreWorkspaceSession({
      sessionProject: '/proj',
      openTabs: ['/proj/a.feature'],
      untitledTabs: [{ path: untitled, content: 'Feature: Draft\n  Scenario: S' }],
      activeTab: untitled,
    } as never)

    const snap = tabsStore.snapshot()
    expect(snap.tabs.map((tab) => tab.path)).toEqual(['/proj/a.feature', untitled])
    expect(snap.welcomeTabVisible).toBe(false)
    expect(snap.activeTab).toBe(untitled)
    expect(applyEditorText).toHaveBeenCalledWith(
      'Feature: Draft\n  Scenario: S',
      expect.objectContaining({ switchTab: true, tabPath: untitled }),
    )
  })

  it('restores untitled-only session without a saved project', async () => {
    const tabsStore = createTabsStore('__welcome__')
    const applyEditorText = vi.fn()
    const openProject = vi.fn()
    const controller = createWorkspaceSessionController({
      stores: {
        settingsStore: createSettingsStore(),
        recorderPrefsStore: createRecorderPrefsStore(),
        uiPrefsStore: createUiPrefsStore(),
        recentsStore: createRecentsStore(),
        dialogBinds: createDialogBindController({
          settingsStore: createSettingsStore(),
          uiPrefsStore: createUiPrefsStore(),
          recorderPrefsStore: createRecorderPrefsStore(),
          recordFormStore: createRecordFormStore(),
          validateDialogStore: createValidateDialogStore(),
          featureDialogStore: createFeatureDialogStore(),
          projectReplaceStore: createProjectReplaceStore(),
          testClientStore: createTestClientStore(),
          vanessaRunStore: createVanessaRunStore(),
          pluginRunStore: createPluginRunStore(),
          runFormStore: createRunFormStore(),
          settingsDialogStore: createSettingsDialogStore(),
        }),
        tabsStore,
        testClientStore: createTestClientStore(),
      },
      welcomeKey: '__welcome__',
      getProjectPath: () => '',
      getTabs: () => tabsStore.snapshot().tabs,
      getActiveTab: () => tabsStore.snapshot().activeTab,
      getEditorText: () => '',
      syncActiveTabContent: vi.fn(),
      isUntitled: (path) => path.startsWith('__untitled__:'),
      saveSettings: vi.fn(),
      saveFeatureDraft: vi.fn(),
      openProject,
      resolveProjectPath: vi.fn(async (path: string) => path),
      applyProjectScan: vi.fn(),
      listTestClients: vi.fn().mockResolvedValue([]),
      loadFeature: vi.fn(),
      applyEditorText,
      trimTabsMemory: vi.fn(),
      appendLog: vi.fn(),
      setStatus: vi.fn(),
      tr: (key) => key,
    })

    await controller.restoreWorkspaceSession({
      sessionProject: '',
      openTabs: [],
      untitledTabs: [{ path: '__untitled__:1/novyy-scenariy.feature', content: 'Feature: Draft' }],
      activeTab: '__untitled__:1/novyy-scenariy.feature',
    } as never)

    expect(openProject).not.toHaveBeenCalled()
    expect(tabsStore.snapshot().tabs).toEqual([
      {
        path: '__untitled__:1/novyy-scenariy.feature',
        content: 'Feature: Draft',
        draft: 'Feature: Draft',
        dirty: true,
      },
    ])
    expect(applyEditorText).toHaveBeenCalled()
  })

  function createUntitledRestoreHarness() {
    const settingsStore = createSettingsStore()
    const tabsStore = createTabsStore('__welcome__')
    const applyEditorText = vi.fn()
    const loadFeature = vi.fn(async (path: string) => {
      tabsStore.appendTab({ path, content: `Feature: ${path}`, dirty: false })
    })
    const controller = createWorkspaceSessionController({
      stores: {
        settingsStore,
        recorderPrefsStore: createRecorderPrefsStore(),
        uiPrefsStore: createUiPrefsStore(),
        recentsStore: createRecentsStore(),
        dialogBinds: createDialogBindController({
          settingsStore,
          uiPrefsStore: createUiPrefsStore(),
          recorderPrefsStore: createRecorderPrefsStore(),
          recordFormStore: createRecordFormStore(),
          validateDialogStore: createValidateDialogStore(),
          featureDialogStore: createFeatureDialogStore(),
          projectReplaceStore: createProjectReplaceStore(),
          testClientStore: createTestClientStore(),
          vanessaRunStore: createVanessaRunStore(),
          pluginRunStore: createPluginRunStore(),
          runFormStore: createRunFormStore(),
          settingsDialogStore: createSettingsDialogStore(),
        }),
        tabsStore,
        testClientStore: createTestClientStore(),
      },
      welcomeKey: '__welcome__',
      getProjectPath: () => '',
      getTabs: () => tabsStore.snapshot().tabs,
      getActiveTab: () => tabsStore.snapshot().activeTab,
      getEditorText: () => '',
      syncActiveTabContent: vi.fn(),
      isUntitled: (path) => path.startsWith('__untitled__:'),
      saveSettings: vi.fn(),
      saveFeatureDraft: vi.fn(),
      openProject: vi.fn(),
      resolveProjectPath: vi.fn(async (path: string) => path),
      applyProjectScan: vi.fn(),
      listTestClients: vi.fn().mockResolvedValue([]),
      loadFeature,
      applyEditorText,
      trimTabsMemory: vi.fn(),
      appendLog: vi.fn(),
      setStatus: vi.fn(),
      tr: (key) => key,
    })
    return { controller, tabsStore, applyEditorText, loadFeature }
  }

  it('restores multiple untitled tabs with independent editable drafts', async () => {
    const { controller, tabsStore, applyEditorText } = createUntitledRestoreHarness()
    await controller.restoreWorkspaceSession({
      sessionProject: '',
      openTabs: [
        '__untitled__:1/one.feature',
        '__untitled__:2/two.feature',
        '__untitled__:3/three.feature',
      ],
      untitledTabs: [
        { path: '__untitled__:1/one.feature', content: 'Feature: One' },
        { path: '__untitled__:2/two.feature', content: 'Feature: Two' },
        { path: '__untitled__:3/three.feature', content: 'Feature: Three' },
      ],
      activeTab: '__untitled__:2/two.feature',
    } as never)

    expect(tabsStore.snapshot().activeTab).toBe('__untitled__:2/two.feature')
    expect(tabsStore.snapshot().tabs).toEqual([
      { path: '__untitled__:1/one.feature', content: 'Feature: One', draft: 'Feature: One', dirty: true },
      { path: '__untitled__:2/two.feature', content: 'Feature: Two', draft: 'Feature: Two', dirty: true },
      { path: '__untitled__:3/three.feature', content: 'Feature: Three', draft: 'Feature: Three', dirty: true },
    ])
    expect(applyEditorText).toHaveBeenCalledWith(
      'Feature: Two',
      expect.objectContaining({ switchTab: true, tabPath: '__untitled__:2/two.feature' }),
    )
  })

  it('falls back to a valid untitled tab when saved active tab is invalid', async () => {
    const { controller, tabsStore, applyEditorText } = createUntitledRestoreHarness()
    await controller.restoreWorkspaceSession({
      sessionProject: '',
      openTabs: ['__untitled__:1/one.feature', '__untitled__:2/two.feature'],
      untitledTabs: [
        { path: '__untitled__:1/one.feature', content: 'Feature: One' },
        { path: '__untitled__:2/two.feature', content: 'Feature: Two' },
      ],
      activeTab: '__untitled__:999/missing.feature',
    } as never)

    expect(tabsStore.snapshot().activeTab).toBe('__untitled__:2/two.feature')
    expect(tabsStore.snapshot().tabs).toHaveLength(2)
    expect(applyEditorText).toHaveBeenCalledWith(
      'Feature: Two',
      expect.objectContaining({ tabPath: '__untitled__:2/two.feature' }),
    )
  })

  it('deduplicates duplicate persisted untitled paths without losing non-empty text', async () => {
    const { controller, tabsStore, applyEditorText } = createUntitledRestoreHarness()
    await controller.restoreWorkspaceSession({
      sessionProject: '',
      openTabs: ['__untitled__:1/one.feature', '__untitled__:1/one.feature'],
      untitledTabs: [
        { path: '__untitled__:1/one.feature', content: 'Feature: One' },
        { path: '__untitled__:1/one.feature', content: '' },
      ],
      activeTab: '__untitled__:1/one.feature',
    } as never)

    expect(tabsStore.snapshot().tabs).toEqual([
      { path: '__untitled__:1/one.feature', content: 'Feature: One', draft: 'Feature: One', dirty: true },
    ])
    expect(applyEditorText).toHaveBeenCalledWith(
      'Feature: One',
      expect.objectContaining({ tabPath: '__untitled__:1/one.feature' }),
    )
  })
})
