# JobPilot

JobPilot is a safety-first job-search assistant. It discovers and evaluates opportunities, prepares truthful application materials from a verified candidate profile, and always leaves final submission to the candidate.

## Phase 1 status

This repository contains the working Phase 1 foundation:

- Go HTTP API with health, Prometheus-compatible metrics, dashboard, and candidate-profile CRUD.
- An editable seeded candidate profile. Unknown personal, dates, work-authorization, and responsibility details are intentionally marked `TODO` rather than invented.
- MySQL migration with the core JobPilot tables, foreign keys, duplicate fingerprint constraint, and query indexes.
- A lightweight React dashboard served by the API. It presents dashboard metrics and lets a user complete the profile.
- Versioned, file-based prompts and environment-variable configuration template.
- Docker Compose with MySQL and backend health ordering.

Phase 2 is deliberately not implemented: job sources, analysis, matching, and application automation must follow after this foundation is deployed with persistent storage.

## Architecture decision

The MVP is a modular Go monolith. The `candidate`, `dashboard`, and HTTP delivery code have separate package boundaries; the MySQL schema models the future jobs, matching, artifacts, audit, and application workflows. This minimizes operational complexity while retaining clean extraction points. Development currently uses an in-memory candidate repository so the API can be exercised without credentials; implementing the MySQL repository is the first Phase 2 persistence task.

## Run locally

1. Copy `.env.example` to `.env` and populate only the credentials you choose to enable.
2. Run `make dev` (or `go run ./backend/cmd/api`).
3. Open `http://localhost:8080` and verify `http://localhost:8080/api/health`.

For the database-backed container stack, start Docker Desktop, then run `make docker-up`. MySQL applies `migrations/001_initial.sql` only when its volume is first created. Run `make docker-down` to stop it.

## API

| Method | Endpoint | Purpose |
| --- | --- | --- |
| GET | `/api/health` | Liveness response |
| GET | `/api/metrics` | Prometheus metrics |
| GET | `/api/dashboard` | Dashboard metrics |
| GET | `/api/candidate` | Read candidate profile |
| PUT | `/api/candidate` | Replace candidate profile and create a new in-memory version |

Example profile update:

```powershell
$profile = Invoke-RestMethod http://localhost:8080/api/candidate
$profile.full_name = 'Your verified name'
$profile | ConvertTo-Json -Depth 8 | Invoke-RestMethod http://localhost:8080/api/candidate -Method Put -ContentType 'application/json'
```

## Security and approval model

Secrets belong only in `.env` or a secret manager; they are never committed. The future browser service must stop on CAPTCHA, MFA, authentication, access restrictions, unknown questions, and high-risk questions. It will not submit applications automatically: the path is **prepare → notify → human review → manual final submission**.

## Testing

Run `make test`. The current unit tests cover candidate profile versioning and API health/profile retrieval. Expand them in each phase with the required matching, validation, state-machine, database, and mock-browser tests.
