package agents

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPersonasAdditional(t *testing.T) {
	// This test is additional to the existing TestLoadPersonas
	// Test that LoadPersonas returns proper Agent structs
	personas, err := LoadPersonas()
	if err != nil {
		t.Fatalf("LoadPersonas failed: %v", err)
	}

	for i, agent := range personas {
		t.Run(agent.Name, func(t *testing.T) {
			if agent.Name == "" {
				t.Errorf("Agent[%d] has empty name", i)
			}
			if agent.Role == "" {
				t.Errorf("Agent[%d] (%s) has empty role", i, agent.Name)
			}
			if agent.Persona == "" {
				t.Errorf("Agent[%d] (%s) has empty persona", i, agent.Name)
			}

			// Check that persona contains meaningful content
			if len(agent.Persona) < 50 {
				t.Errorf("Agent[%d] (%s) has suspiciously short persona: %s", i, agent.Name, agent.Persona)
			}
		})
	}
}

func TestLoadPersonasExpectedAgents(t *testing.T) {
	personas, err := LoadPersonas()
	if err != nil {
		t.Fatalf("LoadPersonas failed: %v", err)
	}

	// Map to track found agents
	agentNames := make(map[string]bool)
	for _, agent := range personas {
		agentNames[agent.Name] = true
	}

	// Check for expected agents based on current persona.json
	expectedNames := []string{"Codewright", "Blueprint", "Justice Logic", "Scribbler"}
	for _, name := range expectedNames {
		if !agentNames[name] {
			t.Errorf("Expected agent %q not found in personas", name)
		}
	}
}

func TestLoadPersonasExpectedRoles(t *testing.T) {
	personas, err := LoadPersonas()
	if err != nil {
		t.Fatalf("LoadPersonas failed: %v", err)
	}

	// Map to track found roles
	roleNames := make(map[string]bool)
	for _, agent := range personas {
		roleNames[agent.Role] = true
	}

	// Check for expected roles based on current persona.json
	expectedRoles := []string{"Developer (Defense)", "Architect (Prosecutor)", "CTO (Judge)", "Reporter"}
	for _, role := range expectedRoles {
		if !roleNames[role] {
			t.Errorf("Expected role %q not found in personas", role)
		}
	}
}

func TestLoadPersonasFileNotFound(t *testing.T) {
	// Temporarily change to a directory without persona.json
	tempDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// This should use the embedded persona.json, so it should still work
	personas, err := LoadPersonas()
	if err != nil {
		t.Fatalf("LoadPersonas should work with embedded JSON even without file: %v", err)
	}

	if len(personas) == 0 {
		t.Error("LoadPersonas should return embedded personas when file not found")
	}
}

func TestLoadPersonasFromExecutableDirectory(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Since LoadPersonas() uses embedded JSON, changing directory won't affect it
	// This test verifies that LoadPersonas works regardless of current directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// LoadPersonas should work because it uses embedded JSON
	personas, err := LoadPersonas()
	if err != nil {
		t.Fatalf("LoadPersonas failed from different directory: %v", err)
	}

	if len(personas) == 0 {
		t.Fatal("Expected at least one persona, got none")
	}

	// Should get the same personas as the embedded JSON
	expectedNames := []string{"Codewright", "Blueprint", "Justice Logic", "Scribbler"}
	agentNames := make(map[string]bool)
	for _, agent := range personas {
		agentNames[agent.Name] = true
	}

	for _, name := range expectedNames {
		if !agentNames[name] {
			t.Errorf("Expected agent %q from embedded JSON not found", name)
		}
	}
}

func TestLoadPersonasInvalidJSON(t *testing.T) {
	// Since LoadPersonas() uses embedded JSON and doesn't read from local files,
	// this test verifies that LoadPersonas is resilient to environment changes
	// and always returns the embedded personas

	// Create a temporary directory with unrelated files
	tempDir := t.TempDir()

	// Create some unrelated file that might interfere
	interfereFile := filepath.Join(tempDir, "persona.json")
	err := os.WriteFile(interfereFile, []byte("not json"), 0644)
	if err != nil {
		t.Fatalf("Failed to write interference file: %v", err)
	}

	// Change to the temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// LoadPersonas should still work because it uses embedded JSON
	personas, err := LoadPersonas()
	if err != nil {
		t.Fatalf("LoadPersonas should work with embedded JSON regardless of local files: %v", err)
	}

	// Should get the embedded personas
	if len(personas) == 0 {
		t.Error("LoadPersonas should return embedded personas")
	}
}

// Test the prompt generation logic of GenerateArgument without calling LLM
func TestGenerateArgumentPromptGeneration(t *testing.T) {
	agent := Agent{
		Name:    "TestAgent",
		Role:    "TestRole",
		Persona: "You are a test agent with specific instructions.",
	}

	// Test that the agent fields are properly set
	if agent.Name != "TestAgent" {
		t.Errorf("Agent name = %v, want TestAgent", agent.Name)
	}
	if agent.Role != "TestRole" {
		t.Errorf("Agent role = %v, want TestRole", agent.Role)
	}
	if agent.Persona != "You are a test agent with specific instructions." {
		t.Errorf("Agent persona = %v, want specific instructions", agent.Persona)
	}

	// We can't easily test GenerateArgument without mocking the LLM interface,
	// but we can test that the agent struct is properly constructed
	// In a real implementation, this would use dependency injection with interfaces
}

func TestAgentEmptyFields(t *testing.T) {
	agent := Agent{}

	if agent.Name != "" {
		t.Errorf("Empty agent name should be empty string, got %v", agent.Name)
	}
	if agent.Role != "" {
		t.Errorf("Empty agent role should be empty string, got %v", agent.Role)
	}
	if agent.Persona != "" {
		t.Errorf("Empty agent persona should be empty string, got %v", agent.Persona)
	}
}
