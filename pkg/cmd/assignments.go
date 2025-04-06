package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Reisender/canvas-cli-v2/pkg/api"
	"github.com/Reisender/canvas-cli-v2/pkg/ui"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
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
		newAssignmentsAddCmd(),
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
		Run:   runAssignmentsView,
	}
}

func newAssignmentsAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add [course-id]",
		Short: "Add a new assignment to a course",
		Long:  `Create a new assignment in a Canvas course with interactive form input.`,
		Args:  cobra.ExactArgs(1),
		Run:   runAssignmentsAdd,
	}
}

// AssignmentForm represents the data collected from the form
type AssignmentForm struct {
	Name            string
	Description     string
	PointsPossible  float64
	DueDate         string
	UnlockDate      string
	LockDate        string
	GradingType     string
	SubmissionTypes []string
	Published       bool
}

// AssignmentDetailModel represents a model for viewing assignment details
type AssignmentDetailModel struct {
	assignment   *api.Assignment
	viewport     viewport.Model
	ready        bool
	width        int
	height       int
	courseID     string
	assignmentID string
}

// Initialize the assignment detail model
func NewAssignmentDetailModel(courseID, assignmentID string) AssignmentDetailModel {
	return AssignmentDetailModel{
		courseID:     courseID,
		assignmentID: assignmentID,
	}
}

// Init initializes the assignment detail model
func (m AssignmentDetailModel) Init() tea.Cmd {
	return func() tea.Msg {
		client := api.NewClient()
		assignment, err := client.GetAssignment(m.courseID, m.assignmentID)
		if err != nil {
			return AssignmentDetailErrorMsg{err}
		}
		return AssignmentDetailLoadedMsg{assignment}
	}
}

// Messages for the assignment detail model
type AssignmentDetailLoadedMsg struct {
	assignment *api.Assignment
}

type AssignmentDetailErrorMsg struct {
	err error
}

// Update updates the assignment detail model
func (m AssignmentDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "enter":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		if !m.ready {
			// If we haven't initialized yet, do so now
			m.viewport = viewport.New(m.width, m.height-4) // leave room for header/footer
			m.viewport.Style = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")).
				PaddingRight(2)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 4
		}

		if m.assignment != nil {
			m.viewport.SetContent(m.formatAssignmentDetails())
		}

	case AssignmentDetailLoadedMsg:
		m.assignment = msg.assignment
		if m.ready {
			m.viewport.SetContent(m.formatAssignmentDetails())
		}

	case AssignmentDetailErrorMsg:
		return m, tea.Quit
	}

	// Handle viewport updates
	if m.ready {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the assignment detail model
func (m AssignmentDetailModel) View() string {
	if !m.ready {
		return "Loading..."
	}

	if m.assignment == nil {
		return "Error loading assignment details"
	}

	// Create header style
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		MarginBottom(1).
		PaddingLeft(2)

	// Create footer style
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		PaddingTop(1).
		PaddingLeft(2)

	// Combine all the parts with header and footer
	return headerStyle.Render("Assignment Details") + "\n" +
		m.viewport.View() + "\n" +
		footerStyle.Render("q/esc/enter: return to list")
}

// formatAssignmentDetails formats the assignment details as a styled string
func (m AssignmentDetailModel) formatAssignmentDetails() string {
	assignment := m.assignment

	// Create styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		MarginBottom(1).
		Width(m.width - 4)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true).
		Width(18)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Width(m.width - 24)

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		Bold(true).
		MarginTop(1).
		MarginBottom(1)

	// Build the content with styled components
	var content strings.Builder

	// Title section
	content.WriteString(titleStyle.Render(assignment.Name) + "\n\n")

	// Basic information section
	content.WriteString(sectionStyle.Render("Basic Information") + "\n")

	content.WriteString(labelStyle.Render("ID:") + valueStyle.Render(fmt.Sprintf("%d", assignment.ID)) + "\n")
	content.WriteString(labelStyle.Render("Points:") + valueStyle.Render(fmt.Sprintf("%.1f", assignment.PointsPossible)) + "\n")

	// Dates section
	content.WriteString(sectionStyle.Render("Dates") + "\n")

	dueDate := "Not set"
	if !assignment.DueAt.IsZero() {
		dueDate = assignment.DueAt.Format("Jan 2, 2006 3:04 PM")
	}
	content.WriteString(labelStyle.Render("Due Date:") + valueStyle.Render(dueDate) + "\n")

	unlockDate := "Not set"
	if !assignment.UnlockAt.IsZero() {
		unlockDate = assignment.UnlockAt.Format("Jan 2, 2006 3:04 PM")
	}
	content.WriteString(labelStyle.Render("Available From:") + valueStyle.Render(unlockDate) + "\n")

	lockDate := "Not set"
	if !assignment.LockAt.IsZero() {
		lockDate = assignment.LockAt.Format("Jan 2, 2006 3:04 PM")
	}
	content.WriteString(labelStyle.Render("Available Until:") + valueStyle.Render(lockDate) + "\n")

	// Configuration section
	content.WriteString(sectionStyle.Render("Configuration") + "\n")

	content.WriteString(labelStyle.Render("Grading Type:") + valueStyle.Render(assignment.GradingType) + "\n")
	content.WriteString(labelStyle.Render("Submission Types:") + valueStyle.Render(strings.Join(assignment.SubmissionTypes, ", ")) + "\n")

	publishedStatus := "No"
	if assignment.Published {
		publishedStatus = "Yes"
	}
	content.WriteString(labelStyle.Render("Published:") + valueStyle.Render(publishedStatus) + "\n")

	// Metadata section
	content.WriteString(sectionStyle.Render("Metadata") + "\n")

	content.WriteString(labelStyle.Render("Created:") + valueStyle.Render(assignment.CreatedAt.Format("Jan 2, 2006")) + "\n")

	updatedAt := "Same as creation date"
	if !assignment.UpdatedAt.IsZero() && !assignment.CreatedAt.Equal(assignment.UpdatedAt) {
		updatedAt = assignment.UpdatedAt.Format("Jan 2, 2006")
	}
	content.WriteString(labelStyle.Render("Last Updated:") + valueStyle.Render(updatedAt) + "\n")

	// Description section
	content.WriteString(sectionStyle.Render("Description") + "\n")

	// Wrap description text to fit viewport
	descriptionStyle := lipgloss.NewStyle().Width(m.width - 6)
	content.WriteString(descriptionStyle.Render(assignment.Description) + "\n")

	return content.String()
}

// runAssignmentsView displays detailed information about a specific assignment
func runAssignmentsView(cmd *cobra.Command, args []string) {
	courseID := args[0]
	assignmentID := args[1]

	// Initialize the assignment detail model
	model := NewAssignmentDetailModel(courseID, assignmentID)

	// Run the program
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running assignment detail view: %v\n", err)
		return
	}
}

// runAssignmentsAdd runs the add assignment command
func runAssignmentsAdd(cmd *cobra.Command, args []string) {
	courseID := args[0]

	// Available submission types
	submissionTypes := []string{
		"online_text_entry",
		"online_url",
		"online_upload",
		"media_recording",
		"none",
	}

	// Available grading types
	gradingTypes := []string{
		"points",
		"pass_fail",
		"percent",
		"letter_grade",
		"gpa_scale",
	}

	// Create the form data structure
	form := AssignmentForm{
		GradingType:     "points",
		SubmissionTypes: []string{"online_text_entry"},
		Published:       true,
	}

	// Build the form with huh
	formUI := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Create New Assignment").
				Description("Enter the details for the new assignment"),

			huh.NewInput().
				Title("Name").
				Prompt("> ").
				Placeholder("Enter assignment name").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}).
				Value(&form.Name),

			huh.NewText().
				Title("Description").
				Placeholder("Enter assignment description").
				Editor("vi").
				CharLimit(1000).
				Value(&form.Description),

			huh.NewInput().
				Title("Points Possible").
				Prompt("> ").
				Placeholder("Enter the maximum points (e.g. 100)").
				Validate(func(s string) error {
					if s == "" {
						return nil
					}
					val, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return fmt.Errorf("points must be a number")
					}
					if val < 0 {
						return fmt.Errorf("points cannot be negative")
					}
					form.PointsPossible = val
					return nil
				}),

			huh.NewInput().
				Title("Due Date").
				Prompt("> ").
				Placeholder("Format: YYYY-MM-DD HH:MM").
				Validate(func(s string) error {
					if s == "" {
						return nil // optional
					}
					_, err := time.Parse("2006-01-02 15:04", s)
					if err != nil {
						return fmt.Errorf("invalid date format")
					}
					form.DueDate = s
					return nil
				}),

			huh.NewInput().
				Title("Unlock Date").
				Prompt("> ").
				Placeholder("Format: YYYY-MM-DD HH:MM (optional)").
				Validate(func(s string) error {
					if s == "" {
						return nil // optional
					}
					_, err := time.Parse("2006-01-02 15:04", s)
					if err != nil {
						return fmt.Errorf("invalid date format")
					}
					form.UnlockDate = s
					return nil
				}),

			huh.NewInput().
				Title("Lock Date").
				Prompt("> ").
				Placeholder("Format: YYYY-MM-DD HH:MM (optional)").
				Validate(func(s string) error {
					if s == "" {
						return nil // optional
					}
					_, err := time.Parse("2006-01-02 15:04", s)
					if err != nil {
						return fmt.Errorf("invalid date format")
					}
					form.LockDate = s
					return nil
				}),

			huh.NewSelect[string]().
				Title("Grading Type").
				Options(
					huh.NewOptions(gradingTypes...)...,
				).
				Value(&form.GradingType),

			huh.NewMultiSelect[string]().
				Title("Submission Types").
				Options(
					huh.NewOptions(submissionTypes...)...,
				).
				Value(&form.SubmissionTypes),

			huh.NewConfirm().
				Title("Published").
				Description("Make the assignment visible to students").
				Value(&form.Published),
		),
	).WithTheme(huh.ThemeBase16())

	// Run the form UI
	err := formUI.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error with form: %v\n", err)
		return
	}

	// Create the assignment object
	assignment := &api.Assignment{
		Name:            form.Name,
		Description:     form.Description,
		PointsPossible:  form.PointsPossible,
		GradingType:     form.GradingType,
		Published:       form.Published,
		SubmissionTypes: form.SubmissionTypes,
	}

	// Parse dates if provided
	if form.DueDate != "" {
		dueDate, _ := time.Parse("2006-01-02 15:04", form.DueDate)
		assignment.DueAt = dueDate
	}

	if form.UnlockDate != "" {
		unlockDate, _ := time.Parse("2006-01-02 15:04", form.UnlockDate)
		assignment.UnlockAt = unlockDate
	}

	if form.LockDate != "" {
		lockDate, _ := time.Parse("2006-01-02 15:04", form.LockDate)
		assignment.LockAt = lockDate
	}

	// Call the API
	client := api.NewClient()
	newAssignment, err := client.CreateAssignment(courseID, assignment)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating assignment: %v\n", err)
		return
	}

	// Show a success message
	fmt.Println("\n✅ Assignment created successfully!")
	fmt.Printf("ID: %d\n", newAssignment.ID)
	fmt.Printf("Name: %s\n", newAssignment.Name)
	fmt.Printf("Points: %.1f\n", newAssignment.PointsPossible)

	// Format and display the dates
	if !newAssignment.DueAt.IsZero() {
		fmt.Printf("Due Date: %s\n", newAssignment.DueAt.Format("2006-01-02 15:04"))
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
	m.Title = fmt.Sprintf("Assignments for Course %s", courseID)
	m.Help = "↑/↓: Navigate • enter: View Assignment • q: Quit"

	// Set up the selection callback to view assignment details
	m.OnSelect = func(row table.Row) {
		// Get assignment ID from the selected row
		assignmentID := row[0]

		// Create the command args for viewing the assignment
		viewArgs := []string{courseID, assignmentID}

		// Clear screen
		fmt.Print("\033[H\033[2J")

		// Run the view command immediately
		runAssignmentsView(nil, viewArgs)

		// After returning from detail view, restart list view
		runAssignmentsList(nil, args)
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
