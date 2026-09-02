package jobs

import "testing"

func TestDuplicateFingerprintIsStable(t *testing.T) {
	if Fingerprint("Stripe", "Backend Engineer", "https://example.com/job") != Fingerprint(" stripe ", " backend engineer ", "https://example.com/job") {
		t.Fatal("equivalent job values should have identical fingerprints")
	}
}
func TestAnalyzeFindsKnownKeywords(t *testing.T) {
	a := Analyze("Build Go microservices with Kafka and MySQL.")
	if len(a.RequiredSkills) < 3 || len(a.Architecture) != 1 {
		t.Fatalf("unexpected analysis: %#v", a)
	}
}
