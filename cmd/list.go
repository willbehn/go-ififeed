package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/willbehn/go-ififeed/internal"
	"github.com/willbehn/go-ififeed/models"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List alle emner i config",
	RunE: func(cmd *cobra.Command, args []string) error {
		courses, err := internal.ReadCourses()
		if err != nil {
			return err
		}
		printList(courses)
		return nil
	},
}

func printList(courses models.Courses) {
	header := lipgloss.NewStyle().Bold(true)
	fmt.Printf("%s\n", prefix)
	fmt.Printf("%s  %s  %s\n", header.Render("emne"), header.Render("semester"), header.Render("tittel"))
	for _, c := range courses.Courses {
		fmt.Printf("%s  %s  %s\n", accent.Render(c.Code), subtle.Render(c.Semester), c.Title)
	}
}
