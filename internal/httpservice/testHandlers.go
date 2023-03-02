package httpservice

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *HTTPService) addTestHandlers() {
	s.mux.HandleFunc("/test/start", s.startTestHandler)
	s.mux.HandleFunc("/test/next", s.testHandler)
	s.mux.HandleFunc("/test/assert", s.assertHandler)
	s.mux.HandleFunc("/test/stop", s.stopHandler)
}

func (s *HTTPService) startTestHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(writer, "Method must be POST", http.StatusBadRequest)
		return
	}
	session := s.testSessions.InitSession()
	bytes, err := json.Marshal(session.SessionID)
	if err != nil {
		http.Error(writer, "Marshal session ID failed", http.StatusInternalServerError)
		return
	}
	_, err = writer.Write(bytes)
	if err != nil {
		http.Error(writer, "Write session ID failed", http.StatusInternalServerError)
		return
	}
	return
}

func (s *HTTPService) testHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(writer, "Method not supported", http.StatusBadRequest)
		return
	}

	err := req.ParseForm()
	if err != nil {
		http.Error(writer, "Unable to parse form", http.StatusInternalServerError)
		return
	}

	idStr := req.FormValue("sessionID")
	sessionID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	session, found := s.testSessions.GetSession(int64(sessionID))
	if !found {
		http.Error(writer, "Session not found", http.StatusUnauthorized)
		return
	}

	problemID, found := session.GetNextUnanswered()
	if !found {
		writer.WriteHeader(204)
		return
	}

	problem, found := s.testSessions.GetProblem(problemID)
	if !found {
		http.Error(writer, "Problem not found", http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(problem)
	if err != nil {
		http.Error(writer, "Marshal error", http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(bytes)
	if err != nil {
		http.Error(writer, "Write error", http.StatusInternalServerError)
		return
	}
	return
}

func (s *HTTPService) assertHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(writer, "Method not supported", http.StatusBadRequest)
		return
	}
	err := req.ParseForm()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID, errSess := strconv.Atoi(req.FormValue("sessionID"))
	problemID, errProb := strconv.Atoi(req.FormValue("problemID"))
	answer := req.FormValue("answer")
	if errSess != nil || errProb != nil {
		http.Error(writer, "sessionID and/or problemID not found", http.StatusBadRequest)
		return
	}

	_, err = s.testSessions.Assert(int64(sessionID), int64(problemID), answer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func (s *HTTPService) stopHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(writer, "Method not supported", http.StatusBadRequest)
		return
	}

	err := req.ParseForm()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	idStr := req.FormValue("sessionID")
	sessionID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	score, err := s.testSessions.CloseSession(int64(sessionID))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusForbidden)
		return
	}

	bytes, err := json.Marshal(score)
	if err != nil {
		http.Error(writer, "Marshal error", http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(bytes)
	if err != nil {
		http.Error(writer, "Write error", http.StatusInternalServerError)
		return
	}
	return
}
