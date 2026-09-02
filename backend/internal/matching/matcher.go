package matching

import (
	"github.com/jobpilot/jobpilot/backend/internal/candidate"
	"github.com/jobpilot/jobpilot/backend/internal/jobs"
	"strings"
)

type Result struct {
	OverallScore      float64  `json:"overall_score"`
	Recommendation    string   `json:"recommendation"`
	MatchedSkills     []string `json:"matched_skills"`
	UnmatchedSkills   []string `json:"unmatched_skills"`
	Strengths         []string `json:"strengths"`
	Risks             []string `json:"risks"`
	TechnicalMatch    float64  `json:"technical_match"`
	ExperienceMatch   float64  `json:"experience_match"`
	ArchitectureMatch float64  `json:"architecture_match"`
	LocationMatch     float64  `json:"location_match"`
}

func Score(profile candidate.Profile, job jobs.Job) Result {
	known := map[string]bool{}
	for _, skill := range profile.Skills {
		known[strings.ToLower(skill)] = true
	}
	result := Result{}
	for _, skill := range job.Analysis.RequiredSkills {
		if known[skill] || (skill == "golang" && known["go"]) || (skill == "go" && known["golang"]) {
			result.MatchedSkills = append(result.MatchedSkills, skill)
		} else {
			result.UnmatchedSkills = append(result.UnmatchedSkills, skill)
		}
	}
	if len(job.Analysis.RequiredSkills) == 0 {
		result.TechnicalMatch = 50
	} else {
		result.TechnicalMatch = 100 * float64(len(result.MatchedSkills)) / float64(len(job.Analysis.RequiredSkills))
	}
	title := strings.ToLower(job.Title)
	for _, role := range profile.TargetRoles {
		if strings.Contains(title, "backend") || strings.Contains(title, "software engineer") || strings.Contains(title, strings.ToLower(role)) {
			result.ExperienceMatch = 100
			break
		}
	}
	if result.ExperienceMatch == 0 {
		result.ExperienceMatch = 50
	}
	if len(job.Analysis.Architecture) == 0 {
		result.ArchitectureMatch = 50
	} else {
		for _, a := range job.Analysis.Architecture {
			if a == "microservices" || strings.Contains(a, "event") {
				result.ArchitectureMatch = 100
				break
			}
		}
		if result.ArchitectureMatch == 0 {
			result.ArchitectureMatch = 50
		}
	}
	location := strings.ToLower(job.Location)
	for _, preferred := range profile.PreferredLocations {
		if strings.Contains(location, strings.ToLower(preferred)) || strings.Contains(location, "remote") {
			result.LocationMatch = 100
			break
		}
	}
	if result.LocationMatch == 0 {
		result.LocationMatch = 50
	}
	result.OverallScore = 0.30*result.TechnicalMatch + 0.25*result.ExperienceMatch + 0.10*result.ArchitectureMatch + 0.05*result.LocationMatch + 0.15*50 + 0.10*50 + 0.05*50
	for _, skill := range result.MatchedSkills {
		result.Strengths = append(result.Strengths, "Verified skill: "+skill)
	}
	result.Risks = append(result.Risks, job.Analysis.Concerns...)
	switch {
	case result.OverallScore >= 90:
		result.Recommendation = "STRONG APPLY"
	case result.OverallScore >= 80:
		result.Recommendation = "APPLY"
	case result.OverallScore >= 70:
		result.Recommendation = "REVIEW"
	default:
		result.Recommendation = "SKIP"
	}
	return result
}
