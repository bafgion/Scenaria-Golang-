package update

import (
	"strings"

	"github.com/bafgion/scenaria-golang/internal/brand"
)

// InstallMode describes how the app was deployed.
type InstallMode string

const (
	InstallModePortable InstallMode = "portable"
	InstallModeSetup    InstallMode = "setup"
)

// ApplyKind is the update mechanism used for a release asset.
type ApplyKind string

const (
	ApplyKindPortable ApplyKind = "portable"
	ApplyKindSetup    ApplyKind = "setup"
)

func pickAssetForMode(mode InstallMode, assets []ReleaseAsset) *ReleaseAsset {
	switch mode {
	case InstallModeSetup:
		if asset := matchAsset(assets, func(name string) bool {
			return name == strings.ToLower(brand.SetupExe)
		}); asset != nil {
			return asset
		}
		return pickPortableAsset(assets)
	default:
		if asset := pickPortableAsset(assets); asset != nil {
			return asset
		}
		return matchAsset(assets, func(name string) bool { return name == strings.ToLower(brand.SetupExe) })
	}
}

func pickPortableAsset(assets []ReleaseAsset) *ReleaseAsset {
	for _, pred := range []func(string) bool{
		func(name string) bool { return name == strings.ToLower(brand.PortableZip) },
		func(name string) bool { return name == strings.ToLower(brand.UpdateZip) },
		func(name string) bool { return strings.Contains(strings.ToLower(name), "portable") && strings.HasSuffix(strings.ToLower(name), ".zip") },
	} {
		if asset := matchAsset(assets, pred); asset != nil {
			return asset
		}
	}
	return nil
}

func applyKindForAsset(name string) ApplyKind {
	lower := strings.ToLower(strings.TrimSpace(name))
	if strings.HasSuffix(lower, ".zip") {
		return ApplyKindPortable
	}
	return ApplyKindSetup
}
