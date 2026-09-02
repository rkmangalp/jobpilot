package matching

import (
	"github.com/jobpilot/jobpilot/backend/internal/candidate"
	"github.com/jobpilot/jobpilot/backend/internal/jobs"
	"testing"
)

func TestScoreRecognizesVerifiedSkills(t *testing.T) {
	result := Score(candidate.Seed(), jobs.Job{Title: "Backend Software Engineer", Location: "Remote US", Analysis: jobs.Analysis{RequiredSkills: []string{"golang", "mysql", "kafka"}, Architecture: []string{"microservices"}}})
	if len(result.MatchedSkills) != 3 || result.OverallScore < 70 {
		t.Fatalf("unexpected result: %#v", result)
	}
}
