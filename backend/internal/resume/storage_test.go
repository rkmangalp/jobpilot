package resume

import (
	"strings"
	"testing"
	"time"
)

func TestPreparePathsCreatesCompanyFolder(t *testing.T) {
	paths, err := PreparePaths(t.TempDir(), "Stripe, Inc.", "Backend Engineer", time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(paths.Directory, "Stripe_Inc") || !strings.HasSuffix(paths.ResumePDF, "Stripe_Inc_Backend_Engineer_20260901_Resume.pdf") {
		t.Fatalf("unexpected paths: %#v", paths)
	}
}

func TestPreparePathsRejectsUnsafeBlankName(t *testing.T) {
	if _, err := PreparePaths(t.TempDir(), "../", "Backend", time.Now()); err == nil {
		t.Fatal("expected error")
	}
}
