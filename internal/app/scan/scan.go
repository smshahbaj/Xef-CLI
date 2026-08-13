// Package scan provides project scanning and reporting commands.
package scan

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type Finding struct {
	File     string `json:"file"`
	Type     string `json:"type"`
	Detail   string `json:"detail"`
	Severity string `json:"severity"`
}

type Result struct {
	Root      string    `json:"root"`
	Files     int       `json:"files"`
	Findings  []Finding `json:"findings"`
	Generated string    `json:"generated"`
}

var patterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bAKIA[0-9A-Z]{16}\b`),
	regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9_]{20,}\b`),
	regexp.MustCompile(`\bsk_(?:live|test)_[A-Za-z0-9]{16,}\b`),
	regexp.MustCompile(`\bxox[baprs]-[A-Za-z0-9-]{20,}\b`),
	regexp.MustCompile(`-----BEGIN (?:RSA |EC |OPENSSH |DSA )?PRIVATE KEY-----`),
	regexp.MustCompile(`(?i)\b(?:api[_-]?key|secret[_-]?key|access[_-]?key|token|password)\s*[:=]\s*["']?[A-Za-z0-9_./+=-]{12,}`),
}

var placeholders = []string{"placeholder", "example", "changeme", "change-me", "dummy", "fake", "sample", "super_secret", "your_", "replace_me"}

func NewCommand() *cobra.Command {
	var root string
	var jsonOut bool
	var max int64
	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan project files for common issues",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			r, err := runScan(root, max, true)
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(r)
			}
			fmt.Printf("Project Scan: %s\nFiles scanned: %d\nFindings: %d\n", r.Root, r.Files, len(r.Findings))
			printFindings(r.Findings)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project directory to scan")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "output machine-readable JSON")
	cmd.Flags().Int64Var(&max, "max-file-size", 10*1024*1024, "report files larger than this many bytes")
	return cmd
}

func NewSecretCommand() *cobra.Command {
	parent := &cobra.Command{Use: "secret", Short: "Security scanning tools"}
	var root string
	var jsonOut bool
	child := &cobra.Command{
		Use:   "scan",
		Short: "Scan project files for exposed credentials",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			r, err := runScan(root, 0, false)
			if err != nil {
				return err
			}
			filtered := r
			filtered.Findings = make([]Finding, 0)
			for _, f := range r.Findings {
				if f.Type == "secret" {
					filtered.Findings = append(filtered.Findings, f)
				}
			}
			if jsonOut {
				return printJSON(filtered)
			}
			fmt.Printf("Secret Scan: %s\nFiles scanned: %d\nSecrets found: %d\n", filtered.Root, filtered.Files, len(filtered.Findings))
			printFindings(filtered.Findings)
			if len(filtered.Findings) > 0 {
				return fmt.Errorf("secret scan found %d potential credential(s)", len(filtered.Findings))
			}
			return nil
		},
	}
	child.Flags().StringVar(&root, "root", ".", "project directory to scan")
	child.Flags().BoolVar(&jsonOut, "json", false, "output machine-readable JSON")
	parent.AddCommand(child)
	return parent
}

func NewReportCommand() *cobra.Command {
	var root, output string
	var openJSON bool
	cmd := &cobra.Command{Use: "report", Short: "Generate an HTML project health report", Args: cobra.NoArgs, RunE: func(_ *cobra.Command, _ []string) error {
		r, err := runScan(root, 10*1024*1024, true)
		if err != nil {
			return err
		}
		if openJSON {
			return printJSON(r)
		}
		if output == "" {
			output = filepath.Join(root, "xef-report.html")
		}
		return writeReport(output, r)
	}}
	cmd.Flags().StringVar(&root, "root", ".", "project directory to report")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output HTML file")
	cmd.Flags().BoolVar(&openJSON, "json", false, "output report data as JSON")
	return cmd
}

func runScan(root string, max int64, includeQuality bool) (Result, error) {
	root = filepath.Clean(root)
	info, err := os.Stat(root)
	if err != nil {
		return Result{}, fmt.Errorf("project directory: %w", err)
	}
	if !info.IsDir() {
		return Result{}, fmt.Errorf("project path is not a directory: %s", root)
	}
	r := Result{Root: root, Generated: time.Now().UTC().Format(time.RFC3339), Findings: []Finding{}}
	err = filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if info.IsDir() {
			if path != root && shouldSkipDir(root, path) {
				return filepath.SkipDir
			}
			return nil
		}
		if shouldSkipFile(root, path) {
			return nil
		}
		r.Files++
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		text := string(data)
		if containsSecret(text) {
			r.Findings = append(r.Findings, Finding{rel(root, path), "secret", "credential-like pattern detected", "HIGH"})
		}
		if includeQuality {
			if info.Size() > max && max > 0 {
				r.Findings = append(r.Findings, Finding{rel(root, path), "size", fmt.Sprintf("file is larger than %d bytes", max), "MEDIUM"})
			}
			if containsTodoComment(text) {
				r.Findings = append(r.Findings, Finding{rel(root, path), "maintenance", "TODO/FIXME comment found", "LOW"})
			}
		}
		return nil
	})
	if err != nil {
		return Result{}, err
	}
	sort.Slice(r.Findings, func(i, j int) bool {
		if r.Findings[i].File == r.Findings[j].File {
			return r.Findings[i].Type < r.Findings[j].Type
		}
		return r.Findings[i].File < r.Findings[j].File
	})
	return r, nil
}

func shouldSkipDir(root, path string) bool {
	base := strings.ToLower(filepath.Base(path))
	switch base {
	case ".git", "node_modules", "dist", "build", "bin", ".idea", ".vscode":
		return true
	}
	return isIgnored(root, path)
}
func shouldSkipFile(root, path string) bool {
	base := strings.ToLower(filepath.Base(path))
	switch filepath.Ext(base) {
	case ".exe", ".dll", ".so", ".dylib", ".test", ".out":
		return true
	}
	if base == "coverage.out" || base == "coverage.html" {
		return true
	}
	return isIgnored(root, path)
}
func isIgnored(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	rel = filepath.ToSlash(rel)

	// Prefer Git's own matcher when this is a real Git repository.
	if _, err := os.Stat(filepath.Join(root, ".git")); err == nil {
		cmd := exec.Command("git", "-C", root, "check-ignore", "--quiet", "--", rel)
		if cmd.Run() == nil {
			return true
		}
	}

	// Fall back to a small, dependency-free .gitignore matcher. This also
	// makes scans deterministic for temporary/test projects with a .gitignore
	// but no initialized Git index.
	data, err := os.ReadFile(filepath.Join(root, ".gitignore"))
	if err != nil {
		return false
	}
	ignored := false
	for _, raw := range strings.Split(string(data), "\n") {
		pattern := strings.TrimSpace(raw)
		if pattern == "" || strings.HasPrefix(pattern, "#") {
			continue
		}
		negated := strings.HasPrefix(pattern, "!")
		if negated {
			pattern = strings.TrimPrefix(pattern, "!")
		}
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		pattern = strings.TrimPrefix(pattern, "/")
		matched := false
		if strings.Contains(pattern, "/") {
			matched, _ = filepath.Match(pattern, rel)
		} else {
			matched, _ = filepath.Match(pattern, filepath.Base(rel))
		}
		if matched {
			ignored = !negated
		}
	}
	return ignored
}
func rel(root, path string) string {
	r, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.Base(path)
	}
	return filepath.ToSlash(r)
}
func containsSecret(text string) bool {
	for _, p := range patterns {
		for _, m := range p.FindAllString(text, -1) {
			low := strings.ToLower(m)
			skip := false
			for _, term := range placeholders {
				if strings.Contains(low, term) {
					skip = true
					break
				}
			}
			if !skip {
				return true
			}
		}
	}
	return false
}
func containsTodoComment(text string) bool {
	for _, line := range strings.Split(text, "\n") {
		t := strings.TrimSpace(line)
		for _, p := range []string{"//", "#", ";", "--", "<!--", "*"} {
			if strings.HasPrefix(t, p) {
				rest := strings.TrimSpace(strings.TrimPrefix(t, p))
				u := strings.ToUpper(rest)
				if strings.HasPrefix(u, "TODO") || strings.HasPrefix(u, "FIXME") {
					return true
				}
			}
		}
	}
	return false
}
func printFindings(fs []Finding) {
	for _, f := range fs {
		fmt.Printf("[%s] %s: %s (%s)\n", f.Severity, f.File, f.Detail, f.Type)
	}
	if len(fs) == 0 {
		fmt.Println("✓ No findings")
	}
}
func printJSON(v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

var reportTemplate = template.Must(template.New("report").Parse(`<!doctype html><html><head><meta charset="utf-8"><title>XefCLI Project Report</title><style>body{font-family:system-ui,sans-serif;max-width:1000px;margin:40px auto;padding:0 20px}table{width:100%;border-collapse:collapse}th,td{padding:10px;border-bottom:1px solid #ddd;text-align:left}.ok{color:#16803c}.high{color:#b42318}.medium{color:#b54708}.low{color:#667085}</style></head><body><h1>XefCLI Project Report</h1><p><b>Root:</b> {{.Root}}</p><p><b>Files scanned:</b> {{.Files}} &nbsp; <b>Generated:</b> {{.Generated}}</p>{{if .Findings}}<table><tr><th>Severity</th><th>File</th><th>Type</th><th>Detail</th></tr>{{range .Findings}}<tr><td>{{.Severity}}</td><td>{{.File}}</td><td>{{.Type}}</td><td>{{.Detail}}</td></tr>{{end}}</table>{{else}}<h2 class="ok">✓ No findings</h2>{{end}}</body></html>`))

func writeReport(path string, r Result) error {
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if err = reportTemplate.Execute(f, r); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err = f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			os.Remove(tmp)
			return fmt.Errorf("replace report: %w", err)
		}
	} else if !os.IsNotExist(err) {
		os.Remove(tmp)
		return err
	}
	if err = os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return err
	}
	fmt.Printf("✓ Report written to %s\n", path)
	return nil
}
