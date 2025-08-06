package config

import "github.com/charmbracelet/lipgloss"

var (
	pinkColor       = lipgloss.Color("#923FAD")
	brightPinkColor = lipgloss.Color("#B266D4")
	textColor       = lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#ddd"}
	dimTextColor    = lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777"}

	titleStyle = lipgloss.NewStyle().Margin(1, 0, 1, 1).Padding(0, 2).Background(pinkColor)
	labelStyle = lipgloss.NewStyle().Margin(0, 1, 0, 1).Foreground(brightPinkColor)
	inputStyle = lipgloss.NewStyle().Foreground(textColor)
	errStyle   = lipgloss.NewStyle().Margin(0, 0, 0, 1).Foreground(dimTextColor)
)

func (m model) View() string {
	sections := make([]string, 0, 2+len(m.inputs))

	{
		title := titleStyle.Render("jfsh")
		sections = append(sections, title)
	}

	{
		label := labelStyle.Render("Host")
		input := inputStyle.Render(m.inputs[hostInput].View())
		err := ""
		if e := m.inputs[hostInput].Err; e != nil {
			err = e.Error()
		}
		err = errStyle.Render(err)
		sections = append(sections, lipgloss.JoinHorizontal(lipgloss.Top, label, input), err)
	}

	{
		label := labelStyle.Render("Username")
		input := inputStyle.Render(m.inputs[usernameInput].View())
		err := ""
		if e := m.inputs[usernameInput].Err; e != nil {
			err = e.Error()
		}
		err = errStyle.Render(err)
		sections = append(sections, lipgloss.JoinHorizontal(lipgloss.Top, label, input), err)
	}

	{
		label := labelStyle.Render("Password")
		input := inputStyle.Render(m.inputs[passwordInput].View())
		err := ""
		if e := m.inputs[passwordInput].Err; e != nil {
			err = e.Error()
		}
		err = errStyle.Render(err)
		sections = append(sections, lipgloss.JoinHorizontal(lipgloss.Top, label, input), err)
	}

	{
		if m.err != nil {
			err := errStyle.Render(m.err.Error())
			sections = append(sections, err)
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	content = lipgloss.NewStyle().Width(m.width/2 + 10).Height(m.height / 2).Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
