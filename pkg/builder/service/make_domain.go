package builder

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"oasis/pkg/config"
	cli_template "oasis/pkg/template"
	"oasis/pkg/utils"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func domainFileExists(domainPath string, snakeCaseModuleName string) bool {
	domainFilePath := fmt.Sprintf("%s/%s.go", domainPath, snakeCaseModuleName)
	if _, err := os.Stat("/" + utils.NormalizePath(domainFilePath)); err == nil {
		return true
	} else {
		return false
	}
}

func parseDomainTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName string, numFields int, fieldNames []string, fieldTypes []string, isFiltered []bool) (string, error) {
	var fields []Field
	var criteriaFields []CriteriaField
	for i := 0; i < numFields; i++ {
		// Process the field names
		titleCaseFieldName, snakeCaseFieldName, _ := utils.ProcessString(fieldNames[i])

		// Append the field to the fields slice
		fields = append(fields, Field{
			Name:    titleCaseFieldName,
			Type:    fieldTypes[i],
			JsonTag: snakeCaseFieldName,
		})

		// Check if the field is filtered
		if isFiltered[i] {
			// Append the field to the criteria fields slice
			adjustedCriteriaType := adjustFieldTypeForCriteria(fieldTypes[i])
			criteriaFields = append(criteriaFields, CriteriaField{
				Name:    titleCaseFieldName,
				Type:    adjustedCriteriaType,
				JsonTag: snakeCaseFieldName,
			})
		}
	}
	// Prepare the data
	domainTemplateData := struct {
		UcFirstName     string
		SmallName       string
		SnakeCaseName   string
		SmallPluralName string
		Fields          []Field
		CriteriaFields  []CriteriaField
	}{
		UcFirstName:     titleCaseModuleName,
		SmallName:       camelCaseModuleName,
		SnakeCaseName:   snakeCaseModuleName,
		SmallPluralName: utils.ToPlural(snakeCaseModuleName),
		Fields:          fields,
		CriteriaFields:  criteriaFields,
	}

	// Read the contents of the file
	sourceContent := cli_template.Domain

	// Create a new template and parse the template string
	parsedTemplate, err := template.New("domainTemplate").Parse(string(sourceContent))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return "", err
	}

	var buf bytes.Buffer
	// Execute the template with the data
	err = parsedTemplate.Execute(&buf, domainTemplateData)
	if err != nil {
		fmt.Println("Error executing parsed template:", err)
		return "", err
	}

	return buf.String(), nil
}

func generateDomainFile(domainPath string, snakeCaseModuleName string, templateString string) error {

	domainFileName := fmt.Sprintf("%s/%s.go", domainPath, snakeCaseModuleName)
	// Write the code to the file
	err := ioutil.WriteFile("/"+utils.NormalizePath(domainFileName), []byte(templateString), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	// Execute the `go fmt` command
	goFmtCmd := exec.Command("go", "fmt", "/"+utils.NormalizePath(domainFileName))
	goFmtCmd.Stdout = os.Stdout
	goFmtCmd.Stderr = os.Stderr
	err = goFmtCmd.Run()
	if err != nil {
		fmt.Println("Error generating domain:", err)
		return err
	}
	return nil
}

func MakeDomain(cmd *cobra.Command, moduleName string, numFields int, fieldNames []string, fieldTypes []string, isFiltered []bool) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	rootPath, _, err := utils.GetGoModFile(cwd)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = utils.CreateDirectory(rootPath, config.Paths().DomainPath)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	MakeTimestamp(rootPath)

	titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName := utils.ProcessString(moduleName)

	if domainFileExists("/"+utils.NormalizePath(fmt.Sprintf("%s/%s", rootPath, config.Paths().DomainPath)), snakeCaseModuleName) {
		fmt.Println("Domain Already Exists With Name:", snakeCaseModuleName)
		return
	}

	templateString, err := parseDomainTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName, numFields, fieldNames, fieldTypes, isFiltered)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = generateDomainFile("/"+utils.NormalizePath(fmt.Sprintf("%s/%s", rootPath, config.Paths().DomainPath)), snakeCaseModuleName, templateString)
	if err != nil {
		fmt.Println("Error generating domain:", err)
		return
	}

	fmt.Println(fmt.Sprintf("Domain created and formatted successfully! with name: %s", snakeCaseModuleName))
}
