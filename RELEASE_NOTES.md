# Release v1.0.2

## Highlights

- New `xef doctor` project health diagnostics.
- Machine-readable `xef doctor --json` output.
- Safer environment JSON serialization.
- Stronger HTTP benchmark validation.
- Additional regression coverage for diagnostics.
- Safe `doctor --fix` repair mode for missing README/.gitignore files.
- `doctor --strict` for CI enforcement.
- Stronger credential-pattern detection without printing secret contents.
- Test fixture credentials no longer trigger doctor secret false positives.
- Generated executables and coverage artifacts no longer pollute content diagnostics; ignored files remain excluded from repository-size checks.
- Safer filesystem traversal with build-artifact and symlink exclusions.

## Verification

- Release version: `v1.0.2`.
- All Go sources parse successfully with the standard-library Go parser.
- All Go sources are `gofmt`-clean.
- Shell scripts pass `bash -n`.
- Repository was checked for generated artifacts and AI-attribution footprint.
- The release CI workflow runs `go vet`, the full test suite, race tests, and a production build on GitHub runners.
- The isolated packaging runner cannot complete dependency-backed `go test ./...` because external Go module downloads are unavailable; no local result is represented as a false test pass.

## Diagnostic hardening
- `xef doctor` now checks module hygiene, Go formatting, environment-file protection, and configurable repository size limits.
- `xef doctor --fix` can safely add `.env` to an existing `.gitignore` without touching `.env` contents.
