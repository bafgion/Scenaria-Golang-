package player

import (
	"math/rand"
	"regexp"
	"testing"
)

func TestNormalizeGeneratorAliases(t *testing.T) {
	cases := map[string]string{
		"телефон":        "phone",
		"имя":            "first_name",
		"фамилия":        "last_name",
		"отчество":       "patronymic",
		"ИНН":            "inn",
		"расчётный счёт": "bank_account",
		"огрнип":         "ogrnip",
	}
	for raw, want := range cases {
		got, ok := NormalizeGeneratorName(raw)
		if !ok || got != want {
			t.Fatalf("NormalizeGeneratorName(%q) = %q, %v want %q", raw, got, ok, want)
		}
	}
}

func TestRunContextReusesValuesWithinRun(t *testing.T) {
	ctx := NewRunContext(nil, 1, "")
	first, err := ctx.generate("phone")
	if err != nil {
		t.Fatal(err)
	}
	second, err := ctx.generate("phone")
	if err != nil || first != second {
		t.Fatalf("phone should be stable within run: %q %q", first, second)
	}
	resolved, err := ctx.ResolveText("{{phone}}")
	if err != nil || resolved != first {
		t.Fatalf("resolve {{phone}} = %q want %q", resolved, first)
	}
}

func TestRunContextNewValueEachGenerator(t *testing.T) {
	ctx := NewRunContext(nil, 2, "")
	phone, _ := ctx.generate("phone")
	first, _ := ctx.generate("first_name")
	if phone == first {
		t.Fatalf("different generators should differ: %q %q", phone, first)
	}
}

func TestPersonFieldsAreGenderConsistent(t *testing.T) {
	ctx := NewRunContext(nil, 3, "")
	first, _ := ctx.generate("first_name")
	last, _ := ctx.generate("last_name")
	patronymic, _ := ctx.generate("patronymic")
	femaleLast := stringsHasSuffix(last, "а")
	femalePat := stringsHasSuffix(patronymic, "на")
	malePat := stringsHasSuffix(patronymic, "ич")
	if first == "" || last == "" || patronymic == "" {
		t.Fatal("person fields should be non-empty")
	}
	if femaleLast != femalePat && !(femaleLast == false && malePat) {
		t.Fatalf("gender mismatch: last=%q patronymic=%q", last, patronymic)
	}
}

func TestGeneratedFormats(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	phone, _ := generateCanonicalForTest(rng, "phone")
	inn := generateINN(rng)
	account := "40817810" + randomDigits(rng, 12)
	ogrnip := generateOGRNIP(rng)

	phoneRE := regexp.MustCompile(`^\+79\d{9}$`)
	if !phoneRE.MatchString(phone) {
		t.Fatalf("phone format: %q", phone)
	}
	if len(inn) != 12 {
		t.Fatalf("inn length: %q", inn)
	}
	if len(account) != 20 {
		t.Fatalf("account length: %q", account)
	}
	if len(ogrnip) != 15 || ogrnip[0] != '3' {
		t.Fatalf("ogrnip format: %q", ogrnip)
	}
}

func TestUnknownGeneratorRaises(t *testing.T) {
	ctx := NewRunContext(nil, 0, "")
	_, err := ctx.generate("unknown")
	if err == nil {
		t.Fatal("expected error for unknown generator")
	}
}

func generateCanonicalForTest(rng *rand.Rand, canonical string) (string, error) {
	ctx := NewRunContext(nil, rng.Int63(), "")
	ctx.rng = rng
	return ctx.generateCanonical(canonical)
}

func stringsHasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
