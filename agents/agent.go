package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var personaJSON []byte

// GenerateArgument generates an argument based on the conversation history and provided PDF text.
// Returns the generated argument as a string.
// GenerateArgument asks the LLM to return JSON with {"argument":"...","tone":"..."}.
// It returns the argument and error.
func (agent *Agent) GenerateArgument(ctx context.Context, llm *LLM, history []Message, pdfText string) (Argument, error) {
	var conv strings.Builder
	for _, msg := range history {
		fmt.Fprintf(&conv, "%s (%s): %s\n", msg.Sender, msg.Time.Format(time.Kitchen), msg.Argument)
	}

	prompt := fmt.Sprintf(`%s
You are participating in a design debate. Use simple conversational English.
Conversation so far:
%s

Additional context from document:
%s

Return exactly one JSON object with two fields: {"content": "...", "tone": "..."}. Keep the content concise, ideally under 50 words.
Tone should be a short label like: neutral, persuasive, conciliatory, critical, heated, sarcastic, analytical, insulting, etc.
If however, you are a reporter agent, generate a very detailed report/content (3-4 paragraphs) summarizing the debate, key tradeoffs, and a final verdict.
If you find --- END OF DEBATE --- in the conversation history, summarise the history and give the final verdict as the argument.
Only return the JSON object and nothing else.
`, agent.Persona, conv.String(), pdfText)

	out, err := llm.GetArgument(ctx, prompt)
	if err != nil {
		return Argument{}, err
	}
	return out, nil
}

func LoadPersonas() ([]Agent, error) {
	// Try multiple paths to find persona.json
	possiblePaths := []string{
		"agents/persona.json",    // From project root
		"../agents/persona.json", // Original path for cmd directory
		"persona.json",           // Same directory as executable
		"./persona.json",         // Current working directory
	}

	// If we can determine the executable path, also try looking relative to it
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		possiblePaths = append(possiblePaths, filepath.Join(execDir, "persona.json"))
	}

	var personaJSON []byte
	var err error
	var foundPath string

	for _, path := range possiblePaths {
		personaJSON, err = os.ReadFile(path)
		if err == nil {
			foundPath = path
			break
		}
	}

	if foundPath == "" {
		return nil, fmt.Errorf("failed to find persona.json in any of the expected locations: %v", possiblePaths)
	}

	var agents Agents
	if err := json.Unmarshal(personaJSON, &agents); err != nil {
		return nil, fmt.Errorf("failed to parse persona.json from %s: %w", foundPath, err)
	}
	return agents.Agents, nil
}
