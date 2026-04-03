package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/willbehn/go-ififeed/feed"
	"github.com/willbehn/go-ififeed/models"
)

func Run(courses models.Courses) {
	ch := make(chan feed.Message, 50)
	go feed.FetchStream(courses, ch)

	p := tea.NewProgram(
		newModel(ch),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	p.Run()
}
