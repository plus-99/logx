# Contributing to LogX

Thank you for your interest in contributing to LogX! This document provides guidelines and information for contributors.

## Development Setup

### Prerequisites

- Go 1.20 or higher
- Git

### Getting Started

1. **Fork the repository**
   ```bash
   git clone https://github.com/your-username/logx.git
   cd logx
   ```

2. **Install dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Run tests**
   ```bash
   go test ./...
   ```

4. **Run benchmarks**
   ```bash
   go test -bench=. tests/bench_test.go
   ```

5. **Run examples**
   ```bash
   go run examples/main.go
   go run examples/redaction_demo.go
   ```

## Project Structure

```
/
├── examples/           # Example applications and demos
├── internal/          # Internal implementation code
│   ├── encoding/      # Log formatter implementations
│   ├── hooks/         # Hook implementations for remote logging
│   ├── redaction/     # Redaction and sensitive data protection
│   ├── context.go     # Context integration
│   ├── encoder.go     # Public encoder interfaces
│   ├── hooks.go       # Public hook interfaces
│   ├── logger.go      # Core logging implementation
│   └── redaction.go   # Public redaction interfaces
├── tests/             # Test files and benchmarks
├── logx.go            # Public API (re-exports from internal)
├── README.md          # Documentation
└── CONTRIBUTING.md    # This file
```

## Building

### Build Examples
```bash
go build -o bin/demo examples/main.go
go build -o bin/redaction-demo examples/redaction_demo.go
```

### Cross-compilation
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o bin/logx-linux-amd64 examples/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/logx-windows-amd64.exe examples/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o bin/logx-darwin-amd64 examples/main.go
```

## Testing

### Run All Tests
```bash
go test ./...
```

### Run Tests with Coverage
```bash
go test -cover ./...
```

### Run Benchmarks
```bash
go test -bench=. tests/bench_test.go
```

### Performance Comparison
The benchmarks compare LogX against popular logging libraries:
- Logrus
- Zerolog

Run benchmarks to verify LogX maintains its performance advantages.

## Code Guidelines

### Code Style
- Follow standard Go conventions and `gofmt` formatting
- Use meaningful variable and function names
- Add comments for public APIs and complex logic
- Keep functions focused and concise

### Package Organization
- **Public API**: Export only necessary types and functions in `logx.go`
- **Internal packages**: Place implementation details in `internal/`
- **Examples**: Add new examples to `examples/` directory
- **Tests**: Place tests in `tests/` directory

### Adding New Features

1. **Hooks**: Add new remote logging hooks in `internal/hooks/`
2. **Encoders**: Add new formatters in `internal/encoding/`
3. **Redaction**: Add new patterns to `internal/redaction/patterns.go`

## Contribution Process

### 1. Create an Issue
Before starting work, create an issue describing:
- The problem you're solving
- Your proposed solution
- Any breaking changes

### 2. Create a Branch
```bash
git checkout -b feature/your-feature-name
```

### 3. Make Changes
- Follow the code guidelines above
- Add tests for new functionality
- Update documentation as needed

### 4. Test Your Changes
```bash
# Run all tests
go test ./...

# Run benchmarks
go test -bench=. tests/bench_test.go

# Test examples
go run examples/main.go
go run examples/redaction_demo.go
```

### 5. Update Documentation
- Update README.md for new features
- Add examples demonstrating new functionality
- Update this CONTRIBUTING.md if needed

### 6. Submit a Pull Request
- Provide a clear description of changes
- Reference any related issues
- Ensure all tests pass

## Adding New Redaction Patterns

To add new sensitive data patterns:

1. Add the regex pattern to `internal/redaction/patterns.go`:
   ```go
   var BuiltInPatterns = []string{
       // ... existing patterns ...
       `your-new-pattern-here`,
   }
   ```

2. Test the pattern in `examples/redaction_demo.go`

3. Update the README with the new pattern type

## Performance Considerations

LogX prioritizes performance:
- Minimize memory allocations
- Use object pooling where appropriate
- Benchmark new features against existing libraries
- Profile code for bottlenecks

## Questions?

- Open an issue for bugs or feature requests
- Check existing issues before creating new ones
- Be respectful and constructive in discussions

Thank you for contributing to LogX!