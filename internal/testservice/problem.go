package testservice

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/maps"
	"log"
	"strings"
)

type ProblemStatus int

const (
	ProblemBatchSize                = 2
	StatusNotAnswered ProblemStatus = iota
	StatusCorrectAnswer
	StatusIncorrectAnswer
)

type Problems struct {
	problemCatalog map[int64]*Problem
	lastID         int64
	lg             *log.Logger
}

type Problem struct {
	ID        int64  `json:"id"`
	Question  string `json:"question"`
	ImageLink string `json:"image"`
	Answer    string `json:"answer"`
}

func (s *Problems) ParseProblems(data []byte) error {
	var taskList []*Problem
	err := json.Unmarshal(data, &taskList)
	if err != nil {
		return fmt.Errorf("failed to parse problemCatalog: %w", err)
	}

	s.problemCatalog = make(map[int64]*Problem, len(taskList))

	for _, task := range taskList {
		if task.ID < 1 {
			s.lg.Printf("skipping question with negative ID: %d\n", task.ID)
			continue
		}
		s.problemCatalog[task.ID] = task
	}

	return nil
}

func (s *Problems) Assert(questionID int64, userAnswer string) (bool, error) {
	questionById, found := s.problemCatalog[questionID]
	if !found {
		return false, fmt.Errorf("question %d not found", questionID)
	}
	corr := strings.ToLower(questionById.Answer)
	user := strings.ToLower(userAnswer)
	return corr == user, nil
}

func (s *Problems) GetProblem(problemID int64) (p *Problem, found bool) {
	p, found = s.problemCatalog[problemID]
	return
}

func (s *Problems) generateProblemBatch() (res map[int64]ProblemStatus) {
	res = make(map[int64]ProblemStatus)

	if ProblemBatchSize >= len(s.problemCatalog) {
		for _, id := range maps.Keys(s.problemCatalog) {
			res[id] = StatusNotAnswered
		}
	} else {
		for _, id := range maps.Keys(s.problemCatalog)[:ProblemBatchSize] {
			res[id] = StatusNotAnswered
		}
	}
	return
}
