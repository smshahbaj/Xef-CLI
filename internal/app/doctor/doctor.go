// Package doctor provides deterministic, local project health diagnostics.
package doctor

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/smshahbaj/Xef-CLI/internal/pkg/tui"
	"github.com/spf13/cobra"
)

type finding struct {
	Check   string `json:"check"`
	Status  string `json:"status"`
	Detail  string `json:"detail"`
	Fixable bool   `json:"fixable"`
}

type fixResult struct {
	Action string `json:"action"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}

var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bAKIA[0-9A-Z]{16}\b`),
	regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9_]{20,}\b`),
	regexp.MustCompile(`\bsk_(?:live|test)_[A-Za-z0-9]{16,}\b`),
	regexp.MustCompile(`\bxox[baprs]-[A-Za-z0-9-]{20,}\b`),
	regexp.MustCompile(`-----BEGIN (?:RSA |EC |OPENSSH |DSA )?PRIVATE KEY-----`),
	regexp.MustCompile(`(?i)\b(?:api[_-]?key|secret[_-]?key|access[_-]?key|token|password)\s*[:=]\s*["']?[^\s"']{8,}`),
}

var placeholderSecretTerms = []string{
	"placeholder", "example", "changeme", "change-me", "dummy", "fake", "sample", "super_secret", "your_", "replace_me",
}

// NewCommand creates the project health doctor command.
func NewCommand() *cobra.Command {
	var jsonOut bool
	var root string
	var fix bool
	var strict bool
	var maxFileSize int64

	cmd := &cobra.Command{
		Use:     "doctor",
		Short:   "Diagnose project health and common development problems",
		Example: "  xef doctor\n  xef doctor --json\n  xef doctor --fix\n  xef doctor --strict --root ./project",
		RunE: func(_ *cobra.Command, _ []string) error {
			var fixes []fixResult
			if fix {
				results, err := applySafeFixes(root)
				if err != nil {
					return err
				}
				fixes = results
				if !jsonOut {
					tui.PrintTitle("Safe Fixes")
					for _, result := range results {
						if result.Status == "PASS" {
							tui.PrintSuccess(result.Detail)
						} else {
							fmt.Printf("- %s: %s\n", result.Status, result.Detail)
						}
					}
				}
			}

			findings, score := diagnoseWithOptions(root, maxFileSize)
			if jsonOut {
				payload := struct {
					Score    int         `json:"score"`
					Findings []finding   `json:"findings"`
					Fixes    []fixResult `json:"fixes,omitempty"`
				}{score, findings, fixes}
				data, err := json.MarshalIndent(payload, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
			} else {
				tui.PrintTitle("Xef Project Doctor")
				for _, f := range findings {
					switch f.Status {
					case "PASS":
						tui.PrintSuccess(f.Check + ": " + f.Detail)
					case "WARN":
						fmt.Printf("⚠ %s: %s\n", f.Check, f.Detail)
					default:
						fmt.Printf("✗ %s: %s\n", f.Check, f.Detail)
					}
				}
				fmt.Printf("\nHealth Score: %d/100\n", score)
			}

			if strict {
				for _, f := range findings {
					if f.Status == "FAIL" || f.Status == "WARN" {
						return errors.New("doctor found issues in strict mode")
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "output machine-readable JSON")
	cmd.Flags().StringVar(&root, "root", ".", "project directory to inspect")
	cmd.Flags().BoolVar(&fix, "fix", false, "apply only safe, non-destructive fixes before diagnosis")
	cmd.Flags().BoolVar(&strict, "strict", false, "exit with an error when any warning or failure is found")
	cmd.Flags().Int64Var(&maxFileSize, "max-file-size", 10*1024*1024, "warn about files larger than this many bytes")
	return cmd
}

func applySafeFixes(root string) ([]fixResult, error) {
	root = filepath.Clean(root)
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("project directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("project path is not a directory: %s", root)
	}

	results := make([]fixResult, 0, 3)
	createIfMissing := func(name, contents string, mode os.FileMode) error {
		path := filepath.Join(root, name)
		if _, err := os.Stat(path); err == nil {
			results = append(results, fixResult{"check", "SKIP", name + " already exists"})
			return nil
		} else if !os.IsNotExist(err) {
			return err
		}
		if err := os.WriteFile(path, []byte(contents), mode); err != nil {
			return fmt.Errorf("create %s: %w", name, err)
		}
		results = append(results, fixResult{"create", "PASS", "created " + name})
		return nil
	}

	if err := createIfMissing("README.md", "# Project\n\nProject documentation.\n", 0o644); err != nil {
		return nil, err
	}
	if err := createIfMissing(".gitignore", "bin/\ndist/\nbuild/\ncoverage.out\n.env\n", 0o644); err != nil {
		return nil, err
	}
	gitignorePath := filepath.Join(root, ".gitignore")
	if envInfo, err := os.Stat(filepath.Join(root, ".env")); err == nil && !envInfo.IsDir() {
		data, readErr := os.ReadFile(gitignorePath)
		if readErr != nil {
			return nil, fmt.Errorf("read .gitignore: %w", readErr)
		}
		if !gitignoreContainsEnv(string(data)) {
			updated := strings.TrimRight(string(data), "\n") + "\n.env\n"
			if err := os.WriteFile(gitignorePath, []byte(updated), 0o644); err != nil {
				return nil, fmt.Errorf("protect .env in .gitignore: %w", err)
			}
			results = append(results, fixResult{"protect", "PASS", "added .env to .gitignore"})
		}
	}
	return results, nil
}

func diagnose(root string) ([]finding, int) {
	return diagnoseWithOptions(root, 10*1024*1024)
}

func diagnoseWithOptions(root string, maxFileSize int64) ([]finding, int) {
	root = filepath.Clean(root)
	findings := make([]finding, 0, 9)
	add := func(check, status, detail string, fixable bool) {
		findings = append(findings, finding{Check: check, Status: status, Detail: detail, Fixable: fixable})
	}

	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		return []finding{{Check: "Project", Status: "FAIL", Detail: "directory does not exist", Fixable: false}}, 0
	}

	if _, err := os.Stat(filepath.Join(root, ".git")); err == nil {
		add("Git", "PASS", "repository detected", false)
		if tracked, err := exec.Command("git", "-C", root, "ls-files", "--error-unmatch", ".env").CombinedOutput(); err == nil && strings.TrimSpace(string(tracked)) != "" {
			add("Git secrets", "FAIL", ".env is tracked by Git", false)
		} else {
			add("Git secrets", "PASS", ".env is not tracked by Git", false)
		}
	} else {
		add("Git", "WARN", "no .git directory detected", false)
	}
	if _, err := os.Stat(filepath.Join(root, "README.md")); err == nil {
		add("Documentation", "PASS", "README.md present", false)
	} else {
		add("Documentation", "WARN", "README.md is missing", true)
	}
	if _, err := os.Stat(filepath.Join(root, ".gitignore")); err == nil {
		add("Git hygiene", "PASS", ".gitignore present", false)
	} else {
		add("Git hygiene", "WARN", ".gitignore is missing", true)
	}

	if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
		if _, err := exec.LookPath("go"); err == nil {
			add("Go", "PASS", "go toolchain available", false)
		} else {
			add("Go", "FAIL", "go toolchain not found", false)
		}
	}

	secrets, todos, large := 0, 0, 0
	walkErr := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			switch info.Name() {
			case ".git", ".idea", ".vscode", "vendor", "node_modules", "dist", "build", "bin":
				return filepath.SkipDir
			}
			return nil
		}
		if isIgnoredPath(root, path) {
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		if maxFileSize > 0 && info.Size() > maxFileSize {
			large++
			return nil
		}
		if isGeneratedArtifact(path) {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !isScannableExtension(ext) && filepath.Base(path) != ".env" {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil || len(data) > 2*1024*1024 {
			return nil
		}
		text := string(data)
		if containsTodoComment(text) {
			todos++
		}
		// Test files commonly contain deliberately fake credentials used by
		// regression tests. They are not deployable configuration, so skip
		// them for credential findings to avoid noisy false positives.
		if !isTestFixtureFile(path) && containsLikelySecret(text) {
			secrets++
		}
		return nil
	})
	if walkErr != nil {
		add("Filesystem", "WARN", "some files could not be inspected", false)
	}

	if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
		if _, err := os.Stat(filepath.Join(root, "go.sum")); err != nil {
			add("Go modules", "WARN", "go.mod exists but go.sum is missing", true)
		} else {
			add("Go modules", "PASS", "go.mod and go.sum present", false)
		}
		if _, err := exec.LookPath("gofmt"); err == nil {
			if files := gofmtFiles(root); len(files) == 0 {
				add("Go formatting", "PASS", "all Go files are gofmt-clean", false)
			} else {
				add("Go formatting", "WARN", fmt.Sprintf("%d Go file(s) need gofmt", len(files)), true)
			}
		}
	}
	if envInfo, err := os.Stat(filepath.Join(root, ".env")); err == nil && !envInfo.IsDir() {
		data, readErr := os.ReadFile(filepath.Join(root, ".gitignore"))
		if readErr == nil && gitignoreContainsEnv(string(data)) {
			add("Environment hygiene", "PASS", ".env is ignored by Git", false)
		} else {
			add("Environment hygiene", "WARN", ".env exists but is not protected by .gitignore", true)
		}
	}

	if secrets == 0 {
		add("Secrets", "PASS", "no obvious credential patterns found", false)
	} else {
		add("Secrets", "FAIL", fmt.Sprintf("%d file(s) contain secret-like patterns", secrets), false)
	}
	if todos == 0 {
		add("Maintenance", "PASS", "no TODO/FIXME markers found", false)
	} else {
		add("Maintenance", "WARN", fmt.Sprintf("%d file(s) contain TODO/FIXME markers", todos), false)
	}
	if large == 0 {
		add("Repository size", "PASS", fmt.Sprintf("no file exceeds %s", formatBytes(maxFileSize)), false)
	} else {
		add("Repository size", "WARN", fmt.Sprintf("%d file(s) exceed %s", large, formatBytes(maxFileSize)), false)
	}

	sort.Slice(findings, func(i, j int) bool { return findings[i].Check < findings[j].Check })
	score := 100
	for _, f := range findings {
		switch f.Status {
		case "FAIL":
			score -= 20
		case "WARN":
			score -= 5
		}
	}
	if score < 0 {
		score = 0
	}
	return findings, score
}

func isTestFixtureFile(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	return strings.HasSuffix(base, "_test.go") || strings.Contains(base, ".test.")
}

func isGeneratedArtifact(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	switch filepath.Ext(base) {
	case ".exe", ".dll", ".so", ".dylib", ".test", ".out":
		return true
	}
	return base == "coverage.out" || base == "coverage.html"
}

func containsLikelySecret(text string) bool {
	for _, pattern := range secretPatterns {
		for _, match := range pattern.FindAllString(text, -1) {
			if isPlaceholderSecret(match) {
				continue
			}
			return true
		}
	}
	return false
}

func isPlaceholderSecret(value string) bool {
	lower := strings.ToLower(value)
	for _, term := range placeholderSecretTerms {
		if strings.Contains(lower, term) {
			return true
		}
	}
	return false
}

func containsTodoComment(text string) bool {
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		for _, prefix := range []string{"//", "#", ";", "--", "<!--", "*"} {
			if strings.HasPrefix(trimmed, prefix) {
				rest := strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
				upper := strings.ToUpper(rest)
				if strings.HasPrefix(upper, "TODO") || strings.HasPrefix(upper, "FIXME") {
					return true
				}
			}
		}
	}
	return false
}

func isIgnoredPath(root, path string) bool {
	if _, err := os.Stat(filepath.Join(root, ".git")); err != nil {
		return false
	}
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	rel = filepath.ToSlash(rel)
	cmd := exec.Command("git", "-C", root, "check-ignore", "-q", "--no-index", "--", rel)
	return cmd.Run() == nil
}

func gitignoreContainsEnv(data string) bool {
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == ".env" || line == "*.env" || line == "/.env" {
			return true
		}
	}
	return false
}

func gofmtFiles(root string) []string {
	files := make([]string, 0)
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			switch info.Name() {
			case ".git", "vendor", "node_modules", "dist", "build", "bin":
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) == ".go" {
			files = append(files, path)
		}
		return nil
	})
	if len(files) == 0 {
		return nil
	}
	args := append([]string{"-l"}, files...)
	out, err := exec.Command("gofmt", args...).Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return nil
	}
	sort.Strings(lines)
	return lines
}

func formatBytes(n int64) string {
	if n <= 0 {
		return "disabled"
	}
	units := []string{"B", "KiB", "MiB", "GiB"}
	v := float64(n)
	i := 0
	for v >= 1024 && i < len(units)-1 {
		v /= 1024
		i++
	}
	if i == 0 {
		return fmt.Sprintf("%d %s", n, units[i])
	}
	return fmt.Sprintf("%.1f %s", v, units[i])
}

func isScannableExtension(ext string) bool {
	switch ext {
	case ".go", ".py", ".js", ".ts", ".tsx", ".jsx", ".java", ".kt", ".rs", ".c", ".h", ".cpp", ".hpp", ".cs", ".php", ".rb", ".swift", ".yaml", ".yml", ".json", ".toml", ".ini", ".xml", ".sh", ".ps1", ".env":
		return true
	default:
		return false
	}
}
