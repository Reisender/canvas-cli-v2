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

// NewAssignmentsCmd creates a new command for managing assignments
func NewAssignmentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assignments",
		Short: "Manage Canvas assignments",
		Long:  `List, view, and interact with Canvas assignments.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add subcommands
	cmd.AddCommand(
		newAssignmentsListCmd(),
		newAssignmentsViewCmd(),
	)

	return cmd
}

func newAssignmentsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list [course-id]",
		Short: "List assignments for a course",
		Long:  `List all assignments for a specific course in Canvas.`,
		Args:  cobra.ExactArgs(1),
		Run:   runAssignmentsList,
	}
}

func newAssignmentsViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view [course-id] [assignment-id]",
		Short: "View a Canvas assignment",
		Long:  `View details about a specific Canvas assignment.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("View assignment:", args[1], "for course:", args[0])
			// TODO: Implement assignment view
		},
	}
}

func runAssignmentsList(cmd *cobra.Command, args []string) {
	courseID := args[0]
	client := api.NewClient()
	assignments, err := client.GetAssignments(courseID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching assignments: %v\n", err)
		return
	}

	// Create a table for assignments
	columns := []table.Column{
		{Title: "ID", Width: 10},
		{Title: "Name", Width: 40},
		{Title: "Due Date", Width: 20},
		{Title: "Points", Width: 10},
	}

	rows := []table.Row{}
	for _, assignment := range assignments {
		dueDate := ""
		if !assignment.DueAt.IsZero() {
			dueDate = assignment.DueAt.Format("Jan 2, 2006 3:04 PM")
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", assignment.ID),
			assignment.Name,
			dueDate,
			fmt.Sprintf("%.1f", assignment.PointsPossible),
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
	m.Title = fmt.Sprintf("Assignments for Course %s", courseID)
	m.Help = "↑/↓: Navigate • enter: Select • q: Quit"

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
