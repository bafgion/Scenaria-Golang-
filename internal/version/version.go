package version

const (
	AppName = "Scenaria"
	Version = "0.12.0"
	Module  = "github.com/bafgion/scenaria-golang"
)

func String() string {
	return AppName + " " + Version
}
