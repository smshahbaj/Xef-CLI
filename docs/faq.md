# Frequently Asked Questions

## General

### What is XefCLI?
XefCLI is a production-grade, cross-platform CLI toolkit for developers built in Go.

### Is it free?
Yes, XefCLI is open source under the MIT license.

## Installation

### How do I install XefCLI?
```bash
go install github.com/smshahbaj/Xef-CLI/cmd/xefcli@latest
```

### What Go version is required?
Go 1.22 or later.

## Usage

### How do I organize files by date?
```bash
xef file organize ./downloads --by date
```

### How do I benchmark an API?
```bash
xef http benchmark https://api.example.com -n 1000 -c 50
```

### How do I generate a secure password?
```bash
xef crypto password --length 32
```

## Troubleshooting

### "not a git repository" error
Make sure you run git commands from inside a Git repository.

### Permission denied errors
Check file permissions and ensure you have write access to the target directory.
