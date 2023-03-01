package testservice

import (
	"fmt"
	"log"
)

type TestService struct {
	sessionMap map[int64]*Session
	lastID     int64
	problems   Problems
	lg         *log.Logger
}

type Session struct {
	SessionID  int64
	unanswered map[int64]ProblemStatus
	answered   map[int64]ProblemStatus
	score      int
}

func NewTestService() *TestService {
	s := &TestService{
		sessionMap: make(map[int64]*Session),
		lastID:     0,
	}
	return s
}

func (s *TestService) GetSession(id int64) (*Session, bool) {
	el, ok := s.sessionMap[id]
	return el, ok
}

func (s *TestService) InitSession() *Session {
	newID := s.lastID + 1
	s.lastID = newID
	newSession := &Session{
		SessionID: newID,
		answered:  make(map[int64]ProblemStatus),
	}
	newSession.unanswered = s.problems.generateProblemBatch()
	s.sessionMap[newID] = newSession
	return newSession
}

// CloseSession returns correct answers and total answered questions
func (s *TestService) CloseSession(id int64) (int, int) {
	session, found := s.sessionMap[id]
	if !found {
		s.lg.Printf("can't find session %d for delete\n", id)
		return -1, -1
	}
	score := session.score
	all := len(session.answered)
	delete(s.sessionMap, id)
	return score, all
}

func (s *Session) GetNextUnanswered() (int64, bool) {
	for key, val := range s.unanswered {
		if val == StatusNotAnswered {
			return key, true
		}
	}
	return -1, false
}

func (s *TestService) GetProblem(problemID int64) (*Problem, bool) {
	return s.problems.GetProblem(problemID)
}

func (s *TestService) Assert(sessionID int64, problemID int64, answer string) (bool, error) {
	session, found := s.GetSession(sessionID)
	if !found {
		return false, fmt.Errorf("session %d not found", sessionID)
	}

	_, found = session.unanswered[problemID]
	if !found {
		return false, fmt.Errorf("no such problem for this session (sID: %d, qID: %d)", sessionID, problemID)
	}

	correct, err := s.problems.Assert(problemID, answer)
	if err != nil {
		return false, fmt.Errorf("assert failed: %w", err)
	}

	if correct {
		session.answered[problemID] = StatusCorrectAnswer
		session.score++
	} else {
		session.answered[problemID] = StatusIncorrectAnswer
	}
	delete(session.unanswered, problemID)

	return correct, nil
}
