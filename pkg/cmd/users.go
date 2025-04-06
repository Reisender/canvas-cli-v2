package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Reisender/canvas-cli-v2/pkg/api"
	"github.com/Reisender/canvas-cli-v2/pkg/ui"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// NewUsersCmd creates a new command for managing users
func NewUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Manage Canvas users",
		Long:  `List, view, and manage users in Canvas courses.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add subcommands
	cmd.AddCommand(
		newUsersListCmd(),
		newUsersViewCmd(),
		newEnrollmentsCmd(),
		newUsersRemoveCmd(),
	)

	return cmd
}

func newUsersListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list [course-id]",
		Short: "List users in a course",
		Long:  `List all users enrolled in a specific Canvas course.`,
		Args:  cobra.ExactArgs(1),
		Run:   runUsersList,
	}
}

func newUsersViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view [user-id]",
		Short: "View a Canvas user",
		Long:  `View details about a specific Canvas user.`,
		Args:  cobra.ExactArgs(1),
		Run:   runUsersView,
	}
}

func newUsersRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove [course-id] [user-id]",
		Short: "Remove a user from a course",
		Long:  `Remove a user from a Canvas course using the user ID.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			courseID := args[0]
			userID := args[1]

			client := api.NewClient()
			if err := client.RemoveUserByID(courseID, userID); err != nil {
				fmt.Fprintf(os.Stderr, "Error removing user: %v\n", err)
				return
			}

			fmt.Printf("Successfully removed user %s from course %s\n", userID, courseID)
		},
	}
}

func newEnrollmentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enrollments",
		Short: "Manage course enrollments",
		Long:  `List, add, and remove course enrollments.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add enrollment subcommands
	cmd.AddCommand(
		newEnrollmentsListCmd(),
		newEnrollmentsAddCmd(),
		newEnrollmentsRemoveCmd(),
	)

	return cmd
}

func newEnrollmentsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list [course-id]",
		Short: "List enrollments for a course",
		Long:  `List all enrollments for a specific Canvas course.`,
		Args:  cobra.ExactArgs(1),
		Run:   runEnrollmentsList,
	}
}

func newEnrollmentsAddCmd() *cobra.Command {
	var enrollmentType string
	var notify bool

	cmd := &cobra.Command{
		Use:   "add [course-id] [user-id]",
		Short: "Add a user to a course",
		Long:  `Enroll a user in a Canvas course with the specified role.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			courseID := args[0]
			userID := args[1]

			client := api.NewClient()
			enrollment, err := client.AddUserToCourse(courseID, userID, enrollmentType, notify)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error enrolling user: %v\n", err)
				return
			}

			fmt.Printf("Successfully enrolled user %d in course %d with role %s\n",
				enrollment.UserID, enrollment.CourseID, enrollment.Role)
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&enrollmentType, "type", "t", "StudentEnrollment",
		"Enrollment type (StudentEnrollment, TeacherEnrollment, TaEnrollment, ObserverEnrollment, DesignerEnrollment)")
	cmd.Flags().BoolVarP(&notify, "notify", "n", false, "Send enrollment notification to the user")

	return cmd
}

func newEnrollmentsRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove [course-id] [enrollment-id]",
		Short: "Remove an enrollment",
		Long:  `Remove a user's enrollment from a Canvas course.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			courseID := args[0]
			enrollmentID := args[1]

			client := api.NewClient()
			if err := client.RemoveUserFromCourse(courseID, enrollmentID); err != nil {
				fmt.Fprintf(os.Stderr, "Error removing enrollment: %v\n", err)
				return
			}

			fmt.Printf("Successfully removed enrollment %s from course %s\n", enrollmentID, courseID)
		},
	}
}

func runUsersList(cmd *cobra.Command, args []string) {
	courseID := args[0]
	client := api.NewClient()
	users, err := client.GetUsers(courseID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching users: %v\n", err)
		return
	}

	// Create a table for users
	columns := []table.Column{
		{Title: "ID", Width: 10},
		{Title: "Name", Width: 30},
		{Title: "Email", Width: 30},
		{Title: "Login ID", Width: 15},
	}

	rows := []table.Row{}
	for _, user := range users {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", user.ID),
			user.Name,
			user.Email,
			user.LoginID,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
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
	m.Title = fmt.Sprintf("Users in Course %s", courseID)
	m.Help = "↑/↓: Navigate • enter: Select • q: Quit"

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

func runUsersView(cmd *cobra.Command, args []string) {
	userID := args[0]
	client := api.NewClient()
	user, err := client.GetUserDetails(userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching user details: %v\n", err)
		return
	}

	// Display user information
	fmt.Println("User Details:")
	fmt.Println("-------------")
	fmt.Printf("ID:           %d\n", user.ID)
	fmt.Printf("Name:         %s\n", user.Name)
	fmt.Printf("SortableName: %s\n", user.SortableName)
	fmt.Printf("ShortName:    %s\n", user.ShortName)
	fmt.Printf("Email:        %s\n", user.Email)
	fmt.Printf("Login ID:     %s\n", user.LoginID)
	fmt.Printf("SIS User ID:  %s\n", user.SISUserID)
	if user.Avatar != "" {
		fmt.Printf("Avatar URL:   %s\n", user.Avatar)
	}
	if user.Locale != "" {
		fmt.Printf("Locale:       %s\n", user.Locale)
	}
}

func runEnrollmentsList(cmd *cobra.Command, args []string) {
	courseID := args[0]
	client := api.NewClient()
	enrollments, err := client.GetEnrollments(courseID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching enrollments: %v\n", err)
		return
	}

	// Create a table for enrollments
	columns := []table.Column{
		{Title: "Enrollment ID", Width: 12},
		{Title: "User ID", Width: 10},
		{Title: "User Name", Width: 25},
		{Title: "Role", Width: 15},
		{Title: "Status", Width: 10},
	}

	rows := []table.Row{}
	for _, enrollment := range enrollments {
		rows = append(rows, table.Row{
			strconv.Itoa(enrollment.ID),
			strconv.Itoa(enrollment.UserID),
			enrollment.User.Name,
			enrollment.Role,
			enrollment.EnrollmentState,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
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
	m.Title = fmt.Sprintf("Enrollments for Course %s", courseID)
	m.Help = "↑/↓: Navigate • enter: Select • q: Quit"

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
