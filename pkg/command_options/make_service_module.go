package command_options

import (
	"errors"

	"github.com/manifoldco/promptui"
)

var moduleName string
var ValidateModuleName = func(input string) error {
	if len(input) < 3 {
		return errors.New("Username must have more than 3 characters")
	}
	return nil
}
var PromptModuleName = promptui.Prompt{
	Label:    "Module Name",
	Validate: ValidateModuleName,
	Default:  moduleName,
	Templates: &promptui.PromptTemplates{
		Valid:   "\u2192 {{ . | green }} ",
		Invalid: "{{ . | faint }} ",
		Success: "\u2713 {{ . | green }} ",
		Prompt:  "\u21AA {{ . | bold }} ",
	},
}
