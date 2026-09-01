package candidate

import "time"

type Answer struct {
	Question        string `json:"question,omitempty"`
	CanonicalAnswer string `json:"canonical_answer,omitempty"`
	Confidence      string `json:"confidence,omitempty"`
	Source          string `json:"source,omitempty"`
}
type Employment struct {
	Company string `json:"company,omitempty"`
	Title   string `json:"title,omitempty"`
	Dates   string `json:"dates,omitempty"`
	Summary string `json:"summary,omitempty"`
}
type Profile struct {
	ID                     string       `json:"id"`
	Version                int          `json:"version"`
	FullName               string       `json:"full_name"`
	Email                  string       `json:"email"`
	Phone                  string       `json:"phone"`
	Location               string       `json:"location"`
	LinkedInURL            string       `json:"linkedin_url"`
	GitHubURL              string       `json:"github_url"`
	PortfolioURL           string       `json:"portfolio_url"`
	TargetRoles            []string     `json:"target_roles"`
	YearsExperience        string       `json:"years_experience"`
	Skills                 []string     `json:"skills"`
	Employment             []Employment `json:"employment_history"`
	WorkAuthorization      string       `json:"work_authorization"`
	SponsorshipRequirement string       `json:"sponsorship_requirement"`
	PreferredLocations     []string     `json:"preferred_locations"`
	Answers                []Answer     `json:"answers"`
	UpdatedAt              time.Time    `json:"updated_at"`
}

func Seed() Profile {
	return Profile{ID: "default", Version: 1, FullName: "TODO: add full name", Email: "TODO: add email", Phone: "TODO: add phone", Location: "TODO: add location", LinkedInURL: "TODO: add LinkedIn URL", GitHubURL: "TODO: add GitHub URL", PortfolioURL: "TODO: add portfolio URL", TargetRoles: []string{"Backend Software Engineer"}, YearsExperience: "TODO: verify years of experience", Skills: []string{"Golang", "Go", "MySQL", "Kafka", "Redis", "Docker", "REST APIs", "gRPC", "Microservices", "SQL", "Kubernetes", "Splunk", "Kibana", "Prometheus", "Grafana", "Git", "GitHub"}, Employment: []Employment{{Company: "Tesla", Title: "Backend Software Engineer / Contractor", Dates: "TODO: add verified dates", Summary: "Material Flow System (MFS); verified achievement: improved Order Monitor API performance by approximately 30% through SQL query optimization and indexing strategies."}, {Company: "American Express", Title: "TODO: verify title", Dates: "TODO: add verified dates", Summary: "Network Modernization Program (NEMO), Core-Switch / Narwhals; details require verification before use."}, {Company: "Global Payments", Title: "Golang Developer", Dates: "TODO: add verified dates", Summary: "TODO: add verified responsibilities."}, {Company: "Cardinal Health", Title: "Golang Developer", Dates: "TODO: add verified dates", Summary: "TODO: add verified responsibilities."}}, WorkAuthorization: "TODO: verified answer required", SponsorshipRequirement: "TODO: verified answer required", PreferredLocations: []string{"Bay Area", "San Jose", "Fremont", "Sunnyvale", "Mountain View", "Palo Alto", "California", "Remote US"}, UpdatedAt: time.Now().UTC()}
}
