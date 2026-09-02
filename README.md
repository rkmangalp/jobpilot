# JobPilot

JobPilot is a safety-first assistant for backend engineering job applications. It helps collect job descriptions, evaluate fit against a verified candidate profile, prepare truthful application materials, and track progress while keeping the candidate in control of final submission.

## What JobPilot does

```text
Discover -> Analyze -> Match -> Tailor -> Prepare -> Notify -> Human review -> Manual final submission -> Track
```

The current MVP supports the first part of that flow:

- Store an editable candidate profile seeded with supplied, verified experience.
- Add a job manually with its company, role, location, application URL, and description.
- Detect duplicate jobs using a normalized SHA-256 fingerprint.
- Extract known technical and architecture keywords from the description.
- Calculate a transparent, deterministic match score from the verified profile.
- Expose health and Prometheus-compatible metrics endpoints.
- Prepare safe company-specific paths for future tailored resumes and cover letters.

## Safety principles

JobPilot never adds an employer, date, title, skill, achievement, certification, education record, or work-authorization answer unless it exists in your verified candidate profile.

A job keyword missing from your current resume can be used only after you add it as a truthful, verified skill or experience item to the profile. Missing or ambiguous information remains `UNKNOWN` / `NEEDS_REVIEW`.

The future browser workflow will never bypass CAPTCHA, MFA, bot detection, authentication, access restrictions, or rate limits. It will stop for unknown or high-risk questions and will never submit an application automatically.

## Technology stack

| Area | Technology |
| --- | --- |
| Backend API and workflow logic | Go + Gin |
| Database schema and repositories | MySQL 8 + Pop |
| Frontend | React (TypeScript migration planned) |
| Local services | Docker Compose |
| Future async workflows | Kafka + franz-go |
| Future caching and locks | Redis |
| Browser preparation | Playwright (future, compliant use only) |
| Tests | Go test + Ginkgo/Gomega |
| Observability | Structured logs and Prometheus metrics |

JobPilot is a modular monolith: packages have clear boundaries today and can be extracted into services only when that complexity is justified.

Gin owns HTTP routing and middleware. Pop is used for future MySQL connection/repository work; Buffalo is intentionally not used because its full-stack conventions overlap with this separate Go API and React frontend. franz-go is isolated behind an event publisher interface and does not connect to Kafka until `KAFKA_BROKERS` is configured.

## Current implementation status

### Completed

- Candidate profile API with versioned in-memory development storage.
- Seed profile for the supplied Tesla, American Express, Global Payments, and Cardinal Health context. Unverified details are intentionally marked `TODO`.
- MySQL migration covering profiles, jobs, matches, applications, documents, audit records, notifications, and browser sessions.
- Manual job ingestion API and duplicate detection.
- Keyword analysis for common backend technologies and architecture terms.
- Deterministic job matching with `STRONG APPLY`, `APPLY`, `REVIEW`, or `SKIP` recommendations.
- Resume/cover-letter output-path preparation under a company-specific folder.
- Go unit tests for candidate versioning, HTTP endpoints, duplicate fingerprints, job analysis, matching, and output paths.

### Planned next

- MySQL-backed repositories instead of development-only in-memory storage.
- A proper React + TypeScript frontend build.
- Configurable scoring weights and job-search preferences in the UI.
- ATS-friendly DOCX/PDF resume and cover-letter generation with claim validation.
- A local mock application site and Playwright field mapping.
- Approved notification-provider integration, scheduling, Redis, and Kafka workflows.

## Run locally

### Run the API directly

1. Copy `.env.example` to `.env`.
2. Update `RESUME_OUTPUT_DIR` if you want a different output location. Its default is `C:\Users\Rk\Documents\Resumes\_FT`.
3. Run `make dev`.
4. Open `http://localhost:8080`.

Check the API:

```powershell
Invoke-RestMethod http://localhost:8080/api/health
```

### Run with Docker and MySQL

Start Docker Desktop, then run:

```powershell
make docker-up
```

MySQL runs `migrations/001_initial.sql` when its data volume is first created. Stop the stack with `make docker-down`.

## API reference

| Method | Endpoint | Purpose |
| --- | --- | --- |
| GET | `/api/health` | Liveness response |
| GET | `/api/metrics` | Prometheus metrics |
| GET | `/api/dashboard` | Dashboard metrics |
| GET | `/api/candidate` | Read the candidate profile |
| PUT | `/api/candidate` | Update the candidate profile |
| GET | `/api/jobs` | List manually added jobs |
| POST | `/api/jobs` | Add a job description |
| POST | `/api/jobs/{id}/analyze` | Extract known job keywords |
| POST | `/api/jobs/{id}/match` | Calculate the deterministic match result |

Example job creation:

```powershell
$job = @{ company = 'ExampleCo'; title = 'Backend Software Engineer'; location = 'Remote US'; application_url = 'https://jobs.example.com/backend-engineer'; description = 'Build Go microservices using Kafka, MySQL, Docker, and REST APIs.' } | ConvertTo-Json
Invoke-RestMethod http://localhost:8080/api/jobs -Method Post -ContentType 'application/json' -Body $job
```

## Match scoring

Scoring is deterministic; an LLM does not choose the score.

| Dimension | Weight |
| --- | --- |
| Required technical skills | 30% |
| Relevant experience/role | 25% |
| Architecture/domain | 10% |
| Location/work arrangement | 5% |
| Responsibilities | 15% |
| Years of experience | 10% |
| Education | 5% |

Thresholds: 90–100 `STRONG APPLY`, 80–89 `APPLY`, 70–79 `REVIEW`, below 70 `SKIP`.

## Generated application files

Set `RESUME_OUTPUT_DIR` in `.env`. The default creates this layout:

```text
C:\Users\Rk\Documents\Resumes\_FT\
  Company_Name\
    Company_Name_Backend_Engineer_YYYYMMDD_Resume.docx
    Company_Name_Backend_Engineer_YYYYMMDD_Resume.pdf
    Company_Name_Backend_Engineer_YYYYMMDD_Cover_Letter.docx
    Company_Name_Backend_Engineer_YYYYMMDD_Cover_Letter.pdf
```

The storage path is implemented. Actual document generation is the next Phase 3 feature and will validate each claim against the candidate profile before saving a file.

## Testing

Run `make test` for the Go suite and `go vet ./...` for static checks.

## Security

- Never commit `.env`, credentials, session cookies, tokens, or API keys.
- Use environment variables or a secret manager for configuration.
- Do not log passwords, authentication data, or sensitive application answers.
- Preserve human review before every final application submission.
