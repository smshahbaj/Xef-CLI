# Changelog

## v1.0.2 — Final release

### Final hardening pass
- Restored and hardened `scan`, `secret scan`, and `report` commands with deterministic output and safe generated-file/ignore handling.
- Hardened repeated HTML report generation for Windows file replacement semantics.
- Added regression coverage for secret scanning, maintenance detection, ignored artifacts, and report generation.
- Fixed project doctor false positives from test placeholders and detector source strings.
- Maintenance scanning now reports TODO/FIXME markers only when they appear in comment lines.
- Repository-size diagnostics now respect Git ignore rules, preventing generated binaries from lowering project health.
- Expanded doctor regression coverage for secret detection, maintenance scanning, and ignored generated files.
- Test-only credential fixtures are excluded from secret findings to prevent false positives while production configuration remains scanned.
- Generated executables and coverage artifacts are excluded from content diagnostics while still respecting repository-size checks unless Git ignores them.
- Hardened project doctor secret detection with credential-pattern matching without exposing secret values
- Added safe `doctor --fix` repair mode and CI-friendly `doctor --strict` behavior
- Expanded diagnostics to skip build artifacts and symlinks safely
- Added regression tests for safe repair and secret detection
- Added structured CLI release wiring and version injection
- Hardened cross-platform behavior and command scaffolding
- Improved build and test readiness for release
- Hardened generated Go scaffolds and added parser-level regression coverage
- Environment output now redacts values by default; use `--show-values` explicitly
- HTTP downloads now finalize atomically to avoid leaving partial destination files
- JSON diff output is deterministic with sorted object keys
- Doctor now detects tracked `.env` files and reports safe-fix actions in JSON mode

### Diagnostic hardening
- Added configurable `doctor --max-file-size` thresholds.
- Added Go module (`go.sum`) health checks and gofmt cleanliness checks.
- Added `.env` Git hygiene diagnostics and a safe repair path.
- Expanded doctor regression coverage for environment protection and size formatting.
