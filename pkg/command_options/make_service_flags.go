package command_options

import (
	"strings"

	"github.com/manifoldco/promptui"
)

type flag struct {
	Name  string
	Usage string
}

var Flags = []flag{
	{Name: "all", Usage: "Create all components"},
	{Name: "domain", Usage: "Create domain"},
	{Name: "repo", Usage: "Create repository"},
	{Name: "usecase", Usage: "Create use case"},
	{Name: "transform", Usage: "Create response transformer"},
	{Name: "delivery", Usage: "Create handler for http and grpc"},
	{Name: "migration", Usage: "Create migration"},
	{Name: "seed", Usage: "Create seeder"},
}
var FlagTemplates = &promptui.SelectTemplates{
	Label:    "{{ . }}?",
	Active:   "\u21AA {{ .Name | cyan }}",
	Inactive: "  {{ .Name | faint }}",
	Selected: "\u2713 {{ .Name | green | cyan }}",
	Details: `
	--------- Flags ----------
	{{ "Name:" | faint }}	{{ .Name }}
	{{ "Usage:" | faint }}	{{ .Usage }}`,
}
var FlagSearcher = func(input string, index int) bool {
	flag := Flags[index]
	name := strings.Replace(strings.ToLower(flag.Name), " ", "", -1)
	input = strings.Replace(strings.ToLower(input), " ", "", -1)

	return strings.Contains(name, input)
}

var FlagPrompt = promptui.Select{
	Label:     "Select flags to create",
	Items:     Flags,
	Templates: FlagTemplates,
	Size:      4,
	Searcher:  FlagSearcher,
}
