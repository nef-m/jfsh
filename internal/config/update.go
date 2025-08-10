package config

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hacel/jfsh/internal/jellyfin"
	"github.com/spf13/viper"
)

func (m *model) initClient() tea.Cmd {
	host, username, password := m.inputs[hostInput].Value(), m.inputs[usernameInput].Value(), m.inputs[passwordInput].Value()
	device, deviceID, clientVersion, token, userID := viper.GetString("device"), viper.GetString("device_id"), viper.GetString("client_version"), viper.GetString("token"), viper.GetString("user_id")
	return func() tea.Msg {
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
		if err != nil {
			return err
		}
		return client
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		m.err = msg
		return m, nil

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		return m, nil

	case *jellyfin.Client:
		host, username, password := m.inputs[hostInput].Value(), m.inputs[usernameInput].Value(), m.inputs[passwordInput].Value()
		viper.Set("host", host)
		viper.Set("username", username)
		viper.Set("password", password)
		viper.Set("user_id", msg.UserID)
		viper.Set("token", msg.Token)
		if err := viper.WriteConfig(); err != nil {
			if err := viper.SafeWriteConfig(); err != nil {
				panic(err)
			}
		}
		m.client = msg
		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyEnter:
			if m.currentInput == len(m.inputs)-1 {
				valid := true
				if m.inputs[hostInput].Err != nil || m.inputs[hostInput].Value() == "" {
					valid = false
				}
				if m.inputs[usernameInput].Err != nil || m.inputs[usernameInput].Value() == "" {
					valid = false
				}
				if m.inputs[passwordInput].Err != nil {
					valid = false
				}
				if valid {
					return m, m.initClient()
				}
			}
			m.currentInput = (m.currentInput + 1) % len(m.inputs)

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyShiftTab, tea.KeyCtrlP, tea.KeyUp:
			m.currentInput--
			if m.currentInput < 0 {
				m.currentInput = len(m.inputs) - 1
			}

		case tea.KeyTab, tea.KeyCtrlN, tea.KeyDown:
			m.currentInput = (m.currentInput + 1) % len(m.inputs)
		}

		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.currentInput].Focus()
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}
