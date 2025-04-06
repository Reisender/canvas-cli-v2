package cmd

import (
	"fmt"
	"os"

	"github.com/Reisender/canvas-cli-v2/pkg/api"
	"github.com/Reisender/canvas-cli-v2/pkg/ui"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// NewCoursesCmd creates a new command for managing courses
func NewCoursesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "courses",
		Short: "Manage Canvas courses",
		Long:  `List, view, and interact with your Canvas courses.`,
		Run:   runCoursesList,
	}

	// Add subcommands
	cmd.AddCommand(
		newCoursesListCmd(),
		newCoursesViewCmd(),
	)

	return cmd
}

func newCoursesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List Canvas courses",
		Long:  `List all courses you have access to in Canvas.`,
		Run:   runCoursesList,
	}
}

func newCoursesViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view [course-id]",
		Short: "View a Canvas course",
		Long:  `View details about a specific Canvas course.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("View course:", args[0])
			// TODO: Implement course view
		},
	}
}

func runCoursesList(cmd *cobra.Command, args []string) {
	client := api.NewClient()
	courses, err := client.GetCourses()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching courses: %v\n", err)
		return
	}

	// Create a table for courses
	columns := []table.Column{
		{Title: "ID", Width: 10},
		{Title: "Course Code", Width: 15},
		{Title: "Name", Width: 40},
	}

	rows := []table.Row{}
	for _, course := range courses {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", course.ID),
			course.CourseCode,
			course.Name,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)
	t.SetStyles(s)

	m := ui.NewTableModel(t)
	m.Title = "Canvas Courses"
	m.Help = "↑/↓: Navigate • enter: Select • q: Quit"

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
