package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/willbehn/go-ifi-feed/feed"
)

func joinMessages(courses []string) string {
	var contents []string
	for _, message := range feed.Fetch(courses) {
		contents = append(contents, message.Content)
	}
	return strings.Join(contents, "")
}

func runTui(courses []string) {
	p := tea.NewProgram(Model{content: joinMessages(courses)},
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

		courses := []string{"IN1020", "IN2040", "IN3050", "IN5060"}

		runTui(courses)
	}

}
