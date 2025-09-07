
# design-on-trial

Your HLDs and LLDs, put on trial. Why? 

Because design reviews are stressful. Why not turn them into courtroom drama and comics?

## ✨ Features
- AI-Powered Debate Simulation: multiple personas argue about your HLD/LLD, all in your terminal.
- Newspaper-Style Report: two-column summary with pros, cons, and verdict, exported as HTML.
- Supported file formats:
    - `.txt`
    - `.md`

## 🛠️ Why It’s Useful (and Fun)

- Structured Feedback: Developers get a clear report of pros, cons, and final verdict on their design choices.
- Bias-Free Debate: The Judge persona ensures balanced reasoning instead of hand-wavy arguments.
- Sharable Artifacts: Teams can use the reports and comics in docs, retros, or presentations. Bring storytelling to life! 

## 🎬 Trailer

1. Gemini 2.5 Pro (Slow but more reliable)

https://github.com/user-attachments/assets/e1eadf35-6a17-46d2-af94-e84361006728

2. Gemini 2.0 Flash (Fast but less reliable)
   
https://github.com/user-attachments/assets/d0f8e3fb-b79f-4fef-b1f3-abdd1f8b6e65

## 📝 The Fine Print
This project uses Gemini 2.5 Flash Image to generate consistent multi-panel courtroom characters (Judge, Prosecutor, Defendant) in the form of a comic strip that can be used in retros, presentations or light-hearted Slack messages.

## 🏁 Getting Started

### Using the pre-built binary

1. Download the appropriate release from GitHub Releases:
    - Visit https://github.com/RitikaPahwa4444/design-on-trial/releases and download the pre-built binary

2. Export your Gemini API key:
```bash
export GEMINI_API_KEY="your-key-here"
```

3. Make the binary executable (if needed) and run:

```bash
chmod +x design-on-trial
./design-on-trial --file ../../sample_hld.md --duration 1m
```

Replace ../../sample_hld.md with the path to your design doc.

---

### Building from source

Prerequisites:
- Go 1.24+
- Google GenAI Go SDK 

1. Clone this repository:

```bash
git clone https://github.com/RitikaPahwa4444/design-on-trial.git
cd design-on-trial
```

2. Install dependencies (example for GenAI SDK):

```bash
go get google.golang.org/genai
```

3. Export your Gemini API key:
```bash
export GEMINI_API_KEY="your-key-here"
```

4. Run the trial:

```bash
cd server/cmd
go run . --file ../../sample_hld.md --duration 1m
```

---

### Usage

Flags:
- `--file` (required): path to the design document (HLD/LLD). Example: `--file ../../sample_hld.md`. The path can be an absolute or relative. The CLI writes the report next to the input file and prints the report path.
- `--turns`: maximum number of argument turns (each agent counts as one turn).
- `--duration`: max duration for the debate (e.g. `30s`, `2m`).
- `--model`: optional LLM model name for text.
- `--image-model`: optional LLM model name for images.

Example trial run for 1m:

```bash
cd server/cmd
go run . --file ../../sample_hld.md --duration 1m --model gemini-2.5-flash --image-model gemini-2.5-flash-image-preview
```

### 👩🏻‍💻 Development notes

- Personas live in `server/agents/persona.json` (edit to add or tweak personas, not everyone loves a discussion as intense as a courtroom 😉).

## License

MIT. Because sharing is caring.
