package agents

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// RunDebate runs a debate between the provided participants in alternating turns.
// The caller decides which participants take part; judge and reporter may be nil.
// RunDebate streams to stdout and returns the message history.
func RunDebate(ctx context.Context, llm *LLM, participants []Agent, judge *Agent, docText string, turns int, duration time.Duration) ([]Message, error) {
	if llm == nil {
		return nil, fmt.Errorf("llm is nil")
	}
	if len(participants) < 2 {
		return nil, fmt.Errorf("need at least two participants")
	}

	log.Println("The trial is now in session...")

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	history := make([]Message, 0)
	started := time.Now()

	for t := 0; ; t++ {
		if turns == 0 && duration == 0 {
			break // nothing to limit the debate
		}

		if turns > 0 && t >= turns {
			break
		}
		if duration > 0 && time.Since(started) > duration {
			break
		}

		speaker := &participants[t%len(participants)]

		arg, _ := speaker.GenerateArgument(ctx, llm, history, docText)

		msg := Message{Sender: speaker.Name, Argument: arg, Time: time.Now()}
		history = append(history, msg)

		fmt.Printf("[%s] %s\n%s\n\n", msg.Time.Format(time.Kitchen), msg.Sender, msg.Argument.Content)

		if judge != nil && rnd.Intn(100) < 15 {
			jarg, err := judge.GenerateArgument(ctx, llm, history, docText)
			if err == nil {
				jmsg := Message{Sender: judge.Name, Argument: jarg, Time: time.Now()}
				history = append(history, jmsg)
				fmt.Printf("[%s] %s\n%s\n\n", jmsg.Time.Format(time.Kitchen), jmsg.Sender, jmsg.Argument.Content)
			}
		}
	}

	if judge != nil {
		history = append(history, Message{Sender: "User", Argument: Argument{Content: "--- END OF DEBATE ---", Tone: "info"}, Time: time.Now()})
		v, err := judge.GenerateArgument(ctx, llm, history, docText)
		if err == nil {
			vmsg := Message{Sender: judge.Name, Argument: v, Time: time.Now()}
			history = append(history, vmsg)
			fmt.Printf("[%s] %s\n%s\n\n", vmsg.Time.Format(time.Kitchen), vmsg.Sender, vmsg.Argument.Content)
		}
	}

	return history, nil
}
