package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/hacel/jfsh/internal/jellyfin"
)

// KeyMap defines keybindings. It satisfies to the help.KeyMap interface, which is used to render the menu.
type KeyMap struct {
	// Keybindings used when browsing the list.
	CursorUp      key.Binding
	CursorDown    key.Binding
	PageUp        key.Binding
	PageDown      key.Binding
	NextTab       key.Binding
	PrevTab       key.Binding
	GoToStart     key.Binding
	GoToEnd       key.Binding
	Search        key.Binding
	ClearSearch   key.Binding
	Filter        key.Binding
	ClearFilter   key.Binding
	Select        key.Binding
	Back          key.Binding
	ToggleWatched key.Binding
	Refresh       key.Binding

	// Keybindings used when searching.
	CancelWhileSearching key.Binding
	AcceptWhileSearching key.Binding

	// Keybindings used when filtering.
	CancelWhileFiltering key.Binding
	AcceptWhileFiltering key.Binding

	// Help toggle keybindings.
	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding

	// The quit keybinding. This won't be caught when searching.
	Quit key.Binding

	// The quit-no-matter-what keybinding. This will be caught when searching.
	ForceQuit key.Binding
}

func defaultKeyMap() KeyMap {
	return KeyMap{
		// Browsing.
		CursorUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "b", "u"),
			key.WithHelp("pgup/b/u", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdn", "f", "d"),
			key.WithHelp("pgdn/f/d", "page down"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "prev tab"),
		),
		NextTab: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "next tab"),
		),
		GoToStart: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "go to start"),
		),
		GoToEnd: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "go to end"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		ClearSearch: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "clear"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		ClearFilter: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "clear"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter", "space"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "backspace"),
			key.WithHelp("esc", "back"),
			key.WithDisabled(),
		),
		ToggleWatched: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "toggle watched"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),

		// Searching.
		CancelWhileSearching: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		AcceptWhileSearching: key.NewBinding(
			key.WithKeys("enter", "tab", "shift+tab", "ctrl+k", "up", "ctrl+j", "down"),
			key.WithHelp("enter", "apply"),
		),

		// Filtering.
		CancelWhileFiltering: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		AcceptWhileFiltering: key.NewBinding(
			key.WithKeys("enter", "tab", "shift+tab", "ctrl+k", "up", "ctrl+j", "down"),
			key.WithHelp("enter", "apply"),
		),

		// Toggle help.
		ShowFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		CloseFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "close help"),
		),

		// Quitting.
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
	}
}

// FullHelp satisifies the help.KeyMap interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return append(
		[][]key.Binding{},
		[]key.Binding{
			k.CursorUp,
			k.CursorDown,
			k.PageUp,
			k.PageDown,
			k.GoToStart,
			k.GoToEnd,
		},
		[]key.Binding{
			k.NextTab,
			k.PrevTab,
			k.Refresh,
			k.Select,
			k.Search,
			k.ClearSearch,
			k.Filter,
			k.ClearFilter,
		},
		[]key.Binding{
			k.ToggleWatched,
			k.Back,
			k.Quit,
			k.CloseFullHelp,
		})
}

// ShortHelp satisifies the help.KeyMap interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Back,
		k.ToggleWatched,

		k.Search,
		k.ClearSearch,
		k.CancelWhileSearching,
		k.AcceptWhileSearching,

		k.Filter,
		k.ClearFilter,
		k.CancelWhileFiltering,
		k.AcceptWhileFiltering,

		k.ShowFullHelp,
		k.Quit,
	}
}

// updateKeys handles enabling and disabling of all keybinds based on UI state
func (m *model) updateKeys() {
	switch {
	case m.filterInput.Focused():
		m.keyMap.CursorUp.SetEnabled(false)
		m.keyMap.CursorDown.SetEnabled(false)
		m.keyMap.NextTab.SetEnabled(false)
		m.keyMap.PrevTab.SetEnabled(false)
		m.keyMap.GoToStart.SetEnabled(false)
		m.keyMap.GoToEnd.SetEnabled(false)
		m.keyMap.Search.SetEnabled(false)
		m.keyMap.ClearSearch.SetEnabled(false)
		m.keyMap.Filter.SetEnabled(false)
		m.keyMap.ClearFilter.SetEnabled(false)
		m.keyMap.Select.SetEnabled(false)
		m.keyMap.Back.SetEnabled(false)
		m.keyMap.ToggleWatched.SetEnabled(false)
		m.keyMap.Refresh.SetEnabled(false)
		m.keyMap.CancelWhileSearching.SetEnabled(false)
		m.keyMap.AcceptWhileSearching.SetEnabled(false)
		m.keyMap.CancelWhileFiltering.SetEnabled(true)
		m.keyMap.AcceptWhileFiltering.SetEnabled(true)
		m.keyMap.ShowFullHelp.SetEnabled(false)
		m.keyMap.CloseFullHelp.SetEnabled(false)
		m.keyMap.Quit.SetEnabled(false)
		m.keyMap.ForceQuit.SetEnabled(true)

	case m.playing != nil:
		m.keyMap.CursorUp.SetEnabled(false)
		m.keyMap.CursorDown.SetEnabled(false)
		m.keyMap.NextTab.SetEnabled(false)
		m.keyMap.PrevTab.SetEnabled(false)
		m.keyMap.GoToStart.SetEnabled(false)
		m.keyMap.GoToEnd.SetEnabled(false)
		m.keyMap.Search.SetEnabled(false)
		m.keyMap.ClearSearch.SetEnabled(false)
		m.keyMap.Filter.SetEnabled(false)
		m.keyMap.ClearFilter.SetEnabled(false)
		m.keyMap.Select.SetEnabled(false)
		m.keyMap.Back.SetEnabled(false)
		m.keyMap.ToggleWatched.SetEnabled(false)
		m.keyMap.Refresh.SetEnabled(false)
		m.keyMap.CancelWhileSearching.SetEnabled(false)
		m.keyMap.AcceptWhileSearching.SetEnabled(false)
		m.keyMap.CancelWhileFiltering.SetEnabled(false)
		m.keyMap.AcceptWhileFiltering.SetEnabled(false)
		m.keyMap.ShowFullHelp.SetEnabled(false)
		m.keyMap.CloseFullHelp.SetEnabled(false)
		m.keyMap.Quit.SetEnabled(false)
		m.keyMap.ForceQuit.SetEnabled(false)

	case m.currentSeries != nil:
		m.keyMap.CursorUp.SetEnabled(true)
		m.keyMap.CursorDown.SetEnabled(true)
		m.keyMap.NextTab.SetEnabled(false)
		m.keyMap.PrevTab.SetEnabled(false)
		m.keyMap.GoToStart.SetEnabled(true)
		m.keyMap.GoToEnd.SetEnabled(true)
		m.keyMap.Search.SetEnabled(false)
		m.keyMap.ClearSearch.SetEnabled(false)
		m.keyMap.Filter.SetEnabled(true)
		m.keyMap.ClearFilter.SetEnabled(m.filterActive)
		m.keyMap.Select.SetEnabled(len(m.items) > 0 && m.currentItem < len(m.items))
		m.keyMap.Back.SetEnabled(true)
		m.keyMap.ToggleWatched.SetEnabled(len(m.items) > 0 && m.currentItem < len(m.items) && !jellyfin.IsSeries(m.items[m.currentItem]))
		m.keyMap.Refresh.SetEnabled(true)
		m.keyMap.CancelWhileSearching.SetEnabled(false)
		m.keyMap.AcceptWhileSearching.SetEnabled(false)
		m.keyMap.CancelWhileFiltering.SetEnabled(false)
		m.keyMap.AcceptWhileFiltering.SetEnabled(false)
		m.keyMap.ShowFullHelp.SetEnabled(!m.help.ShowAll)
		m.keyMap.CloseFullHelp.SetEnabled(m.help.ShowAll)
		m.keyMap.Quit.SetEnabled(true)
		m.keyMap.ForceQuit.SetEnabled(true)

	case m.currentTab != Search:
		m.keyMap.CursorUp.SetEnabled(true)
		m.keyMap.CursorDown.SetEnabled(true)
		m.keyMap.NextTab.SetEnabled(true)
		m.keyMap.PrevTab.SetEnabled(true)
		m.keyMap.GoToStart.SetEnabled(true)
		m.keyMap.GoToEnd.SetEnabled(true)
		m.keyMap.Search.SetEnabled(false)
		m.keyMap.ClearSearch.SetEnabled(false)
		m.keyMap.Filter.SetEnabled(true)
		m.keyMap.ClearFilter.SetEnabled(m.filterActive)
		m.keyMap.Select.SetEnabled(len(m.items) > 0 && m.currentItem < len(m.items))
		m.keyMap.Back.SetEnabled(false)
		m.keyMap.ToggleWatched.SetEnabled(len(m.items) > 0 && m.currentItem < len(m.items) && !jellyfin.IsSeries(m.items[m.currentItem]))
		m.keyMap.Refresh.SetEnabled(true)
		m.keyMap.CancelWhileSearching.SetEnabled(false)
		m.keyMap.AcceptWhileSearching.SetEnabled(false)
		m.keyMap.CancelWhileFiltering.SetEnabled(false)
		m.keyMap.AcceptWhileFiltering.SetEnabled(false)
		m.keyMap.ShowFullHelp.SetEnabled(!m.help.ShowAll)
		m.keyMap.CloseFullHelp.SetEnabled(m.help.ShowAll)
		m.keyMap.Quit.SetEnabled(true)
		m.keyMap.ForceQuit.SetEnabled(true)

	case m.currentTab == Search && !m.searchInput.Focused():
		m.keyMap.CursorUp.SetEnabled(true)
		m.keyMap.CursorDown.SetEnabled(true)
		m.keyMap.NextTab.SetEnabled(true)
		m.keyMap.PrevTab.SetEnabled(true)
		m.keyMap.GoToStart.SetEnabled(true)
		m.keyMap.GoToEnd.SetEnabled(true)
		m.keyMap.Search.SetEnabled(true)
		m.keyMap.ClearSearch.SetEnabled(m.searchInput.Value() != "")
		m.keyMap.Filter.SetEnabled(false)
		m.keyMap.ClearFilter.SetEnabled(false)
		m.keyMap.Select.SetEnabled(len(m.items) > 0 && m.currentItem < len(m.items))
		m.keyMap.Back.SetEnabled(false)
		m.keyMap.ToggleWatched.SetEnabled(len(m.items) > 0 && m.currentItem < len(m.items) && !jellyfin.IsSeries(m.items[m.currentItem]))
		m.keyMap.Refresh.SetEnabled(true)
		m.keyMap.CancelWhileSearching.SetEnabled(false)
		m.keyMap.AcceptWhileSearching.SetEnabled(false)
		m.keyMap.CancelWhileFiltering.SetEnabled(false)
		m.keyMap.AcceptWhileFiltering.SetEnabled(false)
		m.keyMap.ShowFullHelp.SetEnabled(!m.help.ShowAll)
		m.keyMap.CloseFullHelp.SetEnabled(m.help.ShowAll)
		m.keyMap.Quit.SetEnabled(true)
		m.keyMap.ForceQuit.SetEnabled(true)

	case m.currentTab == Search && m.searchInput.Focused():
		m.keyMap.CursorUp.SetEnabled(false)
		m.keyMap.CursorDown.SetEnabled(false)
		m.keyMap.NextTab.SetEnabled(false)
		m.keyMap.PrevTab.SetEnabled(false)
		m.keyMap.GoToStart.SetEnabled(false)
		m.keyMap.GoToEnd.SetEnabled(false)
		m.keyMap.Search.SetEnabled(false)
		m.keyMap.ClearSearch.SetEnabled(false)
		m.keyMap.Filter.SetEnabled(false)
		m.keyMap.ClearFilter.SetEnabled(false)
		m.keyMap.Select.SetEnabled(false)
		m.keyMap.Back.SetEnabled(false)
		m.keyMap.ToggleWatched.SetEnabled(false)
		m.keyMap.Refresh.SetEnabled(false)
		m.keyMap.CancelWhileSearching.SetEnabled(true)
		m.keyMap.AcceptWhileSearching.SetEnabled(true)
		m.keyMap.CancelWhileFiltering.SetEnabled(false)
		m.keyMap.AcceptWhileFiltering.SetEnabled(false)
		m.keyMap.ShowFullHelp.SetEnabled(false)
		m.keyMap.CloseFullHelp.SetEnabled(false)
		m.keyMap.Quit.SetEnabled(false)
		m.keyMap.ForceQuit.SetEnabled(true)
	}
}
