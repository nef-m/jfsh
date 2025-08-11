package main

import "github.com/charmbracelet/bubbles/key"

// FullHelp satisifies the help.KeyMap interface.
func (m model) FullHelp() [][]key.Binding {
	return append(
		[][]key.Binding{},
		[]key.Binding{
			m.keyMap.CursorUp,
			m.keyMap.CursorDown,
			m.keyMap.NextTab,
			m.keyMap.PrevTab,
			m.keyMap.GoToStart,
			m.keyMap.GoToEnd,
		},
		[]key.Binding{
			m.keyMap.Select,
			m.keyMap.Search,
			m.keyMap.ClearSearch,
		},
		[]key.Binding{
			m.keyMap.ToggleWatched,
			m.keyMap.Back,
			m.keyMap.Quit,
			m.keyMap.CloseFullHelp,
		})
}

// ShortHelp satisifies the help.KeyMap interface.
func (m model) ShortHelp() []key.Binding {
	return []key.Binding{
		m.keyMap.Search,
		m.keyMap.ClearSearch,
		m.keyMap.Back,
		m.keyMap.ToggleWatched,

		m.keyMap.CancelWhileSearching,
		m.keyMap.AcceptWhileSearching,

		m.keyMap.ShowFullHelp,
		m.keyMap.Quit,
	}
}
