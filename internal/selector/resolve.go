package selector

import (
	"regexp"
	"strings"

	playwright "github.com/mxschmitt/playwright-go"
)

var hasTextRE = regexp.MustCompile(`:has-text\("((?:[^"\\]|\\.)*)"\)`)

var containerTags = map[string]struct{}{
	"div": {}, "section": {}, "span": {}, "li": {}, "nav": {}, "ul": {},
}

func firstHasText(selector string) string {
	match := hasTextRE.FindStringSubmatch(selector)
	if match == nil {
		return ""
	}
	return strings.ReplaceAll(match[1], `\"`, `"`)
}

func rootTag(selector string) string {
	part := selector
	if idx := strings.IndexAny(part, ":.#"); idx >= 0 {
		part = part[:idx]
	}
	return strings.ToLower(strings.TrimSpace(part))
}

// HoverLocatorCandidates returns ordered hover targets for menu-style selectors.
func HoverLocatorCandidates(selector string) []string {
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return nil
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, 8)
	add := func(candidate string) {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			return
		}
		if _, ok := seen[candidate]; ok {
			return
		}
		seen[candidate] = struct{}{}
		out = append(out, candidate)
	}

	text := firstHasText(selector)
	root := rootTag(selector)
	if text != "" {
		add(`a:has-text("` + text + `")`)
		add(`button:has-text("` + text + `")`)
		add(`nav a:has-text("` + text + `")`)
		add(`nav button:has-text("` + text + `")`)
		add(`[role="menuitem"]:has-text("` + text + `")`)
	}
	if _, isContainer := containerTags[root]; isContainer || strings.Contains(selector, ":has-text(") {
		if text != "" {
			add(selector + ` >> a:has-text("` + text + `")`)
			add(selector + ` >> button:has-text("` + text + `")`)
		}
		add(selector + " >> a")
		add(selector + " >> button")
	}
	add(selector)
	return out
}

// ResolveHoverLocator picks the first visible hover target.
func ResolveHoverLocator(page playwright.Page, selector string) playwright.Locator {
	for _, candidate := range HoverLocatorCandidates(selector) {
		locator := page.Locator(candidate).First()
		count, err := locator.Count()
		if err != nil || count == 0 {
			continue
		}
		if err := locator.WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(1500),
		}); err == nil {
			return locator
		}
	}
	return page.Locator(selector).First()
}

// ResolveChainedLocator picks the innermost container with exactly one target.
func ResolveChainedLocator(page playwright.Page, selector string) playwright.Locator {
	if !strings.Contains(selector, " >> ") {
		return page.Locator(selector)
	}
	parts := strings.SplitN(selector, " >> ", 2)
	containerSel := strings.TrimSpace(parts[0])
	targetSel := strings.TrimSpace(parts[1])
	if containerSel == "" || targetSel == "" {
		return page.Locator(selector)
	}

	containers := page.Locator(containerSel)
	count, err := containers.Count()
	if err != nil || count == 0 {
		return page.Locator(selector)
	}

	var bestTarget playwright.Locator
	var bestArea float64
	hasBest := false

	for i := 0; i < count; i++ {
		container := containers.Nth(i)
		box, err := container.BoundingBox()
		if err != nil || box == nil {
			continue
		}
		target, ok := resolveChainedTarget(container, targetSel)
		if !ok {
			continue
		}
		area := box.Width * box.Height
		if !hasBest || area < bestArea {
			hasBest = true
			bestArea = area
			bestTarget = target
		}
	}
	if hasBest {
		return bestTarget
	}
	return page.Locator(selector)
}

func resolveChainedTarget(container playwright.Locator, targetSel string) (playwright.Locator, bool) {
	text := firstHasText(targetSel)
	if rootTag(targetSel) == "button" && longSelectorText(text) {
		if target, ok := firstVisibleLocatorBySelectors(container, []string{
			"a[href]",
			"a",
			`[role="link"]`,
		}); ok {
			return target, true
		}
	}

	targets := container.Locator(targetSel)
	targetCount, err := targets.Count()
	if err == nil && targetCount == 1 {
		return targets.First(), true
	}
	if err == nil && targetCount > 1 {
		if target, ok := firstVisibleLocator(targets, targetCount); ok {
			return target, true
		}
	}

	if text == "" {
		return nil, false
	}
	for _, candidate := range []string{
		rootTag(targetSel),
		"button",
		"a",
		`[role="button"]`,
		`[role="link"]`,
		`[role="menuitem"]`,
	} {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			continue
		}
		locator := container.Locator(candidate)
		count, err := locator.Count()
		if err != nil || count == 0 {
			continue
		}
		if target, ok := firstVisibleLocator(locator, count); ok {
			return target, true
		}
	}
	return nil, false
}

func longSelectorText(text string) bool {
	return len([]rune(strings.TrimSpace(text))) >= 24
}

func firstVisibleLocatorBySelectors(container playwright.Locator, selectors []string) (playwright.Locator, bool) {
	for _, selector := range selectors {
		selector = strings.TrimSpace(selector)
		if selector == "" {
			continue
		}
		locator := container.Locator(selector)
		count, err := locator.Count()
		if err != nil || count == 0 {
			continue
		}
		if target, ok := firstVisibleLocator(locator, count); ok {
			return target, true
		}
	}
	return nil, false
}

func firstVisibleLocator(locator playwright.Locator, count int) (playwright.Locator, bool) {
	for i := 0; i < count; i++ {
		item := locator.Nth(i)
		box, err := item.BoundingBox()
		if err != nil || box == nil {
			continue
		}
		if box.Width <= 0 || box.Height <= 0 {
			continue
		}
		return item, true
	}
	return nil, false
}

// IsChained reports whether selector uses Playwright chain syntax.
func IsChained(selector string) bool {
	return strings.Contains(selector, " >> ")
}
