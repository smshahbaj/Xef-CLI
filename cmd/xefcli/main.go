// Package main is the entry point for XefCLI.
package main

import (
	"os"

	"github.com/xef/xefcli/internal/app"
)

var version = "v1.0.1"

func main() {
	a := app.New(version)
	if err := a.Execute(); err != nil {
		os.Exit(1)
	}
}
