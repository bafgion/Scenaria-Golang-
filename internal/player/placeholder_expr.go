package player

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func isPlaceholderExpression(key string) bool {
	for _, r := range key {
		switch r {
		case '+', '-', '*', '/', '(', ')':
			return true
		}
	}
	return false
}

func looksLikeArithmeticExpression(text string) bool {
	text = strings.TrimSpace(text)
	if text == "" || strings.Contains(text, "{{") || !strings.ContainsAny(text, "+-*/") {
		return false
	}
	// Phone numbers like +79161234567 are not arithmetic.
	if strings.HasPrefix(text, "+") && !strings.ContainsAny(text[1:], "+-*/") {
		return false
	}
	return true
}

func (c *RunContext) evaluatePlaceholderExpression(expr string) (string, error) {
	v, err := c.evaluateArithmeticExpression(expr)
	if err != nil {
		return "", fmt.Errorf("placeholder expression %q: %w", expr, err)
	}
	return formatExprInt(v), nil
}

func (c *RunContext) evaluateArithmeticExpression(expr string) (int64, error) {
	p := &exprParser{c: c, input: strings.TrimSpace(expr)}
	v, err := p.parseExpr()
	if err != nil {
		return 0, err
	}
	p.skipSpace()
	if p.pos < len(p.input) {
		return 0, fmt.Errorf("unexpected trailing input %q", p.input[p.pos:])
	}
	return v, nil
}

func formatExprInt(n int64) string {
	return strconv.FormatInt(n, 10)
}

type exprParser struct {
	c     *RunContext
	input string
	pos   int
}

func (p *exprParser) parseExpr() (int64, error) {
	return p.parseAddSub()
}

func (p *exprParser) parseAddSub() (int64, error) {
	left, err := p.parseMulDiv()
	if err != nil {
		return 0, err
	}
	for {
		p.skipSpace()
		switch {
		case p.match('+'):
			right, err := p.parseMulDiv()
			if err != nil {
				return 0, err
			}
			left += right
		case p.match('-'):
			right, err := p.parseMulDiv()
			if err != nil {
				return 0, err
			}
			left -= right
		default:
			return left, nil
		}
	}
}

func (p *exprParser) parseMulDiv() (int64, error) {
	left, err := p.parseUnary()
	if err != nil {
		return 0, err
	}
	for {
		p.skipSpace()
		switch {
		case p.match('*'):
			right, err := p.parseUnary()
			if err != nil {
				return 0, err
			}
			left *= right
		case p.match('/'):
			right, err := p.parseUnary()
			if err != nil {
				return 0, err
			}
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			left /= right
		default:
			return left, nil
		}
	}
}

func (p *exprParser) parseUnary() (int64, error) {
	p.skipSpace()
	if p.match('-') {
		v, err := p.parseUnary()
		if err != nil {
			return 0, err
		}
		return -v, nil
	}
	if p.match('+') {
		return p.parseUnary()
	}
	return p.parsePrimary()
}

func (p *exprParser) parsePrimary() (int64, error) {
	p.skipSpace()
	if p.pos >= len(p.input) {
		return 0, fmt.Errorf("unexpected end of expression")
	}
	if p.match('(') {
		v, err := p.parseExpr()
		if err != nil {
			return 0, err
		}
		p.skipSpace()
		if !p.match(')') {
			return 0, fmt.Errorf("expected )")
		}
		return v, nil
	}
	if unicode.IsDigit(rune(p.input[p.pos])) {
		return p.parseNumberLiteral()
	}
	if unicode.IsLetter(rune(p.input[p.pos])) || p.input[p.pos] == '_' {
		return p.parseIdentifierValue()
	}
	return 0, fmt.Errorf("unexpected character %q", p.input[p.pos])
}

func (p *exprParser) parseNumberLiteral() (int64, error) {
	start := p.pos
	for p.pos < len(p.input) && unicode.IsDigit(rune(p.input[p.pos])) {
		p.pos++
	}
	n, err := strconv.ParseInt(p.input[start:p.pos], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q", p.input[start:p.pos])
	}
	return n, nil
}

func (p *exprParser) parseIdentifierValue() (int64, error) {
	start := p.pos
	for p.pos < len(p.input) {
		r := rune(p.input[p.pos])
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			p.pos++
			continue
		}
		break
	}
	name := p.input[start:p.pos]
	raw, err := p.c.resolvePlainPlaceholderKey(name)
	if err != nil {
		return 0, fmt.Errorf("variable %q: %w", name, err)
	}
	return parseNumericOperand(raw)
}

func parseNumericOperand(raw string) (int64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, fmt.Errorf("empty numeric value")
	}
	digits := make([]byte, 0, len(raw))
	negative := false
	for i, r := range raw {
		switch {
		case r == '-' && i == 0:
			negative = true
		case unicode.IsDigit(r):
			digits = append(digits, byte(r))
		case r == ' ' || r == '\u00a0' || r == '₽' || r == ',' || r == '.':
			continue
		default:
			return 0, fmt.Errorf("non-numeric value %q", raw)
		}
	}
	if len(digits) == 0 {
		return 0, fmt.Errorf("non-numeric value %q", raw)
	}
	n, err := strconv.ParseInt(string(digits), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("numeric overflow for %q", raw)
	}
	if negative {
		n = -n
	}
	return n, nil
}

func (p *exprParser) skipSpace() {
	for p.pos < len(p.input) && (p.input[p.pos] == ' ' || p.input[p.pos] == '\t') {
		p.pos++
	}
}

func (p *exprParser) match(ch byte) bool {
	p.skipSpace()
	if p.pos < len(p.input) && p.input[p.pos] == ch {
		p.pos++
		return true
	}
	return false
}
