package report

import (
	"fmt"
	"io"
	"strings"
)

// RenderMarkdown writes a Markdown deployment readiness report.
func RenderMarkdown(w io.Writer, auditReport AuditReport) error {
	if _, err := fmt.Fprintln(w, "# ShipCheck Deployment Readiness Report"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}

	summaryRows := [][2]string{
		{"Path", auditReport.Path},
		{"Git repository", formatBool(auditReport.GitRepository)},
		{"Files scanned", fmt.Sprintf("%d", auditReport.FilesScanned)},
		{"Directories scanned", fmt.Sprintf("%d", auditReport.DirectoriesScanned)},
		{"Score", fmt.Sprintf("%d/100", auditReport.Score)},
		{"Passed", fmt.Sprintf("%d", auditReport.PassedCount)},
		{"Warnings", fmt.Sprintf("%d", auditReport.WarningCount)},
		{"Failed", fmt.Sprintf("%d", auditReport.FailedCount)},
	}

	if _, err := fmt.Fprintln(w, "## Summary"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| Metric | Value |"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| --- | --- |"); err != nil {
		return err
	}
	for _, row := range summaryRows {
		if _, err := fmt.Fprintf(w, "| %s | %s |\n", escapeMarkdownTable(row[0]), escapeMarkdownTable(row[1])); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "## Findings"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| Status | ID | Category | Severity | Message | Remediation |"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| --- | --- | --- | --- | --- | --- |"); err != nil {
		return err
	}
	for _, finding := range auditReport.Findings {
		if _, err := fmt.Fprintf(
			w,
			"| %s | %s | %s | %s | %s | %s |\n",
			escapeMarkdownTable(string(finding.Status)),
			escapeMarkdownTable(finding.ID),
			escapeMarkdownTable(finding.Category),
			escapeMarkdownTable(string(finding.Severity)),
			escapeMarkdownTable(finding.Message),
			escapeMarkdownTable(finding.Remediation),
		); err != nil {
			return err
		}
	}

	return nil
}

func formatBool(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func escapeMarkdownTable(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "|", "\\|")
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "\r", " ")
	return value
}
