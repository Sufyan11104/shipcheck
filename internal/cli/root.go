package cli

import (
	"flag"
	"fmt"

	"github.com/Sufyan11104/shipcheck/internal/report"
	"github.com/Sufyan11104/shipcheck/internal/scanner"
	"github.com/Sufyan11104/shipcheck/internal/version"
)

// Run processes CLI commands
func Run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no command provided")
	}

	command := args[0]

	switch command {
	case "version":
		fmt.Println(version.Version)
		return nil

	case "audit":
		return handleAudit(args[1:])

	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func handleAudit(args []string) error {
	flagSet := flag.NewFlagSet("audit", flag.ExitOnError)
	flagSet.Parse(args)

	positionalArgs := flagSet.Args()
	if len(positionalArgs) < 1 {
		return fmt.Errorf("audit requires a path argument")
	}

	path := positionalArgs[0]

	// Scan the directory
	result := scanner.Scan(path)

	// Print the report
	return report.PrintTextReport(result)
}
