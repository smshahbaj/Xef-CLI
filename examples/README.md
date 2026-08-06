# XefCLI Examples

## File Organization

```bash
# Organize downloads by file extension
xef file organize ~/Downloads --by extension

# Organize by month
xef file organize ~/Downloads --by date

# Dry run first
xef file organize ~/Downloads --by extension --dry-run
```

## JSON Processing

```bash
# Format JSON file
xef json format data.json

# Compact JSON
xef json format data.json --compact

# Validate from stdin
cat data.json | xef json validate

# Diff two files
xef json diff file1.json file2.json
```

## Cryptography

```bash
# Hash a file
xef crypto sha256 document.pdf

# Hash a string
xef crypto sha512 "hello world"

# Generate password
xef crypto password --length 32 --no-special

# Bcrypt password
xef crypto bcrypt mypassword --cost 12

# Verify bcrypt
xef crypto bcrypt mypassword --compare '$2a$12$...'

# Generate UUIDs
xef crypto uuid --count 5

# Base64 encode/decode
xef crypto base64 "hello"
xef crypto base64 "aGVsbG8=" --decode
```

## HTTP Operations

```bash
# GET request
xef http get https://api.github.com/users/octocat

# With headers
xef http get https://api.example.com -H "Authorization:Bearer token"

# Download file
xef http download https://example.com/file.zip -o output.zip

# Benchmark API
xef http benchmark https://api.example.com -n 1000 -c 50 --timeout 10s
```

## System Monitoring

```bash
xef system cpu
xef system memory
xef system disk
```

## Git Analysis

```bash
xef git stats
xef git branches
```

## Development

```bash
# Create Go project
xef dev project create myservice --lang go

# Create Python project
xef dev project create myapp --lang python

# Show env vars as JSON
xef dev env --format json
```
