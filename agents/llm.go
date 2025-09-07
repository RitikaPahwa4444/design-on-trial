package agents

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"

	"google.golang.org/genai"
)

type LLM struct {
	client *genai.Client
	Models Models
}

type Models struct {
	Text  string
	Image string
}

func NewLLM(ctx context.Context, model string, imageModel string) (*LLM, error) {
	c, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}
	if model == "" {
		model = "gemini-2.5-flash"
	}
	if imageModel == "" {
		imageModel = "gemini-2.5-flash-image-preview"
	}
	return &LLM{client: c, Models: Models{Text: model, Image: imageModel}}, nil
}

// GetArgument sends the prompt to the configured model and returns the generated text.
func (l *LLM) GetArgument(ctx context.Context, prompt string) (Argument, error) {
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"content": {Type: genai.TypeString},
					"tone":    {Type: genai.TypeString},
				},
				PropertyOrdering: []string{"content", "tone"},
			},
		},
	}
	result, err := l.client.Models.GenerateContent(
		ctx,
		l.Models.Text,
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return Argument{}, err
	}

	var responses []Argument
	if err := json.Unmarshal([]byte(result.Text()), &responses); err != nil {
		log.Fatalf("failed to unmarshal: %v\nraw: %s", err, result.Text())
	}

	return responses[0], nil
}

func (l *LLM) GetComicStrip(ctx context.Context, prompt string) ([]string, error) {
	result, err := l.client.Models.GenerateContent(
		ctx,
		l.Models.Image,
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return nil, err
	}

	var b64imgs []string
	for _, part := range result.Candidates[0].Content.Parts {
		if part.InlineData != nil {
			b64 := base64.StdEncoding.EncodeToString(part.InlineData.Data)
			b64imgs = append(b64imgs, "data:image/png;base64,"+b64)
		}
	}
	return b64imgs, nil
}
