package app

import (
	"encoding/json"
	"github.com/jobpilot/jobpilot/backend/internal/candidate"
	"github.com/jobpilot/jobpilot/backend/internal/dashboard"
	"github.com/jobpilot/jobpilot/backend/internal/jobs"
	"github.com/jobpilot/jobpilot/backend/internal/matching"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	candidates candidate.Store
	jobs       jobs.Store
	log        *slog.Logger
	started    time.Time
}

func New(c candidate.Store, l *slog.Logger) *Server {
	return &Server{candidates: c, jobs: jobs.NewMemoryStore(), log: l, started: time.Now()}
}
func (s *Server) Routes() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("GET /api/health", s.health)
	m.HandleFunc("GET /api/candidate", s.getCandidate)
	m.HandleFunc("PUT /api/candidate", s.putCandidate)
	m.HandleFunc("GET /api/dashboard", s.getDashboard)
	m.HandleFunc("GET /api/jobs", s.listJobs)
	m.HandleFunc("POST /api/jobs", s.createJob)
	m.HandleFunc("POST /api/jobs/{id}/analyze", s.analyzeJob)
	m.HandleFunc("POST /api/jobs/{id}/match", s.matchJob)
	m.HandleFunc("GET /api/metrics", s.metrics)
	m.HandleFunc("GET /", s.dashboard)
	return s.logging(m)
}

func (s *Server) listJobs(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusOK, s.jobs.List())
}
func (s *Server) createJob(w http.ResponseWriter, r *http.Request) {
	var job jobs.Job
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 2<<20)).Decode(&job); err != nil {
		respond(w, 400, map[string]string{"error": "invalid job: " + err.Error()})
		return
	}
	if job.Company == "" || job.Title == "" || job.ApplicationURL == "" || job.Description == "" {
		respond(w, 400, map[string]string{"error": "company, title, application_url, and description are required"})
		return
	}
	job.Fingerprint = jobs.Fingerprint(job.Company, job.Title, job.ApplicationURL)
	saved, duplicate := s.jobs.Create(job)
	if duplicate {
		respond(w, http.StatusConflict, map[string]any{"error": "duplicate job", "job": saved})
		return
	}
	respond(w, http.StatusCreated, saved)
}
func (s *Server) analyzeJob(w http.ResponseWriter, r *http.Request) {
	job, ok := s.jobs.Get(r.PathValue("id"))
	if !ok {
		respond(w, 404, map[string]string{"error": "job not found"})
		return
	}
	job.Analysis = jobs.Analyze(job.Description)
	saved := s.jobs.Update(job)
	respond(w, 200, saved)
}
func (s *Server) matchJob(w http.ResponseWriter, r *http.Request) {
	job, ok := s.jobs.Get(r.PathValue("id"))
	if !ok {
		respond(w, 404, map[string]string{"error": "job not found"})
		return
	}
	if len(job.Analysis.RequiredSkills) == 0 {
		job.Analysis = jobs.Analyze(job.Description)
	}
	result := matching.Score(s.candidates.Get(), job)
	if result.Recommendation == "SKIP" {
		job.Status = jobs.Skipped
	} else {
		job.Status = jobs.Matched
	}
	s.jobs.Update(job)
	respond(w, 200, result)
}
func respond(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	respond(w, 200, map[string]any{"status": "ok", "service": "jobpilot-api", "uptime_seconds": int(time.Since(s.started).Seconds())})
}
func (s *Server) getCandidate(w http.ResponseWriter, r *http.Request) {
	respond(w, 200, s.candidates.Get())
}
func (s *Server) putCandidate(w http.ResponseWriter, r *http.Request) {
	var p candidate.Profile
	d := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		respond(w, 400, map[string]string{"error": "invalid candidate profile: " + err.Error()})
		return
	}
	p.UpdatedAt = time.Now().UTC()
	respond(w, 200, s.candidates.Update(p))
}
func (s *Server) getDashboard(w http.ResponseWriter, r *http.Request) {
	respond(w, 200, dashboard.Empty())
}
func (s *Server) metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	_, _ = w.Write([]byte("# HELP jobs_discovered_total Jobs discovered by JobPilot\n# TYPE jobs_discovered_total counter\njobs_discovered_total 0\n# HELP applications_prepared_total Prepared applications\n# TYPE applications_prepared_total counter\napplications_prepared_total 0\n"))
}
func (s *Server) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.log.Info("http request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String())
	})
}
