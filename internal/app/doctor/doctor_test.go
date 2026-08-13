package doctor

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiagnoseHealthyProject(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("# test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("bin/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	findings, score := diagnose(root)
	if score < 80 {
		t.Fatalf("expected healthy score >= 80, got %d (%v)", score, findings)
	}
}

func TestDoctorJSONShape(t *testing.T) {
	root := t.TempDir()
	findings, score := diagnose(root)
	payload := struct {
		Score    int       `json:"score"`
		Findings []finding `json:"findings"`
	}{score, findings}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if _, ok := decoded["score"]; !ok {
		t.Fatal("missing score")
	}
	if _, ok := decoded["findings"]; !ok {
		t.Fatal("missing findings")
	}
}

func TestApplySafeFixesCreatesOnlyMissingFiles(t *testing.T) {
	root := t.TempDir()
	if _, err := applySafeFixes(root); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"README.md", ".gitignore"} {
		if _, err := os.Stat(filepath.Join(root, name)); err != nil {
			t.Fatalf("expected %s: %v", name, err)
		}
	}
	before, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := applySafeFixes(root); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("safe fix overwrote an existing file")
	}
}

func TestDiagnoseDetectsCredentialPatternWithoutLeakingIt(t *testing.T) {
	root := t.TempDir()
	secret := "api_key=sk_live_" + strings.Repeat("A", 24)
	if err := os.WriteFile(filepath.Join(root, "config.env"), []byte(secret), 0o600); err != nil {
		t.Fatal(err)
	}
	findings, _ := diagnose(root)
	for _, f := range findings {
		if f.Check == "Secrets" {
			if f.Status != "FAIL" {
				t.Fatalf("expected secret failure, got %+v", f)
			}
			if strings.Contains(f.Detail, secret) || strings.Contains(f.Detail, "super_secret") {
				t.Fatal("secret content leaked into diagnostic detail")
			}
			return
		}
	}
	t.Fatal("secret finding not found")
}

func TestDiagnoseWarnsWhenEnvIsNotIgnored(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("bin/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".env"), []byte("SAFE_PLACEHOLDER=true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	findings, _ := diagnose(root)
	for _, f := range findings {
		if f.Check == "Environment hygiene" {
			if f.Status != "WARN" || !f.Fixable {
				t.Fatalf("expected fixable environment warning, got %+v", f)
			}
			return
		}
	}
	t.Fatal("environment hygiene finding not found")
}

func TestApplySafeFixesProtectsExistingEnv(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("bin/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".env"), []byte("SECRET=placeholder\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := applySafeFixes(root); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	if !gitignoreContainsEnv(string(data)) {
		t.Fatal("expected .env to be protected")
	}
	env, err := os.ReadFile(filepath.Join(root, ".env"))
	if err != nil {
		t.Fatal(err)
	}
	if string(env) != "SECRET=placeholder\n" {
		t.Fatal("safe fix modified .env contents")
	}
}

func TestFormatBytes(t *testing.T) {
	cases := map[int64]string{
		0:        "disabled",
		512:      "512 B",
		1024:     "1.0 KiB",
		10 << 20: "10.0 MiB",
	}
	for input, want := range cases {
		if got := formatBytes(input); got != want {
			t.Fatalf("formatBytes(%d) = %q, want %q", input, got, want)
		}
	}
}

func TestContainsLikelySecretIgnoresPlaceholders(t *testing.T) {
	if containsLikelySecret(`api_key=super_secret_value_123456789`) {
		t.Fatal("placeholder secret should not be reported")
	}
	if !containsLikelySecret("api_key=sk_live_" + strings.Repeat("A", 24)) {
		t.Fatal("realistic secret pattern should be reported")
	}
}

func TestContainsTodoCommentIgnoresStringLiterals(t *testing.T) {
	if containsTodoComment(`fmt.Println("TODO: example")`) {
		t.Fatal("TODO inside a string literal should not be reported")
	}
	if !containsTodoComment(`// TODO: finish this`) {
		t.Fatal("TODO comment should be reported")
	}
}

func TestIgnoredGeneratedFilesDoNotAffectRepositorySize(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("*.exe\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// A real Git repository is required for check-ignore semantics.
	cmd := exec.Command("git", "-C", root, "init")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("git unavailable: %v (%s)", err, out)
	}
	large := filepath.Join(root, "xef.exe")
	if err := os.WriteFile(large, make([]byte, 11*1024*1024), 0o644); err != nil {
		t.Fatal(err)
	}
	findings, _ := diagnoseWithOptions(root, 10*1024*1024)
	for _, f := range findings {
		if f.Check == "Repository size" && f.Status != "PASS" {
			t.Fatalf("ignored generated file affected repository size: %+v", f)
		}
	}
}

func TestGeneratedArtifactsAreIgnoredByDoctor(t *testing.T) {
	for _, name := range []string{"xef.exe", "coverage.out", "program.test"} {
		if !isGeneratedArtifact(name) {
			t.Fatalf("expected generated artifact %q to be ignored", name)
		}
	}
	if isGeneratedArtifact("important.bin") {
		t.Fatal("unexpectedly ignored arbitrary binary")
	}
}

func TestTestFixturesAreExcludedFromSecretFindings(t *testing.T) {
	if !isTestFixtureFile("internal/app/doctor/doctor_test.go") {
		t.Fatal("expected Go test file to be treated as a test fixture")
	}
	if isTestFixtureFile("internal/app/doctor/doctor.go") {
		t.Fatal("production Go source should not be treated as a test fixture")
	}
}
