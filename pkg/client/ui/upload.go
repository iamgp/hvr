package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type uploadModel struct {
	name    string
	version string
	file    string
	step    int
	err     error
}

func (m uploadModel) Init() tea.Cmd {
	return nil
}

func (m uploadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.step < 3 {
				m.step++
			} else {
				return m, tea.Quit
			}
		case "backspace":
			if m.step == 0 && len(m.name) > 0 {
				m.name = m.name[:len(m.name)-1]
			} else if m.step == 1 && len(m.version) > 0 {
				m.version = m.version[:len(m.version)-1]
			} else if m.step == 2 && len(m.file) > 0 {
				m.file = m.file[:len(m.file)-1]
			}
		default:
			if m.step == 0 {
				m.name += msg.String()
			} else if m.step == 1 {
				m.version += msg.String()
			} else if m.step == 2 {
				m.file += msg.String()
			}
		}
	}
	return m, nil
}

func (m uploadModel) View() string {
	var b strings.Builder

	b.WriteString("Upload a library\n\n")

	switch m.step {
	case 0:
		b.WriteString(fmt.Sprintf("Enter library name: %s", m.name))
		b.WriteString("\n\nPress Enter to continue, Ctrl+C to quit")
	case 1:
		b.WriteString(fmt.Sprintf("Library name: %s\n", m.name))
		b.WriteString(fmt.Sprintf("Enter library version: %s", m.version))
		b.WriteString("\n\nPress Enter to continue, Ctrl+C to quit")
	case 2:
		b.WriteString(fmt.Sprintf("Library name: %s\n", m.name))
		b.WriteString(fmt.Sprintf("Library version: %s\n", m.version))
		b.WriteString(fmt.Sprintf("Enter file path: %s", m.file))
		b.WriteString("\n\nPress Enter to upload, Ctrl+C to quit")
	case 3:
		b.WriteString(fmt.Sprintf("Uploading %s version %s from file %s\n", m.name, m.version, m.file))
		b.WriteString("\nPress any key to exit")
	}

	return b.String()
}

func RunUploadTUI() (string, string, string) {
	m := uploadModel{}
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running upload TUI: %v\n", err)
		return "", "", ""
	}

	if finalModel, ok := finalModel.(uploadModel); ok {
		return finalModel.name, finalModel.version, finalModel.file
	}

	return "", "", ""
}
