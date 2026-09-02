package resume

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const defaultOutputDir = `C:\Users\Rk\Documents\Resumes\_FT`

type ArtifactPaths struct {
	Directory  string `json:"directory"`
	ResumeDOCX string `json:"resume_docx"`
	ResumePDF  string `json:"resume_pdf"`
	LetterDOCX string `json:"cover_letter_docx"`
	LetterPDF  string `json:"cover_letter_pdf"`
}

func OutputDir() string {
	if configured := strings.TrimSpace(os.Getenv("RESUME_OUTPUT_DIR")); configured != "" {
		return configured
	}
	return defaultOutputDir
}

func PreparePaths(baseDir, company, role string, createdAt time.Time) (ArtifactPaths, error) {
	companyName := safeName(company)
	roleName := safeName(role)
	if companyName == "" || roleName == "" {
		return ArtifactPaths{}, fmt.Errorf("company and role are required")
	}
	directory := filepath.Join(baseDir, companyName)
	if err := os.MkdirAll(directory, 0o750); err != nil {
		return ArtifactPaths{}, fmt.Errorf("create artifact directory: %w", err)
	}
	prefix := fmt.Sprintf("%s_%s_%s", companyName, roleName, createdAt.UTC().Format("20060102"))
	return ArtifactPaths{
		Directory:  directory,
		ResumeDOCX: filepath.Join(directory, prefix+"_Resume.docx"),
		ResumePDF:  filepath.Join(directory, prefix+"_Resume.pdf"),
		LetterDOCX: filepath.Join(directory, prefix+"_Cover_Letter.docx"),
		LetterPDF:  filepath.Join(directory, prefix+"_Cover_Letter.pdf"),
	}, nil
}

var unsafeName = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

func safeName(value string) string {
	return strings.Trim(unsafeName.ReplaceAllString(strings.TrimSpace(value), "_"), "._-")
}
