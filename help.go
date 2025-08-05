package main

import "github.com/charmbracelet/bubbles/key"

// FullHelp satisifies the help.KeyMap interface.
func (m model) FullHelp() [][]key.Binding {
	kb := [][]key.Binding{{
		m.keyMap.CursorUp,
		m.keyMap.CursorDown,
		m.keyMap.NextTab,
		m.keyMap.PrevTab,
		m.keyMap.GoToStart,
		m.keyMap.GoToEnd,
	}}
	listLevelBindings := []key.Binding{
		m.keyMap.Search,
		m.keyMap.ClearSearch,
		m.keyMap.AcceptWhileSearching,
		m.keyMap.CancelWhileSearching,
	}
	return append(kb,
		listLevelBindings,
		[]key.Binding{
			m.keyMap.Quit,
			m.keyMap.CloseFullHelp,
		})
}

// ShortHelp satisifies the help.KeyMap interface.
func (m model) ShortHelp() []key.Binding {
	searching := m.searchInput.Focused()
	if searching {
		return []key.Binding{
			m.keyMap.CancelWhileSearching,
			m.keyMap.AcceptWhileSearching,
		}
	}
	kb := []key.Binding{
		m.keyMap.CursorUp,
		m.keyMap.CursorDown,
	}
	if m.currentTab == Search {
		kb = append(kb, m.keyMap.Search)
		kb = append(kb, m.keyMap.ClearSearch)
	}
	return append(kb, m.keyMap.Quit, m.keyMap.ShowFullHelp)
}
