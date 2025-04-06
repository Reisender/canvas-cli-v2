package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
	var multiSelect bool

	cmd := &cobra.Command{
		Use:   "list [course-id]",
		Short: "List users in a course",
		Long:  `List all users enrolled in a specific Canvas course.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runUsersList(args[0], multiSelect)
		},
	}

	cmd.Flags().BoolVarP(&multiSelect, "multi", "m", false, "Enable multi-selection mode")
	return cmd
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

// UserActionModel represents the model for the user action selection screen
type UserActionModel struct {
	courseID  string
	userID    string
	userName  string
	choices   []string
	cursor    int
	client    *api.Client
	completed bool
	result    string
}

func (m UserActionModel) Init() tea.Cmd {
	return nil
}

func (m UserActionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor == 0 {
				// View user details
				user, err := m.client.GetUserDetails(m.userID)
				if err != nil {
					m.result = fmt.Sprintf("Error fetching user details: %v", err)
				} else {
					// Format user details in the same way as runUsersView
					var details strings.Builder
					details.WriteString("\nUser Details:\n")
					details.WriteString("-------------\n")
					details.WriteString(fmt.Sprintf("ID:           %d\n", user.ID))
					details.WriteString(fmt.Sprintf("Name:         %s\n", user.Name))
					details.WriteString(fmt.Sprintf("SortableName: %s\n", user.SortableName))
					details.WriteString(fmt.Sprintf("ShortName:    %s\n", user.ShortName))
					details.WriteString(fmt.Sprintf("Email:        %s\n", user.Email))
					details.WriteString(fmt.Sprintf("Login ID:     %s\n", user.LoginID))
					details.WriteString(fmt.Sprintf("SIS User ID:  %s\n", user.SISUserID))
					if user.Avatar != "" {
						details.WriteString(fmt.Sprintf("Avatar URL:   %s\n", user.Avatar))
					}
					if user.Locale != "" {
						details.WriteString(fmt.Sprintf("Locale:       %s\n", user.Locale))
					}
					m.result = details.String()
				}
				m.completed = true
				return m, tea.Quit
			} else if m.cursor == 1 {
				// Remove user
				err := m.client.RemoveUserByID(m.courseID, m.userID)
				if err != nil {
					m.result = fmt.Sprintf("Error removing user: %v", err)
				} else {
					m.result = fmt.Sprintf("Successfully removed user %s (%s) from course %s",
						m.userID, m.userName, m.courseID)
				}
				m.completed = true
				return m, tea.Quit
			} else {
				// Cancel
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m UserActionModel) View() string {
	if m.completed {
		return m.result
	}

	s := fmt.Sprintf("\nUser: %s (ID: %s)\n\n", m.userName, m.userID)
	s += "What would you like to do?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quit.\n"
	return s
}

// MultiActionModel represents the model for bulk actions on selected users
type MultiActionModel struct {
	courseID      string
	selectedUsers []table.Row
	choices       []string
	cursor        int
	client        *api.Client
	completed     bool
	result        string
	progress      int
	total         int
	success       int
	failed        int
}

func (m MultiActionModel) Init() tea.Cmd {
	return nil
}

func (m MultiActionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor == 0 {
				// Remove all selected users
				m.total = len(m.selectedUsers)

				var results strings.Builder
				results.WriteString(fmt.Sprintf("\nRemoving %d users from course %s...\n\n", m.total, m.courseID))

				for _, row := range m.selectedUsers {
					userID := row[0]
					userName := row[1]

					err := m.client.RemoveUserByID(m.courseID, userID)
					if err != nil {
						results.WriteString(fmt.Sprintf("❌ Failed to remove %s (%s): %v\n", userName, userID, err))
						m.failed++
					} else {
						results.WriteString(fmt.Sprintf("✅ Removed %s (%s)\n", userName, userID))
						m.success++
					}
					m.progress++
				}

				results.WriteString(fmt.Sprintf("\nSummary: %d/%d users removed successfully\n", m.success, m.total))
				m.result = results.String()
				m.completed = true
				return m, tea.Quit
			} else {
				// Cancel
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m MultiActionModel) View() string {
	if m.completed {
		return m.result
	}

	s := fmt.Sprintf("\n%d users selected in course %s\n\n", len(m.selectedUsers), m.courseID)
	s += "What would you like to do?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quit.\n"
	return s
}

func runUsersList(courseID string, multiSelect bool) {
	client := api.NewClient()

	// Initialize variables for pagination
	var allUsers []api.User
	page := 1
	perPage := 50
	moreUsers := true

	// Fetch users with pagination
	for moreUsers {
		users, err := client.GetUsers(courseID, page, perPage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching users: %v\n", err)
			return
		}

		// Add users to our collection
		allUsers = append(allUsers, users...)

		// If we got fewer users than requested, we've reached the end
		if len(users) < perPage {
			moreUsers = false
		} else {
			page++
		}
	}

	// If no users found
	if len(allUsers) == 0 {
		fmt.Println("No users found for this course.")
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
	for _, user := range allUsers {
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
	m.Title = fmt.Sprintf("Users in Course %s (%d users total)", courseID, len(allUsers))

	if multiSelect {
		m.EnableMultiSelect()

		// Set up the multi-selection callback
		m.OnMultiSelect = func(selectedRows []table.Row) {
			// Clear screen
			fmt.Print("\033[H\033[2J")

			// Create a new model for bulk actions
			actionModel := MultiActionModel{
				courseID:      courseID,
				selectedUsers: selectedRows,
				choices:       []string{"Remove all selected users", "Cancel"},
				client:        client,
			}

			// Run the action program
			p := tea.NewProgram(actionModel)
			result, err := p.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running action program: %v\n", err)
				return
			}

			// Get the final model
			finalModel, ok := result.(MultiActionModel)
			if ok && finalModel.completed {
				fmt.Println(finalModel.result)
			}
		}
	} else {
		// Single selection mode
		m.Help = "↑/↓: Navigate • enter: Select • q: Quit"

		// Set up the selection callback
		m.OnSelect = func(row table.Row) {
			// Clear screen
			fmt.Print("\033[H\033[2J")

			userID := row[0]
			userName := row[1]

			// Create a new model for user actions
			actionModel := UserActionModel{
				courseID: courseID,
				userID:   userID,
				userName: userName,
				choices:  []string{"View user details", "Remove user from course", "Cancel"},
				client:   client,
			}

			// Run the action program
			p := tea.NewProgram(actionModel)
			result, err := p.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running action program: %v\n", err)
				return
			}

			// Get the final model
			finalModel, ok := result.(UserActionModel)
			if ok && finalModel.completed {
				fmt.Println(finalModel.result)
			}
		}
	}

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
