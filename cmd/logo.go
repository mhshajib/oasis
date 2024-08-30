package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	logoCmd = &cobra.Command{
		Use:   "logo",
		Short: "logo prints oasis logo",
		Long:  "logo prints oasis logo",
		Run:   printLogo,
	}
)

func init() {
	rootCmd.AddCommand(logoCmd)
}

func printLogo(cmd *cobra.Command, args []string) {
	fmt.Printf(logo)
}
