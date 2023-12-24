package builder

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"oasis/pkg/config"
	cli_template "oasis/pkg/template"
	"oasis/pkg/utils"

	"github.com/spf13/cobra"
)

func transformerFileExists(servicePath string, snakeCaseModuleName string) bool {
	transformerFilePath := fmt.Sprintf("%s/%s/transformer/%s_transformer.go", servicePath, snakeCaseModuleName, snakeCaseModuleName)
	if _, err := os.Stat(transformerFilePath); err == nil {
		return true
	} else {
		return false
	}
}

func parseTransformerTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName string) (string, error) {
	// Prepare the data
	transformerTemplateData := struct {
		UcFirstName     string
		SmallName       string
		SnakeCaseName   string
		SmallPluralName string
		ModuleName      string
	}{
		UcFirstName:     titleCaseModuleName,
		SmallName:       camelCaseModuleName,
		SnakeCaseName:   snakeCaseModuleName,
		SmallPluralName: utils.ToPlural(snakeCaseModuleName),
		ModuleName:      config.Paths().ModuleName,
	}

	// Read the contents of the template
	sourceContent := cli_template.Transformer
	// Create a new template and parse the template string
	parsedTemplate, err := template.New("transformerTemplate").Parse(string(sourceContent))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return "", err
	}

	var buf bytes.Buffer
	// Execute the template with the data
	err = parsedTemplate.Execute(&buf, transformerTemplateData)
	if err != nil {
		fmt.Println("Error executing parsed template:", err)
		return "", err
	}

	return buf.String(), nil
}

func generateTransformerFile(servicePath string, snakeCaseModuleName string, templateString string) error {
	// Create the directory path
	directoryPath := fmt.Sprintf("%s/%s/transformer", servicePath, snakeCaseModuleName)
	err := os.MkdirAll(directoryPath, os.ModePerm) // os.ModePerm is 0777
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

	// Create the file path
	transformerFileName := fmt.Sprintf("%s/%s_transformer.go", directoryPath, snakeCaseModuleName)

	// Write the code to the file
	err = ioutil.WriteFile(transformerFileName, []byte(templateString), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	// Execute the `go fmt` command
	goFmtCmd := exec.Command("go", "fmt", transformerFileName)
	goFmtCmd.Stdout = os.Stdout
	goFmtCmd.Stderr = os.Stderr
	err = goFmtCmd.Run()
	if err != nil {
		fmt.Println("Error formatting transformer file:", err)
		return err
	}

	return nil
}

func MakeTransformer(cmd *cobra.Command, args []string) {
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

	moduleName := args[0]
	titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName := utils.ProcessString(moduleName)

	if transformerFileExists(fmt.Sprintf("%s/%s", rootPath, config.Paths().ServicePath), snakeCaseModuleName) {
		fmt.Println("Usecase Already Exists With Name:", strings.ToLower(moduleName))
		return
	}

	templateString, err := parseUsecaseTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = generateTransformerFile(fmt.Sprintf("%s/%s", rootPath, config.Paths().ServicePath), snakeCaseModuleName, templateString)
	if err != nil {
		fmt.Println("Error generating transformer:", err)
		return
	}

	fmt.Println(fmt.Sprintf("Usecase created and formatted successfully! with name: %s", snakeCaseModuleName))
}