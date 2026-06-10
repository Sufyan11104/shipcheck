package report

import (
	"encoding/json"
	"io"
)

// RenderJSON writes a stable machine-readable audit report.
func RenderJSON(w io.Writer, auditReport AuditReport) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(auditReport)
}
