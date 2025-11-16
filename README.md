📦 GoLeet – A Smart LeetCode CLI Suggestion Tool (Built in Go)

A lightweight, fast, intelligent CLI tool that suggests LeetCode problems based on your history.

🌟 Features

⚡ Daily smart problem suggestions

🤖 AI-powered (Gemini) recommendations

📊 Track solved problems + difficulty stats

🔥 Streak tracking (current & longest)

📝 History of last 10 suggested questions

💾 Local JSON storage (no internet needed except AI calls)

🎯 Topic & difficulty filters

🚀 Self-contained binary (powered by go:embed)

📥 Installation
Windows

Download the latest goleet.exe from Releases and place it anywhere.

Linux / macOS

Download the appropriate binary:

chmod +x goleet
./goleet

Or install using Go:
go install github.com/YOUR_USERNAME/goleet@latest

🚀 Getting Started
1) Initialize the tool

This will create your data folder, config file, and ask for Gemini API key.

goleet init

2) Get a suggested problem
goleet suggest

3) Mark a problem as solved
goleet done 1

4) View your stats
goleet stats

5) View previous suggestions

Show last suggestion:

goleet prev


Show last 3:

goleet prev 3

📚 Commands Overview
Command	Description
goleet init	Setup config + embed problems.json
goleet suggest	Suggest a new LeetCode problem
goleet suggest --difficulty Easy	Filter by difficulty
goleet suggest --topic array	Filter by topic
goleet done <id>	Mark a problem solved
goleet stats	Total solved, difficulty stats, streaks
goleet prev [n]	View previous suggestions (max 10)
goleet update	(Coming soon) Auto-update the CLI
📊 Example Output
Stats Screen
╔══ 📊 STATS ═════════════════════════╗
║ Total Solved        : 27            ║
║ Easy / Med / Hard   : 15 / 10 / 2   ║
║ 🔥 Current Streak    : 5             ║
║ 🏆 Longest Streak     : 7             ║
╚════════════════════════════════════╝

Suggestion Example
🧠 Today's Suggested Problem:
1. Two Sum (Easy)
Topics: [Array Hash Table]
Link: https://leetcode.com/problems/two-sum/

🛠️ Tech Stack

Go 1.22+

Cobra – CLI framework

Go Embed – Static asset embedding

JSON Storage – Local problem db + history + solved

Gemini API – AI-powered suggestions

📁 Project Structure
goleet/
│
├── cmd/               # CLI commands
├── data/              # Embedded problems.json
├── internal/          # Core logic (storage, model, AI)
├── main.go            # Entry point
└── go.mod

🔮 Roadmap / Upcoming Features

✔ Auto-update CLI (goleet update)

✔ Weekly report summaries

⏳ GitHub Gist Sync (cloud history backup)

⏳ AI-based difficulty progression

⏳ Topic mastery analytics

⏳ VS Code extension

🤝 Contributing

Contributions are welcome!
Feel free to open issues or submit PRs.

📄 License

MIT License © 2025
Your Name (Chhand Kunal Chaughule)

⭐ Support

If you like the project, consider giving it a ⭐ on GitHub — it motivates development!