package gui

import "github.com/bafgion/scenaria-golang/internal/featurehighlight"

type HighlightSpan struct {
	Text string `json:"text"`
	Kind string `json:"kind"`
}

func (s *Service) HighlightFeature(text string) []HighlightSpan {
	return highlightFeature(text)
}

func highlightFeature(text string) []HighlightSpan {
	spans := featurehighlight.Highlight(text)
	out := make([]HighlightSpan, 0, len(spans))
	for _, span := range spans {
		out = append(out, HighlightSpan{
			Text: span.Text,
			Kind: highlightKindName(span.Kind),
		})
	}
	return out
}

func (s *Service) RefactorUpdateStartURLs(text, newURL string) RefactorResult {
	return UpdateStartURLs(text, newURL)
}

func (s *Service) RefactorNormalizeIndents(text string) string {
	return NormalizeStepIndents(text)
}

func (s *Service) RefactorCollapseBlankLines(text string) string {
	return CollapseBlankLinesBetweenSteps(text)
}

func (s *Service) RefactorFormatFeature(text string) string {
	return FormatFeature(text)
}

func (s *Service) RefactorReplaceInText(text, find, replace string, caseSensitive bool) RefactorResult {
	return ReplaceInText(text, find, replace, caseSensitive)
}

func highlightKindName(kind featurehighlight.Kind) string {
	switch kind {
	case featurehighlight.KindComment:
		return "comment"
	case featurehighlight.KindTag:
		return "tag"
	case featurehighlight.KindGherkinKeyword:
		return "gherkin"
	case featurehighlight.KindStepKeyword:
		return "step"
	case featurehighlight.KindBlockKeyword:
		return "block"
	case featurehighlight.KindString:
		return "string"
	case featurehighlight.KindTestClient:
		return "testclient"
	case featurehighlight.KindError:
		return "error"
	default:
		return "text"
	}
}
