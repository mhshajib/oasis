package command_options

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// Make an exportable function which takes these three fields as arguments by reference so that the values can be modified here and used in the main function.
func GenerateFields(fieldNames *[]string, fieldTypes *[]string, isFiltered *[]bool, numFields int) {
	// Initialize the slices
	*fieldNames = make([]string, numFields)
	*fieldTypes = make([]string, numFields)
	*isFiltered = make([]bool, numFields)

	dataTypes := []string{"string", "[]string", "int", "[]int", "int32", "[]int32", "int64", "[]int64", "float32", "[]float32", "float64", "[]float64", "bool", "[]bool", "byte", "[]byte", "rune", "[]rune", "time.Time"}
	for i := 0; i < numFields; i++ {
		fmt.Printf("Enter name for field #%d: ", i+1)
		var fieldName string
		fmt.Scan(&fieldName)
		(*fieldNames)[i] = strings.TrimSpace(fieldName)

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

		(*fieldNames)[i] = fieldType

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
			(*isFiltered)[i] = true
		} else {
			(*isFiltered)[i] = false
		}
	}
}
