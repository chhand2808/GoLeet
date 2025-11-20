package cmd

import (
	"fmt"
	"strings"

	"github.com/chhand2808/goleet/internal/data"
	"github.com/chhand2808/goleet/internal/gemini"
	utils "github.com/chhand2808/goleet/internal/util"
	"github.com/spf13/cobra"
)

var suggestCmd = &cobra.Command{
	Use:   "suggest",
	Short: "Suggests a new LeetCode problem using Gemini AI",
	Run: func(cmd *cobra.Command, args []string) {
		runAISuggest(cmd)
	},
}

func init() {
	rootCmd.AddCommand(suggestCmd)

	// add debug flag
	suggestCmd.Flags().Bool("debug", false, "Enable debug logging")
}

func runAISuggest(cmd *cobra.Command) {
	debugFlag, _ := cmd.Flags().GetBool("debug")
	if debugFlag {
		utils.DebugEnabled = true
		utils.IsProduction = false
	} else {
		utils.IsProduction = true
	}

	utils.Info("Starting AI suggestion flow...")

	store := data.NewStore()

	// Load problems
	problems, err := store.LoadProblems()
	if err != nil {
		utils.Error("Failed to load problems: %v", err)
		return
	}
	utils.Debug("Loaded %d problems", len(problems))

	// Load history + solved
	history, _ := store.LoadHistory()
	solved, _ := store.LoadSolved()
	utils.Debug("Loaded %d solved, %d history", len(solved), len(history))

	// seriousness = how strict the AI should be
	seriousness := 1
	var final []data.AISuggestion

	for attempt := 1; attempt <= 3; attempt++ {
		utils.Info("Gemini Call Attempt %d (seriousness=%d)", attempt, seriousness)

		prompt := gemini.BuildPrompt(solved, history, problems, seriousness)
		utils.Debug("PROMPT SENT TO GEMINI:\n%s", prompt)

		// Start spinner ONLY in production mode
		var stop chan bool
		if utils.IsProduction && !utils.DebugEnabled {
			stop = utils.StartSpinner()
		}

		ai, err := gemini.GetSuggestions(prompt)

		// Stop spinner
		if utils.IsProduction && !utils.DebugEnabled {
			utils.StopSpinner(stop)
			fmt.Println() // move to next line
		}

		if err != nil {
			utils.Warn("Gemini error: %v", err)
			seriousness++
			continue
		}

		utils.Debug("Gemini returned %d suggestions", len(ai))

		// filter invalid ones
		valid := filterAISuggestions(ai, solved, history, problems)
		utils.Debug("%d suggestions valid after filtering", len(valid))

		if len(valid) > 0 {
			final = valid
			utils.Info("Found valid AI suggestions.")
			break
		}

		utils.Warn("No valid suggestions, retrying with higher seriousness...")
		seriousness++
	}

	if len(final) == 0 {
		utils.Error("No AI suggestions available after retries")
		fmt.Println("⚠️ No AI suggestions available.")
		return
	}

	// pick the first suggestion
	chosen := final[0]

	utils.Info("Chosen suggestion: %d - %s", chosen.Number, chosen.Title)

	fmt.Println("🧠 AI Suggested:")
	fmt.Printf("%d. %s\n", chosen.Number, chosen.Title)
	fmt.Println("Topics:", chosen.Topics)
	fmt.Printf(
		"Link: https://leetcode.com/problems/%s/\n",
		strings.ToLower(strings.ReplaceAll(chosen.Title, " ", "-")),
	)

	// Save history
	err = store.AppendHistory(
		data.NewHistoryEntry(
			fmt.Sprint(chosen.Number),
			chosen.Title,
			"",
		),
		10,
	)

	if err != nil {
		utils.Warn("Failed to update history: %v", err)
	} else {
		utils.Info("History updated successfully.")
	}
}

func filterAISuggestions(
	ai []data.AISuggestion,
	solved []data.SolvedProblem,
	history []data.HistoryEntry,
	problems []data.Problem,
) []data.AISuggestion {

	block := map[string]bool{}

	// block solved
	for _, s := range solved {
		block[s.ID] = true
	}

	// block recent history
	for _, h := range history {
		block[h.ID] = true
	}

	out := []data.AISuggestion{}

	for _, p := range ai {
		id := fmt.Sprint(p.Number)

		// skip if solved/history
		if block[id] {
			continue
		}

		// ensure exists in problem db
		if !existsInProblemDB(id, problems) {
			utils.Debug("AI suggested problem %s not found in DB", id)
			continue
		}

		out = append(out, p)
	}

	return out
}

func existsInProblemDB(id string, problems []data.Problem) bool {
	for _, p := range problems {
		if p.ID == id {
			return true
		}
	}
	return false
}
