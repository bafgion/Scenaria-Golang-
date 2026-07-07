package stepdsl

import "regexp"

// English step patterns — mirror core Russian DSL for # language: en features.
var englishStepPatterns = []struct {
	re   *regexp.Regexp
	kind string
	mapFn func(groups []string) Action
}{
	{regexp.MustCompile(`(?i)^i open "` + quoted + `"$`), "goto", one("goto")},
	{regexp.MustCompile(`(?i)^i (?:go to|navigate to|visit) "` + quoted + `"$`), "goto", one("goto")},
	{regexp.MustCompile(`(?i)^i see "` + quoted + `" is enabled$`), "assert-enabled", one("assert-enabled")},
	{regexp.MustCompile(`(?i)^i see "` + quoted + `" is disabled$`), "assert-disabled", one("assert-disabled")},
	{regexp.MustCompile(`(?i)^i see "` + quoted + `" is selected$`), "assert-selected", one("assert-selected")},
	{regexp.MustCompile(`(?i)^i see "` + quoted + `"$`), "assert-visible", one("assert-visible")},
	{regexp.MustCompile(`(?i)^i (?:do not|don't) see "` + quoted + `"$`), "assert-hidden", one("assert-hidden")},
	{regexp.MustCompile(`(?i)^i (?:click|press) "` + quoted + `"$`), "click", one("click")},
	{regexp.MustCompile(`(?i)^i double[- ]click "` + quoted + `"$`), "double-click", one("double-click")},
	{regexp.MustCompile(`(?i)^i hover (?:over |on )?"` + quoted + `"$`), "hover", one("hover")},
	{regexp.MustCompile(`(?i)^i (?:type|enter) "` + quoted + `" (?:into|in) "` + quoted + `"$`), "fill", two("fill")},
	{regexp.MustCompile(`(?i)^i check text "` + quoted + `" in "` + quoted + `"$`), "assert-text", two("assert-text")},
	{regexp.MustCompile(`(?i)^i check text matches regex "` + quoted + `" in "` + quoted + `"$`), "assert-text-regex", two("assert-text-regex")},
	{regexp.MustCompile(`(?i)^i check "` + quoted + `" contains "` + quoted + `"$`), "assert-var-contains", two("assert-var-contains")},
	{regexp.MustCompile(`(?i)^i check "` + quoted + `" equals "` + quoted + `"$`), "assert-var-equals", two("assert-var-equals")},
	{regexp.MustCompile(`(?i)^i check url "` + quoted + `"$`), "assert-url", one("assert-url")},
	{regexp.MustCompile(`(?i)^url contains "` + quoted + `"$`), "assert-url-contains", one("assert-url-contains")},
	{regexp.MustCompile(`(?i)^i wait (\d+) seconds?$`), "wait", waitSeconds},
	{regexp.MustCompile(`(?i)^i wait (\d+) ms$`), "wait", waitMillis},
	{regexp.MustCompile(`(?i)^i wait until "` + quoted + `" is enabled$`), "wait-enabled", one("wait-enabled")},
	{regexp.MustCompile(`(?i)^i wait until "` + quoted + `" is disabled$`), "wait-disabled", one("wait-disabled")},
	{regexp.MustCompile(`(?i)^i wait for "` + quoted + `"$`), "wait-visible", one("wait-visible")},
	{regexp.MustCompile(`(?i)^i wait until "` + quoted + `" disappears$`), "wait-hidden", one("wait-hidden")},
	{regexp.MustCompile(`(?i)^i reload the page$`), "reload", none("reload")},
	{regexp.MustCompile(`(?i)^i go back$`), "go-back", none("go-back")},
	{regexp.MustCompile(`(?i)^i close the browser$`), "close-browser", none("close-browser")},
	{regexp.MustCompile(`(?i)^i scroll to "` + quoted + `"$`), "scroll-to", one("scroll-to")},
	{regexp.MustCompile(`(?i)^i clear "` + quoted + `"$`), "clear", one("clear")},
	{regexp.MustCompile(`(?i)^i select "` + quoted + `" in "` + quoted + `"$`), "select", two("select")},
	{regexp.MustCompile(`(?i)^i check "` + quoted + `"$`), "check", one("check")},
	{regexp.MustCompile(`(?i)^i uncheck "` + quoted + `"$`), "uncheck", one("uncheck")},
	{regexp.MustCompile(`(?i)^i upload file "` + quoted + `" to "` + quoted + `"$`), "upload", two("upload")},
	{regexp.MustCompile(`(?i)^i press the "` + quoted + `" key$`), "press", one("press")},
	{regexp.MustCompile(`(?i)^i press "` + quoted + `" in "` + quoted + `"$`), "press-in", two("press-in")},
	{regexp.MustCompile(`(?i)^i drag "` + quoted + `" to "` + quoted + `"$`), "drag-drop", two("drag-drop")},
	{regexp.MustCompile(`(?i)^i remember text "` + quoted + `" as "` + quoted + `"$`), "remember-text", two("remember-text")},
	{regexp.MustCompile(`(?i)^i remember field "` + quoted + `" as "` + quoted + `"$`), "remember-field", two("remember-field")},
	{regexp.MustCompile(`(?i)^i remember number from "` + quoted + `" as "` + quoted + `"$`), "remember-number", two("remember-number")},
	{regexp.MustCompile(`(?i)^i remember url as "` + quoted + `"$`), "remember-url", one("remember-url")},
	{regexp.MustCompile(`(?i)^i download by clicking "` + quoted + `"$`), "download-click", one("download-click")},
}
