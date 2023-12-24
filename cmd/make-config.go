package cmd

import (
	"fmt"
	builder "oasis/pkg/builder/config"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var (
	makeConfigCmd = &cobra.Command{
		Use:   "make:config",
		Short: "make:config create a new config for project",
		Long:  "make:config create a new config for project",
		Args:  cobra.ExactArgs(1),
		Run:   makeConfig,
	}
)

func init() {
	rootCmd.AddCommand(makeConfigCmd)
}

func makeConfig(cmd *cobra.Command, args []string) {
	builder.MakeConfig(cmd, args)

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
