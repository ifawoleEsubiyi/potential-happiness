# GitHub Copilot Agent for Validatord

This agent helps contributors work with the Validatord project by providing testing guidance, safe refactor suggestions, and repository policy reminders.

## Project Context

- **Tech Stack**: Go 1.24.10
- **Internal Packages**: aggregator, attest, bls, farming, keystore, payment, watcher
- **Testing Pattern**: Table-driven tests (see `*_test.go` files for examples)
- **Build Commands**: `go build -v .`, `go test ./...`, `gofmt -s -w .`
- **Detailed Guidelines**: See `.github/copilot-instructions.md` for comprehensive project conventions

## Usage Examples

- "How do I fix the failing tests in internal/keystore?"
- "Suggest a safe refactor for function X to return errors instead of panicking."
- "What security checks should I run before adding a new dependency?"
- "How do I write a table-driven test for the payment package?"
- "What are the build and test commands for this project?"

## Note

This agent does not store or expose secrets. It will recommend secure practices only.
