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
		Selected: "\U0001F336 {{ .Name | red | cyan }}",
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

	if flags[flagIndex].Name == "all" || flags[flagIndex].Name == "domain" {
		builder.MakeDomain(cmd, moduleName)
	}
	if flags[flagIndex].Name == "all" || flags[flagIndex].Name == "repo" {
		builder.MakeRepository(cmd, moduleName)
	}

	if flags[flagIndex].Name == "all" || flags[flagIndex].Name == "usecase" {
		builder.MakeUsecase(cmd, moduleName)
	}

	if flags[flagIndex].Name == "all" || flags[flagIndex].Name == "transform" {
		builder.MakeTransformer(cmd, moduleName)
	}

	if flags[flagIndex].Name == "all" || flags[flagIndex].Name == "delivery" {
		builder.MakeHttpHandler(cmd, moduleName)
	}

	if flags[flagIndex].Name == "all" || flags[flagIndex].Name == "migration" {
		builder.MakeMigration(cmd, moduleName)
	}

	if flags[flagIndex].Name == "all" || flags[flagIndex].Name == "seed" {
		builder.MakeSeeder(cmd, moduleName)
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
