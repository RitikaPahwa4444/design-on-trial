package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WriteReport writes the HTML report next to baseName in outDir and returns the written path.
func WriteReport(outDir, baseName, html string) (string, error) {
	if outDir == "" {
		outDir = "."
	}
	filename := fmt.Sprintf("%s_report.html", strings.TrimSuffix(baseName, filepath.Ext(baseName)))
	path := filepath.Join(outDir, filename)
	if err := os.WriteFile(path, []byte(html), 0644); err != nil {
		return "", err
	}
	return path, nil
}
