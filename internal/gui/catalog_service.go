package gui

import (
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepcatalog"
)

// CatalogService owns step catalog search and completion helpers.
type CatalogService struct{}

func NewCatalogService() *CatalogService {
	return &CatalogService{}
}

func (s *CatalogService) Search(query string) []StepCatalogEntry {
	entries := stepcatalog.Search(query)
	out := make([]StepCatalogEntry, 0, len(entries))
	for _, entry := range entries {
		out = append(out, StepCatalogEntry{
			Label:       entry.Label,
			Action:      entry.Action,
			Category:    entry.Category,
			Description: entry.Description,
			Template:    entry.Template,
			Example:     entry.Example,
			Parameters:  entry.Parameters,
			Help:        entry.Help,
		})
	}
	return out
}

func (s *CatalogService) CompletionsForLine(line string, column int, language string) StepCompletionsDTO {
	started := time.Now()
	defer logWailsTiming("CompletionsForLine", started)
	lang := strings.TrimSpace(language)
	if lang == "" {
		lang = string(gherkin.LangRU)
	}
	result := stepcatalog.CompletionsForLineLang(line, column, lang)
	out := StepCompletionsDTO{
		Start: result.Start,
		End:   result.End,
		Items: make([]StepCompletionSnippet, 0, len(result.Items)),
	}
	for _, item := range result.Items {
		out.Items = append(out.Items, StepCompletionSnippet{
			Label:       item.Label,
			Insert:      item.Insert,
			Description: item.Description,
		})
	}
	return out
}
