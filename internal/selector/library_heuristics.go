package selector

import (
	"encoding/json"
	"fmt"

	"github.com/bafgion/scenaria-golang/internal/settings"
	playwright "github.com/mxschmitt/playwright-go"
)

// Library pack DOM patterns (catalog for 6.1.1):
//   MUI click: .MuiButton-root, .MuiIconButton-root, .MuiMenuItem-root + data-testid
//   MUI input: .MuiInputBase-input, .MuiInputBase-root[name]
//   Ant click: .ant-btn, .ant-menu-item, .ant-select
//   Ant input: .ant-input[placeholder]
//   Portal warnings: MuiPopover/Menu/Modal, ant-dropdown/select-dropdown

func LibraryHeuristicsMUIEnabled(cfg *settings.AppSettings) bool {
	if cfg == nil || cfg.LibraryHeuristicsMUI == nil {
		return true
	}
	return *cfg.LibraryHeuristicsMUI
}

func LibraryHeuristicsAntEnabled(cfg *settings.AppSettings) bool {
	if cfg == nil || cfg.LibraryHeuristicsAnt == nil {
		return true
	}
	return *cfg.LibraryHeuristicsAnt
}

func LibraryHeuristicsJS(mui, ant bool) string {
	muiJSON, _ := json.Marshal(mui)
	antJSON, _ := json.Marshal(ant)
	return fmt.Sprintf(`() => {
		window.__scenariaLibraryHeuristics = { mui: %s, ant: %s };
	}`, muiJSON, antJSON)
}

func ApplyLibraryHeuristics(page playwright.Page, mui, ant bool) error {
	if page == nil {
		return nil
	}
	_, err := page.Evaluate(LibraryHeuristicsJS(mui, ant))
	return err
}

func ApplyLibraryHeuristicsFromSettings(page playwright.Page, cfg *settings.AppSettings) error {
	return ApplyLibraryHeuristics(page, LibraryHeuristicsMUIEnabled(cfg), LibraryHeuristicsAntEnabled(cfg))
}
