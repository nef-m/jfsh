package main

import (
	"slices"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hacel/jfsh/internal/jellyfin"
	"github.com/hacel/jfsh/internal/mpv"
)

type playbackStopped struct{ error }

func (m *model) playItem() tea.Cmd {
	client := m.client
	item := m.items[m.currentItem]
	if jellyfin.IsEpisode(item) {
		return func() tea.Msg {
			// get all episodes of the series and find the index of selected episode
			items, err := client.GetEpisodes(item)
			if err != nil {
				return err
			}
			idx := slices.IndexFunc(items, func(i jellyfin.Item) bool {
				return item.GetId() == i.GetId()
			})
			idx = max(0, idx) // sanity check
			if err := mpv.Play(client, items, idx); err != nil {
				return playbackStopped{err}
			}
			return playbackStopped{nil}
		}
	}
	return func() tea.Msg {
		if err := mpv.Play(client, []jellyfin.Item{item}, 0); err != nil {
			return playbackStopped{err}
		}
		return playbackStopped{nil}
	}
}

type markAsWatchedResult struct{ error }

func (m *model) markAsWatched() tea.Cmd {
	client := m.client
	item := m.items[m.currentItem]
	return func() tea.Msg {
		if err := client.MarkAsWatched(item); err != nil {
			return markAsWatchedResult{err}
		}
		return markAsWatchedResult{nil}
	}
}

func (m *model) fetchItems() tea.Cmd {
	client := m.client
	if m.currentSeries != nil {
		return func() tea.Msg {
			items, err := client.GetEpisodes(*m.currentSeries)
			if err != nil {
				return err
			}
			return items
		}
	}
	switch m.currentTab {
	case Resume:
		return func() tea.Msg {
			items, err := client.GetResume()
			if err != nil {
				return err
			}
			return items
		}
	case NextUp:
		return func() tea.Msg {
			items, err := client.GetNextUp()
			if err != nil {
				return err
			}
			return items
		}
	case RecentlyAdded:
		return func() tea.Msg {
			items, err := client.GetRecentlyAdded()
			if err != nil {
				return err
			}
			return items
		}
	case Search:
		query := m.searchInput.Value()
		return func() tea.Msg {
			if query == "" {
				return []jellyfin.Item{}
			}
			items, err := client.Search(query)
			if err != nil {
				return err
			}
			return items
		}
	default:
		panic("oops, selected tab is not in switch statement")
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case playbackStopped:
		if msg.error != nil {
			m.err = msg.error
		}
		m.playing = nil
		m.updateKeys()
		return m, m.fetchItems()

	case markAsWatchedResult:
		if msg.error != nil {
			m.err = msg.error
		}
		return m, m.fetchItems()

	case error:
		m.err = msg
		return m, nil

	case []jellyfin.Item:
		m.currentItem = 0
		m.items = msg
		m.updateKeys()
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.ForceQuit) {
			return m, tea.Quit
		}

		if m.searchInput.Focused() {
			switch {
			case key.Matches(msg, m.keyMap.CancelWhileSearching):
				m.searchInput.Blur()
				m.updateKeys()
				return m, nil
			case key.Matches(msg, m.keyMap.AcceptWhileSearching):
				m.searchInput.Blur()
				m.updateKeys()
				return m, m.fetchItems()
			}
			var cmd tea.Cmd
			m.searchInput, cmd = m.searchInput.Update(msg)
			return m, cmd
		}

		switch {
		case key.Matches(msg, m.keyMap.CursorUp):
			if m.currentItem > 0 {
				m.currentItem--
			}
			m.updateKeys()
			return m, nil
		case key.Matches(msg, m.keyMap.CursorDown):
			if m.currentItem < len(m.items)-1 {
				m.currentItem++
			}
			m.updateKeys()
			return m, nil
		case key.Matches(msg, m.keyMap.GoToEnd):
			m.currentItem = len(m.items) - 1
			m.updateKeys()
			return m, nil
		case key.Matches(msg, m.keyMap.GoToStart):
			m.currentItem = 0
			m.updateKeys()
			return m, nil

		case key.Matches(msg, m.keyMap.NextTab):
			if m.currentTab < Search {
				m.currentTab++
			} else {
				m.currentTab = 0
			}
			m.updateKeys()
			return m, m.fetchItems()
		case key.Matches(msg, m.keyMap.PrevTab):
			if m.currentTab > 0 {
				m.currentTab--
			} else {
				m.currentTab = Search
			}
			m.updateKeys()
			return m, m.fetchItems()

		case key.Matches(msg, m.keyMap.Search):
			m.searchInput.Focus()
			m.updateKeys()
			return m, nil
		case key.Matches(msg, m.keyMap.ClearSearch):
			m.searchInput.SetValue("")
			m.updateKeys()
			return m, m.fetchItems()

		case key.Matches(msg, m.keyMap.Select):
			item := m.items[m.currentItem]
			if jellyfin.IsSeries(item) {
				m.currentSeries = &item
				m.updateKeys()
				return m, m.fetchItems()
			}
			m.playing = &item
			m.updateKeys()
			return m, m.playItem()

		case key.Matches(msg, m.keyMap.Back):
			m.currentSeries = nil
			m.updateKeys()
			return m, m.fetchItems()

		case key.Matches(msg, m.keyMap.ShowFullHelp):
			m.help.ShowAll = !m.help.ShowAll
			m.updateKeys()
			return m, nil

		case key.Matches(msg, m.keyMap.CloseFullHelp):
			m.help.ShowAll = !m.help.ShowAll
			m.updateKeys()
			return m, nil

		case key.Matches(msg, m.keyMap.MarkAsWatched):
			return m, m.markAsWatched()

		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		default:
			return m, nil
		}
	default:
		return m, nil
	}
}
