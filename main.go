package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/willbehn/go-ifi-feed/feed"
	"github.com/willbehn/go-ifi-feed/models"
	yaml "gopkg.in/yaml.v3"
)

func joinMessages(courses models.Courses) string {
	var contents []string
	for _, message := range feed.Fetch(courses) {
		contents = append(contents, message.Content)
	}
	return strings.Join(contents, "")
}

func runTui(courses models.Courses) {
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

func readCourses() (models.Courses, error) {
	file, err := os.Open("courses.yaml")

	if err != nil {
		return models.Courses{}, err
	}

	defer file.Close()

	var data models.Courses
	if err := yaml.NewDecoder(file).Decode(&data); err != nil {
		return models.Courses{}, err
	}

	return data, nil
}

func main() {
	args := os.Args

	if len(args) > 1 {
		if args[1] == "subscribe" || args[1] == "-s" {
			cmdSubscribe(args[2:])
		}

	} else {
		courses, err := readCourses()

		if err != nil {
			fmt.Print(err)
			return
		}

		runTui(courses)
	}

}
