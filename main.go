package main

import (
	"fmt"
	"os"
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
	header := "ifi feed 📚\n"
	footer := fmt.Sprintf("\nScroll: %.0f%% — press q to quit", m.vp.ScrollPercent()*100)
	return header + m.vp.View() + footer
}

func runTui(courses []string) {
	p := tea.NewProgram(model{content: strings.Join(feed.Fetch(courses), "")},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion())

	p.Run()
}

func cmdSubscribe(courses []string) {
	for _, course := range courses {
		fmt.Print(course)
	}

}

func main() {

	args := os.Args

	if len(args) > 1 {
		if args[1] == "subscribe" {
			cmdSubscribe(args[2:])
		}

	} else {

		courses := []string{"IN1000", "IN1020"}

		runTui(courses)
	}

}
