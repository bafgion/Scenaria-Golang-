package update

// Test hooks (overridden in tests only).
var (
	detectInstallModeOverride func(string) InstallMode
	launchHiddenBatchHook     func(string) error
)

func detectInstallMode(installDir string) InstallMode {
	if detectInstallModeOverride != nil {
		return detectInstallModeOverride(installDir)
	}
	return DetectInstallMode(installDir)
}
