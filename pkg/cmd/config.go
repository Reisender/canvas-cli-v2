package cmd

import (
	"fmt"

	"github.com/Reisender/canvas-cli-v2/pkg/config"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// ConfigModel represents the config UI model
type ConfigModel struct {
	inputs     []textinput.Model
	focusIndex int
	done       bool
	err        error
	title      string
}

// NewConfigCmd creates a new command for managing configuration
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure Canvas CLI",
		Long:  `Configure your Canvas API key and other settings.`,
		Run:   runConfig,
	}

	// Add subcommands
	cmd.AddCommand(
		newConfigGetCmd(),
		newConfigSetCmd(),
	)

	return cmd
}

func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get Canvas CLI configuration",
		Long:  `Display the current Canvas CLI configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.GetConfig()
			fmt.Println("Current Configuration:")
			fmt.Println("---------------------")
			fmt.Printf("Base URL: %s\n", cfg.BaseURL)

			// Mask API key for security
			apiKey := "[not set]"
			if cfg.APIKey != "" {
				apiKey = "[set]"
			}
			fmt.Printf("API Key: %s\n", apiKey)
		},
	}
}

func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set Canvas CLI configuration",
		Long:  `Set a configuration value for Canvas CLI.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			value := args[1]

			if err := config.UpdateConfig(key, value); err != nil {
				fmt.Printf("Error updating config: %v\n", err)
				return
			}

			fmt.Printf("Successfully updated %s\n", key)
		},
	}
}

func runConfig(cmd *cobra.Command, args []string) {
	cfg := config.GetConfig()

	// Initialize text inputs
	baseURLInput := textinput.New()
	baseURLInput.Placeholder = "https://canvas.instructure.com/api/v1"
	baseURLInput.Focus()
	baseURLInput.Width = 60
	baseURLInput.Prompt = "› "
	baseURLInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("62"))
	baseURLInput.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	baseURLInput.SetValue(cfg.BaseURL)
	baseURLInput.CharLimit = 150

	apiKeyInput := textinput.New()
	apiKeyInput.Placeholder = "your-api-key"
	apiKeyInput.Width = 60
	apiKeyInput.Prompt = "› "
	apiKeyInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("62"))
	apiKeyInput.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	apiKeyInput.SetValue(cfg.APIKey)
	apiKeyInput.CharLimit = 100
	apiKeyInput.EchoMode = textinput.EchoPassword
	apiKeyInput.EchoCharacter = '•'

	inputs := []textinput.Model{baseURLInput, apiKeyInput}

	model := ConfigModel{
		inputs:     inputs,
		focusIndex: 0,
		title:      "Canvas CLI Configuration",
	}

	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Printf("Error running config: %v\n", err)
	}
}

// Init initializes the config model
func (m ConfigModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update updates the config model
func (m ConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "shift+tab", "up", "down":
			// Cycle focus between inputs
			if msg.String() == "up" || msg.String() == "shift+tab" {
				m.focusIndex--
				if m.focusIndex < 0 {
					m.focusIndex = len(m.inputs) - 1
				}
			} else {
				m.focusIndex++
				if m.focusIndex >= len(m.inputs) {
					m.focusIndex = 0
				}
			}

			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					cmds = append(cmds, m.inputs[i].Focus())
				} else {
					m.inputs[i].Blur()
				}
			}

			return m, tea.Batch(cmds...)
		case "enter":
			if m.focusIndex == len(m.inputs)-1 {
				// Save config
				err := config.UpdateConfig("base_url", m.inputs[0].Value())
				if err != nil {
					m.err = err
					return m, nil
				}

				err = config.UpdateConfig("api_key", m.inputs[1].Value())
				if err != nil {
					m.err = err
					return m, nil
				}

				m.done = true
				return m, tea.Quit
			}

			// Move to next input
			m.focusIndex++
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					cmds = append(cmds, m.inputs[i].Focus())
				} else {
					m.inputs[i].Blur()
				}
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Update each textinput
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *ConfigModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	for i := range m.inputs {
		var cmd tea.Cmd
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}

// View renders the config model
func (m ConfigModel) View() string {
	if m.done {
		return "Configuration saved successfully!\n"
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		MarginLeft(2)

	s := titleStyle.Render(m.title) + "\n\n"

	s += "Base URL:" + "\n"
	s += m.inputs[0].View() + "\n\n"

	s += "API Key:" + "\n"
	s += m.inputs[1].View() + "\n\n"

	s += "Press Enter to save, Esc to cancel\n"

	return s
}
