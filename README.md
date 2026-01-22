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
│   ├── llm/        # GitHub Models API client for LLMs
│   ├── milestone/  # Milestone tracking
│   ├── payment/    # Payment handling
│   └── watcher/    # State monitoring
├── scripts/        # Utility scripts
└── tests/          # Test utilities
```

## Features

### GitHub Models API Integration

Validatord includes built-in support for calling GitHub Models APIs to easily run large language models (LLMs). This feature enables AI-powered functionality through GitHub's model inference service.

#### Usage Example

```go
import (
    "fmt"
    "log"
    
    "github.com/dreadwitdastacc-IFA/validatord/internal/app"
)

func main() {
    // Initialize the application
    application, err := app.New(app.DefaultPaystring)
    if err != nil {
        log.Fatal(err)
    }
    
    // Configure GitHub token for LLM API access
    err = application.LLM.SetToken("your-github-token")
    if err != nil {
        log.Fatal(err)
    }
    
    // Simple completion - ask a question
    response, err := application.LLM.SimpleCompletion("What is blockchain validation?")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Response:", response)
    
    // Chat completion with system context
    response, err = application.LLM.ChatCompletion(
        "You are a helpful blockchain expert",
        "Explain BLS signatures",
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Response:", response)
}
```

#### Configuration

The LLM client can be configured with custom settings:

```go
import "github.com/dreadwitdastacc-IFA/validatord/internal/llm"

// Create client with custom configuration
client, err := llm.NewWithConfig(llm.Config{
    APIEndpoint: "https://models.inference.ai.azure.com",
    Token:       "your-github-token",
    Model:       "gpt-4o",
    Timeout:     30 * time.Second,
})
```

#### Supported Models

The client supports various models available through GitHub Models API. Default model is `gpt-4o`. You can change the model:

```go
err := application.LLM.SetModel("gpt-4o-mini")
```

#### Getting a GitHub Token

To use the GitHub Models API:

1. Visit [GitHub Models](https://github.com/marketplace/models)
2. Generate a personal access token with appropriate permissions
3. Configure the token in your application

For more information, see the [GitHub Models Quickstart Guide](https://docs.github.com/github-models/quickstart).

## Payment Information

**Paystring:** ifawoleesubiyi$paystring.crypto.com

## License

See [LICENSE](LICENSE) for details.
