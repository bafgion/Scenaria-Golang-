package stepcatalog

import "regexp"

// stepSnippetsEN mirrors core Russian snippets for # language: en features.
var stepSnippetsEN = []snippetDef{
	{"I open", `I open "https://site.com"`, "Open page (action: goto)"},
	{"I click", `I click "button.submit"`, "Click element (action: click)"},
	{"I double-click", `I double-click ".file-item"`, "Double click (action: double_click)"},
	{"I hover", `I hover on "nav a.services"`, "Hover before submenu click (action: hover)"},
	{"I type", `I type "text" into "input[name=email]"`, "Type into field (action: fill)"},
	{"I clear", `I clear "input#search"`, "Clear input (action: clear)"},
	{"I select", `I select "Value" in "select#country"`, "Select option (action: select)"},
	{"I check", `I check "input#agree"`, "Check checkbox (action: check)"},
	{"I uncheck", `I uncheck "input#newsletter"`, "Uncheck checkbox (action: uncheck)"},
	{"I press", `I press the "Enter" key`, "Key press (action: press)"},
	{"I press in", `I press "Tab" in "input[name=email]"`, "Key in field (action: press)"},
	{"I upload", `I upload file "C:\\data\\doc.pdf" to "input[type=file]"`, "Upload file (action: upload)"},
	{"I remember field", `I remember field "input#email" as "user_email"`, "Remember field value (action: remember_field)"},
	{"I see", `I see "h1.title"`, "Assert visible (action: assert_visible)"},
	{"I don't see", `I don't see ".modal-overlay"`, "Assert hidden (action: assert_hidden)"},
	{"I see enabled", `I see "button#submit" is enabled`, "Assert enabled (action: assert_enabled)"},
	{"I see disabled", `I see "button#submit" is disabled`, "Assert disabled (action: assert_disabled)"},
	{"I see selected", `I see ".delivery-tile.active" is selected`, "Assert selected tile (action: assert_selected)"},
	{"I check text", `I check text "Success" in ".message"`, "Assert text (action: assert_text)"},
	{"I check text regex", `I check text matches regex "Pay.*\\d+" in "[data-testid=pay]"`, "Assert text by regex (action: assert_text_regex)"},
	{"I remember number", `I remember number from ".order-total" as "total_amount"`, "Remember parsed number (action: remember_number)"},
	{"I check vars", `I check "{{actual}}" contains "{{expected}}"`, "Assert variable contains (action: assert_var_contains)"},
	{"I check url", `I check url "https://site.com/profile"`, "Assert URL (action: assert_url)"},
	{"I wait", "I wait 2 seconds", "Pause (action: wait)"},
	{"I wait ms", "I wait 500 ms", "Short pause in ms (action: wait)"},
	{"I wait for", `I wait for "button.ready"`, "Wait for element (action: wait_for)"},
	{"I wait until enabled", `I wait until "button#pay" is enabled`, "Wait until enabled (action: wait_for_enabled)"},
	{"I wait until disabled", `I wait until "button#pay" is disabled`, "Wait until disabled (action: wait_for_disabled)"},
	{"I wait until", `I wait until ".spinner" disappears`, "Wait until hidden (action: wait_for_hidden)"},
	{"I reload", "I reload the page", "Reload page (action: reload)"},
	{"I go back", "I go back", "Browser back (action: go_back)"},
	{"I close browser", "I close the browser", "Close browser (action: close_browser)"},
	{"I scroll to", `I scroll to "section#contacts"`, "Scroll to element (action: scroll_to)"},
	{"If I see", "If I see \".cookie-banner\"\n\tI click \"button.accept\"", "Conditional block (action: if)"},
	{"If enabled", "If \"button#pay\" is enabled\n\tI click \"button#pay\"", "Conditional block when enabled (action: if)"},
	{"Repeat", "Repeat 3 times\n\tI click \"button.add\"", "Fixed repeat loop (action: repeat)"},
	{"While", "While I see \"button.load-more\"\n\tI click \"button.load-more\"", "While loop (action: while)"},
	{"For each", "For each \".product-card\" as \"card\"\n\tI click \"{{card}} .buy-btn\"", "For-each loop (action: for_each)"},
}

var headerSnippetsEN = []CompletionSnippet{
	{Label: "Feature:", Insert: "Feature: UI scenario", Description: "Feature file title"},
	{Label: "Background:", Insert: "Background:\n\tGiven I connect TestClient \"ClientName\"", Description: "Background: named TestClient before scenarios"},
	{Label: "Scenario:", Insert: "Scenario: Scenario name", Description: "Scenario title"},
}

var completionKeywordsEN = []string{"Given", "When", "Then", "And", "But"}

var headerLineEN = regexp.MustCompile(`(?i)^\s*(feature|scenario|functionality)\s*:?`)
var stepKeywordEN = regexp.MustCompile(`(?i)^(?:(Given|When|Then|And|But)\s+)?(.*)$`)
