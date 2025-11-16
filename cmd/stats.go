package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/chhand2808/goleet/internal/data"

	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show your solving stats",
	Run: func(cmd *cobra.Command, args []string) {
		showStats()
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

func showStats() {
	// Load solved problems
	solved, err := loadSolved()
	if err != nil {
		fmt.Println("❌ Failed to load solved problems:", err)
		return
	}

	// Load all problems (for difficulty count)
	allProblems, err := loadAllProblems()
	if err != nil {
		fmt.Println("❌ Failed to load problems.json:", err)
		return
	}

	totalSolved := len(solved)
	easy, medium, hard := 0, 0, 0

	// Count difficulty
	for _, s := range solved {
		if p, ok := allProblems[s.ID]; ok {
			switch p.Difficulty {
			case "Easy":
				easy++
			case "Medium":
				medium++
			case "Hard":
				hard++
			}
		}
	}

	currentStreak, longestStreak := calculateStreak(solved)

	drawBoxedStats(totalSolved, easy, medium, hard, currentStreak, longestStreak)
}

func loadSolved() ([]data.SolvedProblem, error) {
	file, err := os.ReadFile("data/solved.json")
	if err != nil {
		// No file yet? Return empty list
		return []data.SolvedProblem{}, nil
	}

	var solved []data.SolvedProblem
	err = json.Unmarshal(file, &solved)
	return solved, err
}

func loadAllProblems() (map[string]data.Problem, error) {
	file, err := os.ReadFile("data/problems.json")
	if err != nil {
		return nil, err
	}

	var problems []data.Problem
	err = json.Unmarshal(file, &problems)
	if err != nil {
		return nil, err
	}

	problemMap := make(map[string]data.Problem)
	for _, p := range problems {
		problemMap[p.ID] = p
	}
	return problemMap, nil
}

func calculateStreak(solved []data.SolvedProblem) (int, int) {
	if len(solved) == 0 {
		return 0, 0
	}

	// Extract unique dates
	dateMap := map[string]bool{}
	for _, s := range solved {
		dateMap[s.Date] = true
	}

	// Convert to sorted slice
	var dates []time.Time
	for d := range dateMap {
		t, err := time.Parse("2006-01-02", d)
		if err == nil {
			dates = append(dates, t)
		}
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	currentStreak, longestStreak := 1, 1

	for i := 1; i < len(dates); i++ {
		diff := dates[i].Sub(dates[i-1]).Hours() / 24

		if diff == 1 {
			currentStreak++
		} else {
			if currentStreak > longestStreak {
				longestStreak = currentStreak
			}
			currentStreak = 1
		}
	}

	if currentStreak > longestStreak {
		longestStreak = currentStreak
	}

	// If last solve wasn't yesterday or today, streak is broken
	today := time.Now().Truncate(24 * time.Hour)
	lastSolve := dates[len(dates)-1]

	if lastSolve.Before(today.Add(-24 * time.Hour)) {
		currentStreak = 0
	}

	return currentStreak, longestStreak
}

func drawBoxedStats(total, easy, medium, hard, current, longest int) {
	top := "╔══ 📊 STATS ═════════════════════════╗"
	btm := "╚════════════════════════════════════╝"

	fmt.Println()
	fmt.Println(top)
	fmt.Printf("║ Total Solved        : %-12d ║\n", total)
	fmt.Printf("║ Easy / Med / Hard   : %d / %d / %d      ║\n", easy, medium, hard)
	fmt.Printf("║ 🔥 Current Streak    : %-12d ║\n", current)
	fmt.Printf("║ 🏆 Longest Streak     : %-12d ║\n", longest)
	fmt.Println(btm)
	fmt.Println()
}
