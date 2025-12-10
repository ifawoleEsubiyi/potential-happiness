# Contributing to Validatord

Thank you for your interest in contributing to Validatord! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Running Tests](#running-tests)
- [Linting](#linting)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Pull Request Process](#pull-request-process)

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## Getting Started

### Prerequisites

- Go 1.24 or later
- Git

### Setup

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/fluffy-umbrella.git
   cd fluffy-umbrella
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/dreadwitdastacc-IFA/fluffy-umbrella.git
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```
5. Verify setup:
   ```bash
   make build
   make test
   ```

## Development Workflow

1. Create a new branch from `main`:
   ```bash
   git checkout main
   git pull upstream main
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following the coding conventions

3. Run tests and linters:
   ```bash
   make test
   make lint
   ```

4. Commit your changes (see [Commit Message Guidelines](#commit-message-guidelines))

5. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

6. Open a Pull Request

## Running Tests

```bash
# Run all tests
make test

# Run tests with race detector
go test -race ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out
```

## Linting

We use `golangci-lint` for linting. Install it and run:

```bash
# Install golangci-lint (if not already installed)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
make lint

# Or directly
golangci-lint run

# Auto-format code
make fmt
```

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Keep functions small and focused
- Write meaningful variable and function names
- Add comments for exported types and functions
- Handle errors explicitly

## Commit Message Guidelines

We follow conventional commit messages:

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, semicolons, etc.)
- `refactor`: Code refactoring without feature changes
- `test`: Adding or updating tests
- `chore`: Maintenance tasks (dependencies, CI, etc.)
- `perf`: Performance improvements
- `ci`: CI/CD changes

### Examples

```
feat(payment): add paystring validation

fix(bls): correct signature verification

docs: update contributing guidelines

chore(deps): update Go dependencies
```

### Guidelines

- Use imperative mood ("Add feature" not "Added feature")
- First line should be 72 characters or less
- Reference issues when applicable (e.g., "Fixes #123")

## Pull Request Process

### Before Submitting

- [ ] Tests pass: `make test`
- [ ] Code is formatted: `make fmt`
- [ ] Linter passes: `make lint`
- [ ] No security vulnerabilities: `make vet`
- [ ] Documentation is updated if needed

### PR Checklist

1. **Title**: Use a clear, descriptive title following commit conventions
2. **Description**: Fill out the PR template completely
3. **Related Issues**: Link any related issues
4. **Tests**: Include tests for new functionality
5. **Documentation**: Update docs if needed

### Review Process

1. PRs require at least one approval
2. All CI checks must pass
3. Requested changes should be addressed
4. Squash and merge is preferred for clean history

### After Merge

- Delete your branch
- Update your local main:
  ```bash
  git checkout main
  git pull upstream main
  ```

## Security

If you discover a security vulnerability, please follow our [Security Policy](SECURITY.md).

**Do not** open a public issue for security vulnerabilities.

## Questions?

Feel free to open an issue for questions or discussions.

Thank you for contributing!
