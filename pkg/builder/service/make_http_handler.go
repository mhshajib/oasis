package builder

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"

	"oasis/pkg/config"
	cli_template "oasis/pkg/template"
	"oasis/pkg/utils"

	"github.com/spf13/cobra"
)

func httpHandlerFileExists(servicePath string, snakeCaseModuleName string) bool {
	httpHandlerFilePath := fmt.Sprintf("%s/%s/delivery/http/%s_handler.go", servicePath, snakeCaseModuleName, snakeCaseModuleName)
	if _, err := os.Stat("/" + utils.NormalizePath(httpHandlerFilePath)); err == nil {
		return true
	} else {
		return false
	}
}

func parseHttpHandlerTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName string, numFields int, fieldNames []string, fieldTypes []string, isFiltered []bool) (string, error) {
	var fields []Field
	var criteriaFields []CriteriaField
	for i := 0; i < numFields; i++ {
		// Process the field names
		titleCaseFieldName, snakeCaseFieldName, _ := utils.ProcessString(fieldNames[i])

		// Append the field to the fields slice
		fields = append(fields, Field{
			Name:               titleCaseFieldName,
			Type:               fieldTypes[i],
			JsonTag:            snakeCaseFieldName,
			OperatorLessThan:   template.HTML("<"),
			OperatorGreterThan: template.HTML(">"),
		})

		// Check if the field is filtered
		if isFiltered[i] {
			// Append the field to the criteria fields slice
			adjustedCriteriaType := adjustFieldTypeForCriteria(fieldTypes[i])
			criteriaFields = append(criteriaFields, CriteriaField{
				Name:               titleCaseFieldName,
				Type:               adjustedCriteriaType,
				JsonTag:            snakeCaseFieldName,
				OperatorLessThan:   template.HTML("<"),
				OperatorGreterThan: template.HTML(">"),
			})
		}
	}
	// Prepare the data
	httpHandlerTemplateData := struct {
		UcFirstName     string
		SmallName       string
		SnakeCaseName   string
		SmallPluralName string
		ModuleName      string
		DomainPath      string
		TransformerPath string
		Fields          []Field
		CriteriaFields  []CriteriaField
	}{
		UcFirstName:     titleCaseModuleName,
		SmallName:       camelCaseModuleName,
		SnakeCaseName:   snakeCaseModuleName,
		SmallPluralName: utils.ToPlural(snakeCaseModuleName),
		ModuleName:      config.Paths().ModuleName,
		DomainPath:      utils.NormalizePath(fmt.Sprintf("%s/%s", config.Paths().ModuleName, config.Paths().DomainPath)),
		TransformerPath: utils.NormalizePath(fmt.Sprintf("%s/%s/%s/transformer", config.Paths().ModuleName, config.Paths().ServicePath, camelCaseModuleName)),
		Fields:          fields,
		CriteriaFields:  criteriaFields,
	}

	// Read the contents of the file
	sourceContent := cli_template.HttpHandler

	// Create a new template and parse the template string
	parsedTemplate, err := template.New("httpHandlerTemplate").Parse(string(sourceContent))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return "", err
	}

	var buf bytes.Buffer
	// Execute the template with the data
	err = parsedTemplate.Execute(&buf, httpHandlerTemplateData)
	if err != nil {
		fmt.Println("Error executing parsed template:", err)
		return "", err
	}

	return buf.String(), nil
}

func generateHttpHandlerFile(servicePath string, snakeCaseModuleName string, templateString string) error {
	// Create the directory path
	directoryPath := fmt.Sprintf("%s/%s/delivery/http", servicePath, snakeCaseModuleName)
	err := os.MkdirAll("/"+utils.NormalizePath(directoryPath), os.ModePerm) // os.ModePerm is 0777
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

	// Create the file path
	httpHandlerFileName := fmt.Sprintf("%s/%s_handler.go", directoryPath, snakeCaseModuleName)

	// Write the code to the file
	err = ioutil.WriteFile("/"+utils.NormalizePath(httpHandlerFileName), []byte(templateString), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	// Execute the `go fmt` command
	goFmtCmd := exec.Command("go", "fmt", "/"+utils.NormalizePath(httpHandlerFileName))
	goFmtCmd.Stdout = os.Stdout
	goFmtCmd.Stderr = os.Stderr
	err = goFmtCmd.Run()
	if err != nil {
		fmt.Println("Error formatting http handler file:", err)
		return err
	}

	return nil
}

func MakeHttpHandler(cmd *cobra.Command, moduleName string, numFields int, fieldNames []string, fieldTypes []string, isFiltered []bool) {
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

	titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName := utils.ProcessString(moduleName)
	filePath := "/" + utils.NormalizePath(fmt.Sprintf("%s/%s", rootPath, config.Paths().ServicePath))
	if httpHandlerFileExists(filePath, snakeCaseModuleName) {
		fmt.Println("Http Handler Already Exists With Name:", snakeCaseModuleName)
		return
	}

	templateString, err := parseHttpHandlerTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName, numFields, fieldNames, fieldTypes, isFiltered)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = generateHttpHandlerFile(filePath, snakeCaseModuleName, templateString)
	if err != nil {
		fmt.Println("Error generating http handler:", err)
		return
	}

	fmt.Println(fmt.Sprintf("Http Handler created and formatted successfully! with name: %s", snakeCaseModuleName))
}
