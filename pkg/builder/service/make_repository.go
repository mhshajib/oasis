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

func repositoryFileExists(servicePath string, snakeCaseModuleName string) bool {
	repositoryFilePath := fmt.Sprintf("%s/%s/repository/%s_mongo.go", servicePath, snakeCaseModuleName, snakeCaseModuleName)
	if _, err := os.Stat("/" + utils.NormalizePath(repositoryFilePath)); err == nil {
		return true
	} else {
		return false
	}
}

func parseRepositoryTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName string) (string, error) {
	// Prepare the data
	repositoryTemplateData := struct {
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
	sourceContent := cli_template.Repository

	// Create a new template and parse the template string
	parsedTemplate, err := template.New("repositoryTemplate").Parse(string(sourceContent))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return "", err
	}

	var buf bytes.Buffer
	// Execute the template with the data
	err = parsedTemplate.Execute(&buf, repositoryTemplateData)
	if err != nil {
		fmt.Println("Error executing parsed template:", err)
		return "", err
	}

	return buf.String(), nil
}

func generateRepositoryFile(servicePath string, snakeCaseModuleName string, templateString string) error {

	// Create the directory path
	directoryPath := fmt.Sprintf("%s/%s/repository", servicePath, snakeCaseModuleName)
	err := os.MkdirAll("/"+utils.NormalizePath(directoryPath), os.ModePerm) // os.ModePerm is 0777
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

	// Create the file path
	repositoryFileName := fmt.Sprintf("%s/%s_mongo.go", directoryPath, snakeCaseModuleName)

	// Write the code to the file
	err = ioutil.WriteFile("/"+utils.NormalizePath(repositoryFileName), []byte(templateString), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	// Execute the `go fmt` command
	goFmtCmd := exec.Command("go", "fmt", "/"+utils.NormalizePath(repositoryFileName))
	goFmtCmd.Stdout = os.Stdout
	goFmtCmd.Stderr = os.Stderr
	err = goFmtCmd.Run()
	if err != nil {
		fmt.Println("Error formatting repository file:", err)
		return err
	}

	return nil
}

func MakeRepository(cmd *cobra.Command, moduleName string) {
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

	if repositoryFileExists("/"+utils.NormalizePath(fmt.Sprintf("%s/%s", rootPath, config.Paths().ServicePath)), snakeCaseModuleName) {
		fmt.Println("Domain Already Exists With Name:", snakeCaseModuleName)
		return
	}

	templateString, err := parseRepositoryTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = generateRepositoryFile("/"+utils.NormalizePath(fmt.Sprintf("%s/%s", rootPath, config.Paths().ServicePath)), snakeCaseModuleName, templateString)
	if err != nil {
		fmt.Println("Error generating repository:", err)
		return
	}

	fmt.Println(fmt.Sprintf("Repository created and formatted successfully! with name: %s", snakeCaseModuleName))
}
