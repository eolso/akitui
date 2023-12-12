package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eolso/akiapi"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Quit}}
}

type loadingModel struct {
	spinner spinner.Model
	help    help.Model
	err     error
}

func newLoadingModel() loadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return loadingModel{spinner: s, help: help.New()}
}

func (l loadingModel) Init() tea.Cmd {
	initClient := func() tea.Msg {
		akiapi.SetHttpClient(&http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}})
		game, err := akiapi.NewGame(akiapi.GameOptions{Theme: akiapi.CharactersTheme, Language: akiapi.English})
		if err != nil {
			panic(err)
		}

		items := []list.Item{
			item("Yes"),
			item("No"),
			item("Not Sure"),
			item("Probably"),
			item("Probably Not"),
		}

		l := list.New(items, newGameList(), 30, 10)
		l.Title = fmt.Sprintf("%d) %s", len(game.Responses())+1, game.Question())
		l.SetWidth(len(game.Question()) + 14)
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.Styles.Title = titleStyle
		l.Styles.PaginationStyle = paginationStyle
		l.Styles.HelpStyle = helpStyle
		l.SetShowHelp(false)

		return gameUpdate{
			game: game,
			list: l,
		}
	}

	return tea.Batch(l.spinner.Tick, initClient)
}

func (l loadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case gameUpdate:
		return &gameModel{gameUpdate: msg, spinner: l.spinner.Spinner, help: help.New()}, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return l, tea.Quit
		default:
			return l, nil
		}
	case error:
		l.err = msg
		return l, nil
	default:
		var cmd tea.Cmd
		l.spinner, cmd = l.spinner.Update(msg)
		return l, cmd
	}
}

func (l loadingModel) View() string {
	if l.err != nil {
		return l.err.Error()
	}

	str := fmt.Sprintf("\n  %s Loading game...\n\n", l.spinner.View())
	str += l.help.View(keyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q", "quit"),
		)},
	)

	return str
}
