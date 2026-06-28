package version

const (
	AppName = "Scenaria"
	Version = "0.1.0-go"
	Module  = "github.com/bafgion/scenaria-golang"
)

func String() string {
	return AppName + " " + Version
}
