package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/willbehn/go-ififeed/internal"
	"github.com/willbehn/go-ififeed/models"
)

var addCmd = &cobra.Command{
	Use:   "add <emnekode> <semester> <tittel>",
	Short: "Legger til et nytt emne du vil få beskjeder fra",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		course := models.Course{
			Code:     args[0],
			Semester: args[1],
			Title:    args[2],
		}

		courses, err := internal.ReadCourses()
		if err != nil {
			return err
		}

		courses.Courses = append(courses.Courses, course)

		if err := internal.WriteCourses(courses); err != nil {
			return err
		}

		fmt.Printf("%s la til %s %s\n", prefix, accent.Render(course.Code), subtle.Render(course.Semester))
		return nil
	},
}
