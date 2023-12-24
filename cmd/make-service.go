package cmd

import (
	"fmt"
	"oasis/pkg/builder"
	"os/exec"
	"time"

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
	makeBlockCmd.Flags().BoolVar(&repoFlag, "repo", false, "Create use case")
	makeBlockCmd.Flags().BoolVar(&useCaseFlag, "usecase", false, "Create use case")
	makeBlockCmd.Flags().BoolVar(&transformFlag, "transform", false, "Create transformer")
	rootCmd.AddCommand(makeBlockCmd)
}

func makeBlock(cmd *cobra.Command, args []string) {
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

	// Start the animation in a separate goroutine
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\rExecuting 'go mod tidy' %s", animationFrame())
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

// animationFrame returns a string representing the current frame of the animation
func animationFrame() string {
	frames := []string{"-", "\\", "|", "/"}
	return frames[time.Now().UnixMilli()/100%int64(len(frames))]
}
