package main

import (
	"fmt"
	"os"

	"github.com/cicd/internal/runner"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: cicd <pipeline.yaml>\n")
		os.Exit(1)
	}

	if err := runner.Run(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "❌ %v\n", err)
		os.Exit(1)
	}
}
