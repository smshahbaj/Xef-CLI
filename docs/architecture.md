# XefCLI Architecture

## Overview
XefCLI follows Clean Architecture principles with clear separation between:
- **Core**: Business logic, interfaces, errors, config, logger
- **Infrastructure**: Concrete implementations (filesystem, HTTP, crypto, system)
- **App**: CLI commands organized by domain
- **Pkg**: Shared utilities and TUI helpers

## Dependency Flow
```
cmd/xefcli -> internal/app -> internal/core/interfaces -> internal/infrastructure
```

## Key Design Decisions
1. **Interfaces First**: All external dependencies are abstracted behind interfaces
2. **Dependency Injection**: App dependencies are wired in `app.New()`
3. **No Global State**: All state is injected through constructors
4. **Context Propagation**: All I/O operations accept `context.Context`
5. **Structured Logging**: Zap-based logger with field-based API
6. **Error Wrapping**: Every error includes context and exit codes
7. **Path Traversal Protection**: File system operations sanitize paths
8. **Worker Pools**: Concurrent operations use bounded goroutine pools
9. **Graceful Degradation**: Operations continue on partial failures

## Folder Structure
```
cmd/xefcli/          # Entry point
internal/
  app/               # CLI commands (file, json, crypto, http, git, system, dev)
  core/
    config/          # Configuration management
    errors/          # Custom error types and exit codes
    interfaces/      # Core abstractions
    logger/          # Structured logging
  infrastructure/
    crypto/          # Hashing implementations
    filesystem/      # OS file system with security
    network/         # HTTP client
    systeminfo/      # System info provider
  pkg/
    tui/             # Terminal UI helpers
    utils/           # Common utilities
pkg/xeflib/          # Public library API
configs/             # Configuration files
docs/                # Documentation
examples/            # Usage examples
scripts/             # Build scripts
assets/              # Screenshots and media
```
