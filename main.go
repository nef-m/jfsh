package main

import (
	"io"
	"log"
	"log/slog"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hacel/jfsh/config"

	"github.com/adrg/xdg"
	"github.com/spf13/pflag"
)

const (
	clientVersion = "0.1.0"
)

func main() {
	cfgPath := pflag.StringP("config", "c", filepath.Join(xdg.ConfigHome, "jfsh", "jfsh.yaml"), "config file path")
	debug := pflag.StringP("debug", "d", "", "debug log file path (enables debug logging)")
	printVersion := pflag.BoolP("version", "v", false, "show version")
	help := pflag.BoolP("help", "h", false, "show help")
	pflag.Parse()

	if *help {
		println("Usage:  jfsh [OPTIONS]")
		println()
		println("Options:")
		pflag.PrintDefaults()
		return
	}

	if *printVersion {
		println(clientVersion)
		return
	}

	if *debug != "" {
		f, err := tea.LogToFile(*debug, "")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		slog.Info("enabled debug logging")
	} else {
		log.SetOutput(io.Discard)
	}

	// first off, run a side bubbletea model that takes care of configuration and initializing the api client
	client := config.Run(clientVersion, *cfgPath)
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
