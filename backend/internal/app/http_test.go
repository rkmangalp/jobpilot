package app

import (
	"bytes"
	"encoding/json"
	"github.com/jobpilot/jobpilot/backend/internal/candidate"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCandidateCRUD(t *testing.T) {
	s := New(candidate.NewMemoryStore(), slog.Default()).Routes()
	r := httptest.NewRequest(http.MethodGet, "/api/candidate", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("GET %d", w.Code)
	}
	r = httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w = httptest.NewRecorder()
	s.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("health %d", w.Code)
	}
	var profile candidate.Profile
	payload := []byte(`{"full_name":"Verified Candidate","skills":["Go"]}`)
	r = httptest.NewRequest(http.MethodPut, "/api/candidate", bytes.NewReader(payload))
	w = httptest.NewRecorder()
	s.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("PUT %d: %s", w.Code, w.Body.String())
	}
	if err := json.Unmarshal(w.Body.Bytes(), &profile); err != nil || profile.FullName != "Verified Candidate" || profile.Version != 2 {
		t.Fatalf("unexpected update: %#v, %v", profile, err)
	}
}
