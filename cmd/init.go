package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chhand2808/goleet/data"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Setup GoLeet for first-time use",
	Long:  "Initializes GoLeet by storing your Gemini API key and creating required data files",
	Run: func(cmd *cobra.Command, args []string) {
		err := InitConfig()
		if err != nil {
			fmt.Println("❌ Failed to initialize:", err)
		} else {
			fmt.Println("✅ GoLeet successfully initialized!")
		}
	},
}

func InitConfig() error {
	configDir := "data"

	// Create data directory if missing
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.Mkdir(configDir, 0755)
		if err != nil {
			return err
		}
	}

	// Prompt for API Key
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your Gemini API Key: ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)

	// Save config.json
	configPath := filepath.Join(configDir, "config.json")
	err := os.WriteFile(configPath, []byte(fmt.Sprintf(`{"api_key": "%s"}`, apiKey)), 0644)
	if err != nil {
		return err
	}

	// ✅ Write embedded problems.json (only if not exists)
	problemsPath := filepath.Join(configDir, "problems.json")
	if _, err := os.Stat(problemsPath); os.IsNotExist(err) {
		err = os.WriteFile(problemsPath, data.EmbeddedProblems, 0644)
		if err != nil {
			return err
		}
	}

	// Create solved.json if missing
	solvedPath := filepath.Join(configDir, "solved.json")
	if _, err := os.Stat(solvedPath); os.IsNotExist(err) {
		err = os.WriteFile(solvedPath, []byte("[]"), 0644)
		if err != nil {
			return err
		}
	}

	// Create history.json if missing
	historyPath := filepath.Join(configDir, "history.json")
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		err = os.WriteFile(historyPath, []byte("[]"), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
