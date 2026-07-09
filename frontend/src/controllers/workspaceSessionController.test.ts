import { describe, expect, it, vi } from 'vitest'
import { createWorkspaceSessionController } from './workspaceSessionController'
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
      applyProjectScan: vi.fn(),
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
})
