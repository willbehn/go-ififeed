package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/willbehn/go-ifi-feed/feed"
)

type model struct {
	vp      viewport.Model
	content string
	ready   bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:

		if !m.ready {
			m.vp = viewport.New(msg.Width, msg.Height-2) // leave room for header/footer
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

func (m model) View() string {
	if !m.ready {
		return "loading...\n"
	}
	header := "Scrollable output (viewport)\n"
	footer := fmt.Sprintf("\nScroll: %.0f%% — press q to quit", m.vp.ScrollPercent()*100)
	return header + m.vp.View() + footer
}

func main() {
	p := tea.NewProgram(model{content: strings.Join(feed.Fetch("IN1000"), "")},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion())

	p.Run()
}
