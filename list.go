package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle      = lipgloss.NewStyle().MarginLeft(2)
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle   = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// Responses don't need to be filtered, so create a generic type that implements list.Item
type item string

func (i item) FilterValue() string { return "" }

type gameList struct {
	itemStyle     lipgloss.Style
	selectedStyle lipgloss.Style
}

func newGameList() gameList {
	return gameList{
		itemStyle:     lipgloss.NewStyle().PaddingLeft(4),
		selectedStyle: lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("141")),
	}
}

func (d gameList) Height() int                             { return 1 }
func (d gameList) Spacing() int                            { return 0 }
func (d gameList) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d gameList) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)
	if index == m.Index() {
		fmt.Fprintf(w, d.selectedStyle.Render("> ", str))
	} else {
		fmt.Fprintf(w, d.itemStyle.Render(str))
	}
}
