package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type searchModel struct {
	query string
	step  int
	err   error
}

func (m searchModel) Init() tea.Cmd {
	return nil
}

func (m searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.step < 1 {
				m.step++
			} else {
				return m, tea.Quit
			}
		case "backspace":
			if m.step == 0 && len(m.query) > 0 {
				m.query = m.query[:len(m.query)-1]
			}
		default:
			if m.step == 0 {
				m.query += msg.String()
			}
		}
	}
	return m, nil
}

func (m searchModel) View() string {
	var b strings.Builder

	b.WriteString("Search for libraries\n\n")

	switch m.step {
	case 0:
		b.WriteString(fmt.Sprintf("Enter search query: %s", m.query))
		b.WriteString("\n\nPress Enter to search, Ctrl+C to quit")
	case 1:
		b.WriteString(fmt.Sprintf("Searching for: %s\n", m.query))
		b.WriteString("\nPress any key to exit")
	}

	return b.String()
}

func RunSearchTUI() string {
	m := searchModel{}
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running search TUI: %v\n", err)
		return ""
	}

	if finalModel, ok := finalModel.(searchModel); ok {
		return finalModel.query
	}

	return ""
}
