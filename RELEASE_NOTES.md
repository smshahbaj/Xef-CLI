# Release v1.0.1

Summary:
- Fix: resolved filesystem `WalkDir` deadlock that prevented `file clean` from removing files on Windows.
- Fix: prevented Cobra command flag leakage in `file clean` and made tests deterministic.
- Chore: bumped CLI version to `v1.0.1`.

Artifacts produced locally:
- Local git tag: `v1.0.1` (created locally)

Sanity checks performed:
- `gofmt` applied to source
- `go vet ./...` - no issues
- `go test ./...` - all packages pass

Publish steps (requires GitHub remote and credentials):
1. Push commits and tag:

```bash
git push origin HEAD
git push origin v1.0.1
```

2. Create GitHub release using `gh` (or use GitHub web UI):

```bash
gh release create v1.0.1 --title "v1.0.1" --notes-file RELEASE_NOTES.md
```

Optional: run GoReleaser to build binaries and publish release artifacts:

```bash
goreleaser release --rm-dist
```

If you want, I can push and create the release for you (I will need your permission and access to your Git credentials or a remote URL). Otherwise run the commands above locally or allow me to continue and I will attempt to push if you provide the remote URL and confirm.
