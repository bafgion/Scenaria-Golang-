export const confirm = {
  deleteFeature: {
    title: 'Delete scenario',
    message: 'Delete «{name}» permanently? It will not be moved to Recycle Bin or Trash.',
    confirmLabel: 'Delete',
  },
  openOtherProject: {
    title: 'Open another project',
    messageDirty: '{count} unsaved file(s). Close the current project and open a new one?',
    messageClean: 'Close the current project and open a new one?',
    confirmLabel: 'Open',
  },
  closeProject: {
    title: 'Close project',
    message: '{count} unsaved file(s). Close the project without saving?',
    confirmLabel: 'Close',
  },
  diskChanged: {
    title: 'File changed on disk',
    message:
      '«{name}» was modified outside the editor. Reload from disk? Unsaved editor changes will be lost.',
    confirmLabel: 'Reload',
  },
  recordingActive: {
    title: 'Recording active',
    message:
      'Steps are being recorded into «{from}». Switch tab? Further steps will be recorded into «{to}».',
    confirmLabel: 'Switch',
    dontAskAgainLabel: "Don't ask again",
  },
  recordingTargetTab: {
    title: 'Recording target tab',
    messagePaused: '«{name}» is the recording target (paused). Close the tab?',
    messageRecording: 'Recording into «{name}». Pause recording or close the tab deliberately.',
    confirmLabelClose: 'Close',
    confirmLabelForceClose: 'Close anyway',
  },
  headless: {
    title: 'Headless',
    message:
      'Changing between headed and headless mode applies to the current session and may restart the browser. Cookies, session data, and page state can be cleared. Continue?',
    confirmLabel: 'Apply',
  },
  generic: {
    title: 'Confirm',
    confirmLabelDelete: 'Delete',
  },
  closeBrowser: {
    message: 'There are unsaved changes. Close the browser?',
  },
} as const
