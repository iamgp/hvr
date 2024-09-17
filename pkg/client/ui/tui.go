package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type mainModel struct {
	choices  []string
	cursor   int
	selected string
}

func initialModel() mainModel {
	return mainModel{
		choices: []string{"Upload", "Download", "Search", "Quit"},
	}
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.selected = m.choices[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m mainModel) View() string {
	s := "Hamilton Venus Registry\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quit.\n"
	return s
}

func RunMainTUI() string {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running main TUI: %v\n", err)
		return ""
	}

	if m, ok := m.(mainModel); ok {
		return m.selected
	}
	return ""
}
