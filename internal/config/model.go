package config

import (
	"errors"
	"net/url"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hacel/jfsh/internal/jellyfin"
	"github.com/spf13/viper"
)

// form fields
const (
	hostInput = iota
	usernameInput
	passwordInput
)

type model struct {
	client *jellyfin.Client
	err    error

	height int
	width  int

	inputs       []textinput.Model
	currentInput int
}

func initialModel() model {
	form := make([]textinput.Model, 3)

	form[hostInput] = textinput.New()
	form[hostInput].Focus()
	form[hostInput].Prompt = ""
	form[hostInput].SetValue(viper.GetString("host"))
	form[hostInput].Validate = func(s string) error {
		u, err := url.Parse(s)
		if err != nil {
			return errors.New("invalid format")
		}
		if u.Scheme == "" {
			return errors.New("must include scheme (http:// or https://)")
		}
		if u.Host == "" {
			return errors.New("URL must include host")
		}
		return nil
	}

	form[usernameInput] = textinput.New()
	form[usernameInput].Prompt = ""
	form[usernameInput].SetValue(viper.GetString("username"))

	form[passwordInput] = textinput.New()
	form[passwordInput].Prompt = ""
	form[passwordInput].EchoMode = textinput.EchoPassword
	form[passwordInput].SetValue(viper.GetString("password"))

	return model{inputs: form}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}
