package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RitikaPahwa4444/design-on-trial/server/agents"
)

func main() {
	filePath := flag.String("file", "", "path to the design doc (HLD/LLD)")
	turns := flag.Int("turns", 0, "maximum number of argument turns (each agent counts as one turn)")
	duration := flag.Duration("duration", 0, "max duration for the debate (e.g. 30s, 2m). 0 means no limit")
	model := flag.String("model", "", "LLM model to use for text(optional)")
	imageModel := flag.String("image-model", "", "LLM model to use for images (optional)")

	flag.Parse()

	if *filePath == "" {
		fmt.Fprintln(os.Stderr, "usage: --file /path/to/doc [--turns N] [--duration 30s]")
		os.Exit(2)
	}

	data, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("error reading %s: %v", *filePath, err)
	}
	docText := string(data)

	ctx := context.Background()

	llm, err := agents.NewLLM(ctx, *model, *imageModel)
	if err != nil {
		log.Fatalf("failed to initialize LLM: %v", err)
	}

	personas, err := agents.LoadPersonas()
	if err != nil {
		log.Fatalf("failed to load personas: %v", err)
	}

	// pick judge and reporter from personas by role (caller assigns roles)
	var judge *agents.Agent
	var reporter *agents.Agent
	participants := make([]agents.Agent, 0, len(personas))
	for i := range personas {
		r := strings.ToLower(personas[i].Role)

		if judge == nil && strings.Contains(r, "judge") {
			judge = &personas[i]
			continue
		}
		if reporter == nil && (strings.Contains(r, "report") || strings.Contains(r, "reporter")) {
			reporter = &personas[i]
			continue
		}

		// not judge/reporter (or they were already set), keep as participant
		participants = append(participants, personas[i])
	}
	if len(participants) < 2 {
		// fallback: use all personas if too few participants
		participants = personas
	}

	history, err := agents.RunDebate(ctx, llm, participants, judge, docText, *turns, *duration)
	if err != nil {
		log.Fatalf("debate run failed: %v", err)
	}

	html, err := agents.BuildReportFromLLM(ctx, llm, reporter, docText, history)
	if err != nil {
		log.Fatalf("failed to build report: %v", err)
	}

	outPath, err := agents.WriteReport(filepath.Dir(*filePath), filepath.Base(*filePath), html)
	if err != nil {
		log.Fatalf("failed to write report: %v", err)
	}
	fmt.Printf("Report written to %s\n", outPath)
}
