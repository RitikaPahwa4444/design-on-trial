package agents

import (
	"encoding/json"
	"testing"
	"time"
)

func TestAgentStruct(t *testing.T) {
	agent := Agent{
		Name:    "Test Agent",
		Role:    "Test Role",
		Persona: "Test Persona Description",
	}

	if agent.Name != "Test Agent" {
		t.Errorf("Expected name 'Test Agent', got %s", agent.Name)
	}
	if agent.Role != "Test Role" {
		t.Errorf("Expected role 'Test Role', got %s", agent.Role)
	}
	if agent.Persona != "Test Persona Description" {
		t.Errorf("Expected persona 'Test Persona Description', got %s", agent.Persona)
	}
}

func TestAgentsStruct(t *testing.T) {
	agents := Agents{
		Agents: []Agent{
			{Name: "Agent1", Role: "Role1", Persona: "Persona1"},
			{Name: "Agent2", Role: "Role2", Persona: "Persona2"},
		},
	}

	if len(agents.Agents) != 2 {
		t.Errorf("Expected 2 agents, got %d", len(agents.Agents))
	}
	if agents.Agents[0].Name != "Agent1" {
		t.Errorf("Expected first agent name 'Agent1', got %s", agents.Agents[0].Name)
	}
}

func TestArgumentStruct(t *testing.T) {
	arg := Argument{
		Content: "Test content",
		Tone:    "neutral",
	}

	if arg.Content != "Test content" {
		t.Errorf("Expected content 'Test content', got %s", arg.Content)
	}
	if arg.Tone != "neutral" {
		t.Errorf("Expected tone 'neutral', got %s", arg.Tone)
	}
}

func TestMessageStruct(t *testing.T) {
	now := time.Now()
	arg := Argument{Content: "Test message", Tone: "friendly"}
	msg := Message{
		Sender:   "Test Sender",
		Argument: arg,
		Time:     now,
	}

	if msg.Sender != "Test Sender" {
		t.Errorf("Expected sender 'Test Sender', got %s", msg.Sender)
	}
	if msg.Argument.Content != "Test message" {
		t.Errorf("Expected argument content 'Test message', got %s", msg.Argument.Content)
	}
	if !msg.Time.Equal(now) {
		t.Errorf("Expected time %v, got %v", now, msg.Time)
	}
}

func TestAgentsJSONSerialization(t *testing.T) {
	agents := Agents{
		Agents: []Agent{
			{Name: "Agent1", Role: "Role1", Persona: "Persona1"},
		},
	}

	// Test marshaling
	jsonData, err := json.Marshal(agents)
	if err != nil {
		t.Fatalf("Failed to marshal agents: %v", err)
	}

	// Test unmarshaling
	var unmarshaledAgents Agents
	err = json.Unmarshal(jsonData, &unmarshaledAgents)
	if err != nil {
		t.Fatalf("Failed to unmarshal agents: %v", err)
	}

	if len(unmarshaledAgents.Agents) != 1 {
		t.Errorf("Expected 1 agent after JSON round-trip, got %d", len(unmarshaledAgents.Agents))
	}
	if unmarshaledAgents.Agents[0].Name != "Agent1" {
		t.Errorf("Expected agent name 'Agent1' after JSON round-trip, got %s", unmarshaledAgents.Agents[0].Name)
	}
}

func TestArgumentJSONSerialization(t *testing.T) {
	arg := Argument{Content: "Test content", Tone: "neutral"}

	// Test marshaling
	jsonData, err := json.Marshal(arg)
	if err != nil {
		t.Fatalf("Failed to marshal argument: %v", err)
	}

	// Test unmarshaling
	var unmarshaledArg Argument
	err = json.Unmarshal(jsonData, &unmarshaledArg)
	if err != nil {
		t.Fatalf("Failed to unmarshal argument: %v", err)
	}

	if unmarshaledArg.Content != "Test content" {
		t.Errorf("Expected content 'Test content' after JSON round-trip, got %s", unmarshaledArg.Content)
	}
	if unmarshaledArg.Tone != "neutral" {
		t.Errorf("Expected tone 'neutral' after JSON round-trip, got %s", unmarshaledArg.Tone)
	}
}

func TestMessageJSONSerialization(t *testing.T) {
	now := time.Now()
	arg := Argument{Content: "Test message", Tone: "friendly"}
	msg := Message{
		Sender:   "Test Sender",
		Argument: arg,
		Time:     now,
	}

	// Test marshaling
	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	// Test unmarshaling
	var unmarshaledMsg Message
	err = json.Unmarshal(jsonData, &unmarshaledMsg)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	if unmarshaledMsg.Sender != "Test Sender" {
		t.Errorf("Expected sender 'Test Sender' after JSON round-trip, got %s", unmarshaledMsg.Sender)
	}
	if unmarshaledMsg.Argument.Content != "Test message" {
		t.Errorf("Expected argument content 'Test message' after JSON round-trip, got %s", unmarshaledMsg.Argument.Content)
	}
	// Note: Time precision may be lost in JSON round-trip, so we check within a second
	if unmarshaledMsg.Time.Unix() != now.Unix() {
		t.Errorf("Expected time %v after JSON round-trip, got %v", now, unmarshaledMsg.Time)
	}
}
