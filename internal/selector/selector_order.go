package selector

import (
	"encoding/json"
	"fmt"

	playwright "github.com/mxschmitt/playwright-go"
)

func SelectorOrderJS(click, input []string) string {
	click = NormalizeClickStrategies(click)
	input = NormalizeInputStrategies(input)
	clickJSON, _ := json.Marshal(click)
	inputJSON, _ := json.Marshal(input)
	return fmt.Sprintf(`() => {
		window.__scenariaSelectorOrder = { click: %s, input: %s };
	}`, clickJSON, inputJSON)
}

func ApplySelectorOrder(page playwright.Page, click, input []string) error {
	if page == nil {
		return nil
	}
	_, err := page.Evaluate(SelectorOrderJS(click, input))
	return err
}
