package agents

import (
	"context"
	"testing"
	"time"
)

func TestRunDebateNilLLM(t *testing.T) {
	ctx := context.Background()
	participants := []Agent{
		{Name: "Agent1", Role: "Role1", Persona: "Persona1"},
		{Name: "Agent2", Role: "Role2", Persona: "Persona2"},
	}

	_, err := RunDebate(ctx, nil, participants, nil, "test doc", 1, 0)
	if err == nil {
		t.Error("RunDebate() expected error for nil LLM, got nil")
	}
	if err.Error() != "llm is nil" {
		t.Errorf("RunDebate() error = %v, want 'llm is nil'", err)
	}
}

func TestRunDebateInsufficientParticipants(t *testing.T) {
	ctx := context.Background()
	// We can't create a real LLM for testing, but we can test the validation logic
	// by passing a non-nil pointer (even though it won't work for actual calls)
	dummyLLM := &LLM{} // This will fail later, but we're only testing parameter validation

	// Test with no participants
	_, err := RunDebate(ctx, dummyLLM, []Agent{}, nil, "test doc", 1, 0)
	if err == nil {
		t.Error("RunDebate() expected error for no participants, got nil")
	}
	if err.Error() != "need at least two participants" {
		t.Errorf("RunDebate() error = %v, want 'need at least two participants'", err)
	}

	// Test with one participant
	participants := []Agent{{Name: "Agent1", Role: "Role1", Persona: "Persona1"}}
	_, err = RunDebate(ctx, dummyLLM, participants, nil, "test doc", 1, 0)
	if err == nil {
		t.Error("RunDebate() expected error for one participant, got nil")
	}
	if err.Error() != "need at least two participants" {
		t.Errorf("RunDebate() error = %v, want 'need at least two participants'", err)
	}
}

func TestRunDebateZeroTurnsZeroDuration(t *testing.T) {
	ctx := context.Background()
	dummyLLM := &LLM{} // We won't actually call LLM methods in this test
	participants := []Agent{
		{Name: "Agent1", Role: "Role1", Persona: "Persona1"},
		{Name: "Agent2", Role: "Role2", Persona: "Persona2"},
	}

	history, err := RunDebate(ctx, dummyLLM, participants, nil, "test doc", 0, 0)
	if err != nil {
		t.Fatalf("RunDebate() error = %v", err)
	}

	// Should exit immediately with no turns or duration limit
	if len(history) != 0 {
		t.Errorf("RunDebate() with 0 turns and 0 duration should return empty history, got %d messages", len(history))
	}
}

func TestRunDebateParameterValidation(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		llm          *LLM
		participants []Agent
		expectedErr  string
	}{
		{
			name:         "nil LLM",
			llm:          nil,
			participants: []Agent{{Name: "A1"}, {Name: "A2"}},
			expectedErr:  "llm is nil",
		},
		{
			name:         "no participants",
			llm:          &LLM{},
			participants: []Agent{},
			expectedErr:  "need at least two participants",
		},
		{
			name:         "one participant",
			llm:          &LLM{},
			participants: []Agent{{Name: "A1"}},
			expectedErr:  "need at least two participants",
		},
		{
			name:         "valid parameters",
			llm:          &LLM{},
			participants: []Agent{{Name: "A1"}, {Name: "A2"}},
			expectedErr:  "", // No error expected for parameter validation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RunDebate(ctx, tt.llm, tt.participants, nil, "test doc", 0, 0)

			if tt.expectedErr == "" {
				// For valid parameters, we expect no parameter validation error
				// (though LLM calls will fail later, but that's not what we're testing)
				if err != nil && err.Error() == "llm is nil" {
					t.Errorf("RunDebate() parameter validation failed: %v", err)
				}
				if err != nil && err.Error() == "need at least two participants" {
					t.Errorf("RunDebate() parameter validation failed: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("RunDebate() expected error %q, got nil", tt.expectedErr)
				} else if err.Error() != tt.expectedErr {
					t.Errorf("RunDebate() error = %v, want %v", err, tt.expectedErr)
				}
			}
		})
	}
}

func TestRunDebateEarlyExit(t *testing.T) {
	ctx := context.Background()
	dummyLLM := &LLM{}
	participants := []Agent{
		{Name: "Agent1", Role: "Role1", Persona: "Persona1"},
		{Name: "Agent2", Role: "Role2", Persona: "Persona2"},
	}

	// Test zero turns, zero duration (should exit immediately)
	history, err := RunDebate(ctx, dummyLLM, participants, nil, "test doc", 0, 0)
	if err != nil {
		t.Fatalf("RunDebate() error = %v", err)
	}
	if len(history) != 0 {
		t.Errorf("RunDebate() with 0 turns and 0 duration should return empty history, got %d messages", len(history))
	}

	// Test with very short duration (should exit due to timeout)
	history, err = RunDebate(ctx, dummyLLM, participants, nil, "test doc", 0, 1*time.Nanosecond)
	if err != nil {
		t.Fatalf("RunDebate() error = %v", err)
	}
	// Should have very few or no messages due to immediate timeout
	if len(history) > 1 {
		t.Errorf("RunDebate() with very short duration should return at most 1 message, got %d", len(history))
	}
}
