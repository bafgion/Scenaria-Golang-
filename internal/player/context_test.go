package player

import "testing"

func TestGenerateByKind(t *testing.T) {
	ctx := NewRunContext(nil, 42, "")
	phone, err := ctx.GenerateByKind("phone")
	if err != nil || phone == "" {
		t.Fatalf("phone generator failed: %v %q", err, phone)
	}
	inn, err := ctx.GenerateByKind("inn")
	if err != nil || len(inn) != 10 {
		t.Fatalf("inn generator failed: %v %q", err, inn)
	}
}

func TestEmailCodePromptCallback(t *testing.T) {
	EmailCodePrompt = func(email string) (string, error) {
		if email != "" && email != "a@b.c" {
			t.Fatalf("unexpected email: %q", email)
		}
		return "654321", nil
	}
	defer func() { EmailCodePrompt = nil }()
	ctx := NewRunContext(nil, 1, "")
	code, err := ctx.EmailCode()
	if err != nil || code != "654321" {
		t.Fatalf("unexpected code: %v %q", err, code)
	}
}

func TestEmailCodeFromEnv(t *testing.T) {
	EmailCodePrompt = nil
	t.Setenv("SCENARIA_EMAIL_CODE", "123456")
	ctx := NewRunContext(nil, 1, "")
	code, err := ctx.EmailCode()
	if err != nil || code != "123456" {
		t.Fatalf("unexpected email code: %v %q", err, code)
	}
}
