package stepcatalog

import "strings"

func normalizeActionKind(kind string) string {
	return strings.ReplaceAll(kind, "-", "_")
}

// LookupByAction finds a catalog entry for a parsed stepdsl action kind.
func LookupByAction(kind string) (Entry, bool) {
	if kind == "" {
		return Entry{}, false
	}
	norm := normalizeActionKind(kind)
	for _, entry := range Entries() {
		if entry.Action == kind || normalizeActionKind(entry.Action) == norm {
			return entry, true
		}
	}
	return Entry{}, false
}

// LookupByStepText resolves catalog help for a Gherkin step body (without keyword).
func LookupByStepText(stepText string) (Entry, bool) {
	stepText = strings.TrimSpace(stepText)
	if stepText == "" {
		return Entry{}, false
	}

	if entry, ok := lookupByLabelPrefix(stepText); ok {
		return entry, true
	}

	entries := Search(stepText)
	if len(entries) > 0 {
		return entries[0], true
	}
	return Entry{}, false
}

func lookupByLabelPrefix(stepText string) (Entry, bool) {
	lower := stringsToLower(stepText)
	var best Entry
	bestLen := 0
	found := false
	for _, entry := range Entries() {
		lbl := stringsToLower(entry.Label)
		if lbl == "" || !strings.HasPrefix(lower, lbl) {
			continue
		}
		rest := lower[len(lbl):]
		if rest != "" && rest[0] != ' ' && rest[0] != '"' && rest[0] != '\'' {
			continue
		}
		if !found || len(lbl) < bestLen {
			best = entry
			bestLen = len(lbl)
			found = true
		}
	}
	if found {
		return best, true
	}
	return Entry{}, false
}
