package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SelectionCallback is a function called when a row is selected
type SelectionCallback func(row table.Row)

// MultiSelectionCallback is a function called with multiple selected rows
type MultiSelectionCallback func(rows []table.Row)

// TableModel represents a table UI model
type TableModel struct {
	table           table.Model
	selectionTable  table.Model // Separate table for multi-selection view
	Title           string
	Help            string
	OnSelect        SelectionCallback
	OnMultiSelect   MultiSelectionCallback
	selectedRows    map[int]bool
	multiSelectMode bool
}

// NewTableModel creates a new table model
func NewTableModel(t table.Model) *TableModel {
	return &TableModel{
		table:           t,
		Title:           "Table",
		Help:            "↑/↓: Navigate • enter: Select • q: Quit",
		selectedRows:    make(map[int]bool),
		multiSelectMode: false,
	}
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			MarginLeft(2)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(2).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("63")).
			Bold(true)

	selectionIndicator = lipgloss.NewStyle().
				Foreground(lipgloss.Color("170")).
				Bold(true).
				Render("✓ ")

	noSelectionIndicator = "  "
)

// Init initializes the table model
func (m TableModel) Init() tea.Cmd {
	return nil
}

// IsRowSelected checks if a row is selected
func (m TableModel) IsRowSelected(index int) bool {
	return m.selectedRows[index]
}

// ToggleRow toggles selection status of the current row
func (m *TableModel) ToggleRow() {
	currentIndex := m.table.Cursor()
	if m.selectedRows[currentIndex] {
		delete(m.selectedRows, currentIndex)
	} else {
		m.selectedRows[currentIndex] = true
	}

	// Update selection table if it exists
	if m.multiSelectMode {
		m.updateSelectionTable()
	}
}

// updateSelectionTable creates or updates the selection table
func (m *TableModel) updateSelectionTable() {
	rows := m.table.Rows()
	cursorPos := m.table.Cursor()

	// Build columns with an extra selection column
	columns := []table.Column{
		{Title: "", Width: 2},
	}
	for _, col := range m.table.Columns() {
		columns = append(columns, col)
	}

	// Create new rows with checkmarks
	newRows := make([]table.Row, len(rows))
	for i, row := range rows {
		// If selected, add a checkmark as the first element
		indicator := ""
		if m.IsRowSelected(i) {
			indicator = "✓"
		}

		newRow := make(table.Row, len(row)+1)
		newRow[0] = indicator
		copy(newRow[1:], row)
		newRows[i] = newRow
	}

	// Create the new table
	newTable := table.New(
		table.WithColumns(columns),
		table.WithRows(newRows),
		table.WithFocused(true),
		table.WithHeight(m.table.Height()),
	)

	// Apply styles
	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	tableStyles.Selected = tableStyles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)
	newTable.SetStyles(tableStyles)

	// Set cursor to match original table
	newTable.SetCursor(cursorPos)

	// Save the selection table
	m.selectionTable = newTable
}

// GetSelectedRows returns all selected rows
func (m TableModel) GetSelectedRows() []table.Row {
	var selected []table.Row
	rows := m.table.Rows()

	for i := range rows {
		if m.selectedRows[i] {
			selected = append(selected, rows[i])
		}
	}

	return selected
}

// SelectAll selects all rows
func (m *TableModel) SelectAll() {
	for i := range m.table.Rows() {
		m.selectedRows[i] = true
	}

	// Update selection table
	if m.multiSelectMode {
		m.updateSelectionTable()
	}
}

// ClearSelections clears all selected rows
func (m *TableModel) ClearSelections() {
	m.selectedRows = make(map[int]bool)

	// Update selection table
	if m.multiSelectMode {
		m.updateSelectionTable()
	}
}

// EnableMultiSelect enables multi-selection mode
func (m *TableModel) EnableMultiSelect() {
	m.multiSelectMode = true
	m.Help = "↑/↓: Navigate • space: Select/Deselect • a: Select All • enter: Perform Action on Selected • q: Quit"
	m.updateSelectionTable()
}

// Update updates the table model
func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case " ":
			if m.multiSelectMode {
				m.ToggleRow()
			}
			return m, nil
		case "a":
			if m.multiSelectMode {
				m.SelectAll()
			}
			return m, nil
		case "enter":
			if m.multiSelectMode && len(m.selectedRows) > 0 && m.OnMultiSelect != nil {
				m.OnMultiSelect(m.GetSelectedRows())
				return m, nil
			} else if !m.multiSelectMode && m.OnSelect != nil && len(m.table.Rows()) > 0 {
				selectedRow := m.table.SelectedRow()
				m.OnSelect(selectedRow)
			}
			return m, nil
		}
	}

	// Update the main table
	m.table, cmd = m.table.Update(msg)

	// If in multi-select mode, also update the selection table's cursor
	if m.multiSelectMode {
		// Transfer cursor position from main table to selection table
		m.selectionTable.SetCursor(m.table.Cursor())
	}

	return m, cmd
}

// View renders the table model
func (m TableModel) View() string {
	result := titleStyle.Render(m.Title) + "\n\n"

	if m.multiSelectMode {
		// For multi-selection mode, show selection count
		if len(m.selectedRows) > 0 {
			result += lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Render(
				fmt.Sprintf("%d items selected", len(m.selectedRows))) + "\n\n"
		}

		// Update the selection table in case anything changed
		m.updateSelectionTable()

		// Use the selection table view
		result += m.selectionTable.View() + "\n\n"
	} else {
		result += m.table.View() + "\n\n"
	}

	result += helpStyle.Render(m.Help)
	return result
}
