package agents

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPersonas(t *testing.T) {
	// Test that LoadPersonas can find persona.json from different locations
	
	// This test depends on the actual persona.json file existing
	// First, check if we can load from the standard relative path
	personas, err := LoadPersonas()
	if err != nil {
		t.Fatalf("LoadPersonas failed: %v", err)
	}
	
	if len(personas) == 0 {
		t.Fatal("Expected at least one persona, got none")
	}
	
	// Verify we have the expected agent roles
	rolesSeen := make(map[string]bool)
	for _, agent := range personas {
		rolesSeen[agent.Role] = true
		if agent.Name == "" {
			t.Errorf("Agent has empty name: %+v", agent)
		}
		if agent.Persona == "" {
			t.Errorf("Agent has empty persona: %+v", agent)
		}
	}
	
	// Check for expected roles based on persona.json content
	expectedRoles := []string{"Developer (Defense)", "Architect (Prosecutor)", "Judge", "Reporter"}
	for _, role := range expectedRoles {
		if !rolesSeen[role] {
			t.Errorf("Expected role %q not found in personas", role)
		}
	}
}

func TestLoadPersonasFromExecutableDir(t *testing.T) {
	// Create a temporary directory to simulate executable deployment
	tempDir := t.TempDir()
	
	// Copy persona.json to the temp directory
	originalPersona, err := os.ReadFile("persona.json")
	if err != nil {
		// Fallback: try to find it in the project
		originalPersona, err = os.ReadFile("../persona.json")
		if err != nil {
			t.Skip("Could not find persona.json for test")
		}
	}
	
	personaPath := filepath.Join(tempDir, "persona.json")
	err = os.WriteFile(personaPath, originalPersona, 0644)
	if err != nil {
		t.Fatalf("Failed to write test persona.json: %v", err)
	}
	
	// Change to the temp directory to simulate running from executable directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldDir)
	
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	
	// Now LoadPersonas should find the persona.json in the current directory
	personas, err := LoadPersonas()
	if err != nil {
		t.Fatalf("LoadPersonas failed from executable directory: %v", err)
	}
	
	if len(personas) == 0 {
		t.Fatal("Expected at least one persona, got none")
	}
}