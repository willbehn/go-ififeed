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

type itemMsg struct{ content string }
type doneMsg struct{}

type Model struct {
	vp      viewport.Model
	content string
	ready   bool
	loading bool
	msgCh   chan feed.Message
}

func waitForMsg(ch chan feed.Message) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return doneMsg{}
		}
		md := feed.ConvertToMarkdown(msg.Content)
		out, err := glamour.Render(md, "dark")
		if err != nil {
			return itemMsg{content: md}
		}
		return itemMsg{content: out}
	}
}

func (m Model) Init() tea.Cmd {
	return waitForMsg(m.msgCh)
}

func banner() string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true).Render(models.Banner)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case itemMsg:
		m.content += msg.content
		if m.ready {
			m.vp.SetContent(banner() + "\n" + m.content)
		}
		return m, waitForMsg(m.msgCh)

	case doneMsg:
		m.loading = false
		return m, nil

	case tea.WindowSizeMsg:
		footerLines := 3
		availHeight := max(msg.Height-footerLines, 1)

		if !m.ready {
			m.vp = viewport.New(msg.Width, availHeight)
			m.vp.MouseWheelEnabled = true
			m.vp.SetContent(banner() + "\n" + m.content)
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
		return "\n  Loading...\n"
	}

	subtle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	separator := subtle.Render(strings.Repeat("─", m.vp.Width))

	status := "  q quit"
	if m.loading {
		status = "  fetching...   q quit"
	}
	footer := subtle.Render(fmt.Sprintf("%s   %.0f%%", status, m.vp.ScrollPercent()*100))

	return m.vp.View() + "\n" + separator + "\n" + footer
}
