package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/willbehn/go-ifi-feed/feed"
	"github.com/willbehn/go-ifi-feed/models"
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
		footerLines := 3

		availHeight := max(msg.Height-footerLines, 1)

		if !m.ready {
			m.vp = viewport.New(msg.Width, availHeight)
			m.vp.MouseWheelEnabled = true

			style := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffffff"))
			colored := style.Render(models.Banner)

			m.vp.SetContent(colored + "\n" + m.content)
			m.ready = true
		} else {
			m.vp.Width = msg.Width
			m.vp.Height = availHeight
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

	footer := fmt.Sprintf("\n\n  Scroll: %.0f%% — press q to quit", m.vp.ScrollPercent()*100)
	return m.vp.View() + footer
}

func CombineMessages(courses models.Courses) string {
	var contents []string
	for _, message := range feed.Fetch(courses) {
		contents = append(contents, message.Content)
	}

	var allMessages = feed.ConvertToMarkdown(strings.Join(contents, ""))
	out, err := glamour.Render(allMessages, "dark")

	if err != nil {
		fmt.Println("Error rendering markdown:", err)
	}

	return out
}
