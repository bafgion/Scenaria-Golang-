package player

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestResolveEmailForCodeFromPriorFill(t *testing.T) {
	ctx := NewRunContext(map[string]string{}, 1, "")
	email, err := ctx.ResolveEmailForCode("", []gherkin.Step{
		{Text: `ввожу "user@example.com" в "#email"`},
	})
	if err != nil || email != "user@example.com" {
		t.Fatalf("unexpected email: %v %q", err, email)
	}
}
