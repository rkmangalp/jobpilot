package app

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jobpilot/jobpilot/backend/internal/candidate"
	"github.com/jobpilot/jobpilot/backend/internal/dashboard"
	"github.com/jobpilot/jobpilot/backend/internal/jobs"
	"github.com/jobpilot/jobpilot/backend/internal/matching"
	"log/slog"
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
func (s *Server) Routes() *gin.Engine {
	r := gin.New()
	r.Use(s.logging(), gin.Recovery())
	r.GET("/api/health", s.health)
	r.GET("/api/candidate", s.getCandidate)
	r.PUT("/api/candidate", s.putCandidate)
	r.GET("/api/dashboard", s.getDashboard)
	r.GET("/api/jobs", s.listJobs)
	r.POST("/api/jobs", s.createJob)
	r.POST("/api/jobs/:id/analyze", s.analyzeJob)
	r.POST("/api/jobs/:id/match", s.matchJob)
	r.GET("/api/metrics", s.metrics)
	r.GET("/", gin.WrapF(s.dashboard))
	return r
}

func (s *Server) listJobs(c *gin.Context) {
	c.JSON(200, s.jobs.List())
}
func (s *Server) createJob(c *gin.Context) {
	var job jobs.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(400, gin.H{"error": "invalid job: " + err.Error()})
		return
	}
	if job.Company == "" || job.Title == "" || job.ApplicationURL == "" || job.Description == "" {
		c.JSON(400, gin.H{"error": "company, title, application_url, and description are required"})
		return
	}
	job.Fingerprint = jobs.Fingerprint(job.Company, job.Title, job.ApplicationURL)
	saved, duplicate := s.jobs.Create(job)
	if duplicate {
		c.JSON(409, gin.H{"error": "duplicate job", "job": saved})
		return
	}
	c.JSON(201, saved)
}
func (s *Server) analyzeJob(c *gin.Context) {
	job, ok := s.jobs.Get(c.Param("id"))
	if !ok {
		c.JSON(404, gin.H{"error": "job not found"})
		return
	}
	job.Analysis = jobs.Analyze(job.Description)
	saved := s.jobs.Update(job)
	c.JSON(200, saved)
}
func (s *Server) matchJob(c *gin.Context) {
	job, ok := s.jobs.Get(c.Param("id"))
	if !ok {
		c.JSON(404, gin.H{"error": "job not found"})
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
	c.JSON(200, result)
}
func (s *Server) health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok", "service": "jobpilot-api", "uptime_seconds": int(time.Since(s.started).Seconds())})
}
func (s *Server) getCandidate(c *gin.Context) { c.JSON(200, s.candidates.Get()) }
func (s *Server) putCandidate(c *gin.Context) {
	var p candidate.Profile
	if err := json.NewDecoder(c.Request.Body).Decode(&p); err != nil {
		c.JSON(400, gin.H{"error": "invalid candidate profile: " + err.Error()})
		return
	}
	p.UpdatedAt = time.Now().UTC()
	c.JSON(200, s.candidates.Update(p))
}
func (s *Server) getDashboard(c *gin.Context) { c.JSON(200, dashboard.Empty()) }
func (s *Server) metrics(c *gin.Context) {
	c.Data(200, "text/plain; version=0.0.4", []byte("# HELP jobs_discovered_total Jobs discovered by JobPilot\n# TYPE jobs_discovered_total counter\njobs_discovered_total 0\n# HELP applications_prepared_total Prepared applications\n# TYPE applications_prepared_total counter\napplications_prepared_total 0\n"))
}
func (s *Server) logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		s.log.Info("http request", "method", c.Request.Method, "path", c.Request.URL.Path, "status", c.Writer.Status(), "duration", time.Since(start).String())
	}
}
