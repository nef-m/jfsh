package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hacel/jfsh/jellyfin"
)

type tab int

const (
	Resume tab = iota
	NextUp
	Latest
	Search
	ResumeTabName = "Resume"
	NextUpTabName = "Next Up"
	LatestTabName = "Latest"
	SearchTabName = "Search"
)

type model struct {
	keyMap KeyMap
	help   help.Model

	width  int
	height int

	client *jellyfin.Client

	currentTab  tab
	searchInput textinput.Model

	items       []jellyfin.Item
	currentItem int

	playing *jellyfin.Item

	err error
}

func initialModel(client *jellyfin.Client) model {
	searchInput := textinput.New()
	searchInput.Prompt = "Search: "

	return model{
		keyMap:      defaultKeyMap(),
		help:        help.New(),
		client:      client,
		searchInput: searchInput,
	}
}

func (m model) Init() tea.Cmd {
	return fetchItems(m.client, m.currentTab, m.searchInput.Value())
}
