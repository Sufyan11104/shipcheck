package main

import (
	"fmt"
	"os"

	"github.com/Sufyan11104/shipcheck/internal/cli"
)

func main() {
	err := cli.Run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
