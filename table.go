package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type tableView struct {
	table table.Model
	game  gameModel
	help  help.Model
}

//func newTableView(game gameModel) tableView {
//	guesses, err := game.session.ListGuesses()
//	if err != nil {
//		panic(err)
//	}
//
//	var longest int
//	var rows []table.Row
//	for _, guess := range guesses {
//		rows = append(rows, table.Row{guess.Name(), fmt.Sprintf("%.2f%%", guess.Probability())})
//		if len(guess.Name()) > longest {
//			longest = len(guess.Name())
//		}
//	}
//
//	columns := []table.Column{
//		{Title: "Name", Width: longest},
//		{Title: "Certainty", Width: 10},
//	}
//
//	t := table.New(
//		table.WithColumns(columns),
//		table.WithRows(rows),
//		table.WithFocused(true),
//		table.WithHeight(7),
//	)
//
//	s := table.DefaultStyles()
//	s.Header = s.Header.
//		BorderStyle(lipgloss.NormalBorder()).
//		BorderForeground(lipgloss.Color("240")).
//		BorderBottom(true).
//		Bold(false)
//	s.Selected = s.Selected.
//		Foreground(lipgloss.Color("229")).
//		Background(lipgloss.Color("57")).
//		Bold(false)
//	t.SetStyles(s)
//
//	return tableView{
//		table: t,
//		game:  game,
//		help:  help.New(),
//	}
//}

func (m tableView) Init() tea.Cmd { return nil }

func (m tableView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "g":
			return m.game, nil
		}
	}

	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m tableView) View() string {
	return baseStyle.Render(m.table.View()) + "\n" + m.help.View(m)
}

func (m tableView) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "return to session"),
		),
		key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (m tableView) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
