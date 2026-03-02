# Contributing to doxctl

Thank you for your interest in contributing to doxctl! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Development Environment Setup](#development-environment-setup)
- [How to Run Tests](#how-to-run-tests)
- [Code Style Guidelines](#code-style-guidelines)
- [Pull Request Process](#pull-request-process)
- [Commit Message Conventions](#commit-message-conventions)

## Development Environment Setup

### Prerequisites

- **Go 1.25.0+** (as specified in go.mod)
- **Make** (for build automation)
- **Git** (for version control)
- **Docker** (optional, for container builds)

### Initial Setup

1. **Fork and Clone the Repository**
   ```bash
   git clone https://github.com/YOUR_USERNAME/doxctl.git
   cd doxctl
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Build the Project**
   ```bash
   make build
   ```
   Or for a local installation:
   ```bash
   make install
   ```

4. **Verify Installation**
   ```bash
   doxctl -h
   ```

### Configuration

The tool requires a configuration file `.doxctl.yaml`. A sample configuration is provided:
```bash
cp .doxctl.yaml_SAMPLE ~/.doxctl.yaml
```

Edit `~/.doxctl.yaml` to customize settings for your environment.

## How to Run Tests

### Running All Tests

```bash
make test
```

This runs the complete test suite with race detection and generates a coverage report.

### Running Tests Manually

```bash
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
```

### Test Coverage

Coverage reports are generated in `coverage.txt`. You can view coverage details with:
```bash
go tool cover -html=coverage.txt
```

## Code Style Guidelines

### Go Code Standards

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` to format your code (this is enforced)
- Run `go vet` to catch common mistakes
- Run `make lint` to run golangci-lint before submitting PRs
- Use meaningful variable and function names

### Linting

The project uses [golangci-lint](https://golangci-lint.run/) for code quality enforcement. 

**Installation:**
```bash
# macOS
brew install golangci-lint

# Linux/WSL
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

**Running locally:**
```bash
make lint
```

The linting configuration is in `.golangci.yml` and includes checks for:
- Unchecked errors (errcheck)
- Code simplification (gosimple)
- Static analysis (staticcheck, govet)
- Ineffectual assignments (ineffassign)
- Unused code (unused)
- Code formatting (gofmt, goimports)
- Spelling errors (misspell)
- Code review issues (revive)
- Security vulnerabilities (gosec)

**Note:** Linting is automatically run in CI/CD on all pull requests.

### Code Organization

- Keep functions focused and single-purpose
- Add comments for exported functions and types
- Group related functionality in the same package
- Follow the existing project structure:
  - `cmd/` - CLI command implementations
  - `internal/` - Internal packages and helpers
  - `model_cmds/` - Model command scripts

### Documentation

- Document all exported functions, types, and packages
- Include examples in documentation where appropriate
- Update README.md if adding new features

## Pull Request Process

### Before Submitting

1. **Create a Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Your Changes**
   - Write clear, concise code
   - Add or update tests as needed
   - Update documentation

3. **Test Your Changes**
   ```bash
   make test
   ```

4. **Lint Your Code**
   ```bash
   make lint
   ```

5. **Commit Your Changes**
   - Follow the commit message conventions (see below)
   - Keep commits atomic and focused

### Submitting a Pull Request

1. **Push to Your Fork**
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Open a Pull Request**
   - Provide a clear title and description
   - Reference any related issues
   - Describe what changes were made and why
   - Include testing steps if applicable

3. **Code Review**
   - Address reviewer feedback promptly
   - Keep discussions focused and constructive
   - Update your PR based on feedback

4. **Merge Requirements**
   - All tests must pass
   - Code must be reviewed and approved
   - No merge conflicts with main branch
   - Follows project coding standards

## Commit Message Conventions

We follow a conventional commit message format to maintain a clear project history.

### Format

```
<type>: <subject>

<body>

<footer>
```

### Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation changes
- **style**: Code style changes (formatting, missing semicolons, etc.)
- **refactor**: Code refactoring (no functional changes)
- **perf**: Performance improvements
- **test**: Adding or updating tests
- **chore**: Maintenance tasks (dependency updates, build changes, etc.)

### Examples

```
feat: add support for custom DNS timeout configuration

Add a new flag --dns-timeout to allow users to configure the DNS
query timeout value. This is useful for slow network environments.

Closes #123
```

```
fix: correct VPN route counting logic

The minimum route count check was incorrectly comparing against
a hardcoded value. Now uses the configured minVpnRoutes setting.
```

```
docs: update installation instructions for macOS

Added Homebrew installation steps and updated prerequisites.
```

### Best Practices

- Keep the subject line under 50 characters
- Use the imperative mood ("add" not "added")
- Separate subject from body with a blank line
- Wrap body at 72 characters
- Reference issues and PRs in the footer

## Getting Help

- **Issues**: Open an issue for bug reports or feature requests
- **Discussions**: Use GitHub Discussions for questions and ideas
- **Security**: Report security vulnerabilities via SECURITY.md

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on what is best for the project and community
- Show empathy towards other community members

Thank you for contributing to doxctl! 🎉
