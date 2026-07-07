package player

import "testing"

func TestResolvePlaceholderExpressionInsideBraces(t *testing.T) {
	ctx := NewRunContext(map[string]string{"total": "7180"}, 1, "")
	got, err := ctx.ResolveText("{{total / 4}}")
	if err != nil || got != "1795" {
		t.Fatalf("ResolveText() = %q, %v", got, err)
	}
}

func TestResolvePlaceholderExpressionAfterSubstitution(t *testing.T) {
	ctx := NewRunContext(map[string]string{"total": "7180"}, 1, "")
	got, err := ctx.ResolveText("{{total}} / 4")
	if err != nil || got != "1795" {
		t.Fatalf("ResolveText() = %q, %v", got, err)
	}
}

func TestResolvePlaceholderExpressionWithVariables(t *testing.T) {
	ctx := NewRunContext(map[string]string{
		"total":  "7200",
		"points": "200",
	}, 1, "")
	got, err := ctx.ResolveText("{{(total - points) / 4}}")
	if err != nil || got != "1750" {
		t.Fatalf("ResolveText() = %q, %v", got, err)
	}
}

func TestResolvePlaceholderExpressionDivisionByZero(t *testing.T) {
	ctx := NewRunContext(map[string]string{"total": "100"}, 1, "")
	_, err := ctx.ResolveText("{{total / 0}}")
	if err == nil {
		t.Fatal("expected division by zero error")
	}
}

func TestEvaluateArithmeticExpression(t *testing.T) {
	ctx := NewRunContext(map[string]string{"total": "7180"}, 1, "")
	v, err := ctx.evaluateArithmeticExpression("total / 4")
	if err != nil || v != 1795 {
		t.Fatalf("evaluateArithmeticExpression() = %d, %v", v, err)
	}
}

func TestResolvePlaceholderExpressionDoesNotManglePhone(t *testing.T) {
	ctx := NewRunContext(nil, 1, "")
	ctx.values["phone"] = "+79161234567"
	got, err := ctx.ResolveText("{{phone}}")
	if err != nil || got != "+79161234567" {
		t.Fatalf("ResolveText() = %q, %v", got, err)
	}
}

func TestParseNumericOperand(t *testing.T) {
	v, err := parseNumericOperand("7 180 ₽")
	if err != nil || v != 7180 {
		t.Fatalf("parseNumericOperand() = %d, %v", v, err)
	}
}
