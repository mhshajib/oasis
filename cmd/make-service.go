package cmd

import (
	"fmt"
	builder "oasis/pkg/builder/service"
	"oasis/pkg/command_options"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var (
	makeServiceCmd = &cobra.Command{
		Use:   "make:service",
		Short: "make:service create a new service for project",
		Long:  "make:service create a new service for project, including a domain, delivery (http), usecase & repository",
		Run:   makeService,
	}
)

// var allFlag, domainFlag, migrationFlag, seedFlag, transformFlag, useCaseFlag, repoFlag, deliveryFlag bool

func init() {
	rootCmd.AddCommand(makeServiceCmd)
}

func makeService(cmd *cobra.Command, args []string) {

	// Create a promptui selector for flags
	flagIndex, _, err := command_options.FlagPrompt.Run()

	// Create a promptui prompt for module name
	moduleName, err := command_options.PromptModuleName.Run()

	fmt.Print("Enter the number of fields: ")
	var numFields int
	fmt.Scan(&numFields)
	if numFields <= 0 {
		fmt.Println("Number of fields must be greater than 0")
		return

	}
	fieldNames := make([]string, numFields)
	fieldTypes := make([]string, numFields)
	isFiltered := make([]bool, numFields)

	// Generate fields
	command_options.GenerateFields(&fieldNames, &fieldTypes, &isFiltered, numFields)

	if command_options.Flags[flagIndex].Name == "all" || command_options.Flags[flagIndex].Name == "domain" {
		builder.MakeDomain(cmd, moduleName, numFields, fieldNames, fieldTypes, isFiltered)
	}
	if command_options.Flags[flagIndex].Name == "all" || command_options.Flags[flagIndex].Name == "repo" {
		builder.MakeRepository(cmd, moduleName, numFields, fieldNames, fieldTypes, isFiltered)
	}

	if command_options.Flags[flagIndex].Name == "all" || command_options.Flags[flagIndex].Name == "usecase" {
		builder.MakeUsecase(cmd, moduleName, numFields, fieldNames, fieldTypes)
	}

	if command_options.Flags[flagIndex].Name == "all" || command_options.Flags[flagIndex].Name == "transform" {
		builder.MakeTransformer(cmd, moduleName, numFields, fieldNames, fieldTypes)
	}

	if command_options.Flags[flagIndex].Name == "all" || command_options.Flags[flagIndex].Name == "delivery" {
		builder.MakeHttpHandler(cmd, moduleName, numFields, fieldNames, fieldTypes, isFiltered)
	}

	if command_options.Flags[flagIndex].Name == "all" || command_options.Flags[flagIndex].Name == "migration" {
		builder.MakeMigration(cmd, moduleName, numFields, fieldNames, fieldTypes)
	}

	if command_options.Flags[flagIndex].Name == "all" || command_options.Flags[flagIndex].Name == "seed" {
		builder.MakeSeeder(cmd, moduleName, numFields, fieldNames, fieldTypes)
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
	err = execCmd.Run()
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
