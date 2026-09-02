package main

import (
	"github.com/jobpilot/jobpilot/backend/internal/app"
	"github.com/jobpilot/jobpilot/backend/internal/candidate"
	"log/slog"
	"os"
)

func main() {
	addr := os.Getenv("JOBPILOT_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	server := app.New(candidate.NewMemoryStore(), log)
	log.Info("JobPilot API started", "addr", addr)
	if err := server.Routes().Run(addr); err != nil {
		log.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
