package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/willbehn/go-ififeed/internal"
)

var removeCmd = &cobra.Command{
	Use:   "remove <emnekode> <semester>",
	Short: "Fjerner et emne du får beskjeder fra",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		code := args[0]
		semester := args[1]

		courses, err := internal.ReadCourses()
		if err != nil {
			return err
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
			return nil
		}

		courses.Courses = filtered

		if err := internal.WriteCourses(courses); err != nil {
			return err
		}

		fmt.Printf("%s fjernet %s %s\n", prefix, accent.Render(code), subtle.Render(semester))
		return nil
	},
}
