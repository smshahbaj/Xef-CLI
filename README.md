# 🚀 XefCLI

**The Ultimate Open-Source Developer Toolkit**

[![CI](https://github.com/smshahbaj/Xef-CLI/actions/workflows/ci.yml/badge.svg)](https://github.com/smshahbaj/Xef-CLI/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/smshahbaj/Xef-CLI)](https://goreportcard.com/report/github.com/smshahbaj/Xef-CLI)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A production-grade, cross-platform CLI toolkit built in Go for developers who demand quality, performance, and reliability.

## ✨ Features

### 📁 File Tools
- `organize` - Organize files by extension or date
- `stats` - Show detailed file/directory statistics
- `duplicates` - Find duplicate files with parallel hashing
- `clean` - Remove temporary and unwanted files

### 📊 JSON Tools
- `format` - Pretty-print or compact JSON
- `validate` - Validate JSON syntax
- `diff` - Compare two JSON files recursively

### 🔐 Crypto Tools
- `sha256` / `sha512` - Compute file/string hashes
- `bcrypt` - Hash and verify passwords
- `uuid` - Generate UUIDs
- `base64` - Encode/decode Base64
- `password` - Generate secure random passwords

### 🌐 HTTP Tools
- `get` - Perform HTTP GET requests
- `download` - Download files safely with atomic finalization
- `benchmark` - Load test HTTP endpoints

### 💻 System Tools
- `cpu` - CPU information and usage
- `memory` - Memory usage statistics
- `disk` - Disk usage by mount point

### 🔧 Git Tools
- `stats` - Repository statistics and contributors
- `branches` - List branches with commit info

### 🛠️ Dev Tools
- `project` - Scaffold Go/Python projects
- `env` - Display environment variables

### 🩺 Project Doctor
- `doctor` - Diagnose project health locally
- `doctor --json` - Emit CI-friendly diagnostics
- `doctor --fix` - Apply only safe, non-destructive missing-file fixes
- `doctor --strict` - Fail when warnings or errors are detected

### 🔎 Security & Project Scanning
- `scan` - Scan project files for security and maintenance findings
- `scan --json` - Emit machine-readable scan results
- `secret scan` - Scan specifically for exposed credential patterns
- `secret scan --json` - Emit machine-readable secret-scan results
- `report` - Generate an HTML project health report

## 🚀 Installation

### From Source
```bash
go install github.com/smshahbaj/Xef-CLI/cmd/xefcli@latest
```

### Pre-built Binaries
Download from [Releases](https://github.com/smshahbaj/Xef-CLI/releases)

## 📖 Quick Start

```bash
# File statistics
xef file stats ./my-project

# Format JSON
xef json format data.json --indent "  "

# Validate JSON
xef json validate data.json

# Generate password
xef crypto password --length 32

# Benchmark API
xef http benchmark https://api.example.com -n 1000 -c 50

# System info
xef system cpu
xef system memory
xef system disk

# Git stats
xef git stats

# Create project
xef dev project myapp --lang go

# Diagnose and safely repair project hygiene
xef doctor
xef doctor --fix
xef doctor --strict --json
xef scan
xef secret scan
xef report
```

## 🏗️ Architecture

XefCLI follows **Clean Architecture** principles:

- **Interfaces First** - All dependencies abstracted
- **Dependency Injection** - No global state
- **Context Propagation** - Graceful cancellation
- **Structured Logging** - Zap-based with fields
- **Error Wrapping** - Contextual errors with exit codes

See [docs/architecture.md](docs/architecture.md) for details.

## 🧪 Testing

```bash
# Run all tests
make test

# Run with race detection
make test-race

# Run linter
make lint
```

## 🔎 Release-quality checks

```bash
# Project health
xef doctor --strict

# Machine-readable health report
xef doctor --json

# Format and static analysis
make fmt
make vet
make lint
```

## 📦 Building

```bash
# Build binary
make build

# Build for all platforms
make snapshot

# Release
make release
```

## 🤝 Contributing

See [docs/contributing.md](docs/contributing.md) for guidelines.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Viper](https://github.com/spf13/viper) - Configuration
- [Zap](https://github.com/uber-go/zap) - Logging
- [gopsutil](https://github.com/shirou/gopsutil) - System information
