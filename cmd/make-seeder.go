package cmd

import (
	"fmt"
	builder "oasis/pkg/builder/service"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var (
	makeSeederCmd = &cobra.Command{
		Use:   "make:seed",
		Short: "make:seed create a new seeder for project",
		Long:  "make:seed create a new seeder for project",
		Args:  cobra.ExactArgs(1),
		Run:   makeSeeder,
	}
)

func init() {
	rootCmd.AddCommand(makeSeederCmd)
}

func makeSeeder(cmd *cobra.Command, args []string) {
	builder.MakeSeeder(cmd, args)

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
