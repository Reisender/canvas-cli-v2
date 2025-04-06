package ui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TableModel represents a table UI model
type TableModel struct {
	table table.Model
	Title string
	Help  string
}

// NewTableModel creates a new table model
func NewTableModel(t table.Model) *TableModel {
	return &TableModel{
		table: t,
		Title: "Table",
		Help:  "↑/↓: Navigate • enter: Select • q: Quit",
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
)

// Init initializes the table model
func (m TableModel) Init() tea.Cmd {
	return nil
}

// Update updates the table model
func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			// TODO: Handle selection
			return m, nil
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View renders the table model
func (m TableModel) View() string {
	s := titleStyle.Render(m.Title) + "\n\n"
	s += m.table.View() + "\n\n"
	s += helpStyle.Render(m.Help)
	return s
}
