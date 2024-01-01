package builder

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"oasis/pkg/config"
	cli_template "oasis/pkg/template"
	"os"
	"os/exec"
	"strings"

	"oasis/pkg/utils"

	"github.com/spf13/cobra"
)

func usecaseFileExists(servicePath string, snakeCaseModuleName string) bool {
	usecaseFilePath := fmt.Sprintf("%s/%s/usecase/%s_usecase.go", servicePath, snakeCaseModuleName, snakeCaseModuleName)
	if _, err := os.Stat(utils.NormalizePath(usecaseFilePath)); err == nil {
		return true
	} else {
		return false
	}
}

func parseUsecaseTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName string) (string, error) {
	usecaseTemplateData := struct {
		UcFirstName     string
		SmallName       string
		SnakeCaseName   string
		SmallPluralName string
		ModuleName      string
		DomainPath      string
	}{
		UcFirstName:     titleCaseModuleName,
		SmallName:       camelCaseModuleName,
		SnakeCaseName:   snakeCaseModuleName,
		SmallPluralName: utils.ToPlural(snakeCaseModuleName),
		ModuleName:      config.Paths().ModuleName,
		DomainPath:      utils.NormalizePath(fmt.Sprintf("%s/%s", config.Paths().ModuleName, config.Paths().DomainPath)),
	}

	// Read the contents of the template
	sourceContent := cli_template.UseCase

	// Create a new template and parse the template string
	parsedTemplate, err := template.New("usecaseTemplate").Parse(string(sourceContent))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return "", err
	}

	var buf bytes.Buffer
	// Execute the template with the data
	err = parsedTemplate.Execute(&buf, usecaseTemplateData)
	if err != nil {
		fmt.Println("Error executing parsed template:", err)
		return "", err
	}

	return buf.String(), nil
}

func generateUsecaseFile(servicePath string, snakeCaseModuleName string, templateString string) error {
	// Create the directory path
	directoryPath := fmt.Sprintf("%s/%s/usecase", servicePath, snakeCaseModuleName)
	err := os.MkdirAll(utils.NormalizePath(directoryPath), os.ModePerm) // os.ModePerm is 0777
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

	// Create the file path
	usecaseFileName := fmt.Sprintf("%s/%s_usecase.go", directoryPath, snakeCaseModuleName)

	// Write the code to the file
	err = ioutil.WriteFile(utils.NormalizePath(usecaseFileName), []byte(templateString), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	// Execute the `go fmt` command
	goFmtCmd := exec.Command("go", "fmt", utils.NormalizePath(usecaseFileName))
	goFmtCmd.Stdout = os.Stdout
	goFmtCmd.Stderr = os.Stderr
	err = goFmtCmd.Run()
	if err != nil {
		fmt.Println("Error formatting usecase file:", err)
		return err
	}

	return nil
}

func MakeUsecase(cmd *cobra.Command, args []string) {
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

	if usecaseFileExists(utils.NormalizePath(fmt.Sprintf("%s/%s", rootPath, config.Paths().ServicePath)), snakeCaseModuleName) {
		fmt.Println("Usecase Already Exists With Name:", strings.ToLower(moduleName))
		return
	}

	templateString, err := parseUsecaseTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = generateUsecaseFile(utils.NormalizePath(fmt.Sprintf("%s/%s", rootPath, config.Paths().ServicePath)), snakeCaseModuleName, templateString)
	if err != nil {
		fmt.Println("Error generating usecase:", err)
		return
	}

	fmt.Println(fmt.Sprintf("Usecase created and formatted successfully! with name: %s", snakeCaseModuleName))
}
