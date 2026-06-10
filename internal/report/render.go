package report

import (
	"fmt"
	"io"
)

// Render writes an audit report in the requested format.
func Render(w io.Writer, auditReport AuditReport, format string) error {
	switch format {
	case FormatText:
		return RenderText(w, auditReport)
	case FormatJSON:
		return RenderJSON(w, auditReport)
	case FormatMarkdown:
		return RenderMarkdown(w, auditReport)
	default:
		return fmt.Errorf("unknown report format %q (valid: text, json, markdown)", format)
	}
}

// IsValidFormat returns whether format is a supported report format.
func IsValidFormat(format string) bool {
	switch format {
	case FormatText, FormatJSON, FormatMarkdown:
		return true
	default:
		return false
	}
}
