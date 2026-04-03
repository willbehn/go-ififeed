package cmd

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/willbehn/go-ififeed/internal"
	"github.com/willbehn/go-ififeed/tui"
)

var (
	accent = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	subtle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	prefix = accent.Render("ififeed") + subtle.Render(":")
)

var rootCmd = &cobra.Command{
	Use:   "ififeed",
	Short: "ifi emnebeskjeder rett i terminalen",
	RunE: func(cmd *cobra.Command, args []string) error {
		courses, err := internal.ReadCourses()
		if err != nil {
			return err
		}
		tui.Run(courses)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(listCmd, addCmd, removeCmd)
}
