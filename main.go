package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/willbehn/go-ififeed/feed"
	"github.com/willbehn/go-ififeed/models"
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
		fmt.Println("bruk: ififeed add \"<emnekode>\" \"<semester>\" \"<emnetittel>\"")
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

	fmt.Printf("la til %s (%s)\n", course.Code, course.Semester)
}

func cmdList(courses models.Courses) {
	for _, c := range courses.Courses {
		fmt.Printf("%s\t%s\t%s\n", c.Code, c.Semester, c.Title)
	}
}

func cmdRemove(args []string) {
	if len(args) < 2 {
		fmt.Println("bruk: ififeed remove \"<emnekode>\" \"<semester>\"")
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
		fmt.Printf("fant ikke %s (%s)\n", code, semester)
		return
	}

	courses.Courses = filtered

	if err := writeCourses(courses); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("fjernet %s (%s)\n", code, semester)
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
