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

func seederFileExists(seederPath string, snakeCaseModuleName string) bool {
	seederFilePath := fmt.Sprintf("%s/%s.go", seederPath, snakeCaseModuleName)
	if _, err := os.Stat("/" + utils.NormalizePath(seederFilePath)); err == nil {
		return true
	} else {
		return false
	}
}

func parseSeederTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName string, numFields int, fieldNames []string, fieldTypes []string) (string, error) {
	var fields []Field
	for i := 0; i < numFields; i++ {
		// Process the field names
		titleCaseFieldName, snakeCaseFieldName, _ := utils.ProcessString(fieldNames[i])

		// Append the field to the fields slice
		fields = append(fields, Field{
			Name:    titleCaseFieldName,
			Type:    fieldTypes[i],
			JsonTag: snakeCaseFieldName,
		})
	}
	// Prepare the data
	seederTemplateData := struct {
		UcFirstName     string
		SmallName       string
		SnakeCaseName   string
		SmallPluralName string
		ModuleName      string
		DomainPath      string
		RepositoryPath  string
		UsecasePath     string
		Fields          []Field
	}{
		UcFirstName:     titleCaseModuleName,
		SmallName:       camelCaseModuleName,
		SnakeCaseName:   snakeCaseModuleName,
		SmallPluralName: utils.ToPlural(snakeCaseModuleName),
		ModuleName:      config.Paths().ModuleName,
		DomainPath:      utils.NormalizePath(fmt.Sprintf("%s/%s", config.Paths().ModuleName, config.Paths().DomainPath)),
		RepositoryPath:  utils.NormalizePath(fmt.Sprintf("%s/%s/%s/repository", config.Paths().ModuleName, config.Paths().ServicePath, camelCaseModuleName)),
		UsecasePath:     utils.NormalizePath(fmt.Sprintf("%s/%s/%s/usecase", config.Paths().ModuleName, config.Paths().ServicePath, camelCaseModuleName)),
		Fields:          fields,
	}

	// Read the contents of the file
	sourceContent := cli_template.SeederMongo

	// Create a new template and parse the template string
	parsedTemplate, err := template.New("seederTemplate").Parse(string(sourceContent))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return "", err
	}

	var buf bytes.Buffer
	// Execute the template with the data
	err = parsedTemplate.Execute(&buf, seederTemplateData)
	if err != nil {
		fmt.Println("Error executing parsed template:", err)
		return "", err
	}

	return buf.String(), nil
}

func generateSeederFile(seederPath string, snakeCaseModuleName string, templateString string) error {
	// Create the directory path
	directoryPath := fmt.Sprintf("%s", seederPath)
	err := os.MkdirAll("/"+utils.NormalizePath(directoryPath), os.ModePerm) // os.ModePerm is 0777
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

	seederFileName := fmt.Sprintf("%s/%s.go", seederPath, snakeCaseModuleName)

	// Write the code to the file
	err = ioutil.WriteFile("/"+utils.NormalizePath(seederFileName), []byte(templateString), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	// Execute the `go fmt` command
	goFmtCmd := exec.Command("go", "fmt", "/"+utils.NormalizePath(seederFileName))
	goFmtCmd.Stdout = os.Stdout
	goFmtCmd.Stderr = os.Stderr
	err = goFmtCmd.Run()
	if err != nil {
		fmt.Println("Error generating seeder:", err)
		return err
	}
	return nil
}

func MakeSeeder(cmd *cobra.Command, moduleName string, numFields int, fieldNames []string, fieldTypes []string) {
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

	if seederFileExists("/"+utils.NormalizePath(fmt.Sprintf("%s/%s", rootPath, config.Paths().SeederPath)), snakeCaseModuleName) {
		fmt.Println("Seeder Already Exists With Name:", snakeCaseModuleName)
		return
	}

	templateString, err := parseSeederTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName, numFields, fieldNames, fieldTypes)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = generateSeederFile("/"+utils.NormalizePath(fmt.Sprintf("%s/%s", rootPath, config.Paths().SeederPath)), snakeCaseModuleName, templateString)
	if err != nil {
		fmt.Println("Error generating seeder:", err)
		return
	}

	fmt.Println(fmt.Sprintf("Seeder created and formatted successfully! with name: %s", snakeCaseModuleName))
}
