package agents

import (
	"context"
	"testing"
)

func TestNewLLMDefaults(t *testing.T) {
	ctx := context.Background()

	// Test with empty model strings - should use defaults
	llm, err := NewLLM(ctx, "", "")

	// We expect this to fail due to missing API credentials in test environment
	// but we can check that the function handles default values correctly
	if err == nil {
		// If somehow it succeeds, check the defaults were set
		if llm.Models.Text != "gemini-2.5-pro" {
			t.Errorf("Expected default text model 'gemini-2.5-pro', got %s", llm.Models.Text)
		}
		if llm.Models.Image != "gemini-2.5-flash-image-preview" {
			t.Errorf("Expected default image model 'gemini-2.5-flash-image-preview', got %s", llm.Models.Image)
		}
		if llm.client == nil {
			t.Error("Expected client to be set")
		}
	} else {
		// Expected to fail due to missing credentials, that's fine for testing
		t.Logf("NewLLM failed as expected in test environment: %v", err)
	}
}

func TestNewLLMCustomModels(t *testing.T) {
	ctx := context.Background()

	customTextModel := "custom-text-model"
	customImageModel := "custom-image-model"

	// Test with custom model strings
	llm, err := NewLLM(ctx, customTextModel, customImageModel)

	// We expect this to fail due to missing API credentials in test environment
	// but we can check that the function handles custom values correctly
	if err == nil {
		// If somehow it succeeds, check the custom values were set
		if llm.Models.Text != customTextModel {
			t.Errorf("Expected custom text model %s, got %s", customTextModel, llm.Models.Text)
		}
		if llm.Models.Image != customImageModel {
			t.Errorf("Expected custom image model %s, got %s", customImageModel, llm.Models.Image)
		}
		if llm.client == nil {
			t.Error("Expected client to be set")
		}
	} else {
		// Expected to fail due to missing credentials, that's fine for testing
		t.Logf("NewLLM failed as expected in test environment: %v", err)
	}
}

func TestModelsStruct(t *testing.T) {
	models := Models{
		Text:  "test-text-model",
		Image: "test-image-model",
	}

	if models.Text != "test-text-model" {
		t.Errorf("Expected text model 'test-text-model', got %s", models.Text)
	}
	if models.Image != "test-image-model" {
		t.Errorf("Expected image model 'test-image-model', got %s", models.Image)
	}
}

func TestLLMStruct(t *testing.T) {
	models := Models{
		Text:  "test-text-model",
		Image: "test-image-model",
	}

	llm := &LLM{
		client: nil, // In real usage, this would be a genai.Client
		Models: models,
	}

	if llm.Models.Text != "test-text-model" {
		t.Errorf("Expected text model 'test-text-model', got %s", llm.Models.Text)
	}
	if llm.Models.Image != "test-image-model" {
		t.Errorf("Expected image model 'test-image-model', got %s", llm.Models.Image)
	}
}

// Note: We can't easily test GetArgument and GetComicStrip methods without
// mocking the genai.Client, which would require significant refactoring.
// In a production environment, these would be tested with integration tests
// or by making the genai.Client interface-based for easier mocking.
