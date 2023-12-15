package utils

import (
	pluralize "github.com/gertd/go-pluralize"
)

func ToPlural(singular string) string {
	pluralize := pluralize.NewClient()
	if pluralize.IsPlural(singular) && !pluralize.IsSingular(singular) {
		return singular
	} else {
		return pluralize.Plural(singular)
	}
}
