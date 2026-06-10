package cli

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Sufyan11104/shipcheck/internal/engine"
	"github.com/Sufyan11104/shipcheck/internal/report"
	"github.com/Sufyan11104/shipcheck/internal/scanner"
	"github.com/Sufyan11104/shipcheck/internal/version"
)

type auditOptions struct {
	path       string
	format     string
	failUnder  int
	categories []string
}

// ExitError carries a CLI exit code without calling os.Exit inside command logic.
type ExitError struct {
	Code    int
	Message string
}

func (e *ExitError) Error() string {
	return e.Message
}

func (e *ExitError) ExitCode() int {
	return e.Code
}

// Run processes CLI commands
func Run(args []string) error {
	return RunWithWriter(args, os.Stdout)
}

// RunWithWriter processes CLI commands and writes command output to w.
func RunWithWriter(args []string, w io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("no command provided")
	}

	command := args[0]

	switch command {
	case "version":
		fmt.Fprintln(w, version.Version)
		return nil

	case "audit":
		return handleAudit(args[1:], w)

	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func handleAudit(args []string, w io.Writer) error {
	options, err := parseAuditArgs(args)
	if err != nil {
		return err
	}

	result := scanner.Scan(options.path)
	if result.Error != nil {
		return result.Error
	}

	eng := engine.NewEngine(result.Path)
	findings, _ := eng.RunChecks(result.IsGitRepository)
	findings = engine.FilterFindingsByCategory(findings, options.categories)
	score := engine.CalculateScore(findings)

	auditReport := report.NewAuditReport(result, findings, score)
	if err := report.Render(w, auditReport, options.format); err != nil {
		return err
	}

	if err := EvaluateFailUnder(score, options.failUnder); err != nil {
		if options.format == report.FormatText {
			fmt.Fprintf(w, "\n%s\n", err.Error())
		}
		return err
	}

	return nil
}

// EvaluateFailUnder returns an ExitError if score is below threshold.
func EvaluateFailUnder(score, threshold int) error {
	if threshold <= 0 || score >= threshold {
		return nil
	}

	return &ExitError{
		Code:    1,
		Message: fmt.Sprintf("Score %d is below fail-under threshold %d", score, threshold),
	}
}

func parseAuditArgs(args []string) (auditOptions, error) {
	options := auditOptions{
		format: report.FormatText,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "--format":
			value, next, err := requireFlagValue(args, i, "--format")
			if err != nil {
				return options, err
			}
			i = next
			if err := setFormat(&options, value); err != nil {
				return options, err
			}
		case strings.HasPrefix(arg, "--format="):
			if err := setFormat(&options, strings.TrimPrefix(arg, "--format=")); err != nil {
				return options, err
			}
		case arg == "--fail-under":
			value, next, err := requireFlagValue(args, i, "--fail-under")
			if err != nil {
				return options, err
			}
			i = next
			if err := setFailUnder(&options, value); err != nil {
				return options, err
			}
		case strings.HasPrefix(arg, "--fail-under="):
			if err := setFailUnder(&options, strings.TrimPrefix(arg, "--fail-under=")); err != nil {
				return options, err
			}
		case arg == "--category":
			value, next, err := requireFlagValue(args, i, "--category")
			if err != nil {
				return options, err
			}
			i = next
			if err := setCategories(&options, value); err != nil {
				return options, err
			}
		case strings.HasPrefix(arg, "--category="):
			if err := setCategories(&options, strings.TrimPrefix(arg, "--category=")); err != nil {
				return options, err
			}
		case strings.HasPrefix(arg, "-"):
			return options, fmt.Errorf("unknown audit flag: %s", arg)
		default:
			if options.path != "" {
				return options, fmt.Errorf("unexpected argument: %s", arg)
			}
			options.path = arg
		}
	}

	if options.path == "" {
		return options, fmt.Errorf("audit requires a path argument")
	}

	return options, nil
}

func requireFlagValue(args []string, index int, name string) (string, int, error) {
	next := index + 1
	if next >= len(args) || strings.HasPrefix(args[next], "-") {
		return "", index, fmt.Errorf("%s requires a value", name)
	}

	return args[next], next, nil
}

func setFormat(options *auditOptions, value string) error {
	format := strings.ToLower(strings.TrimSpace(value))
	if !report.IsValidFormat(format) {
		return fmt.Errorf("unknown report format %q (valid: text, json, markdown)", value)
	}

	options.format = format
	return nil
}

func setFailUnder(options *auditOptions, value string) error {
	threshold, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return fmt.Errorf("invalid --fail-under value %q: must be an integer", value)
	}
	if threshold < 0 {
		return fmt.Errorf("invalid --fail-under value %q: must be 0 or greater", value)
	}

	options.failUnder = threshold
	return nil
}

func setCategories(options *auditOptions, value string) error {
	categories, err := engine.ParseCategoryFilter(value)
	if err != nil {
		return err
	}

	options.categories = categories
	return nil
}
