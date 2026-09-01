# Frontend delivery

Phase 1 delivers the small React dashboard from the Go server to avoid requiring a local Node installation. The source is embedded in `backend/internal/app/dashboard.go` and uses React’s ESM distribution during development. In Phase 2, move it into a Vite TypeScript project and bake static assets into the backend image.
