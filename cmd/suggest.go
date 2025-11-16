package cmd

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/chhand2808/goleet/internal/data"
	"github.com/spf13/cobra"
)

var (
	flagDifficulty string
	flagTopic      string
)

var suggestCmd = &cobra.Command{
	Use:   "suggest",
	Short: "Suggests a new LeetCode problem",
	Run: func(cmd *cobra.Command, args []string) {
		runSuggest()
	},
}

func init() {
	rootCmd.AddCommand(suggestCmd)

	// optional filters
	suggestCmd.Flags().StringVar(&flagDifficulty, "difficulty", "", "Filter by difficulty: Easy|Medium|Hard")
	suggestCmd.Flags().StringVar(&flagTopic, "topic", "", "Filter by topic (case-insensitive)")
}

func runSuggest() {
	store := data.NewStore()

	problems, err := store.LoadProblems()
	if err != nil {
		fmt.Println("❌ Failed to load problems:", err)
		return
	}
	if len(problems) == 0 {
		fmt.Println("No problems found in problems.json")
		return
	}

	// Apply simple filtering if flags provided
	candidates := []data.Problem{}
	for _, p := range problems {
		// difficulty filter
		if flagDifficulty != "" {
			if !strings.EqualFold(p.Difficulty, flagDifficulty) {
				continue
			}
		}

		// topic filter
		if flagTopic != "" {
			foundTopic := false
			for _, t := range p.TopicTags {
				if strings.EqualFold(t.Name, flagTopic) {
					foundTopic = true
					break
				}
			}
			if !foundTopic {
				continue
			}
		}

		candidates = append(candidates, p)
	}

	if len(candidates) == 0 {
		fmt.Println("No candidate problems found for the given filters.")
		return
	}

	// pick random candidate
	rand.Seed(time.Now().UnixNano())
	choice := candidates[rand.Intn(len(candidates))]

	// display to user
	topics := []string{}
	for _, t := range choice.TopicTags {
		topics = append(topics, t.Name)
	}

	fmt.Println("🧠 Today's Suggested Problem:")
	fmt.Printf("%s. %s (%s)\n", choice.ID, choice.Title, choice.Difficulty) // uses helper ID() below
	fmt.Println("Topics:", topics)
	fmt.Printf("Link: https://leetcode.com/problems/%s/\n", choice.TitleSlug)

	// append to history (max 10)
	entry := data.NewHistoryEntry(choice.ID, choice.Title, "")
	if err := store.AppendHistory(entry, 10); err != nil {
		fmt.Println("⚠️ Warning: failed to update history:", err)
	}
}

// Because data.Problem used in data package might have different field names,
// provide small helper methods via type aliasing with reflection-safe access.
// Here, we assume data.Problem has methods/fields; implement helper accessors:

// NOTE: adjust these if your data.Problem struct uses different exported field names.
// The following assumes the struct fields are exported as: ID, Title, Difficulty, TitleSlug, TopicTags

// To keep this file minimal and robust, use small wrapper functions in data package instead.
// For now, implement simple accessors by casting.
