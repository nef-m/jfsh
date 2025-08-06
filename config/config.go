// Package config is a bubbletea model for setting up the configuration file and initializing the jellyfin client
package config

import (
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/hacel/jfsh/jellyfin"
	"github.com/spf13/viper"
)

func Run(clientVersion, path string) *jellyfin.Client {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// auto-create config dir
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		panic(err)
	}

	configModified := false

	// read config file
	if err := viper.ReadInConfig(); err != nil {
		configModified = true
	}

	// get/set client variables
	host := viper.GetString("host")
	username := viper.GetString("username")
	password := viper.GetString("password")
	device := viper.GetString("device")
	if device == "" {
		device, _ = os.Hostname()
		viper.Set("device", device)
		configModified = true
	}
	deviceID := viper.GetString("device_id")
	if deviceID == "" {
		deviceID = uuid.NewString()
		viper.Set("device_id", deviceID)
		configModified = true
	}
	if viper.GetString("client_version") != clientVersion {
		viper.Set("client_version", clientVersion)
		configModified = true
	}
	token := viper.GetString("token")
	userID := viper.GetString("user_id")

	if configModified {
		if err := viper.WriteConfig(); err != nil {
			if err := viper.SafeWriteConfig(); err != nil {
				panic(err)
			}
		}
	}

	// short circuit if we can already make a client
	if host != "" && username != "" && password != "" {
		client, err := jellyfin.NewClient(
			host,
			username,
			password,
			device,
			deviceID,
			clientVersion,
			token,
			userID,
		)
		if err == nil {
			return client
		}
	}

	// run the bubbletea form model otherwise
	m, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run()
	if err != nil {
		panic(err)
	}
	// the model should've created a valid client
	client := m.(model).client
	return client
}
