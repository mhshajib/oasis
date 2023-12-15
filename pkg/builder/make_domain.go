package builder

import (
	"bytes"
	"clean-cli/pkg/config"
	cli_template "clean-cli/pkg/template"
	"clean-cli/pkg/utils"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func domainFileExists(domainPath string, snakeCaseModuleName string) bool {
	domainFilePath := fmt.Sprintf("%s/%s.go", domainPath, snakeCaseModuleName)
	if _, err := os.Stat(domainFilePath); err == nil {
		return true
	} else {
		return false
	}
}

func parseDomainTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName string) (string, error) {
	// Prepare the data
	domainTemplateData := struct {
		UcFirstName     string
		SmallName       string
		SnameName       string
		SmallPluralName string
	}{
		UcFirstName:     titleCaseModuleName,
		SmallName:       camelCaseModuleName,
		SnameName:       snakeCaseModuleName,
		SmallPluralName: utils.ToPlural(snakeCaseModuleName),
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
	err := ioutil.WriteFile(domainFileName, []byte(templateString), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	// Execute the `go fmt` command
	goFmtCmd := exec.Command("go", "fmt", domainFileName)
	goFmtCmd.Stdout = os.Stdout
	goFmtCmd.Stderr = os.Stderr
	err = goFmtCmd.Run()
	if err != nil {
		fmt.Println("Error generating domain:", err)
		return err
	}
	return nil
}

func MakeDomain(cmd *cobra.Command, args []string) {
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

	MakeTimestamp(rootPath)

	err = utils.CreateDirectory(rootPath, config.Paths().DomainPath)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}
	moduleName := args[0]

	titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName := utils.ProcessString(moduleName)

	if domainFileExists(fmt.Sprintf("%s/%s", rootPath, config.Paths().DomainPath), snakeCaseModuleName) {
		fmt.Println("Domain Already Exists With Name:", snakeCaseModuleName)
		return
	}

	templateString, err := parseDomainTemplate(titleCaseModuleName, snakeCaseModuleName, camelCaseModuleName)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = generateDomainFile(fmt.Sprintf("%s/%s", rootPath, config.Paths().DomainPath), snakeCaseModuleName, templateString)
	if err != nil {
		fmt.Println("Error generating domain:", err)
		return
	}

	fmt.Println(fmt.Sprintf("Domain created and formatted successfully! with name: %s", snakeCaseModuleName))
}
