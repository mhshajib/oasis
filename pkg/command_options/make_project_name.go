package command_options

import (
	"errors"
	"strings"

	"github.com/manifoldco/promptui"
)

var projectName string
var ValidateProjectName = func(input string) error {
	if strings.Contains(input, " ") {
		return errors.New("Project name cannot contain spaces")
	}
	if len(input) < 3 {
		return errors.New("Project name must have more than 3 characters")
	}
	return nil
}
var PromptProjectName = promptui.Prompt{
	Label:    "Project Name",
	Validate: ValidateProjectName,
	Default:  projectName,
	Templates: &promptui.PromptTemplates{
		Valid:   "\u2192 {{ . | green }} ",
		Invalid: "{{ . | faint }} ",
		Success: "\u2713 {{ . | green }} ",
		Prompt:  "\u21AA {{ . | bold }} ",
	},
}
