package cmd

import (
	"github.com/Reisender/canvas-cli-v2/pkg/config"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "canvas-cli",
		Short: "A CLI for interacting with the Canvas LMS API",
		Long: `Canvas CLI is a command line interface for interacting with the Canvas LMS API.
It provides commands for managing courses, assignments, grades, and more.
Built with Charm libraries for a delightful terminal experience.`,
	}

	// Initialize config
	config.InitConfig()

	// Add commands
	rootCmd.AddCommand(
		NewCoursesCmd(),
		NewAssignmentsCmd(),
		NewUsersCmd(),
		NewConfigCmd(),
	)

	return rootCmd
}
