package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var (
	makeProjectCmd = &cobra.Command{
		Use:   "make:project",
		Short: "make:project initializes a new project boilerplate on current directory",
		Long:  "make:project initializes a new project boilerplate on current directory",
		Args:  cobra.ExactArgs(1),
		Run:   makeProject,
	}
)

// var allFlag, domainFlag, migrationFlag, seedFlag, transformFlag, useCaseFlag, repoFlag, deliveryFlag bool

func init() {
	rootCmd.AddCommand(makeProjectCmd)
}

func makeProject(cmd *cobra.Command, args []string) {
	projectName := args[0]
	fmt.Printf("Project name: %s\n", projectName)

	// Clone the repository
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Start the animation in a separate goroutine
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\rCreating new Project: %s | %s", projectName, AnimationFrame())
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	cloneCmd := exec.Command("git", "clone", "https://github.com/rtyley/small-test-repo", projectName)
	cloneCmd.Dir = cwd
	err = cloneCmd.Run()
	if err != nil {
		fmt.Println("Failed To Initialize Project :", err)
		done <- true
		return
	}

	done <- true
	fmt.Println("\rProject Initialized Successfully.")
}
