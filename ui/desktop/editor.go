//go:build desktop

package desktop

import (
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2/container"
)

type tabManager struct {
	tabs         *container.AppTabs
	paths        map[string]*container.TabItem
	editors      map[string]*syntaxEditor
	onTextChange func(text string)
}

func newTabManager() *tabManager {
	tabs := container.NewAppTabs()
	return &tabManager{
		tabs:    tabs,
		paths:   map[string]*container.TabItem{},
		editors: map[string]*syntaxEditor{},
	}
}

func (m *tabManager) SetTextChangeHandler(fn func(text string)) {
	m.onTextChange = fn
}

func (m *tabManager) Widget() *container.AppTabs {
	return m.tabs
}

func (m *tabManager) Open(path, content string) {
	if item, ok := m.paths[path]; ok {
		m.tabs.Select(item)
		if m.onTextChange != nil {
			m.onTextChange(m.CurrentText())
		}
		return
	}
	editor := newSyntaxEditor()
	editor.SetText(content)
	editor.SetOnChanged(func(text string) {
		if m.onTextChange != nil {
			m.onTextChange(text)
		}
	})
	title := filepath.Base(path)
	item := container.NewTabItem(title, editor)
	m.tabs.Append(item)
	m.paths[path] = item
	m.editors[path] = editor
	m.tabs.Select(item)
	if m.onTextChange != nil {
		m.onTextChange(content)
	}
}

func (m *tabManager) InsertLine(text string) {
	path := m.CurrentPath()
	if path == "" {
		return
	}
	editor, ok := m.editors[path]
	if !ok {
		return
	}
	current := editor.Text()
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	if current != "" && !strings.HasSuffix(current, "\n") {
		current += "\n"
	}
	editor.SetText(current + text)
}

func (m *tabManager) CurrentPath() string {
	selected := m.tabs.Selected()
	if selected == nil {
		return ""
	}
	for path, item := range m.paths {
		if item == selected {
			return path
		}
	}
	return ""
}

func (m *tabManager) CurrentText() string {
	path := m.CurrentPath()
	if path == "" {
		return ""
	}
	if editor, ok := m.editors[path]; ok {
		return editor.Text()
	}
	return ""
}

func (m *tabManager) CloseCurrent() {
	path := m.CurrentPath()
	if path == "" {
		return
	}
	if item, ok := m.paths[path]; ok {
		m.tabs.Remove(item)
		delete(m.paths, path)
		delete(m.editors, path)
	}
}
