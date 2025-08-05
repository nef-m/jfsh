package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hacel/jfsh/jellyfin"
	"github.com/hacel/jfsh/mpv"
)

type playbackStopped struct{}

func playItem(client *jellyfin.Client, item jellyfin.Item) tea.Cmd {
	return func() tea.Msg {
		mpv.Play(client, item)
		return playbackStopped{}
	}
}

func fetchItems(client *jellyfin.Client, tab tab, searchQuery string) tea.Cmd {
	return func() tea.Msg {
		switch tab {
		case Resume:
			items, err := client.GetResume()
			if err != nil {
				return err
			}
			return items
		case NextUp:
			items, err := client.GetNextUp()
			if err != nil {
				return err
			}
			return items
		case Latest:
			items, err := client.GetLatest()
			if err != nil {
				return err
			}
			return items
		case Search:
			if searchQuery == "" {
				return []jellyfin.Item{}
			}
			items, err := client.Search(searchQuery)
			if err != nil {
				return err
			}
			return items
		default:
			panic("oops, selected tab is not in switch statement")
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []jellyfin.Item:
		m.items = msg
		return m, nil
	}

	switch msg := msg.(type) {
	case error:
		m.err = msg

	case []jellyfin.Item:
		m.items = msg
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case playbackStopped:
		m.playing = nil
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.ForceQuit) {
			return m, tea.Quit
		}

		if m.currentTab == Search {
			searching := m.searchInput.Focused()
			if searching {
				switch {
				case key.Matches(msg, m.keyMap.CancelWhileSearching):
					m.searchInput.Blur()
					return m, nil
				case key.Matches(msg, m.keyMap.AcceptWhileSearching):
					m.searchInput.Blur()
					return m, fetchItems(m.client, m.currentTab, m.searchInput.Value())
				}
				var cmd tea.Cmd
				m.searchInput, cmd = m.searchInput.Update(msg)
				return m, cmd
			}
			switch {
			case key.Matches(msg, m.keyMap.Search):
				m.searchInput.Focus()
				return m, nil
			case key.Matches(msg, m.keyMap.ClearSearch):
				m.searchInput.SetValue("")
				return m, fetchItems(m.client, m.currentTab, m.searchInput.Value())
			}
		}

		switch {
		case key.Matches(msg, m.keyMap.NextTab):
			if m.currentTab < Search {
				m.currentTab++
			} else {
				m.currentTab = 0
			}
			return m, fetchItems(m.client, m.currentTab, m.searchInput.Value())
		case key.Matches(msg, m.keyMap.PrevTab):
			if m.currentTab > 0 {
				m.currentTab--
			} else {
				m.currentTab = Search
			}
			return m, fetchItems(m.client, m.currentTab, m.searchInput.Value())
		case key.Matches(msg, m.keyMap.CursorUp):
			if m.currentItem > 0 {
				m.currentItem--
			}
			return m, nil
		case key.Matches(msg, m.keyMap.CursorDown):
			if m.currentItem < len(m.items)-1 {
				m.currentItem++
			}
			return m, nil
		case key.Matches(msg, m.keyMap.GoToEnd):
			m.currentItem = len(m.items) - 1
			return m, nil
		case key.Matches(msg, m.keyMap.GoToStart):
			m.currentItem = 0
			return m, nil

		case key.Matches(msg, m.keyMap.Select):
			if m.currentItem >= len(m.items) {
				panic("selected item is out of range")
			}
			item := m.items[m.currentItem]
			m.playing = &item
			return m, playItem(m.client, item)

		case key.Matches(msg, m.keyMap.ShowFullHelp):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

		case key.Matches(msg, m.keyMap.CloseFullHelp):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		}
	}
	return m, nil
}
