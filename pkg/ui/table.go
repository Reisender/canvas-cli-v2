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
	baseRows        []table.Row    // Original rows without selection indicators
	baseColumns     []table.Column // Original columns without selection column
	Title           string
	Help            string
	OnSelect        SelectionCallback
	OnMultiSelect   MultiSelectionCallback
	selectedRows    map[int]bool
	multiSelectMode bool
}

// NewTableModel creates a new table model
func NewTableModel(t table.Model) *TableModel {
	// Store original rows and columns
	baseRows := make([]table.Row, len(t.Rows()))
	copy(baseRows, t.Rows())

	baseColumns := make([]table.Column, len(t.Columns()))
	copy(baseColumns, t.Columns())

	return &TableModel{
		table:           t,
		baseRows:        baseRows,
		baseColumns:     baseColumns,
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

	// Update the table rows to reflect selection changes
	if m.multiSelectMode {
		m.updateTableWithSelectionIndicators()
	}
}

// updateTableWithSelectionIndicators updates the main table to show selection indicators
func (m *TableModel) updateTableWithSelectionIndicators() {
	// Keep track of the current cursor position
	cursorPos := m.table.Cursor()

	// Get current table dimensions
	height := m.table.Height()
	height = 25

	// Create new rows with checkmarks
	newRows := make([]table.Row, len(m.baseRows))
	for i, row := range m.baseRows {
		// If selected, add a checkmark as the first element
		indicator := ""
		if m.IsRowSelected(i) {
			indicator = "✓"
		}

		// Create a new row with the selection indicator plus all original data
		newRow := make(table.Row, len(row)+1)
		newRow[0] = indicator
		for j, cell := range row {
			newRow[j+1] = cell
		}
		newRows[i] = newRow
	}

	// Create a columns slice with selection column
	columns := []table.Column{
		{Title: "", Width: 2},
	}
	columns = append(columns, m.baseColumns...)

	// Create a new table with the updated data but preserving other settings
	newTable := table.New(
		table.WithColumns(columns),
		table.WithRows(newRows),
		table.WithFocused(true),
		table.WithHeight(height),
	)

	// Apply default styles since we can't access the existing styles directly
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

	// Replace the existing table
	m.table = newTable
}

// GetSelectedRows returns all selected rows
func (m TableModel) GetSelectedRows() []table.Row {
	var selected []table.Row

	for i, row := range m.baseRows {
		if m.selectedRows[i] {
			selected = append(selected, row)
		}
	}

	return selected
}

// SelectAll selects all rows
func (m *TableModel) SelectAll() {
	for i := range m.baseRows {
		m.selectedRows[i] = true
	}

	// Update the table rows to reflect selection changes
	if m.multiSelectMode {
		m.updateTableWithSelectionIndicators()
	}
}

// ClearSelections clears all selected rows
func (m *TableModel) ClearSelections() {
	m.selectedRows = make(map[int]bool)

	// Update the table rows to reflect selection changes
	if m.multiSelectMode {
		m.updateTableWithSelectionIndicators()
	}
}

// EnableMultiSelect enables multi-selection mode
func (m *TableModel) EnableMultiSelect() {
	m.multiSelectMode = true
	m.Help = "↑/↓: Navigate • space: Select/Deselect • a: Select All • enter: Perform Action on Selected • q: Quit"

	// Create a fixed-height table
	fixedHeight := 25 // Use a consistent height value
	newTable := table.New(
		table.WithColumns(m.table.Columns()),
		table.WithRows(m.table.Rows()),
		table.WithFocused(true),
		table.WithHeight(fixedHeight),
	)

	// Copy styles
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

	// Preserve cursor position
	newTable.SetCursor(m.table.Cursor())

	// Update the main table
	m.table = newTable

	// Then add selection indicators
	m.updateTableWithSelectionIndicators()
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
				// Return only the original row data without selection indicators
				m.OnMultiSelect(m.GetSelectedRows())
				return m, nil
			} else if !m.multiSelectMode && m.OnSelect != nil && len(m.table.Rows()) > 0 {
				selectedRow := m.table.SelectedRow()
				// For single selection, return the raw selected row
				m.OnSelect(selectedRow)
			}
			return m, nil
		}
	}

	// Update the main table
	m.table, cmd = m.table.Update(msg)

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

		// The table already has selection indicators from updateTableWithSelectionIndicators
		result += m.table.View() + "\n\n"
	} else {
		result += m.table.View() + "\n\n"
	}

	result += helpStyle.Render(m.Help)
	return result
}
