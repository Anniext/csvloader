package embed

import (
	"bufio"
	"embed"
	"log"
	"strings"
)

//go:embed template
var embeddedFiles embed.FS

func OpenTemplateFileToString(strPath string) string {
	var stringOut []string
	file, err := embeddedFiles.Open("template/" + strPath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	scaner := bufio.NewScanner(file)
	for scaner.Scan() {
		line := scaner.Text()
		stringOut = append(stringOut, line)
	}
	return strings.Join(stringOut, "\n")
}
