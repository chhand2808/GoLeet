package storage

import (
	"encoding/json"
	"os"

	"github.com/chhand2808/goleet/internal/model"
)

func LoadProblems() ([]model.Problem, error) {
	file, err := os.Open("data/problems.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var problems []model.Problem
	err = json.NewDecoder(file).Decode(&problems)
	return problems, err
}
