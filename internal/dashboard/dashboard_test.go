package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Sufyan11104/shipcheck/internal/report"
)

func TestDashboardHandlerReturnsHTML(t *testing.T) {
	tmpDir := newDashboardFixture(t)

	handler := NewHandler(tmpDir, nil)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "ShipCheck") {
		t.Errorf("expected dashboard title, got %s", body)
	}
	if !strings.Contains(body, "Score") {
		t.Errorf("expected score in dashboard HTML, got %s", body)
	}
	if !strings.Contains(body, "docs.readme_exists") {
		t.Errorf("expected findings in dashboard HTML, got %s", body)
	}
}

func TestDashboardAPIReportReturnsJSON(t *testing.T) {
	tmpDir := newDashboardFixture(t)

	handler := NewHandler(tmpDir, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/report", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var auditReport report.AuditReport
	if err := json.Unmarshal(rr.Body.Bytes(), &auditReport); err != nil {
		t.Fatalf("expected valid JSON report, got %v", err)
	}
	if auditReport.Path != tmpDir {
		t.Errorf("expected path %q, got %q", tmpDir, auditReport.Path)
	}
	if len(auditReport.Findings) == 0 {
		t.Fatal("expected JSON report to include findings")
	}
}

func TestBuildReportInvalidPath(t *testing.T) {
	missingPath := filepath.Join(t.TempDir(), "missing")

	_, err := BuildReport(missingPath, nil)
	if err == nil {
		t.Fatal("expected error for missing path")
	}
	if !strings.Contains(err.Error(), "failed to access path") {
		t.Errorf("expected clear path error, got %v", err)
	}
}

func newDashboardFixture(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	writeDashboardFile(t, filepath.Join(tmpDir, "README.md"), "# Test\n")
	writeDashboardFile(t, filepath.Join(tmpDir, ".gitignore"), ".env\n")
	writeDashboardFile(t, filepath.Join(tmpDir, ".env.example"), "# no secrets required\n")
	return tmpDir
}

func writeDashboardFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create parent directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
}
