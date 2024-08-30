package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	cloneCmd := exec.Command("git", "clone", "https://github.com/mhshajib/oasis_boilerplate", projectName)
	cloneCmd.Dir = cwd
	err = cloneCmd.Run()
	if err != nil {
		fmt.Println("Failed To Initialize Project :", err)
		done <- true
		return
	}

	done <- true
	fmt.Println("\rProject Initialized Successfully.")

	// Replace the text in all files
	err = filepath.Walk(projectName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Read the file content
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// Replace the text
			newContent := strings.ReplaceAll(string(content), "github.com/mhshajib/oasis_boilerplate", fmt.Sprintf("%s/", projectName))
			finalContent := strings.ReplaceAll(string(newContent), "projectName", projectName)

			// Write the updated content back to the file
			err = ioutil.WriteFile(path, []byte(finalContent), info.Mode())
			if err != nil {
				return err
			}

			// copy config.develop.yml to config.yml
			if strings.Contains(path, "config.develop.yml") {
				destination := strings.ReplaceAll(path, "config.develop.yml", "config.yml")
				err = ioutil.WriteFile(destination, []byte(finalContent), info.Mode())
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Failed to replace text in files:", err)
		return
	}
}
