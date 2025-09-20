package main

import (
	"flag"
	"os"
	"testing"
	"time"
)

func TestMainFunctionWithoutArgs(t *testing.T) {
	// Save original args and replace them
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test with no --file argument (should exit)
	os.Args = []string{"cmd"}

	// Since main() calls os.Exit(2) when no file is provided, we can't easily test it
	// Instead, we test the flag parsing logic separately

	// Reset flags for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	filePath := flag.String("file", "", "path to the design doc (HLD/LLD)")
	turns := flag.Int("turns", 0, "maximum number of argument turns (each agent counts as one turn)")
	duration := flag.Duration("duration", 0, "max duration for the debate (e.g. 30s, 2m). 0 means no limit")
	model := flag.String("model", "", "LLM model to use for text(optional)")
	imageModel := flag.String("image-model", "", "LLM model to use for images (optional)")

	// Parse empty args
	err := flag.CommandLine.Parse([]string{})
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	// Check default values
	if *filePath != "" {
		t.Errorf("Expected empty file path, got %s", *filePath)
	}
	if *turns != 0 {
		t.Errorf("Expected turns default to 0, got %d", *turns)
	}
	if *duration != 0 {
		t.Errorf("Expected duration default to 0, got %v", *duration)
	}
	if *model != "" {
		t.Errorf("Expected empty model, got %s", *model)
	}
	if *imageModel != "" {
		t.Errorf("Expected empty imageModel, got %s", *imageModel)
	}
}

func TestFlagParsing(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected map[string]interface{}
	}{
		{
			name: "basic file argument",
			args: []string{"--file", "test.md"},
			expected: map[string]interface{}{
				"file":        "test.md",
				"turns":       0,
				"duration":    time.Duration(0),
				"model":       "",
				"image-model": "",
			},
		},
		{
			name: "all arguments",
			args: []string{
				"--file", "design.md",
				"--turns", "5",
				"--duration", "2m",
				"--model", "custom-model",
				"--image-model", "custom-image-model",
			},
			expected: map[string]interface{}{
				"file":        "design.md",
				"turns":       5,
				"duration":    2 * time.Minute,
				"model":       "custom-model",
				"image-model": "custom-image-model",
			},
		},
		{
			name: "duration parsing",
			args: []string{"--file", "test.md", "--duration", "30s"},
			expected: map[string]interface{}{
				"file":        "test.md",
				"turns":       0,
				"duration":    30 * time.Second,
				"model":       "",
				"image-model": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new flag set for each test
			fs := flag.NewFlagSet("test", flag.ContinueOnError)

			filePath := fs.String("file", "", "path to the design doc (HLD/LLD)")
			turns := fs.Int("turns", 0, "maximum number of argument turns (each agent counts as one turn)")
			duration := fs.Duration("duration", 0, "max duration for the debate (e.g. 30s, 2m). 0 means no limit")
			model := fs.String("model", "", "LLM model to use for text(optional)")
			imageModel := fs.String("image-model", "", "LLM model to use for images (optional)")

			err := fs.Parse(tt.args)
			if err != nil {
				t.Fatalf("Failed to parse flags: %v", err)
			}

			// Check each expected value
			if *filePath != tt.expected["file"] {
				t.Errorf("file = %v, want %v", *filePath, tt.expected["file"])
			}
			if *turns != tt.expected["turns"] {
				t.Errorf("turns = %v, want %v", *turns, tt.expected["turns"])
			}
			if *duration != tt.expected["duration"] {
				t.Errorf("duration = %v, want %v", *duration, tt.expected["duration"])
			}
			if *model != tt.expected["model"] {
				t.Errorf("model = %v, want %v", *model, tt.expected["model"])
			}
			if *imageModel != tt.expected["image-model"] {
				t.Errorf("image-model = %v, want %v", *imageModel, tt.expected["image-model"])
			}
		})
	}
}

func TestInvalidDurationParsing(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Duration("duration", 0, "test duration")

	err := fs.Parse([]string{"--duration", "invalid"})
	if err == nil {
		t.Error("Expected error for invalid duration, got nil")
	}
}

func TestInvalidTurnsParsing(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Int("turns", 0, "test turns")

	err := fs.Parse([]string{"--turns", "invalid"})
	if err == nil {
		t.Error("Expected error for invalid turns, got nil")
	}
}

func TestHelpFlag(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("file", "", "path to the design doc (HLD/LLD)")
	fs.Int("turns", 0, "maximum number of argument turns")
	fs.Duration("duration", 0, "max duration for the debate")
	fs.String("model", "", "LLM model to use for text")
	fs.String("image-model", "", "LLM model to use for images")

	// Redirect output to avoid cluttering test output
	fs.SetOutput(os.Stderr)

	err := fs.Parse([]string{"-h"})
	if err != flag.ErrHelp {
		t.Errorf("Expected ErrHelp, got %v", err)
	}
}
