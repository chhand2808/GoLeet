package cmd

import (
	"fmt"

	"github.com/chhand2808/goleet/internal/data"

	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done [questionID]",
	Short: "Mark a question as solved",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		questionID := args[0]
		store := data.NewStore()

		// Load all problems
		problems, err := store.LoadProblems()
		if err != nil {
			fmt.Println("❌ Failed to load problems:", err)
			return
		}

		// Find problem by ID
		var title string
		found := false
		for _, p := range problems {
			if p.ID == questionID {
				title = p.Title
				found = true
				break
			}
		}

		if !found {
			fmt.Println("⚠️ Problem ID not found:", questionID)
			return
		}

		// Mark as solved
		err = store.MarkSolved(questionID, title)
		if err != nil {
			fmt.Println("❌ Failed to mark as solved:", err)
			return
		}

		fmt.Printf("✅ Marked as solved: %s (%s)\n", title, questionID)
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
