package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type downloadModel struct {
	name    string
	version string
	step    int
	err     error
}

func (m downloadModel) Init() tea.Cmd {
	return nil
}

func (m downloadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.step < 2 {
				m.step++
			} else {
				return m, tea.Quit
			}
		case "backspace":
			if m.step == 0 && len(m.name) > 0 {
				m.name = m.name[:len(m.name)-1]
			} else if m.step == 1 && len(m.version) > 0 {
				m.version = m.version[:len(m.version)-1]
			}
		default:
			if m.step == 0 {
				m.name += msg.String()
			} else if m.step == 1 {
				m.version += msg.String()
			}
		}
	}
	return m, nil
}

func (m downloadModel) View() string {
	var b strings.Builder

	b.WriteString("Download a library\n\n")

	switch m.step {
	case 0:
		b.WriteString(fmt.Sprintf("Enter library name: %s", m.name))
		b.WriteString("\n\nPress Enter to continue, Ctrl+C to quit")
	case 1:
		b.WriteString(fmt.Sprintf("Library name: %s\n", m.name))
		b.WriteString(fmt.Sprintf("Enter library version: %s", m.version))
		b.WriteString("\n\nPress Enter to download, Ctrl+C to quit")
	case 2:
		b.WriteString(fmt.Sprintf("Downloading %s version %s\n", m.name, m.version))
		b.WriteString("\nPress any key to exit")
	}

	return b.String()
}

func RunDownloadTUI() (string, string) {
	m := downloadModel{}
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running download TUI: %v\n", err)
		return "", ""
	}

	if finalModel, ok := finalModel.(downloadModel); ok {
		return finalModel.name, finalModel.version
	}

	return "", ""
}
