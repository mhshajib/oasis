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
)

func timestampFileExists(domainPath string) bool {
	domainFilePath := fmt.Sprintf("%s/timestamp.go", domainPath)
	if _, err := os.Stat(utils.NormalizePath(domainFilePath)); err == nil {
		return true
	} else {
		return false
	}
}

func parseTimestampTemplate() (string, error) {

	// Read the contents of the file
	sourceContent := cli_template.TimeStamp

	// Create a new template and parse the template string
	parsedTemplate, err := template.New("timestampTemplate").Parse(string(sourceContent))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return "", err
	}

	var buf bytes.Buffer
	// Execute the template with the data
	err = parsedTemplate.Execute(&buf, nil)
	if err != nil {
		fmt.Println("Error executing parsed template:", err)
		return "", err
	}

	return buf.String(), nil
}

func generateTimestampFile(domainPath string, templateString string) error {

	timestampFileName := fmt.Sprintf("%s/timestamp.go", domainPath)

	// Write the code to the file
	err := ioutil.WriteFile(utils.NormalizePath(timestampFileName), []byte(templateString), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	// Execute the `go fmt` command
	goFmtCmd := exec.Command("go", "fmt", utils.NormalizePath(timestampFileName))
	goFmtCmd.Stdout = os.Stdout
	goFmtCmd.Stderr = os.Stderr
	err = goFmtCmd.Run()
	if err != nil {
		fmt.Println("Error generating domain:", err)
		return err
	}
	return nil
}

func MakeTimestamp(rootPath string) {
	if timestampFileExists(utils.NormalizePath(fmt.Sprintf("%s/%s", rootPath, config.Paths().DomainPath))) {
		return
	}

	templateString, err := parseTimestampTemplate()
	if err != nil {
		fmt.Println("Error parsing timestamp template:", err)
		return
	}

	err = generateTimestampFile(utils.NormalizePath(fmt.Sprintf("%s/%s", rootPath, config.Paths().DomainPath)), templateString)
	if err != nil {
		fmt.Println("Error generating timestamp:", err)
		return
	}

	fmt.Println("Timestamp created and formatted successfully!")
}
