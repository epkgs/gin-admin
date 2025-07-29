package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// This function creates a CLI command that prints the version number.
func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			// Version will be set by main.go
			fmt.Println(cmd.Root().Version)
		},
	}
}
