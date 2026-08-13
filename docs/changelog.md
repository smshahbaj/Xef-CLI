## [1.0.2] - 2026-08-13
- Hardened doctor diagnostics and safe repair mode
- Added tracked `.env` detection
- Redacted environment values by default
- Made JSON diffs deterministic
- Made HTTP downloads atomic
- Hardened project scaffolding and added parser regression coverage

- Added deterministic project health diagnostics with `xef doctor`.
- Added JSON output for automation and CI integrations.
- Hardened HTTP benchmark flag validation.
- Hardened environment JSON serialization.
- Improved command-level error handling and project diagnostics.

# Changelog

All notable changes to this project will be documented in this file.

## [1.0.0] - 2024-01-15

### Added
- Initial release
- File management: organize, stats, duplicates, clean
- JSON processing: format, validate, diff
- Cryptographic tools: sha256, sha512, bcrypt, uuid, base64, password
- HTTP client: get, download, benchmark
- System monitoring: cpu, memory, disk
- Git utilities: stats, branches
- Development tools: project create, env
- Cross-platform support (Linux, macOS, Windows)
- Structured logging with Zap
- Clean Architecture implementation
- Comprehensive test suite
- CI/CD pipeline with GitHub Actions
- GoReleaser configuration

## [1.0.1] - 2026-08-06

### Fixed
- Resolve filesystem WalkDir deadlock preventing file removal on Windows
- Fix command flag state leakage in `file clean` and improve tests

### Chore
- Bumped CLI version to `v1.0.1` for release
