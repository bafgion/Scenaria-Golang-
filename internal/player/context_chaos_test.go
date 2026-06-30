package player

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestChaosResolveTextRandomPlaceholders(t *testing.T) {
	generators := []string{"phone", "first_name", "last_name", "inn", "ogrnip", "address", "bank_account", "patronymic"}
	const iterations = 200
	for i := 0; i < iterations; i++ {
		seed := int64(i + 1)
		rng := rand.New(rand.NewSource(seed))
		ctx := NewRunContext(map[string]string{"known": "value"}, seed, t.TempDir())

		parts := make([]string, 0, 1+rng.Intn(4))
		for j := 0; j < cap(parts); j++ {
			switch rng.Intn(4) {
			case 0:
				parts = append(parts, "{{known}}")
			case 1:
				parts = append(parts, fmt.Sprintf("{{%s}}", generators[rng.Intn(len(generators))]))
			case 2:
				parts = append(parts, "plain")
			default:
				parts = append(parts, "{{missing_"+fmt.Sprint(rng.Intn(1000))+"}}")
			}
		}
		text := strings.Join(parts, " ")
		out, err := ctx.ResolveText(text)
		if err != nil {
			if strings.Contains(text, "{{missing_") {
				continue
			}
			t.Fatalf("seed %d text %q: %v", seed, text, err)
		}
		if strings.Contains(out, "{{") && strings.Contains(out, "}}") {
			// unknown placeholders should error, not leak
			if !strings.Contains(text, "missing_") {
				t.Fatalf("unresolved placeholders in %q -> %q", text, out)
			}
		}
	}
}

func TestChaosMaxLoopIterationsEnforced(t *testing.T) {
	exec := NewStepExecutor(ExecutorOptions{MaxLoopIterations: 3})
	steps := []gherkin.Step{{
		Block:       gherkin.BlockRepeat,
		RepeatCount: 100,
		Children:    []gherkin.Step{{Text: `жду 1 мс`, Line: 2}},
		Line:        1,
	}}
	err := exec.ExecuteSteps(context.Background(), &browserSession{}, steps, NewRunContext(nil, 1, t.TempDir()))
	if err != nil {
		t.Fatalf("repeat capped execution failed: %v", err)
	}
}
