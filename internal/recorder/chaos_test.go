package recorder

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestChaosCoalesceShuffledEvents(t *testing.T) {
	actions := []string{"click", "fill", "hover", "press", "check", "uncheck"}
	selectors := []string{"#a", "#email", `label:has-text("OK")`, "nav", "input"}
	values := []string{"", "Tab", "on", "off", "hello@x.com", "Ivan"}

	for seed := int64(0); seed < 50; seed++ {
		t.Run(fmt.Sprintf("seed-%d", seed), func(t *testing.T) {
			rng := rand.New(rand.NewSource(seed))
			steps := make([]RecordedStep, 0, 1+rng.Intn(24))
			for i := 0; i < cap(steps); i++ {
				steps, _ = ApplyCoalescedStep(steps, RecordedStep{
					Action:    actions[rng.Intn(len(actions))],
					Selector:  selectors[rng.Intn(len(selectors))],
					Value:     values[rng.Intn(len(values))],
					InputType: randomInputType(rng),
					Text:      "Label",
				})
			}
			for _, step := range steps {
				if step.Action == "" {
					t.Fatal("empty action in coalesced steps")
				}
			}
		})
	}
}

func randomInputType(rng *rand.Rand) string {
	switch rng.Intn(4) {
	case 0:
		return "checkbox"
	case 1:
		return "text"
	default:
		return ""
	}
}
