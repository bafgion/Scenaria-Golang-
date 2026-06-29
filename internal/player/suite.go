package player

import (
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/settings"
)

type RunCase struct {
	FeaturePath string
	Name        string
	Steps       []gherkin.Step
	TestClient  *settings.TestClient
	Variables   map[string]string
	ProjectRoot string
}

type ExecutionPlan struct {
	Cases []RunCase
}

func BuildExecutionPlan(features []FeatureInput, tag string, variables map[string]string) ExecutionPlan {
	return buildExecutionPlan(features, tag, "", variables, "")
}

func BuildExecutionPlanWithTestClient(features []FeatureInput, tag, scenario string, variables map[string]string, testClientOverride string) ExecutionPlan {
	return buildExecutionPlan(features, tag, scenario, variables, testClientOverride)
}

func buildExecutionPlan(features []FeatureInput, tag, scenario string, variables map[string]string, testClientOverride string) ExecutionPlan {
	plan := ExecutionPlan{
		Cases: make([]RunCase, 0),
	}
	for _, input := range features {
		if tag != "" && !gherkin.FeatureHasTag(input.Feature, tag) {
			continue
		}
		projectRoot := paths.InferProjectRoot([]string{input.Path})
		var testClient *settings.TestClient
		clientName := strings.TrimSpace(input.Feature.TestClient)
		if override := strings.TrimSpace(testClientOverride); override != "" {
			clientName = override
		}
		if clientName != "" && projectRoot != "" {
			client, err := loadTestClientForFeature(projectRoot, clientName)
			if err == nil {
				testClient = client
			}
		}
		for _, runnable := range gherkin.ExpandFeature(input.Feature) {
			if tag != "" && !gherkin.TagsInclude(runnable.Tags, tag) {
				continue
			}
			if scenario != "" && !strings.EqualFold(strings.TrimSpace(runnable.Title), strings.TrimSpace(scenario)) {
				continue
			}
			plan.Cases = append(plan.Cases, RunCase{
				FeaturePath: input.Path,
				Name:        runnable.Title,
				Steps:       runnable.Steps,
				TestClient:  testClient,
				Variables:   variables,
				ProjectRoot: projectRoot,
			})
		}
	}
	return plan
}

func SummarizePlan(plan ExecutionPlan) (files int, scenarios int, steps int, results []ScenarioResult) {
	seenFiles := make(map[string]struct{})
	results = make([]ScenarioResult, 0, len(plan.Cases))
	for _, runCase := range plan.Cases {
		seenFiles[runCase.FeaturePath] = struct{}{}
		steps += gherkin.CountLeafSteps(runCase.Steps)
		results = append(results, ScenarioResult{
			FeaturePath: runCase.FeaturePath,
			Scenario:    runCase.Name,
			Status:      "skipped",
			Message:     "dry-run mode",
		})
	}
	return len(seenFiles), len(plan.Cases), steps, results
}
