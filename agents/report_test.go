package agents

import (
	"testing"
	"time"
)

func TestBuildTranscriptSnippet(t *testing.T) {
	// Create test messages
	now := time.Now()
	history := []Message{
		{Sender: "Agent1", Argument: Argument{Content: "First message", Tone: "neutral"}, Time: now},
		{Sender: "Agent2", Argument: Argument{Content: "Second message", Tone: "confident"}, Time: now.Add(time.Minute)},
		{Sender: "Agent3", Argument: Argument{Content: "Third message", Tone: "analytical"}, Time: now.Add(2 * time.Minute)},
	}

	tests := []struct {
		name     string
		history  []Message
		n        int
		expected string
	}{
		{
			name:     "get last 2 messages",
			history:  history,
			n:        2,
			expected: "Agent2: Second message\nAgent3: Third message",
		},
		{
			name:     "get all messages when n is larger than history",
			history:  history,
			n:        5,
			expected: "Agent1: First message\nAgent2: Second message\nAgent3: Third message",
		},
		{
			name:     "get last 1 message",
			history:  history,
			n:        1,
			expected: "Agent3: Third message",
		},
		{
			name:     "zero messages requested",
			history:  history,
			n:        0,
			expected: "",
		},
		{
			name:     "negative n",
			history:  history,
			n:        -1,
			expected: "",
		},
		{
			name:     "empty history",
			history:  []Message{},
			n:        2,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildTranscriptSnippet(tt.history, tt.n)
			if result != tt.expected {
				t.Errorf("buildTranscriptSnippet() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestBuildTranscriptSnippetWithLongMessages(t *testing.T) {
	now := time.Now()
	longMessage := ""
	for i := 0; i < 300; i++ {
		longMessage += "a"
	}

	history := []Message{
		{Sender: "Agent1", Argument: Argument{Content: longMessage, Tone: "verbose"}, Time: now},
	}

	result := buildTranscriptSnippet(history, 1)

	// Should be truncated to 200 characters + "…"
	expected := "Agent1: " + longMessage[:200] + "…"
	if result != expected {
		t.Errorf("buildTranscriptSnippet() with long message = %q, want %q", result, expected)
	}
}

func TestTruncateLocal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "short string not truncated",
			input:    "hello",
			maxLen:   10,
			expected: "hello",
		},
		{
			name:     "exact length not truncated",
			input:    "hello",
			maxLen:   5,
			expected: "hello",
		},
		{
			name:     "long string truncated",
			input:    "this is a very long string that should be truncated",
			maxLen:   10,
			expected: "this is a …",
		},
		{
			name:     "empty string",
			input:    "",
			maxLen:   5,
			expected: "",
		},
		{
			name:     "zero max length",
			input:    "hello",
			maxLen:   0,
			expected: "…",
		},
		{
			name:     "single character truncation",
			input:    "hello",
			maxLen:   1,
			expected: "h…",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateLocal(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncateLocal(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}