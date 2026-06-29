package stepcatalog

import "strings"

type Entry struct {
	Label       string
	Action      string
	Category    string
	Description string
	Template    string
	Example     string
	Parameters  []string
	Help        string // same as Description; kept for older clients
}

func Entries() []Entry {
	return loadEntries()
}

func Search(query string) []Entry {
	query = stringsToLower(query)
	if query == "" {
		return Entries()
	}
	out := make([]Entry, 0)
	for _, entry := range Entries() {
		if entryMatches(entry, query) {
			out = append(out, entry)
		}
	}
	return out
}

func entryMatches(entry Entry, query string) bool {
	parts := []string{
		entry.Label,
		entry.Action,
		entry.Category,
		entry.Description,
		entry.Template,
		entry.Example,
		entry.Help,
	}
	for _, p := range entry.Parameters {
		parts = append(parts, p)
	}
	return containsFold(strings.Join(parts, " "), query)
}

func stringsToLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func containsFold(haystack, needle string) bool {
	return strings.Contains(stringsToLower(haystack), needle)
}
