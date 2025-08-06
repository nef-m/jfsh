package main

import (
	"math"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/hacel/jfsh/jellyfin"
)

var (
	blueColor       = lipgloss.Color("#000B25")
	pinkColor       = lipgloss.Color("#923FAD")
	brightPinkColor = lipgloss.Color("#B266D4")
	textColor       = lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#ddd"}
	dimTextColor    = lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777"}

	tabStyle        = lipgloss.NewStyle().Margin(1, 1, 1, 1).Padding(0, 2).Foreground(lipgloss.Color("#ddd")).Background(blueColor)
	currentTabStyle = tabStyle.Background(pinkColor)

	searchInputStyle = lipgloss.NewStyle().Margin(0, 0, 1, 2).Foreground(textColor)

	titleStyle = lipgloss.NewStyle().Margin(0, 0, 0, 1).Padding(0, 0, 0, 2).Foreground(textColor)
	descStyle  = titleStyle.Margin(0, 0, 1, 1).Foreground(dimTextColor)

	currentTitleStyle = lipgloss.NewStyle().
				Margin(0, 0, 0, 1).
				Padding(0, 0, 0, 1).
				Foreground(brightPinkColor).
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(brightPinkColor).
				Bold(true)
	currentDescStyle = currentTitleStyle.Margin(0, 0, 1, 1).Foreground(pinkColor).UnsetBold()

	scrollbarStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#1a1a1a"))
	scrollbarThumbStyle = lipgloss.NewStyle().Foreground(pinkColor)

	errStyle = lipgloss.NewStyle().Foreground(brightPinkColor)
)

func (m model) View() string {
	if m.playing != nil {
		messageView := lipgloss.NewStyle().Foreground(textColor).Render("Playing")
		title := jellyfin.GetItemTitle(*m.playing)
		titleView := lipgloss.NewStyle().Foreground(pinkColor).Render(title)
		exitView := lipgloss.NewStyle().Foreground(dimTextColor).Render("\nExit mpv to return")
		v := lipgloss.NewStyle().Padding(1, 2).BorderForeground(brightPinkColor).Render(
			lipgloss.JoinVertical(lipgloss.Top, messageView, titleView, exitView),
		)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, v)
	}

	var sections []string
	availHeight := m.height

	{
		if m.err != nil {
			sections = append(sections, errStyle.Render(m.err.Error()))
			availHeight -= lipgloss.Height(errStyle.Render(m.err.Error()))
		}
	}

	{
		if m.currentSeries == nil {
			var tabs []string
			for i, name := range []string{ResumeTabName, NextUpTabName, RecentlyAddedTabName, SearchTabName} {
				if tab(i) == m.currentTab {
					tabs = append(tabs, currentTabStyle.Render(name))
					continue
				}
				tabs = append(tabs, tabStyle.Render(name))
			}
			v := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
			sections = append(sections, v)
			availHeight -= lipgloss.Height(v)

			if m.currentTab == Search {
				v := searchInputStyle.Render(m.searchInput.View())
				sections = append(sections, v)
				availHeight -= lipgloss.Height(v)
			}
		} else {
			title := jellyfin.GetItemTitle(*m.currentSeries)
			v := currentTabStyle.Render(title)
			sections = append(sections, v)
			availHeight -= lipgloss.Height(v)
		}
	}

	var helpView string
	{
		helpView = m.help.View(m)
		helpView = lipgloss.NewStyle().Margin(0, 0, 0, 2).Render(helpView)
		availHeight -= lipgloss.Height(helpView)
	}

	{
		if len(m.items) > 0 {
			itemsPerPage := max(availHeight/3, 1)
			firstItem := max(m.currentItem-itemsPerPage/2, 0)
			if firstItem > len(m.items)-itemsPerPage {
				firstItem = max(len(m.items)-itemsPerPage, 0)
			}
			lastItem := min(firstItem+itemsPerPage, len(m.items))
			var itemViews []string
			for i := firstItem; i < lastItem; i++ {
				item := m.items[i]
				title := jellyfin.GetItemTitle(item)
				desc := jellyfin.GetItemDescription(item)

				// Prevent text from exceeding list width
				textwidth := m.width - 6
				title = ansi.Truncate(title, textwidth, "…")
				desc = ansi.Truncate(desc, textwidth, "…")
				if i == m.currentItem {
					title = currentTitleStyle.Render(title)
					desc = currentDescStyle.Render(desc)
				} else {
					if jellyfin.Watched(item) {
						title = titleStyle.Foreground(dimTextColor).Render(title)
					} else {
						title = titleStyle.Render(title)
					}
					desc = descStyle.Render(desc)
				}
				itemViews = append(itemViews, title, desc)
			}
			listContent := lipgloss.JoinVertical(lipgloss.Left, itemViews...)
			listContent = lipgloss.NewStyle().Width(m.width - 2).Render(listContent)

			scrollbarLines := make([]string, availHeight)
			if len(m.items) > itemsPerPage {
				thumbPosition := int(math.Round(float64(m.currentItem) / float64(len(m.items)-1) * float64(availHeight-1)))
				for i := range availHeight {
					if i == thumbPosition {
						scrollbarLines[i] = scrollbarThumbStyle.Render("█")
					} else {
						scrollbarLines[i] = scrollbarStyle.Render("│")
					}
				}
			}
			scrollbarView := lipgloss.JoinVertical(lipgloss.Left, scrollbarLines...)
			sections = append(sections, lipgloss.JoinHorizontal(lipgloss.Top, listContent, scrollbarView))
		} else {
			sections = append(sections, descStyle.Height(availHeight-1).Render("No items."))
		}
	}

	sections = append(sections, helpView)
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
