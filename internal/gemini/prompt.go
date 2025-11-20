package gemini

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/chhand2808/goleet/internal/data"
)

//
// =============== PUBLIC API ===============
//

func BuildPrompt(
	solved []data.SolvedProblem,
	history []data.HistoryEntry,
	allProblems []data.Problem,
	seriousness int,
) string {

	// 1️⃣ Compute dynamic streak
	currentStreak := computeStreak(solved)

	// 2️⃣ Compute weak topics
	weakTopics := computeWeakTopics(allProblems, solved)

	// 3️⃣ Difficulty guidance driven by streak level
	difficultyAdvice := difficultyBasedOnStreak(currentStreak)

	// 4️⃣ Format solved lines
	solvedLines := []string{}
	for _, s := range solved {
		solvedLines = append(solvedLines, fmt.Sprintf("%s | %s", s.ID, s.Title))
	}

	// 5️⃣ Format recent history IDs
	historyIDs := []string{}
	for _, h := range history {
		historyIDs = append(historyIDs, h.ID)
	}

	// 6️⃣ Build final prompt
	return fmt.Sprintf(`You are GoLeet AI. Suggest exactly 3 new LeetCode problems in JSON ONLY.

USER_SOLVED (never repeat):
%s

RECENT_SUGGESTIONS (avoid repeating):
%s

USER_STATS:
Current_Streak: %d days
Weak_Topics: %s

STREAK_BASED_DIFFICULTY_GUIDANCE:
%s

GOAL:
- Suggest 3 UNSOLVED, NEW problems.
- Follow difficulty guidance above.
- Prefer weak topics moderately.
- Ensure variety (avoid repeating topics too much).
- Avoid all solved and history items.
- NO explanations. NO text outside JSON.

STRICT JSON OUTPUT FORMAT:
[
  {"title": "", "number": 0, "topics": ["",""]},
  {"title": "", "number": 0, "topics": ["",""]},
  {"title": "", "number": 0, "topics": ["",""]}
]

SERIOUSNESS_MODE: %d
(1 = normal, 2 = stronger diversity, 3 = strict filtering)

Return ONLY the JSON array.
`,
		strings.Join(solvedLines, "\n"),
		strings.Join(historyIDs, ", "),
		currentStreak,
		strings.Join(weakTopics, ", "),
		difficultyAdvice,
		seriousness,
	)
}

//
// =============== HELPERS ===============
//

// Computes the consecutive-day streak based on solved.json
func computeStreak(solved []data.SolvedProblem) int {
	if len(solved) == 0 {
		return 0
	}

	// Sort by date descending
	sort.Slice(solved, func(i, j int) bool {
		return solved[i].Date > solved[j].Date
	})

	format := "2006-01-02"
	streak := 1

	lastDate, err := time.Parse(format, solved[0].Date)
	if err != nil {
		return 1
	}

	for i := 1; i < len(solved); i++ {
		curDate, err := time.Parse(format, solved[i].Date)
		if err != nil {
			continue
		}

		diff := lastDate.Sub(curDate).Hours()

		if diff <= 24 && diff >= 0 {
			streak++
			lastDate = curDate
		} else {
			break
		}
	}

	return streak
}

// Difficulty personality based on streak
func difficultyBasedOnStreak(streak int) string {
	switch {
	case streak < 3:
		return "User is early in streak. Prefer EASY problems (some MEDIUM allowed)."
	case streak < 10:
		return "User is mid-streak. Suggest a balanced mix of EASY and MEDIUM."
	default:
		return "User is on a strong streak. Suggest MEDIUM and MEDIUM-HARD challenges."
	}
}

// Computes weak topics (least solved vs most available)
func computeWeakTopics(allProblems []data.Problem, solved []data.SolvedProblem) []string {

	// Count total availability per topic
	totalCount := map[string]int{}
	for _, p := range allProblems {
		for _, t := range p.TopicTags {
			totalCount[t.Name]++
		}
	}

	// Count solved per topic
	solvedCount := map[string]int{}
	for _, s := range solved {
		for _, p := range allProblems {
			if p.ID == s.ID {
				for _, t := range p.TopicTags {
					solvedCount[t.Name]++
				}
			}
		}
	}

	// Calculate weakness score = solved / total
	type pair struct {
		Topic string
		Score float64
	}

	scoreList := []pair{}
	for topic, total := range totalCount {
		solved := solvedCount[topic]
		score := float64(solved) / float64(total) // low = weak
		scoreList = append(scoreList, pair{topic, score})
	}

	// Sort by increasing score (weakest first)
	sort.Slice(scoreList, func(i, j int) bool {
		return scoreList[i].Score < scoreList[j].Score
	})

	// Pick top 3 weakest topics
	limit := 3
	if len(scoreList) < limit {
		limit = len(scoreList)
	}

	weakTopics := []string{}
	for i := 0; i < limit; i++ {
		weakTopics = append(weakTopics, scoreList[i].Topic)
	}

	return weakTopics
}
