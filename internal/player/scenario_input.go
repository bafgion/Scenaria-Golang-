package player

import (
	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/settings"
)
type ScenarioInput struct {
	FeaturePath  string
	ScenarioName string
	Steps        []gherkin.Step
	TestClient   *settings.TestClient
	Variables    map[string]string
	ProjectRoot  string
	StartStep    int
	EndStep      int
}
