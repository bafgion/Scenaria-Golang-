package player

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	playwright "github.com/mxschmitt/playwright-go"
)

const MaxLoopIterations = 100

// EmailCodePrompt is set by desktop/CLI for interactive OTP entry.
var EmailCodePrompt func(email string) (string, error)

var placeholderRE = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_:]+)\s*\}\}`)

type RunContext struct {
	Variables      map[string]string
	values         map[string]string
	rng            *rand.Rand
	person         *personBundle
	page           playwright.Page
	projectRoot    string
	lastDownload   string
	downloadDir    string
	completedSteps []gherkin.Step
	PromptEmailCode func(email string) (string, error)
}

func NewRunContext(variables map[string]string, seed int64, projectRoot string) *RunContext {
	if variables == nil {
		variables = map[string]string{}
	}
	return &RunContext{
		Variables:       variables,
		values:          map[string]string{},
		rng:             rand.New(rand.NewSource(seed)),
		projectRoot:     projectRoot,
		PromptEmailCode: EmailCodePrompt,
	}
}

func (c *RunContext) SetPage(page playwright.Page) {
	c.page = page
}

func (c *RunContext) Page() playwright.Page {
	return c.page
}

func (c *RunContext) Remember(name, value string) {
	c.Variables[name] = value
}

func (c *RunContext) SetLastDownload(path string) {
	c.lastDownload = path
}

func (c *RunContext) LastDownload() string {
	return c.lastDownload
}

func (c *RunContext) DownloadDir() string {
	if c.downloadDir != "" {
		return c.downloadDir
	}
	root := strings.TrimSpace(c.projectRoot)
	if root == "" {
		root = os.TempDir()
	}
	dir := filepath.Join(root, ".scenaria", "downloads")
	_ = os.MkdirAll(dir, 0o755)
	c.downloadDir = dir
	return dir
}

func (c *RunContext) ResolveText(text string) (string, error) {
	if text == "" || !strings.Contains(text, "{{") {
		return text, nil
	}
	var err error
	out := placeholderRE.ReplaceAllStringFunc(text, func(match string) string {
		if err != nil {
			return match
		}
		groups := placeholderRE.FindStringSubmatch(match)
		key := groups[1]
		if strings.HasPrefix(strings.ToLower(key), "env:") {
			envName := key[4:]
			value := os.Getenv(envName)
			if value == "" {
				err = fmt.Errorf("environment variable %q is empty", envName)
				return match
			}
			return value
		}
		if value, ok := c.Variables[key]; ok {
			return value
		}
		if value, ok := c.values[key]; ok {
			return value
		}
		generated, genErr := c.generate(key)
		if genErr != nil {
			err = genErr
			return match
		}
		c.values[key] = generated
		return generated
	})
	if err != nil {
		return "", err
	}
	return out, nil
}

func (c *RunContext) GenerateByKind(kind string) (string, error) {
	canonical, ok := NormalizeGeneratorName(kind)
	if !ok {
		return "", fmt.Errorf("unknown generator %q", kind)
	}
	return c.generateCanonical(canonical)
}

func (c *RunContext) EmailCode() (string, error) {
	if code := strings.TrimSpace(os.Getenv("SCENARIA_EMAIL_CODE")); code != "" {
		return code, nil
	}
	if code := strings.TrimSpace(c.Variables["email_code"]); code != "" {
		return code, nil
	}
	if c.PromptEmailCode != nil {
		email, _ := c.ResolveEmailForCode("", c.PriorSteps())
		return c.PromptEmailCode(email)
	}
	return "", fmt.Errorf("email verification code not set (use SCENARIA_EMAIL_CODE env, --var email_code=..., or interactive prompt)")
}

func (c *RunContext) generate(key string) (string, error) {
	canonical, ok := NormalizeGeneratorName(key)
	if !ok {
		return "", fmt.Errorf("unknown generator %q", key)
	}
	return c.generateCanonical(canonical)
}

func randomDigits(rng *rand.Rand, count int) string {
	if count <= 0 {
		return ""
	}
	out := make([]byte, count)
	for i := range out {
		out[i] = byte('0' + rng.Intn(10))
	}
	return string(out)
}

func (c *RunContext) EvaluateCondition(cond *gherkin.Condition) bool {
	if cond == nil || c.page == nil {
		return false
	}
	switch cond.Type {
	case "visible":
		selector, err := c.ResolveText(cond.Selector)
		if err != nil {
			return false
		}
		visible, err := c.page.Locator(selector).IsVisible()
		return err == nil && visible
	case "hidden":
		selector, err := c.ResolveText(cond.Selector)
		if err != nil {
			return false
		}
		visible, err := c.page.Locator(selector).IsVisible()
		return err == nil && !visible
	case "url_contains":
		value, err := c.ResolveText(cond.Value)
		if err != nil {
			return false
		}
		return strings.Contains(c.page.URL(), value)
	case "page_text":
		value, err := c.ResolveText(cond.Value)
		if err != nil {
			return false
		}
		content, err := c.page.Content()
		return err == nil && strings.Contains(content, value)
	default:
		return false
	}
}

func (c *RunContext) RecordStep(step gherkin.Step) {
	c.completedSteps = append(c.completedSteps, step)
}

func (c *RunContext) PriorSteps() []gherkin.Step {
	return c.completedSteps
}

func (c *RunContext) NowSeed() int64 {
	if c.rng == nil {
		return time.Now().UnixNano()
	}
	return c.rng.Int63()
}
