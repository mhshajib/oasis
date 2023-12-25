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

func migrationFileExists(migrationPath string, snakeCaseModuleName string) bool {
	migrationFilePath := fmt.Sprintf("%s/%s.go", migrationPath, snakeCaseModuleName)
	if _, err := os.Stat(migrationFilePath); err == nil {
		return true
	} else {
		return false
	}
}

func parseMigrationTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName string) (string, error) {
	// Prepare the data
	migrationTemplateData := struct {
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
	sourceContent := cli_template.MigrationMongo

	// Create a new template and parse the template string
	parsedTemplate, err := template.New("migrationTemplate").Parse(string(sourceContent))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return "", err
	}

	var buf bytes.Buffer
	// Execute the template with the data
	err = parsedTemplate.Execute(&buf, migrationTemplateData)
	if err != nil {
		fmt.Println("Error executing parsed template:", err)
		return "", err
	}

	return buf.String(), nil
}

func generateMigrationFile(migrationPath string, snakeCaseModuleName string, templateString string) error {

	migrationFileName := fmt.Sprintf("%s/%s.go", migrationPath, snakeCaseModuleName)

	// Write the code to the file
	err := ioutil.WriteFile(migrationFileName, []byte(templateString), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	// Execute the `go fmt` command
	goFmtCmd := exec.Command("go", "fmt", migrationFileName)
	goFmtCmd.Stdout = os.Stdout
	goFmtCmd.Stderr = os.Stderr
	err = goFmtCmd.Run()
	if err != nil {
		fmt.Println("Error generating migration:", err)
		return err
	}
	return nil
}

func MakeMigration(cmd *cobra.Command, args []string) {
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

	if migrationFileExists(fmt.Sprintf("%s/%s", rootPath, config.Paths().MigrationPath), snakeCaseModuleName) {
		fmt.Println("Migration Already Exists With Name:", snakeCaseModuleName)
		return
	}

	templateString, err := parseMigrationTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = generateMigrationFile(fmt.Sprintf("%s/%s", rootPath, config.Paths().MigrationPath), snakeCaseModuleName, templateString)
	if err != nil {
		fmt.Println("Error generating migration:", err)
		return
	}

	fmt.Println(fmt.Sprintf("Migration created and formatted successfully! with name: %s", snakeCaseModuleName))
}
