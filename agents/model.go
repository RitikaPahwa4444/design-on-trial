package agents

import "time"

type Agents struct {
	Agents []Agent `json:"agents"`
}

type Agent struct {
	Name    string `json:"name"`
	Role    string `json:"role"`
	Persona string `json:"persona"`
}

type Message struct {
	Sender   string    `json:"sender"`
	Argument Argument  `json:"argument"`
	Time     time.Time `json:"time"`
}

type Argument struct {
	Content string `json:"content"`
	Tone    string `json:"tone"`
}
