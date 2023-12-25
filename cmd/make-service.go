package cmd

import (
	"fmt"
	builder "oasis/pkg/builder/service"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var (
	makeServiceCmd = &cobra.Command{
		Use:   "make:service",
		Short: "make:service create a new service for project",
		Long:  "make:service create a new service for project, including a domain, delivery (http), usecase & repository",
		Args:  cobra.ExactArgs(1),
		Run:   makeService,
	}
)
var allFlag, domainFlag, migrationFlag, seedFlag, transformFlag, useCaseFlag, repoFlag, deliveryFlag bool

func init() {
	makeServiceCmd.Flags().BoolVar(&allFlag, "all", false, "Create all components")
	makeServiceCmd.Flags().BoolVar(&domainFlag, "domain", false, "Create domain")
	makeServiceCmd.Flags().BoolVar(&repoFlag, "repo", false, "Create use case")
	makeServiceCmd.Flags().BoolVar(&useCaseFlag, "usecase", false, "Create use case")
	makeServiceCmd.Flags().BoolVar(&transformFlag, "transform", false, "Create transformer")
	makeServiceCmd.Flags().BoolVar(&deliveryFlag, "delivery", false, "Create handler")
	makeServiceCmd.Flags().BoolVar(&migrationFlag, "migration", false, "Create migration")
	rootCmd.AddCommand(makeServiceCmd)
}

func makeService(cmd *cobra.Command, args []string) {
	if allFlag || domainFlag {
		builder.MakeDomain(cmd, args)
	}
	if allFlag || repoFlag {
		builder.MakeRepository(cmd, args)
	}

	if allFlag || useCaseFlag {
		builder.MakeUsecase(cmd, args)
	}

	if allFlag || transformFlag {
		builder.MakeTransformer(cmd, args)
	}

	if allFlag || deliveryFlag {
		builder.MakeHttpHandler(cmd, args)
	}

	if allFlag || migrationFlag {
		builder.MakeMigration(cmd, args)
	}

	// Start the animation in a separate goroutine
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\rExecuting 'go mod tidy' %s", AnimationFrame())
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Execute 'go mod tidy'
	execCmd := exec.Command("go", "mod", "tidy")
	err := execCmd.Run()
	if err != nil {
		fmt.Println("\nError running 'go mod tidy':", err)
		done <- true
		return
	}

	// Stop the animation
	done <- true
	fmt.Println("\r'go mod tidy' executed successfully.")
}

// AnimationFrame returns a string representing the current frame of the animation
func AnimationFrame() string {
	frames := []string{"-", "\\", "|", "/"}
	return frames[time.Now().UnixMilli()/100%int64(len(frames))]
}
