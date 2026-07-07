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

func TestResolveEmailForCodeFromStepField(t *testing.T) {
	ctx := NewRunContext(map[string]string{}, 1, "")
	email, err := ctx.ResolveEmailForCode("qa@test.com", nil)
	if err != nil || email != "qa@test.com" {
		t.Fatalf("unexpected email: %v %q", err, email)
	}
}

func TestEmailCodeForStepUsesStepEmail(t *testing.T) {
	var gotEmail string
	ctx := NewRunContext(map[string]string{}, 1, "", WithPromptEmailCode(func(email string) (string, error) {
		gotEmail = email
		return "123456", nil
	}))
	code, err := ctx.EmailCodeForStep("qa@test.com")
	if err != nil || code != "123456" || gotEmail != "qa@test.com" {
		t.Fatalf("unexpected prompt: err=%v code=%q email=%q", err, code, gotEmail)
	}
}

func TestEmailCodeForStepSkipsPromptWithEnv(t *testing.T) {
	t.Setenv("SCENARIA_EMAIL_CODE", "999999")
	ctx := NewRunContext(map[string]string{}, 1, "", WithPromptEmailCode(func(email string) (string, error) {
		t.Fatal("prompt should not run when env is set")
		return "", nil
	}))
	code, err := ctx.EmailCodeForStep("qa@test.com")
	if err != nil || code != "999999" {
		t.Fatalf("unexpected code: err=%v code=%q", err, code)
	}
}
