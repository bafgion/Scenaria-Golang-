package stepcatalog

import "strings"

type Entry struct {
	Category string
	Template string
	Help     string
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
		if containsFold(entry.Template, query) || containsFold(entry.Help, query) || containsFold(entry.Category, query) {
			out = append(out, entry)
		}
	}
	return out
}

func stringsToLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func containsFold(haystack, needle string) bool {
	return strings.Contains(stringsToLower(haystack), needle)
}
