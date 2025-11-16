# Contributing to NextTrace Exporter

Thank you for your interest in contributing to NextTrace Exporter! This document provides guidelines and instructions for contributing.

## Code of Conduct

Be respectful and constructive in all interactions with the community.

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include:

- A clear and descriptive title
- Detailed steps to reproduce the issue
- Expected vs actual behavior
- Your environment (OS, Go version, nexttrace version)
- Relevant logs or error messages
- Configuration file (sanitized if needed)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- A clear and descriptive title
- Detailed description of the proposed functionality
- Use cases and benefits
- Possible implementation approach (optional)

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes**:
   - Write clear, readable code
   - Follow Go conventions and best practices
   - Add tests for new functionality
   - Update documentation as needed
3. **Ensure the test suite passes**:
   ```bash
   go test ./...
   go vet ./...
   ```
4. **Commit your changes**:
   - Use clear, descriptive commit messages
   - Follow conventional commit format when possible:
     - `feat:` for new features
     - `fix:` for bug fixes
     - `docs:` for documentation changes
     - `test:` for test additions/changes
     - `refactor:` for code refactoring
5. **Push to your fork** and submit a pull request

### Pull Request Guidelines

- Include tests for new functionality
- Update README.md with any new features or changes
- Update CHANGELOG.md following Keep a Changelog format
- Ensure CI checks pass
- Request review from maintainers

## Development Setup

### Prerequisites

- Go 1.21 or later
- nexttrace installed
- make (optional, but recommended)

### Setting Up Development Environment

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/nexttrace_exporter.git
cd nexttrace_exporter

# Install dependencies
go mod download

# Build
make build

# Run tests
make test

# Run locally
./nexttrace_exporter --config.file=examples/config.yml --log.level=debug
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...

# Generate coverage report
make test-coverage
```

### Code Style

- Follow standard Go formatting (`gofmt`)
- Use `go vet` to check for common issues
- Run `golangci-lint` if available
- Write clear comments for exported functions
- Keep functions focused and reasonably sized

### Project Structure

```
nexttrace_exporter/
├── main.go              # Entry point
├── config/              # Configuration handling
├── executor/            # NextTrace execution logic
├── collector/           # Prometheus metrics collection
├── parser/              # JSON parsing
└── examples/            # Example configs and integrations
```

## Documentation

- Update README.md for user-facing changes
- Update code comments for API changes
- Add examples for new features
- Update QUICKSTART.md if setup process changes

## Testing

### Writing Tests

- Write tests for all new functionality
- Use table-driven tests where appropriate
- Mock external dependencies (nexttrace execution)
- Test error paths, not just happy paths

### Test Coverage

Aim for reasonable test coverage:
- Critical paths: 90%+ coverage
- New features: Include tests in PR
- Bug fixes: Add regression tests

## Release Process

(For maintainers)

1. Update CHANGELOG.md
2. Update version in relevant files
3. Create and push a git tag
4. GitHub Actions will build and publish releases

## Questions?

Feel free to open an issue for questions or discussions!

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
