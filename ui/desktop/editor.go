//go:build desktop

package desktop

import (
	"path/filepath"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type tabManager struct {
	tabs    *container.AppTabs
	paths   map[string]*container.TabItem
	editors map[string]*widget.Entry
}

func newTabManager() *tabManager {
	tabs := container.NewAppTabs()
	return &tabManager{
		tabs:    tabs,
		paths:   map[string]*container.TabItem{},
		editors: map[string]*widget.Entry{},
	}
}

func (m *tabManager) Widget() *container.AppTabs {
	return m.tabs
}

func (m *tabManager) Open(path, content string) {
	if item, ok := m.paths[path]; ok {
		m.tabs.Select(item)
		return
	}
	editor := widget.NewMultiLineEntry()
	editor.SetText(content)
	title := filepath.Base(path)
	item := container.NewTabItem(title, editor)
	m.tabs.Append(item)
	m.paths[path] = item
	m.editors[path] = editor
	m.tabs.Select(item)
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
		return editor.Text
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
