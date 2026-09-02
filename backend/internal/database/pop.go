package database

import (
	"fmt"
	"os"

	"github.com/gobuffalo/pop/v6"
)

// NewPopConnection builds the MySQL connection used by persistent repositories.
// It does not connect until Open is called by the application bootstrap.
func NewPopConnection() (*pop.Connection, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return nil, fmt.Errorf("DATABASE_URL is required for persistent storage")
	}
	return pop.NewConnection(&pop.ConnectionDetails{URL: url})
}
