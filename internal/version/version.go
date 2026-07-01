package version

import "github.com/bafgion/scenaria-golang/internal/brand"

const (
	AppName = brand.Name
	Version = "0.18.0"
	Module  = "github.com/bafgion/scenaria-golang"
)

func String() string {
	return AppName + " " + Version
}
