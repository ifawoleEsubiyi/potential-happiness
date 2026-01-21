# Validatord

A validator daemon written in Go that provides infrastructure for blockchain validators, including attestation, BLS signature operations, key management, data aggregation, and state monitoring.

## Quickstart

### Prerequisites

- Go 1.24 or later

### Build

```bash
# Clone the repository
git clone https://github.com/dreadwitdastacc-IFA/fluffy-umbrella.git
cd fluffy-umbrella

# Build the project
make build

# Or directly with Go
go build -v .
```

### Test

```bash
# Run all tests
make test

# Run tests with race detector
make test-race

# Run all checks (format, vet, lint, test)
make check
```

### Run

```bash
./validatord
```

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines, including:

- Setting up your development environment
- Running tests and linters
- Commit message conventions
- Pull request process

## For Students

If you're a student looking to leverage GitHub Education benefits and Copilot in your workflow, check out our comprehensive guide:

- [GitHub Education Benefits & Copilot Workflow](GITHUB_EDUCATION_WORKFLOW.md)

This guide covers project setup, Copilot usage, collaboration best practices, and how to make the most of the GitHub Student Developer Pack.

## Security

Please see our [Security Policy](SECURITY.md) for:

- How to report vulnerabilities
- Security best practices
- Supported versions

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md).

## Project Structure

```
├── cmd/            # Command-line applications
├── internal/       # Private packages
│   ├── aggregator/ # Data aggregation
│   ├── attest/     # Attestation logic
│   ├── bls/        # BLS cryptography
│   ├── farming/    # Farming operations
│   ├── keystore/   # Key management
│   ├── milestone/  # Milestone tracking
│   ├── payment/    # Payment handling
│   └── watcher/    # State monitoring
├── scripts/        # Utility scripts
└── tests/          # Test utilities
```

## Payment Information

**Paystring:** ifawoleesubiyi$paystring.crypto.com

## License

See [LICENSE](LICENSE) for details.
