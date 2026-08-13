package scan

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunScanFindsSecretAndIgnoresPlaceholder(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "safe.txt"), []byte(`api_key=super_secret_value_123456789`), 0o600)
	os.WriteFile(filepath.Join(root, "real.env"), []byte(`api_key=sk_live_`+strings.Repeat("A", 24)), 0o600)
	r, err := runScan(root, 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Findings) != 1 || r.Findings[0].Type != "secret" {
		t.Fatalf("unexpected findings: %+v", r.Findings)
	}
}

func TestRunScanIgnoresGeneratedAndGitIgnoredFiles(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, ".gitignore"), []byte("*.exe\nignored.txt\n"), 0o644)
	os.Mkdir(filepath.Join(root, ".git"), 0o755)
	os.WriteFile(filepath.Join(root, "ignored.txt"), []byte(`api_key=sk_live_`+strings.Repeat("A", 24)), 0o600)
	os.WriteFile(filepath.Join(root, "xef.exe"), []byte(`api_key=sk_live_`+strings.Repeat("A", 24)), 0o600)
	os.WriteFile(filepath.Join(root, "ok.txt"), []byte("hello"), 0o600)
	r, err := runScan(root, 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Findings) != 0 {
		t.Fatalf("unexpected findings: %+v", r.Findings)
	}
}

func TestContainsTodoComment(t *testing.T) {
	if containsTodoComment(`fmt.Println("TODO: example")`) {
		t.Fatal("string literal should not count")
	}
	if !containsTodoComment(`// TODO: fix this`) {
		t.Fatal("comment should count")
	}
}

func TestWriteReport(t *testing.T) {
	root := t.TempDir()
	out := filepath.Join(root, "report.html")
	r := Result{Root: root, Files: 1, Generated: "now", Findings: []Finding{{File: "x.txt", Type: "secret", Detail: "credential-like pattern detected", Severity: "HIGH"}}}
	if err := writeReport(out, r); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "XefCLI Project Report") || !strings.Contains(string(data), "credential-like pattern detected") {
		t.Fatal("report content incomplete")
	}
}

func TestCommandStructure(t *testing.T) {
	if NewCommand().Name() != "scan" {
		t.Fatal("scan command name mismatch")
	}
	secret := NewSecretCommand()
	if secret.Name() != "secret" {
		t.Fatal("secret command name mismatch")
	}
	child := secret.Commands()
	if len(child) != 1 || child[0].Name() != "scan" {
		t.Fatal("secret scan command missing")
	}
	if NewReportCommand().Name() != "report" {
		t.Fatal("report command name mismatch")
	}
}

func TestRunScanJSONFindingsAreNonNil(t *testing.T) {
	root := t.TempDir()
	r, err := runScan(root, 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if r.Findings == nil {
		t.Fatal("findings must be an empty array, not nil")
	}
}
