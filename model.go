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
	RecentlyAdded
	Search
	ResumeTabName        = "Resume"
	NextUpTabName        = "Next Up"
	RecentlyAddedTabName = "Recently Added"
	SearchTabName        = "Search"
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

	currentSeries *jellyfin.Item

	playing *jellyfin.Item

	err error
}

func initialModel(client *jellyfin.Client) model {
	searchInput := textinput.New()
	searchInput.Prompt = "Search: "
	searchInput.Width = 40

	m := model{
		keyMap:      defaultKeyMap(),
		help:        help.New(),
		client:      client,
		searchInput: searchInput,
	}
	m.updateKeys()
	return m
}

func (m model) Init() tea.Cmd {
	return m.fetchItems()
}
