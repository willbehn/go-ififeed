package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	vp      viewport.Model
	content string
	ready   bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:

		if !m.ready {
			m.vp = viewport.New(msg.Width, msg.Height-2)
			m.vp.MouseWheelEnabled = true
			m.vp.SetContent(m.content)
			m.ready = true
		} else {
			m.vp.Width = msg.Width
			m.vp.Height = msg.Height - 2
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if !m.ready {
		return "loading...\n"
	}
	header := "ifi feed 📚\n"
	footer := fmt.Sprintf("\nScroll: %.0f%% — press q to quit", m.vp.ScrollPercent()*100)
	return header + m.vp.View() + footer
}
