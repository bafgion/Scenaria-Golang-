package selector

import _ "embed"

//go:embed picker_script.js
var pickerInstallScriptJS string

// PickerInstallScript installs the in-browser element picker overlay.
var PickerInstallScript = pickerInstallScriptJS

const PickerUninstallScript = `(() => {
  if (typeof window.__shopPickerCleanup === 'function') {
    window.__shopPickerCleanup();
  }
})();`
