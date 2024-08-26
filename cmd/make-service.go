package cmd

import (
	"errors"
	"fmt"
	builder "oasis/pkg/builder/service"
	"os/exec"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	makeServiceCmd = &cobra.Command{
		Use:   "make:service",
		Short: "make:service create a new service for project",
		Long:  "make:service create a new service for project, including a domain, delivery (http), usecase & repository",
		//Args:  cobra.ExactArgs(1),
		Run: makeService,
	}
)

type flag struct {
	Name  string
	Usage string
}

var allFlag, domainFlag, migrationFlag, seedFlag, transformFlag, useCaseFlag, repoFlag, deliveryFlag bool

func init() {
	//makeServiceCmd.Flags().BoolVar(&allFlag, "all", false, "Create all components")
	//makeServiceCmd.Flags().BoolVar(&domainFlag, "domain", false, "Create domain")
	//makeServiceCmd.Flags().BoolVar(&repoFlag, "repo", false, "Create use case")
	//makeServiceCmd.Flags().BoolVar(&useCaseFlag, "usecase", false, "Create use case")
	//makeServiceCmd.Flags().BoolVar(&transformFlag, "transform", false, "Create transformer")
	//makeServiceCmd.Flags().BoolVar(&deliveryFlag, "delivery", false, "Create handler")
	//makeServiceCmd.Flags().BoolVar(&migrationFlag, "migration", false, "Create migration")
	//makeServiceCmd.Flags().BoolVar(&seedFlag, "seed", false, "Create seeder")
	rootCmd.AddCommand(makeServiceCmd)
}

func makeService(cmd *cobra.Command, args []string) {
	flags := []flag{
		{Name: "all", Usage: "Create all components"},
		{Name: "domain", Usage: "Create domain"},
		{Name: "repo", Usage: "Create repository"},
		{Name: "usecase", Usage: "Create use case"},
		{Name: "transform", Usage: "Create response transformer"},
		{Name: "delivery", Usage: "Create handler for http and grpc"},
		{Name: "migration", Usage: "Create migration"},
		{Name: "seed", Usage: "Create seeder"},
	}
	flagTemplates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\u21AA {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\u21AA {{ .Name | red | cyan }}",
		Details: `
		--------- Flags ----------
		{{ "Name:" | faint }}	{{ .Name }}
		{{ "Usage:" | faint }}	{{ .Usage }}`,
	}
	flagSearcher := func(input string, index int) bool {
		flag := flags[index]
		name := strings.Replace(strings.ToLower(flag.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	flagPrompt := promptui.Select{
		Label:     "Select flags to create",
		Items:     flags,
		Templates: flagTemplates,
		Size:      4,
		Searcher:  flagSearcher,
	}

	flagIndex, _, err := flagPrompt.Run()

	var moduleName string
	validateModuleName := func(input string) error {
		if len(input) < 3 {
			return errors.New("Username must have more than 3 characters")
		}
		return nil
	}
	promptModuleName := promptui.Prompt{
		Label:    "Module Name",
		Validate: validateModuleName,
		Default:  moduleName,
	}

	moduleName, err = promptModuleName.Run()

	fmt.Print("Enter the number of fields: ")
	var numFields int
	fmt.Scan(&numFields)
	fieldNames := make([]string, numFields)
	fieldTypes := make([]string, numFields)
	isFiltered := make([]bool, numFields)
	dataTypes := []string{"string", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64", "complex64", "complex128", "bool", "byte", "rune", "time.Time", "*Type", "[]Type", "[n]Type", "map[KeyType]ValueType", "chan Type", "func(Type1, Type2) ReturnType", "error"}
	for i := 0; i < numFields; i++ {
		fmt.Printf("Enter name for field #%d: ", i+1)
		var fieldName string
		fmt.Scan(&fieldName)
		fieldNames[i] = strings.TrimSpace(fieldName)

		// Create a promptui selector for data types
		prompt := promptui.Select{
			Label: "Select Data Type",
			Items: dataTypes,
		}

		_, fieldType, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fieldTypes[i] = fieldType

		addFiltersPrompt := promptui.Select{
			Label: "Add filterable fields within this field?",
			Items: []string{"Yes", "No"},
		}
		_, addFilters, err := addFiltersPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		if addFilters == "Yes" {
			isFiltered[i] = true
		} else {
			isFiltered[i] = false
		}
	}

	if flags[flagIndex].Name == "all" || flags[flagIndex].Name == "domain" {
		builder.MakeDomain(cmd, moduleName, numFields, fieldNames, fieldTypes, isFiltered)
	}
	if flags[flagIndex].Name == "all" || flags[flagIndex].Name == "repo" {
		builder.MakeRepository(cmd, moduleName, numFields, fieldNames, fieldTypes)
	}

	if flags[flagIndex].Name == "all" || flags[flagIndex].Name == "usecase" {
		builder.MakeUsecase(cmd, moduleName, numFields, fieldNames, fieldTypes)
	}

	if flags[flagIndex].Name == "all" || flags[flagIndex].Name == "transform" {
		builder.MakeTransformer(cmd, moduleName, numFields, fieldNames, fieldTypes)
	}

	if flags[flagIndex].Name == "all" || flags[flagIndex].Name == "delivery" {
		builder.MakeHttpHandler(cmd, moduleName, numFields, fieldNames, fieldTypes)
	}

	if flags[flagIndex].Name == "all" || flags[flagIndex].Name == "migration" {
		builder.MakeMigration(cmd, moduleName, numFields, fieldNames, fieldTypes)
	}

	if flags[flagIndex].Name == "all" || flags[flagIndex].Name == "seed" {
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
