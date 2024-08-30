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

	err = filepath.WalkDir(projectName, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip system and hidden files or directories
		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		if shouldSkip(path) {
			return nil
		}

		if !d.IsDir() {
			fmt.Printf("Processing file: %s\n", path)
			err := replaceInFile(path, d, projectName)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func shouldSkip(path string) bool {
	// Add more conditions if needed
	return strings.Contains(path, ".git")
}

func replaceInFile(path string, d os.DirEntry, projectName string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	updatedContent := strings.ReplaceAll(string(content), "github.com/mhshajib/oasis_boilerplate", fmt.Sprintf("%s", projectName))
	updatedContent = strings.ReplaceAll(updatedContent, "projectName", projectName)

	err = ioutil.WriteFile(path, []byte(updatedContent), d.Type())
	if err != nil {
		return err
	}

	if d.Name() == "config.develop.yml" {
		newPath := strings.ReplaceAll(path, "config.develop.yml", "config.yml")
		err := os.Rename(path, newPath)
		if err != nil {
			return err
		}
	}

	return nil
}
