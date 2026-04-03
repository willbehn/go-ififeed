package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
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

func cmdList(courses models.Courses) {
	header := lipgloss.NewStyle().Bold(true)
	fmt.Printf("%s\n", prefix)
	fmt.Printf("%s  %s  %s\n", header.Render("emne"), header.Render("semester"), header.Render("tittel"))
	for _, c := range courses.Courses {
		fmt.Printf("%s  %s  %s\n", accent.Render(c.Code), subtle.Render(c.Semester), c.Title)
	}
}

func cmdAdd(code, semester, title string) {
	course := models.Course{
		Code:     code,
		Semester: semester,
		Title:    title,
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

func cmdRemove(code, semester string) {
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
	rootCmd := &cobra.Command{
		Use:   "ififeed",
		Short: "ifi emnebeskjeder rett i terminalen",
		RunE: func(cmd *cobra.Command, args []string) error {
			courses, err := readCourses()
			if err != nil {
				return err
			}
			runTui(courses)
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List alle emner i config",
		RunE: func(cmd *cobra.Command, args []string) error {
			courses, err := readCourses()
			if err != nil {
				return err
			}
			cmdList(courses)
			return nil
		},
	}

	addCmd := &cobra.Command{
		Use:   "add <emnekode> <semester> <tittel>",
		Short: "Legger til et nytt emne du vil få bskjeder fra",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdAdd(args[0], args[1], args[2])
			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <emnekode> <semester>",
		Short: "Fjerner et emne du får beskjeder fra",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdRemove(args[0], args[1])
			return nil
		},
	}

	rootCmd.AddCommand(listCmd, addCmd, removeCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
