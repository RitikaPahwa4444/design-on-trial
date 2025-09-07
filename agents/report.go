package agents

import (
	"context"
	"fmt"
	"strings"
	"text/template"
	"time"
)

// BuildReportFromLLM orchestrates calls to the LLM to generate a title, comic images
// and the report body, then renders the HTML using the bundled template.
func BuildReportFromLLM(ctx context.Context, llm *LLM, reporter *Agent, docText string, history []Message) (string, error) {
	fmt.Printf("[%s] %s\nGenerating report...\n", time.Now().Format(time.Kitchen), reporter.Name)

	if llm == nil {
		return "", fmt.Errorf("llm is required")
	}

	// Ask the LLM to pick the most important snippet from the transcript.
	// Fall back to a recent-last-n snippet if the request fails.
	fullTranscript := buildTranscriptSnippet(history, len(history))
	chooseSnippetPrompt := fmt.Sprintf("You are an assistant that extracts the most important lines from a debate transcript. Given the full transcript below, return a concise snippet (3-6 lines) that captures the most important arguments and tradeoffs. Return only the snippet, no commentary.\n\nTranscript:\n%s\n", fullTranscript)

	transcriptSnippet := ""
	if out, err := llm.GetArgument(ctx, chooseSnippetPrompt); err == nil {
		transcriptSnippet = strings.TrimSpace(out.Content)
	}

	// 1) Title
	titlePrompt := fmt.Sprintf("You are a concise report title generator. Given the design doc excerpt and the debate transcript, produce a single short title (6-10 words) appropriate for a design trial report.\n\nDesign doc:\n%s\n\nTranscript snippet:\n%s\n", docText, transcriptSnippet)
	title, err := llm.GetArgument(ctx, titlePrompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate title: %w", err)
	}
	titleStr := strings.TrimSpace(title.Content)

	// 2) Comic images: ask for four panels
	comicPrompt := fmt.Sprintf(`
Create a 4-panel comic strip in a satirical courtroom style. 
Requirements:
- Exactly 4 different panels, each 1200x600, arranged for a 2x2 grid.
- Use clean, minimal black-and-white line art.
- Each panel must include clearly readable English speech bubbles (avoid misspellings, random characters, or nonsense text).
- Keep dialogue short (max 6 words per bubble) and directly tied to the debate.
- Context: This is a courtroom debate about a software architecture (High-Level/Low-Level Design).
- Personas: Judge (neutral, witty), Defendant (developer defending design), Prosecutor (architect challenging design).
- Focus on accurate, simple English text inside bubbles over artistic detail.

Design summary (HLD/LLD):
%s

Transcript snippet:
%s
`, docText, transcriptSnippet)
	images, err := llm.GetComicStrip(ctx, comicPrompt)

	if err != nil {
		return "", fmt.Errorf("failed to generate comic images: %w", err)
	}

	// 3) Report body — prefer the reporter agent (scribbler) if provided
	var bodyText string
	if reporter != nil {
		// reporter.GenerateArgument returns an Argument with Content and Tone
		rb, err := reporter.GenerateArgument(ctx, llm, history, docText)
		if err != nil {
			return "", fmt.Errorf("reporter failed to generate report body: %w", err)
		}
		bodyText = rb.Content
	} else {
		reportPrompt := fmt.Sprintf("You are a helpful, structured reporter. Given the title:\n%s\n\nDesign doc:\n%s\n\nFull transcript:\n%s\n\nWrite a concise report (3 short paragraphs) summarizing the debate, key tradeoffs, and a final verdict. Use a neutral professional tone.", title, docText, transcriptSnippet)
		bt, err := llm.GetArgument(ctx, reportPrompt)
		if err != nil {
			return "", fmt.Errorf("failed to generate report body: %w", err)
		}
		bodyText = bt.Content
	}

	var html strings.Builder
	html.WriteString("<html><body style='font-family:sans-serif'>")

	html.WriteString("<h1 style='text-align:center'>" + template.HTMLEscapeString(titleStr) + "</h1>")
	html.WriteString("<p style='font-size:0.9em;color:gray;text-align:center'>Generated on " + time.Now().Format(time.RFC1123) + "</p>")
	html.WriteString("<div style='display:flex;flex-wrap:wrap;justify-content:center;margin-top:1em'>")
	for _, img := range images {
		html.WriteString(fmt.Sprintf("<img src='%s' style='width:45%%;margin:1%%;border:1px solid #ccc;border-radius:8px'/>", img))
	}
	html.WriteString("</div>")
	html.WriteString("<div style='max-width:1000px;margin:2em auto;line-height:1.6;font-size:1em;text-align:justify;column-count:2;column-gap:2em'>")
	html.WriteString(bodyText)
	html.WriteString("</div>")

	html.WriteString("</body></html>")
	return html.String(), nil
}

// buildTranscriptSnippet returns the last n messages as a plain text snippet.
func buildTranscriptSnippet(history []Message, n int) string {
	if n <= 0 {
		return ""
	}
	if len(history) < n {
		n = len(history)
	}
	parts := make([]string, 0, n)
	for i := len(history) - n; i < len(history); i++ {
		if i >= 0 {
			parts = append(parts, fmt.Sprintf("%s: %s", history[i].Sender, truncateLocal(history[i].Argument.Content, 200)))
		}
	}
	return strings.Join(parts, "\n")
}

// local helpers to avoid cross-file dependencies
func truncateLocal(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
