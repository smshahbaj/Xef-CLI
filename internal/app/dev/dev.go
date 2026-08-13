// Package dev provides development helper commands used during development.
package dev

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/smshahbaj/Xef-CLI/internal/core/interfaces"
	"github.com/smshahbaj/Xef-CLI/internal/core/logger"
	"github.com/smshahbaj/Xef-CLI/internal/pkg/tui"
	"github.com/spf13/cobra"
)

// NewCommand creates the dev command group.
func NewCommand(fs interfaces.FileSystem, log logger.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Development tools",
		Long:  "Project scaffolding, environment management, and dev workflows.",
	}

	cmd.AddCommand(newProjectCreateCmd(fs, log))
	cmd.AddCommand(newEnvCmd(log))
	return cmd
}

func newProjectCreateCmd(fs interfaces.FileSystem, _ logger.Logger) *cobra.Command {
	var lang string

	cmd := &cobra.Command{
		Use:     "project [name]",
		Short:   "Create a new project scaffold",
		Example: `  xef dev project myapp --lang go`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			if fs.Exists(name) {
				return fmt.Errorf("directory %s already exists", name)
			}

			if err := fs.MkdirAll(name, 0o755); err != nil {
				return fmt.Errorf("failed to create project: %w", err)
			}

			switch strings.ToLower(lang) {
			case "go":
				if err := scaffoldGo(fs, name); err != nil {
					_ = os.RemoveAll(name)
					return fmt.Errorf("failed to scaffold Go project: %w", err)
				}
			case "python":
				if err := scaffoldPython(fs, name); err != nil {
					_ = os.RemoveAll(name)
					return fmt.Errorf("failed to scaffold Python project: %w", err)
				}
			default:
				_ = os.RemoveAll(name)
				return fmt.Errorf("unsupported language: %s (supported: go, python)", lang)
			}

			tui.PrintSuccess(fmt.Sprintf("Created %s project: %s", lang, name))
			return nil
		},
	}

	cmd.Flags().StringVar(&lang, "lang", "go", "project language (go, python)")
	return cmd
}

func scaffoldGo(fs interfaces.FileSystem, name string) error {
	projectName := filepath.Base(filepath.Clean(name))
	dirs := []string{
		filepath.Join(name, "cmd", projectName),
		filepath.Join(name, "internal"),
		filepath.Join(name, "pkg"),
		filepath.Join(name, "configs"),
		filepath.Join(name, "scripts"),
	}

	for _, d := range dirs {
		if err := fs.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}

	mainGo := fmt.Sprintf(`package main

import "fmt"

func main() {
	fmt.Println("Hello from %s!")
}
`, projectName)

	files := map[string]string{
		filepath.Join(name, "go.mod"):                      fmt.Sprintf("module %s\n\ngo 1.22\n", projectName),
		filepath.Join(name, "cmd", projectName, "main.go"): mainGo,
		filepath.Join(name, "README.md"):                   fmt.Sprintf("# %s\n\nA Go project.\n", toTitle(projectName)),
		filepath.Join(name, "Makefile"):                    ".PHONY: build\nbuild:\n\tgo build -o bin/" + projectName + " ./cmd/" + projectName + "\n",
		filepath.Join(name, ".gitignore"):                  "/bin/\n*.exe\n*.test\n/vendor/\n",
	}

	for path, content := range files {
		if err := fs.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func scaffoldPython(fs interfaces.FileSystem, name string) error {
	projectName := filepath.Base(filepath.Clean(name))
	dirs := []string{
		filepath.Join(name, projectName),
		filepath.Join(name, "tests"),
		filepath.Join(name, "docs"),
	}

	for _, d := range dirs {
		if err := fs.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}

	files := map[string]string{
		filepath.Join(name, "pyproject.toml"): fmt.Sprintf(`[project]
name = "%s"
version = "0.1.0"
description = "A Python project"
requires-python = ">=3.10"
`, projectName),
		filepath.Join(name, projectName, "__init__.py"): "",
		filepath.Join(name, projectName, "main.py"):     "def main():\n    print('Hello, World!')\n\nif __name__ == '__main__':\n    main()\n",
		filepath.Join(name, "README.md"):                fmt.Sprintf("# %s\n\nA Python project.\n", toTitle(projectName)),
		filepath.Join(name, ".gitignore"):               "__pycache__/\n*.pyc\n.venv/\n",
	}

	for path, content := range files {
		if err := fs.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// toTitle converts a string to title case without using deprecated strings.Title.
func toTitle(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func newEnvCmd(_ logger.Logger) *cobra.Command {
	var format string
	var showValues bool

	cmd := &cobra.Command{
		Use:     "env",
		Short:   "Show environment variables",
		Example: `  xef dev env --format json`,
		RunE: func(_ *cobra.Command, _ []string) error {
			envs := os.Environ()
			masked := make([]string, 0, len(envs))
			values := make(map[string]string, len(envs))
			for _, e := range envs {
				parts := strings.SplitN(e, "=", 2)
				if len(parts) != 2 {
					continue
				}
				value := parts[1]
				if !showValues {
					value = "<redacted>"
				}
				masked = append(masked, parts[0]+"="+value)
				values[parts[0]] = value
			}
			switch strings.ToLower(format) {
			case "json":
				data, err := json.MarshalIndent(values, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to encode environment: %w", err)
				}
				fmt.Println(string(data))
			case "export":
				for _, e := range masked {
					fmt.Printf("export %s\n", e)
				}
			case "list":
				for _, e := range masked {
					fmt.Println(e)
				}
			default:
				return fmt.Errorf("unsupported format: %s (use list, json, export)", format)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "list", "output format: list, json, export")
	cmd.Flags().BoolVar(&showValues, "show-values", false, "show environment values (may expose secrets)")
	return cmd
}
