# Validatord

A validator daemon written in Go that provides infrastructure for blockchain validators, including attestation, BLS signature operations, key management, data aggregation, state monitoring, and GitHub Models API integration for AI-powered capabilities.

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
│   ├── models/     # GitHub Models API integration
│   ├── payment/    # Payment handling
│   └── watcher/    # State monitoring
├── scripts/        # Utility scripts
└── tests/          # Test utilities
```

## Features

### GitHub Models API Integration

Validatord includes integration with GitHub Models API, allowing you to easily run LLMs for AI-powered capabilities. The `models` package provides a simple interface for calling GitHub's AI inference API.

**Setup:**

1. Create a GitHub personal access token with the `models` scope: [GitHub Settings > Personal Access Tokens](https://github.com/settings/tokens)

2. Use the Models API in your code:

```go
package main

import (
    "fmt"
    "log"
    "github.com/dreadwitdastacc-IFA/validatord/internal/models"
)

func main() {
    // Initialize with your GitHub token
    m, err := models.NewWithToken("YOUR_GITHUB_PAT")
    if err != nil {
        log.Fatal(err)
    }

    // Simple chat with default model
    response, err := m.Chat("What is the capital of France?")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response)

    // Use a specific model
    response, err = m.ChatWithModel("Explain recursion", "openai/gpt-4")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response)

    // Advanced usage with custom parameters
    req := &models.ChatRequest{
        Model: "openai/gpt-4o-mini",
        Messages: []models.Message{
            {Role: "system", Content: "You are a helpful assistant."},
            {Role: "user", Content: "Hello!"},
        },
        Temperature: 0.7,
        MaxTokens: 1000,
    }
    resp, err := m.CallModel(req)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(resp.Choices[0].Message.Content)
}
```

**Available Models:**

You can use any model from the [GitHub Models marketplace](https://github.com/marketplace/models), including:
- `openai/gpt-4o-mini` (default)
- `openai/gpt-4`
- `openai/gpt-4o`
- And many more from OpenAI, Meta, DeepSeek, and other providers

**API Documentation:**

See the [GitHub Models API documentation](https://docs.github.com/github-models) for more details.

## Payment Information

**Paystring:** ifawoleesubiyi$paystring.crypto.com

## License

See [LICENSE](LICENSE) for details.
