package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/willbehn/go-ifi-feed/feed"
)

// TODO gjør disse velgbare i tui
const (
	AsciiStyle      = "ascii"
	AutoStyle       = "auto"
	DarkStyle       = "dark"
	DraculaStyle    = "dracula"
	TokyoNightStyle = "tokyo-night"
	LightStyle      = "light"
	NoTTYStyle      = "notty"
	PinkStyle       = "pink"
)

type itemMsg struct{ content string }
type doneMsg struct{}

type Model struct {
	vp        viewport.Model
	content   string
	ready     bool
	loading   bool
	msgCh     chan feed.Message
	spin      spinner.Model
	itemCount int
}

func newModel(ch chan feed.Message) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return Model{loading: true, msgCh: ch, spin: s}
}

func waitForMsg(ch chan feed.Message) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return doneMsg{}
		}
		md := feed.ConvertToMarkdown(msg.Content)
		out, err := glamour.Render(md, TokyoNightStyle)
		if err != nil {
			return itemMsg{content: md}
		}
		return itemMsg{content: out}
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spin.Tick,
		waitForMsg(m.msgCh),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case itemMsg:
		m.content += msg.content
		m.itemCount++
		if m.ready {
			m.vp.SetContent(m.content)
		}
		return m, waitForMsg(m.msgCh)

	case doneMsg:
		m.loading = false
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spin, cmd = m.spin.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		reservedLines := 3
		availHeight := max(msg.Height-reservedLines, 1)

		if !m.ready {
			m.vp = viewport.New(msg.Width, availHeight)
			m.vp.MouseWheelEnabled = true
			m.vp.SetContent(m.content)
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
	subtleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	accentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39"))

	if !m.ready {
		return "\n  Loading...\n"
	}

	sep := subtleStyle.Render(strings.Repeat("─", m.vp.Width))

	// Footer
	var footerLeft string
	if m.loading {
		footerLeft = "  " + accentStyle.Render(m.spin.View()) + subtleStyle.Render(fmt.Sprintf(" fetching (%d loaded)", m.itemCount))
	} else {
		footerLeft = subtleStyle.Render(fmt.Sprintf(" %d items", m.itemCount))
	}

	footerMiddle := subtleStyle.Render("ififeed v1")
	footerRight := subtleStyle.Render(fmt.Sprintf("q quit   %.0f%%  ", m.vp.ScrollPercent()*100))

	gap := max((m.vp.Width-lipgloss.Width(footerLeft)-lipgloss.Width(footerMiddle)-lipgloss.Width(footerRight))/2, 0)
	footer := footerLeft + strings.Repeat(" ", gap) + footerMiddle + strings.Repeat(" ", gap) + footerRight

	return m.vp.View() + "\n" + sep + "\n" + footer
}
