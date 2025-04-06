package main

import (
	"fmt"
	"os"

	"github.com/Reisender/canvas-cli-v2/pkg/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
