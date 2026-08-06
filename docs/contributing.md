# Contributing to XefCLI

Thank you for your interest in contributing! This document provides guidelines for contributing to XefCLI.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/xefcli.git`
3. Install dependencies: `make deps`
4. Run tests: `make test`

## Development Workflow

1. Create a feature branch: `git checkout -b feature/my-feature`
2. Make your changes
3. Run linting: `make lint`
4. Run tests: `make test`
5. Commit with clear messages
6. Push and create a Pull Request

## Code Standards

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Every exported function must be documented
- No duplicated logic
- No magic numbers
- Small, focused functions
- Meaningful variable names
- No global mutable state

## Testing Requirements

- Minimum 90% code coverage
- Unit tests for all packages
- Integration tests for commands
- Benchmarks for performance-critical code
- Race detection enabled

## Commit Message Format

```
type(scope): subject

body

footer
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

## Pull Request Process

1. Update documentation for any API changes
2. Ensure CI passes
3. Request review from maintainers
4. Address feedback promptly
5. Squash commits if requested
