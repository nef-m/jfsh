package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hacel/jfsh/config"

	"github.com/spf13/pflag"
)

const (
	clientName    = "jfsh"
	clientVersion = "0.1.0"
)

func main() {
	cfgPath := pflag.StringP("config", "c", "", "override path to configuration file")
	pflag.Parse()

	// first off, run a side bubbletea model that takes care of configuration and initializing the api client
	client := config.Run(clientName, clientVersion, *cfgPath)
	if client == nil {
		// err handling should happen inside the config model, this means the user quit
		return
	}

	// now we can run the main bubbletea model
	p := tea.NewProgram(initialModel(client), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
