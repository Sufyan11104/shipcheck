package report

import (
	"encoding/json"
	"io"
	"path"
	"regexp"
	"strings"
	"unicode"

	"github.com/Sufyan11104/shipcheck/internal/rules"
	"github.com/Sufyan11104/shipcheck/internal/version"
)

const sarifSchemaURL = "https://json.schemastore.org/sarif-2.1.0.json"
const shipCheckInfoURI = "https://github.com/Sufyan11104/shipcheck"

// RenderSARIF writes a SARIF 2.1.0 report for warning and failure findings.
func RenderSARIF(w io.Writer, auditReport AuditReport) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(newSARIFLog(auditReport))
}

type sarifLog struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string                `json:"name"`
	InformationURI string                `json:"informationUri"`
	Version        string                `json:"version"`
	Rules          []sarifRuleDescriptor `json:"rules,omitempty"`
}

type sarifRuleDescriptor struct {
	ID                   string                    `json:"id"`
	Name                 string                    `json:"name,omitempty"`
	ShortDescription     sarifMessage              `json:"shortDescription,omitempty"`
	FullDescription      sarifMessage              `json:"fullDescription,omitempty"`
	DefaultConfiguration sarifDefaultConfiguration `json:"defaultConfiguration,omitempty"`
	Help                 sarifMessage              `json:"help,omitempty"`
	Properties           sarifRuleProperties       `json:"properties,omitempty"`
}

type sarifDefaultConfiguration struct {
	Level string `json:"level,omitempty"`
}

type sarifRuleProperties struct {
	Category          string   `json:"category,omitempty"`
	ShipCheckSeverity string   `json:"shipcheckSeverity,omitempty"`
	Tags              []string `json:"tags,omitempty"`
}

type sarifResult struct {
	RuleID     string                 `json:"ruleId"`
	Level      string                 `json:"level"`
	Message    sarifMessage           `json:"message"`
	Locations  []sarifLocation        `json:"locations,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

func newSARIFLog(auditReport AuditReport) sarifLog {
	results := make([]sarifResult, 0, auditReport.WarningCount+auditReport.FailedCount)
	ruleDescriptors := make([]sarifRuleDescriptor, 0)
	seenRules := make(map[string]bool)

	for _, finding := range auditReport.Findings {
		if !includeSARIFResult(finding) {
			continue
		}

		results = append(results, newSARIFResult(finding))
		if !seenRules[finding.ID] {
			ruleDescriptors = append(ruleDescriptors, newSARIFRuleDescriptor(finding))
			seenRules[finding.ID] = true
		}
	}

	return sarifLog{
		Version: "2.1.0",
		Schema:  sarifSchemaURL,
		Runs: []sarifRun{
			{
				Tool: sarifTool{
					Driver: sarifDriver{
						Name:           "ShipCheck",
						InformationURI: shipCheckInfoURI,
						Version:        version.Version,
						Rules:          ruleDescriptors,
					},
				},
				Results: results,
			},
		},
	}
}

func includeSARIFResult(finding ReportFinding) bool {
	return finding.Status == rules.StatusWarn || finding.Status == rules.StatusFail
}

func newSARIFRuleDescriptor(finding ReportFinding) sarifRuleDescriptor {
	help := finding.Remediation
	if help == "" || help == "N/A" {
		help = finding.Message
	}

	return sarifRuleDescriptor{
		ID:               finding.ID,
		Name:             finding.Title,
		ShortDescription: sarifMessage{Text: finding.Title},
		FullDescription:  sarifMessage{Text: finding.Message},
		DefaultConfiguration: sarifDefaultConfiguration{
			Level: sarifLevel(finding.Severity),
		},
		Help: sarifMessage{Text: help},
		Properties: sarifRuleProperties{
			Category:          finding.Category,
			ShipCheckSeverity: string(finding.Severity),
			Tags:              sarifTags(finding),
		},
	}
}

func newSARIFResult(finding ReportFinding) sarifResult {
	result := sarifResult{
		RuleID:  finding.ID,
		Level:   sarifLevel(finding.Severity),
		Message: sarifMessage{Text: finding.Message},
		Properties: map[string]interface{}{
			"category":          finding.Category,
			"shipcheckStatus":   string(finding.Status),
			"shipcheckSeverity": string(finding.Severity),
		},
	}

	if remediation := sanitizeSARIFText(finding.Remediation); remediation != "" && remediation != "N/A" {
		result.Properties["remediation"] = remediation
	}
	if evidence := sanitizeEvidence(finding.Evidence); evidence != "" {
		result.Properties["evidence"] = evidence
	}
	if uri := sanitizeSARIFURI(finding.Path); uri != "" {
		result.Locations = []sarifLocation{
			{
				PhysicalLocation: sarifPhysicalLocation{
					ArtifactLocation: sarifArtifactLocation{URI: uri},
				},
			},
		}
	}

	return result
}

func sarifLevel(severity rules.Severity) string {
	switch severity {
	case rules.SeverityInfo, rules.SeverityLow:
		return "note"
	case rules.SeverityMedium:
		return "warning"
	case rules.SeverityHigh, rules.Severity("critical"):
		return "error"
	default:
		return "warning"
	}
}

func sarifTags(finding ReportFinding) []string {
	tags := []string{"shipcheck"}
	if finding.Category != "" {
		tags = append(tags, finding.Category)
	}
	if finding.Severity != "" {
		tags = append(tags, string(finding.Severity))
	}
	return tags
}

func sanitizeSARIFURI(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	value = strings.ReplaceAll(value, "\\", "/")
	value = path.Clean(value)
	if value == "." || strings.HasPrefix(value, "/") || strings.Contains(value, ":") {
		return ""
	}
	if value == ".." || strings.HasPrefix(value, "../") {
		return ""
	}

	return value
}

func sanitizeSARIFText(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	var builder strings.Builder
	lastWasSpace := false
	for _, r := range value {
		if unicode.IsSpace(r) {
			if !lastWasSpace {
				builder.WriteRune(' ')
				lastWasSpace = true
			}
			continue
		}
		if r == unicode.ReplacementChar {
			continue
		}
		builder.WriteRune(r)
		lastWasSpace = false
	}

	return strings.TrimSpace(builder.String())
}

func sanitizeEvidence(value string) string {
	value = sanitizeSARIFText(value)
	if value == "" {
		return ""
	}

	parts := strings.Fields(value)
	for i, part := range parts {
		parts[i] = redactSensitiveAssignment(part)
	}
	value = strings.Join(parts, " ")

	const maxEvidenceLength = 240
	if len(value) > maxEvidenceLength {
		value = value[:maxEvidenceLength] + "..."
	}
	return value
}

var sensitiveEvidenceName = regexp.MustCompile(`(?i)(secret|token|password|passwd|api[_-]?key|private[_-]?key|credential)`)

func redactSensitiveAssignment(value string) string {
	for _, separator := range []string{"=", ":"} {
		index := strings.Index(value, separator)
		if index <= 0 {
			continue
		}
		name := strings.Trim(value[:index], `"'`)
		if sensitiveEvidenceName.MatchString(name) {
			return value[:index+1] + "[redacted]"
		}
	}
	return value
}
