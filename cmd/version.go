package cmd

import (
	"fmt"

	"clean-cli/pkg/config"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version print the current build number and version information",
		Long:  `Version print the current build number and version information`,
		Run:   version,
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func version(cmd *cobra.Command, args []string) {
	fmt.Printf("Version: %s\n", config.Version)
}
