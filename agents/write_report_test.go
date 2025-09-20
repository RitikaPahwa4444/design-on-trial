package agents

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteReport(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	testHTML := "<html><body><h1>Test Report</h1></body></html>"

	tests := []struct {
		name     string
		outDir   string
		baseName string
		html     string
		wantErr  bool
	}{
		{
			name:     "basic write to temp dir",
			outDir:   tempDir,
			baseName: "test.md",
			html:     testHTML,
			wantErr:  false,
		},
		{
			name:     "write with no extension",
			outDir:   tempDir,
			baseName: "test",
			html:     testHTML,
			wantErr:  false,
		},
		{
			name:     "write to current dir when outDir empty",
			outDir:   "",
			baseName: "test.txt",
			html:     testHTML,
			wantErr:  false,
		},
		{
			name:     "write with different extension",
			outDir:   tempDir,
			baseName: "design.txt",
			html:     testHTML,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := WriteReport(tt.outDir, tt.baseName, tt.html)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return // Test expects error, we're done
			}

			// Verify the file was written
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("WriteReport() file was not created at %s", path)
				return
			}

			// Verify the content
			content, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("Failed to read written file: %v", err)
				return
			}

			if string(content) != tt.html {
				t.Errorf("WriteReport() content = %v, want %v", string(content), tt.html)
			}

			// Verify the filename format
			expectedName := strings.TrimSuffix(tt.baseName, filepath.Ext(tt.baseName)) + "_report.html"
			if !strings.HasSuffix(path, expectedName) {
				t.Errorf("WriteReport() path = %v, expected to end with %v", path, expectedName)
			}

			// Clean up the file if it's in current directory
			if tt.outDir == "" {
				os.Remove(path)
			}
		})
	}
}

func TestWriteReportPathFormat(t *testing.T) {
	tempDir := t.TempDir()
	testHTML := "<html><body>Test</body></html>"

	tests := []struct {
		baseName     string
		expectedName string
	}{
		{"test.md", "test_report.html"},
		{"design.txt", "design_report.html"},
		{"sample_hld.md", "sample_hld_report.html"},
		{"noextension", "noextension_report.html"},
		{"multiple.dots.in.name.md", "multiple.dots.in.name_report.html"},
	}

	for _, tt := range tests {
		t.Run(tt.baseName, func(t *testing.T) {
			path, err := WriteReport(tempDir, tt.baseName, testHTML)
			if err != nil {
				t.Fatalf("WriteReport() error = %v", err)
			}

			fileName := filepath.Base(path)
			if fileName != tt.expectedName {
				t.Errorf("WriteReport() filename = %v, want %v", fileName, tt.expectedName)
			}
		})
	}
}

func TestWriteReportInvalidDirectory(t *testing.T) {
	// Try to write to a directory that doesn't exist and can't be created
	invalidDir := "/invalid/nonexistent/directory"
	testHTML := "<html><body>Test</body></html>"

	_, err := WriteReport(invalidDir, "test.md", testHTML)
	if err == nil {
		t.Error("WriteReport() expected error for invalid directory, got nil")
	}
}

func TestWriteReportEmptyHTML(t *testing.T) {
	tempDir := t.TempDir()

	path, err := WriteReport(tempDir, "test.md", "")
	if err != nil {
		t.Fatalf("WriteReport() error = %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	if string(content) != "" {
		t.Errorf("WriteReport() with empty HTML should write empty file, got %v", string(content))
	}
}
