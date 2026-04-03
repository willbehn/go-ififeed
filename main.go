package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/willbehn/go-ififeed/feed"
	"github.com/willbehn/go-ififeed/models"
)

var (
	accent = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	subtle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	prefix = accent.Render("ififeed") + subtle.Render(":")
)

func runTui(courses models.Courses) {
	ch := make(chan feed.Message, 50)
	go feed.FetchStream(courses, ch)

	p := tea.NewProgram(
		newModel(ch),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	p.Run()
}

func cmdAdd(args []string) {
	if len(args) < 3 {
		fmt.Printf("%s bruk: ififeed add \"<emnekode>\" \"<semester>\" \"<emnetittel>\"\n", prefix)
		return
	}

	course := models.Course{
		Code:     args[0],
		Semester: args[1],
		Title:    args[2],
	}

	courses, err := readCourses()
	if err != nil {
		fmt.Println(err)
		return
	}

	courses.Courses = append(courses.Courses, course)

	if err := writeCourses(courses); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s la til %s %s\n", prefix, accent.Render(course.Code), subtle.Render(course.Semester))
}

func cmdList(courses models.Courses) {
	header := lipgloss.NewStyle().Bold(true)
	fmt.Printf("%s\n", prefix)
	fmt.Printf("%s  %s  %s\n", header.Render("emne"), header.Render("semester"), header.Render("tittel"))
	for _, c := range courses.Courses {
		fmt.Printf("%s  %s  %s\n", accent.Render(c.Code), subtle.Render(c.Semester), c.Title)
	}
}

func cmdRemove(args []string) {
	if len(args) < 2 {
		fmt.Printf("%s bruk: ififeed remove \"<emnekode>\" \"<semester>\"\n", prefix)
		return
	}

	code := args[0]
	semester := args[1]

	courses, err := readCourses()
	if err != nil {
		fmt.Println(err)
		return
	}

	filtered := courses.Courses[:0]
	removed := false
	for _, c := range courses.Courses {
		if c.Code == code && c.Semester == semester {
			removed = true
		} else {
			filtered = append(filtered, c)
		}
	}

	if !removed {
		fmt.Printf("%s fant ikke %s %s\n", prefix, accent.Render(code), subtle.Render(semester))
		return
	}

	courses.Courses = filtered

	if err := writeCourses(courses); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s fjernet %s %s\n", prefix, accent.Render(code), subtle.Render(semester))
}

func main() {
	args := os.Args

	courses, err := readCourses()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(args) > 1 {
		switch args[1] {
		case "list":
			cmdList(courses)
		case "add":
			cmdAdd(args[2:])
		case "remove":
			cmdRemove(args[2:])
		}
	} else {
		runTui(courses)
	}
}
