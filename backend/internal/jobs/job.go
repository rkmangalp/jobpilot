package jobs

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"sync"
	"time"
)

type Status string

const (
	Discovered Status = "DISCOVERED"
	Matched    Status = "MATCHED"
	Skipped    Status = "SKIPPED"
)

type Job struct {
	ID             string    `json:"id"`
	Company        string    `json:"company"`
	Title          string    `json:"title"`
	Location       string    `json:"location"`
	ApplicationURL string    `json:"application_url"`
	Description    string    `json:"description"`
	Fingerprint    string    `json:"fingerprint"`
	Status         Status    `json:"status"`
	Analysis       Analysis  `json:"analysis,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type Analysis struct {
	RequiredSkills  []string `json:"required_skills"`
	PreferredSkills []string `json:"preferred_skills"`
	Architecture    []string `json:"architecture"`
	Concerns        []string `json:"concerns"`
}

type Store interface {
	Create(Job) (Job, bool)
	List() []Job
	Get(string) (Job, bool)
	Update(Job) Job
}

type MemoryStore struct {
	mu           sync.RWMutex
	jobs         map[string]Job
	fingerprints map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{jobs: map[string]Job{}, fingerprints: map[string]string{}}
}
func (s *MemoryStore) Create(job Job) (Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id, exists := s.fingerprints[job.Fingerprint]; exists {
		return s.jobs[id], true
	}
	job.ID = job.Fingerprint[:16]
	job.CreatedAt = time.Now().UTC()
	job.Status = Discovered
	s.jobs[job.ID] = job
	s.fingerprints[job.Fingerprint] = job.ID
	return job, false
}
func (s *MemoryStore) List() []Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		result = append(result, job)
	}
	return result
}
func (s *MemoryStore) Get(id string) (Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	job, ok := s.jobs[id]
	return job, ok
}
func (s *MemoryStore) Update(job Job) Job {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.ID] = job
	return job
}

func Fingerprint(company, title, url string) string {
	normalized := strings.ToLower(strings.TrimSpace(company) + "|" + strings.TrimSpace(title) + "|" + strings.TrimSpace(url))
	sum := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(sum[:])
}

func Analyze(description string) Analysis {
	lower := strings.ToLower(description)
	a := Analysis{}
	for _, skill := range []string{"golang", "go", "mysql", "kafka", "redis", "docker", "kubernetes", "grpc", "rest", "sql", "prometheus", "grafana", "aws", "java", "python"} {
		if strings.Contains(lower, skill) {
			a.RequiredSkills = append(a.RequiredSkills, skill)
		}
	}
	for _, architecture := range []string{"microservices", "distributed systems", "event-driven", "event driven"} {
		if strings.Contains(lower, architecture) {
			a.Architecture = append(a.Architecture, architecture)
		}
	}
	for _, concern := range []string{"sponsorship", "security clearance", "must be authorized", "on-site", "onsite"} {
		if strings.Contains(lower, concern) {
			a.Concerns = append(a.Concerns, concern)
		}
	}
	return a
}
