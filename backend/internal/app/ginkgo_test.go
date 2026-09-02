package app_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jobpilot/jobpilot/backend/internal/app"
	"github.com/jobpilot/jobpilot/backend/internal/candidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"log/slog"
)

func TestApp(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "JobPilot API Suite")
}

var _ = ginkgo.Describe("Health endpoint", func() {
	ginkgo.It("reports the API as healthy", func() {
		router := app.New(candidate.NewMemoryStore(), slog.Default()).Routes()
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		router.ServeHTTP(response, request)
		gomega.Expect(response.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(response.Body.String()).To(gomega.ContainSubstring(`"status":"ok"`))
	})
})
