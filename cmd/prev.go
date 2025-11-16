package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/chhand2808/goleet/internal/data"
	"github.com/spf13/cobra"
)

var prevCmd = &cobra.Command{
	Use:   "prev [n]",
	Short: "Show previously suggested problems. Default n=1 (max 10)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		n := 1
		if len(args) == 1 {
			if v, err := strconv.Atoi(args[0]); err == nil {
				n = v
			}
		}
		if n <= 0 {
			fmt.Println("Please provide a positive number")
			return
		}
		if n > 10 {
			n = 10
		}
		showPrev(n)
	},
}

func init() {
	rootCmd.AddCommand(prevCmd)
}

func showPrev(n int) {
	store := data.NewStore()
	hist, err := store.LoadHistory()
	if err != nil {
		// If we returned an error but also an empty history, continue; otherwise show message
		fmt.Println("❌ Failed to load history:", err)
		// still try to proceed if hist not nil
	}

	if len(hist) == 0 {
		fmt.Println("No previously suggested problems found. Try running: suggest")
		return
	}

	// We store most recent at the end, so latestFirst = reverse order
	total := len(hist)
	toShow := n
	if toShow > total {
		toShow = total
	}

	fmt.Println("Recent Suggested Problems:")
	// show latest first
	for i := 0; i < toShow; i++ {
		idx := total - 1 - i // latest at end -> show it first
		entry := hist[idx]
		dateStr := formatRelativeDate(entry.Date)
		fmt.Printf("%d. %s (%s)\n", i+1, entry.Title, dateStr)
	}
}

// formatRelativeDate performs D3 formatting:
// Today / Yesterday / DD Mon YYYY
func formatRelativeDate(dateISO string) string {
	if dateISO == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02", dateISO)
	if err != nil {
		return dateISO
	}
	now := time.Now()
	y, m, d := now.Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	y1, m1, d1 := t.Date()
	tt := time.Date(y1, m1, d1, 0, 0, 0, 0, now.Location())

	diff := int(today.Sub(tt).Hours() / 24)
	switch diff {
	case 0:
		return "Today"
	case 1:
		return "Yesterday"
	default:
		return tt.Format("02 Jan 2006")
	}
}
