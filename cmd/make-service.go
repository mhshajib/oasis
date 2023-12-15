package cmd

import (
	"clean-cli/pkg/builder"

	"github.com/spf13/cobra"
)

var (
	makeBlockCmd = &cobra.Command{
		Use:   "make:service",
		Short: "make:service create a new service for project",
		Long:  "make:service create a new service for project, including a domain, delivery (http), usecase & repository",
		Args:  cobra.ExactArgs(1),
		Run:   makeBlock,
	}
)
var allFlag, domainFlag, migrationFlag, seedFlag, transformFlag, useCaseFlag, repoFlag, deliveryFlag bool

func init() {
	makeBlockCmd.Flags().BoolVar(&allFlag, "all", false, "Create all components")
	makeBlockCmd.Flags().BoolVar(&domainFlag, "domain", false, "Create domain")
	rootCmd.AddCommand(makeBlockCmd)
}

func makeBlock(cmd *cobra.Command, args []string) {
	if allFlag || domainFlag {
		builder.MakeDomain(cmd, args)
	}
}
