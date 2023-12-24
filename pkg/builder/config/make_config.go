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

func configFileExists(configPath string, snakeCaseModuleName string) bool {
	configFilePath := fmt.Sprintf("%s/%s.go", configPath, snakeCaseModuleName)
	if _, err := os.Stat(configFilePath); err == nil {
		return true
	} else {
		return false
	}
}

func parseConfigTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName string) (string, error) {
	// Prepare the data
	configTemplateData := struct {
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

	// Read the contents of the file
	sourceContent := cli_template.Config

	// Create a new template and parse the template string
	parsedTemplate, err := template.New("configTemplate").Parse(string(sourceContent))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return "", err
	}

	var buf bytes.Buffer
	// Execute the template with the data
	err = parsedTemplate.Execute(&buf, configTemplateData)
	if err != nil {
		fmt.Println("Error executing parsed template:", err)
		return "", err
	}

	return buf.String(), nil
}

func generateConfigFile(configPath string, snakeCaseModuleName string, templateString string) error {

	configFileName := fmt.Sprintf("%s/%s.go", configPath, snakeCaseModuleName)

	// Write the code to the file
	err := ioutil.WriteFile(configFileName, []byte(templateString), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	// Execute the `go fmt` command
	goFmtCmd := exec.Command("go", "fmt", configFileName)
	goFmtCmd.Stdout = os.Stdout
	goFmtCmd.Stderr = os.Stderr
	err = goFmtCmd.Run()
	if err != nil {
		fmt.Println("Error generating config:", err)
		return err
	}
	return nil
}

func MakeConfig(cmd *cobra.Command, args []string) {
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

	if configFileExists(fmt.Sprintf("%s/%s", rootPath, config.Paths().ConfigPath), snakeCaseModuleName) {
		fmt.Println("Config Already Exists With Name:", snakeCaseModuleName)
		return
	}

	templateString, err := parseConfigTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = generateConfigFile(fmt.Sprintf("%s/%s", rootPath, config.Paths().ConfigPath), snakeCaseModuleName, templateString)
	if err != nil {
		fmt.Println("Error generating config:", err)
		return
	}

	fmt.Println(fmt.Sprintf("Config created and formatted successfully! with name: %s", snakeCaseModuleName))
}
