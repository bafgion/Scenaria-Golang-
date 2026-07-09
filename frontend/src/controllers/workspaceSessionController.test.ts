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
})
