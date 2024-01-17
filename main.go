package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eolso/akiapi"
)

type gameState int

const (
	initializingState gameState = iota
	thinkingState
	questionPromptState
	answerPromptState
	undoState
)

type gameUpdate struct {
	game         akiapi.Game
	list         list.Model
	state        gameState
	answerBuffer int
	err          error
}

type gameModel struct {
	gameUpdate
	spinner spinner.Spinner
	help    help.Model
}

func (m gameModel) Init() tea.Cmd {
	return nil
}

func (m gameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case gameUpdate:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}

		m.gameUpdate = msg

		return m, nil
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.state == thinkingState {
				return m, nil
			}

			update := updateGame(m.state, m)
			m.state = thinkingState

			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)

			m.list.SetSpinner(m.spinner)

			return m, tea.Batch(cmd, m.list.StartSpinner(), update)
		case "u":
			if len(m.game.Responses()) > 0 {
				update := updateGame(undoState, m)
				m.state = thinkingState

				var cmd tea.Cmd
				m.list, cmd = m.list.Update(msg)
				m.list.SetSpinner(m.spinner)

				return m, tea.Batch(cmd, m.list.StartSpinner(), update)
			}
		case "g":
			if m.state == thinkingState {
				return m, nil
			}

			return newTableView(m), nil
		}

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m gameModel) View() string {
	if m.err != nil {
		return quitTextStyle.Render(m.err.Error())
	}

	str := "\n" + m.list.View()
	str += "\n" + lipgloss.NewStyle().Padding(0, 0, 1, 4).Render(m.help.View(m))

	return str
}

func (m gameModel) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "view guesses"),
		),
		key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "undo"),
		),
		key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (m gameModel) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

func main() {
	if _, err := tea.NewProgram(newLoadingModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func updateGame(state gameState, model gameModel) func() tea.Msg {
	return func() tea.Msg {
		update := gameUpdate{game: model.game, list: model.list, answerBuffer: model.answerBuffer}

		switch state {
		case questionPromptState, initializingState:
			update.err = model.game.SelectOption(update.list.Cursor())
			if update.err != nil {
				return update
			}

			if update.game.Progress() > 80.0 && update.answerBuffer <= 0 {
				guess, err := update.game.Guess()
				if err != nil {
					update.err = err
					return update
				}

				update.list.Title = fmt.Sprintf("You're thinking of: %s", guess.Name())
				update.list.SetWidth(len(update.list.Title) + 14)
				update.list.SetItems([]list.Item{item("Yes"), item("No")})
				update.state = answerPromptState

				return update
			}

			update.answerBuffer--
		case undoState:
			update.err = update.game.Undo()
			if update.answerBuffer == 5 {
				update.answerBuffer = -1
			}
		case answerPromptState:
			if update.list.Cursor() == 0 {
				return tea.Quit()
			} else {
				update.err = update.game.Undo()
				update.answerBuffer = 5
			}
		}

		update.list.Title = fmt.Sprintf("%d) %s", len(update.game.Responses())+1, update.game.Question())
		update.list.SetWidth(len(update.list.Title) + 14)
		update.list.SetItems([]list.Item{
			item("Yes"),
			item("No"),
			item("Not Sure"),
			item("Probably"),
			item("Probably Not"),
		})

		update.state = questionPromptState

		return update
	}
}
