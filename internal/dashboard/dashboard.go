package dashboard

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/Sufyan11104/shipcheck/internal/engine"
	"github.com/Sufyan11104/shipcheck/internal/report"
	"github.com/Sufyan11104/shipcheck/internal/rules"
	"github.com/Sufyan11104/shipcheck/internal/scanner"
)

// AuditRunner builds the report shown by the dashboard.
type AuditRunner func() (report.AuditReport, error)

// BuildReport runs the same scan and audit flow used by the audit command.
func BuildReport(path string, categories []string) (report.AuditReport, error) {
	result := scanner.Scan(path)
	if result.Error != nil {
		return report.AuditReport{}, result.Error
	}

	eng := engine.NewEngine(result.Path)
	findings, _ := eng.RunChecksWithCategories(result.IsGitRepository, categories)
	findings = engine.FilterFindingsByCategory(findings, categories)
	score := engine.CalculateScore(findings)

	return report.NewAuditReport(result, findings, score), nil
}

// NewHandler returns an HTTP handler for the dashboard and JSON API.
func NewHandler(path string, categories []string) http.Handler {
	return NewHandlerWithRunner(func() (report.AuditReport, error) {
		return BuildReport(path, categories)
	})
}

// NewHandlerWithRunner returns an HTTP handler using a custom report builder.
func NewHandlerWithRunner(run AuditRunner) http.Handler {
	mux := http.NewServeMux()
	h := &handler{run: run}

	mux.HandleFunc("/", h.handleDashboard)
	mux.HandleFunc("/api/report", h.handleAPIReport)

	return mux
}

type handler struct {
	run AuditRunner
}

func (h *handler) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	auditReport, err := h.run()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := dashboardTemplate.Execute(w, newPageView(auditReport)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) handleAPIReport(w http.ResponseWriter, r *http.Request) {
	auditReport, err := h.run()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := report.RenderJSON(w, auditReport); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type pageView struct {
	Report      report.AuditReport
	ScoreLabel  string
	GeneratedAt string
	Categories  []categoryView
}

type categoryView struct {
	Title    string
	Findings []findingView
}

type findingView struct {
	ID              string
	Status          string
	Severity        string
	Message         string
	Remediation     string
	Symbol          string
	Class           string
	ShowRemediation bool
}

func newPageView(auditReport report.AuditReport) pageView {
	return pageView{
		Report:      auditReport,
		ScoreLabel:  scoreLabel(auditReport.Score),
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05 MST"),
		Categories:  groupFindings(auditReport.Findings),
	}
}

func groupFindings(findings []report.ReportFinding) []categoryView {
	grouped := make(map[string][]findingView)
	for _, finding := range findings {
		key := strings.ToLower(finding.Category)
		grouped[key] = append(grouped[key], newFindingView(finding))
	}

	ordered := []categoryView{
		{Title: "Documentation", Findings: grouped["docs"]},
		{Title: "Repository", Findings: grouped["repo"]},
		{Title: "Environment", Findings: grouped["env"]},
		{Title: "Docker", Findings: grouped["docker"]},
		{Title: "GitHub Actions", Findings: grouped["ci"]},
		{Title: "Kubernetes", Findings: grouped["k8s"]},
		{Title: "Terraform", Findings: grouped["terraform"]},
	}

	return ordered
}

func newFindingView(finding report.ReportFinding) findingView {
	return findingView{
		ID:              finding.ID,
		Status:          string(finding.Status),
		Severity:        string(finding.Severity),
		Message:         finding.Message,
		Remediation:     finding.Remediation,
		Symbol:          statusSymbol(finding.Status),
		Class:           statusClass(finding.Status),
		ShowRemediation: shouldShowRemediation(finding),
	}
}

func shouldShowRemediation(finding report.ReportFinding) bool {
	if finding.Status != rules.StatusWarn && finding.Status != rules.StatusFail {
		return false
	}
	return finding.Remediation != "" && finding.Remediation != "N/A"
}

func scoreLabel(score int) string {
	switch {
	case score >= 90:
		return "Excellent"
	case score >= 70:
		return "Good"
	case score >= 50:
		return "Needs attention"
	default:
		return "High risk"
	}
}

func statusSymbol(status rules.Status) string {
	switch status {
	case rules.StatusPass:
		return "✓"
	case rules.StatusWarn:
		return "!"
	case rules.StatusFail:
		return "✗"
	case rules.StatusSkip:
		return "-"
	default:
		return "?"
	}
}

func statusClass(status rules.Status) string {
	switch status {
	case rules.StatusPass:
		return "pass"
	case rules.StatusWarn:
		return "warn"
	case rules.StatusFail:
		return "fail"
	case rules.StatusSkip:
		return "skip"
	default:
		return "unknown"
	}
}

var dashboardTemplate = template.Must(template.New("dashboard").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>ShipCheck Dashboard</title>
  <style>
    :root {
      color-scheme: light;
      --bg: #f7f8fa;
      --panel: #ffffff;
      --ink: #17202a;
      --muted: #667085;
      --border: #d8dee8;
      --pass: #147a46;
      --warn: #9a5b00;
      --fail: #b42318;
      --skip: #667085;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      background: var(--bg);
      color: var(--ink);
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      line-height: 1.45;
    }
    main {
      max-width: 1120px;
      margin: 0 auto;
      padding: 32px 20px 48px;
    }
    header {
      display: flex;
      justify-content: space-between;
      gap: 24px;
      align-items: flex-start;
      margin-bottom: 24px;
    }
    h1, h2, h3, p { margin-top: 0; }
    h1 {
      font-size: 32px;
      line-height: 1.1;
      margin-bottom: 8px;
    }
    .meta {
      color: var(--muted);
      font-size: 14px;
      margin: 0;
    }
    .score {
      min-width: 220px;
      text-align: right;
    }
    .score strong {
      display: block;
      font-size: 34px;
      line-height: 1;
    }
    .score span {
      color: var(--muted);
      font-weight: 600;
    }
    .cards {
      display: grid;
      grid-template-columns: repeat(4, minmax(0, 1fr));
      gap: 12px;
      margin: 24px 0;
    }
    .card, .section {
      background: var(--panel);
      border: 1px solid var(--border);
      border-radius: 8px;
    }
    .card {
      padding: 16px;
    }
    .card .label {
      color: var(--muted);
      font-size: 13px;
      font-weight: 600;
      text-transform: uppercase;
    }
    .card .value {
      display: block;
      font-size: 28px;
      font-weight: 700;
      margin-top: 4px;
    }
    .section {
      margin-top: 14px;
      overflow: hidden;
    }
    .section h2 {
      margin: 0;
      padding: 14px 16px;
      border-bottom: 1px solid var(--border);
      font-size: 18px;
    }
    .empty {
      color: var(--muted);
      margin: 0;
      padding: 16px;
    }
    .finding {
      display: grid;
      grid-template-columns: 32px 1fr;
      gap: 12px;
      padding: 16px;
      border-top: 1px solid var(--border);
    }
    .finding:first-of-type { border-top: 0; }
    .symbol {
      width: 28px;
      height: 28px;
      border-radius: 999px;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      font-weight: 700;
      border: 1px solid currentColor;
    }
    .finding.pass .symbol { color: var(--pass); }
    .finding.warn {
      border-left: 4px solid var(--warn);
      background: #fff9ec;
    }
    .finding.warn .symbol { color: var(--warn); }
    .finding.fail {
      border-left: 4px solid var(--fail);
      background: #fff4f2;
    }
    .finding.fail .symbol { color: var(--fail); }
    .finding.skip {
      color: var(--skip);
      background: #fafafa;
    }
    .finding.skip .symbol { color: var(--skip); }
    .finding h3 {
      font-size: 15px;
      margin-bottom: 4px;
      word-break: break-word;
    }
    .finding p {
      margin-bottom: 6px;
    }
    .finding .detail {
      color: var(--muted);
      font-size: 13px;
      margin-bottom: 0;
    }
    .fix {
      margin-top: 8px;
      font-size: 14px;
    }
    code {
      font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
      font-size: 0.95em;
    }
    @media (max-width: 720px) {
      main { padding: 24px 14px 40px; }
      header {
        display: block;
      }
      .score {
        text-align: left;
        margin-top: 18px;
      }
      .cards {
        grid-template-columns: repeat(2, minmax(0, 1fr));
      }
    }
  </style>
</head>
<body>
  <main>
    <header>
      <div>
        <h1>ShipCheck</h1>
        <p class="meta">Target: <code>{{.Report.Path}}</code></p>
        <p class="meta">Git repository: {{if .Report.GitRepository}}yes{{else}}no{{end}} · Files: {{.Report.FilesScanned}} · Directories: {{.Report.DirectoriesScanned}}</p>
        <p class="meta">Generated: {{.GeneratedAt}}</p>
      </div>
      <div class="score">
        <span>Score</span>
        <strong>{{.Report.Score}}/100</strong>
        <span>{{.ScoreLabel}}</span>
      </div>
    </header>

    <section class="cards" aria-label="Audit summary">
      <div class="card"><span class="label">Passed</span><span class="value">{{.Report.PassedCount}}</span></div>
      <div class="card"><span class="label">Warnings</span><span class="value">{{.Report.WarningCount}}</span></div>
      <div class="card"><span class="label">Failed</span><span class="value">{{.Report.FailedCount}}</span></div>
      <div class="card"><span class="label">Skipped</span><span class="value">{{.Report.SkippedCount}}</span></div>
    </section>

    {{range .Categories}}
    <section class="section">
      <h2>{{.Title}}</h2>
      {{if .Findings}}
        {{range .Findings}}
        <article class="finding {{.Class}}">
          <div class="symbol" aria-label="{{.Status}}">{{.Symbol}}</div>
          <div>
            <h3><code>{{.ID}}</code></h3>
            <p>{{.Message}}</p>
            <p class="detail">Status: {{.Status}} · Severity: {{.Severity}}</p>
            {{if .ShowRemediation}}<p class="fix"><strong>Fix:</strong> {{.Remediation}}</p>{{end}}
          </div>
        </article>
        {{end}}
      {{else}}
        <p class="empty">No findings in this category.</p>
      {{end}}
    </section>
    {{end}}
  </main>
</body>
</html>`))
