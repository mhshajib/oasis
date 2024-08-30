package command_options

import (
	"errors"
	"strings"

	"github.com/manifoldco/promptui"
)

var packageName string
var ValidatePackageName = func(input string) error {
	if strings.Contains(input, " ") {
		return errors.New("Package name cannot contain spaces")
	}
	if len(input) < 3 {
		return errors.New("Package name must have more than 3 characters")
	}
	return nil
}
var PromptPackageName = promptui.Prompt{
	Label:    "Package Name",
	Validate: ValidatePackageName,
	Default:  packageName,
	Templates: &promptui.PromptTemplates{
		Valid:   "\u2192 {{ . | green }} ",
		Invalid: "{{ . | faint }} ",
		Success: "\u2713 {{ . | green }} ",
		Prompt:  "\u21AA {{ . | bold }} ",
	},
}
