---
description: Instructions for GitHub Copilot coding agent in the Validatord repository
applyTo: "**/*"
---

# Copilot Instructions for Validatord

> **Version:** 1.2 | **Last Updated:** 2025-11-24

## Project Overview

Validatord is a validator daemon written in Go that provides infrastructure for blockchain validators. The project includes components for attestation, BLS signature operations, key management, data aggregation, and state monitoring.

**Tech Stack:**
- Language: Go 1.24.10
- Architecture: Modular internal packages
- Key Components: attestation, BLS cryptography, keystore, aggregator, watcher, payment, farming

**Repository Context:**
- **Repository Name**: `fluffy-umbrella` (GitHub repository)
- **Module Name**: `github.com/dreadwitdastacc-IFA/validatord` (Go module)
- **Primary Purpose**: Blockchain validator infrastructure with secure payment handling
- **Code Size**: ~350 lines of production Go code (~1100 lines including tests) across 7 internal packages
- **Test Coverage**: All packages have test files with table-driven tests
- **Development Stage**: Active development with established patterns

## Important: What NOT to Do

**Security & Safety Boundaries:**
- ❌ NEVER commit secrets, private keys, or credentials to the repository
- ❌ DO NOT modify the payment paystring in README.md or hardcoded values without explicit instruction
- ❌ NEVER remove or weaken cryptographic security measures
- ❌ DO NOT push directly to main branch - always use pull requests
- ❌ NEVER ignore test failures - all tests must pass before merging
- ❌ DO NOT add external dependencies without security review (especially for cryptographic operations)
- ❌ DO NOT remove or modify existing error handling without replacement
- ❌ NEVER use `panic()` in library code; reserve for truly unrecoverable situations
- ❌ DO NOT use `_` to ignore errors without explicit justification in comments

## Coding Guidelines

### General Principles
- Write clean, idiomatic Go code following the [Effective Go](https://go.dev/doc/effective_go) guidelines
- Keep functions small and focused with a single responsibility
- Use meaningful variable and function names that clearly express intent
- Prioritize code readability and maintainability over cleverness

**Example of Good Code Structure:**
```go
// Package example demonstrates proper Go code structure
package example

// Validator represents a blockchain validator
type Validator struct {
    id      string
    pubKey  []byte
    isActive bool
}

// NewValidator creates a new validator instance with validation
func NewValidator(id string, pubKey []byte) (*Validator, error) {
    if id == "" {
        return nil, fmt.Errorf("validator ID cannot be empty")
    }
    if len(pubKey) == 0 {
        return nil, fmt.Errorf("public key cannot be empty")
    }
    return &Validator{
        id:      id,
        pubKey:  pubKey,
        isActive: true,
    }, nil
}
```

### Go Style Conventions
- Use `gofmt` for code formatting (all code should be formatted before commit)
- Follow standard Go naming conventions:
  - Use camelCase for unexported names
  - Use PascalCase for exported names
  - Use ALL_CAPS only for constants that are true constants
- Package names should be short, lowercase, and singular (no underscores)
- Interface names should describe behavior (e.g., `Reader`, `Writer`, `Validator`)

### Code Organization
- Keep related code in the same package under `internal/`
- Each package should have a clear, single purpose
- Export only what's necessary; prefer unexported implementation details
- Place tests in `*_test.go` files alongside the code they test

### Error Handling
- Always handle errors explicitly; never ignore returned errors
- Use descriptive error messages that include context
- Wrap errors with additional context using `fmt.Errorf("context: %w", err)`
- Return errors rather than panicking, except in truly unrecoverable situations

**Example of Proper Error Handling:**
```go
func processData(data []byte) error {
    if len(data) == 0 {
        return fmt.Errorf("data cannot be empty")
    }
    
    result, err := validateData(data)
    if err != nil {
        return fmt.Errorf("failed to validate data: %w", err)
    }
    
    if err := storeResult(result); err != nil {
        return fmt.Errorf("failed to store result: %w", err)
    }
    
    return nil
}
```

### Comments and Documentation
- Write package-level documentation for each package in `doc.go` or at the top of the main file
- Document all exported types, functions, methods, and constants
- Start comments with the name of the thing being documented
- Keep comments concise but informative
- Use TODO comments to mark incomplete implementations: `// TODO: Add implementation`

### Security Considerations
- Never hardcode sensitive data (keys, passwords, secrets)
- Validate all inputs, especially from external sources
- Use cryptographically secure random number generation for security-critical operations
- Handle cryptographic keys securely (proper zeroing, safe storage)
- Follow the principle of least privilege in access control
- Sanitize all user inputs to prevent injection attacks
- Use constant-time comparison for sensitive data (e.g., `subtle.ConstantTimeCompare`)
- **Avoid timing attacks**: Don't use early returns in security-critical comparisons
- **Log security events**: Track authentication attempts, validation failures, and suspicious activities
- **Rate limiting**: Implement rate limiting for external-facing APIs
- **Input size limits**: Always enforce maximum input sizes to prevent DoS attacks
- **Fail securely**: Default to deny, fail closed rather than open
- **Defense in depth**: Layer security measures, don't rely on single point of validation

**Example: Secure Input Validation**
```go
func sanitizePaystring(input string) (string, error) {
    // Trim whitespace
    input = strings.TrimSpace(input)
    
    // Validate length to prevent DoS
    if len(input) > 256 {
        return "", fmt.Errorf("paystring exceeds maximum length")
    }
    
    // Validate format before processing
    if err := ValidatePaystring(input); err != nil {
        return "", fmt.Errorf("invalid paystring format: %w", err)
    }
    
    return input, nil
}
```

**Example: Constant-Time Comparison for Secrets**
```go
import "crypto/subtle"

func verifyToken(provided, expected []byte) bool {
    // Use constant-time comparison to prevent timing attacks
    // Note: subtle.ConstantTimeCompare handles length differences securely
    return subtle.ConstantTimeCompare(provided, expected) == 1
}
```

**Example: Secure Key Handling**
```go
func processKey(key []byte) error {
    // Use defer to ensure key is zeroed after use
    defer func() {
        // Zero out sensitive data
        for i := range key {
            key[i] = 0
        }
    }()
    
    // Process the key...
    return nil
}
```

**Security Checklist for Code Review:**
- [ ] No hardcoded credentials or API keys
- [ ] All user inputs are validated and sanitized
- [ ] Cryptographic operations use approved libraries
- [ ] Sensitive data is not logged or exposed in errors
- [ ] Rate limiting is implemented where needed
- [ ] Error messages don't leak sensitive information
- [ ] Proper authentication and authorization checks
- [ ] Secure random number generation for cryptographic purposes

## Testing Guidelines

### Test Organization
- Write table-driven tests where appropriate
- Test files should be named `*_test.go`
- Use meaningful test function names: `TestFunctionName_Scenario`
- Group related test cases using subtests with `t.Run()`

**Example of Table-Driven Test:**
```go
func TestValidatePaystring(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
        errMsg  string
    }{
        {
            name:    "valid paystring",
            input:   "user$domain.com",
            wantErr: false,
        },
        {
            name:    "empty paystring",
            input:   "",
            wantErr: true,
            errMsg:  "paystring cannot be empty",
        },
        {
            name:    "missing dollar sign",
            input:   "userdomain.com",
            wantErr: true,
            errMsg:  "paystring must contain exactly one '$' separator",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidatePaystring(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidatePaystring() error = %v, wantErr %v", err, tt.wantErr)
            }
            if tt.wantErr && err != nil && err.Error() != tt.errMsg {
                t.Errorf("ValidatePaystring() error = %v, want %v", err.Error(), tt.errMsg)
            }
        })
    }
}
```

### Test Quality
- Write tests that are independent and can run in any order
- Test both success cases and error conditions
- Use test helpers to reduce duplication
- Mock external dependencies appropriately
- Aim for meaningful test coverage (>80% for critical packages)
- Test edge cases and boundary conditions
- Use descriptive assertion messages for easier debugging

### Running Tests
```bash
go test ./...                    # Run all tests
go test -v ./...                 # Run with verbose output
go test -race ./...              # Run with race detector
go test -cover ./...             # Show test coverage
go test -coverprofile=coverage.out ./...  # Generate coverage report
go tool cover -html=coverage.out          # View coverage in browser
go test -bench=. ./...           # Run benchmarks
go test -short ./...             # Skip long-running tests
```

### Advanced Testing Patterns
- **Parallel Tests**: Use `t.Parallel()` for independent tests to speed up test runs
- **Benchmarks**: Add benchmark functions with `Benchmark` prefix for performance tracking
- **Test Fixtures**: Use `testdata/` directory for test input files
- **Build Tags**: Use build tags to separate integration tests: `//go:build integration`
- **Context Usage**: Pass context to functions and test timeout handling

**Example: Benchmark Test**
```go
func BenchmarkValidatePaystring(b *testing.B) {
    paystring := "user$domain.com"
    b.ResetTimer() // Reset after setup
    for i := 0; i < b.N; i++ {
        _ = ValidatePaystring(paystring)
    }
}
```

**Example: Parallel Test**
```go
func TestMultipleValidations(t *testing.T) {
    tests := []string{"user1$domain.com", "user2$example.org"}
    for _, tc := range tests {
        tc := tc // Capture range variable
        t.Run(tc, func(t *testing.T) {
            t.Parallel() // Run in parallel
            err := ValidatePaystring(tc)
            if err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })
    }
}
```

## Building and Development

### Build Commands
```bash
go build -v .                    # Build the project
go run main.go                   # Run directly
go mod tidy                      # Clean up dependencies
go mod verify                    # Verify dependencies
```

### Code Quality
```bash
gofmt -s -w .                    # Format all Go files
go vet ./...                     # Run Go vet
golint ./...                     # Run golint (if available)
```

### CI/CD Workflow
- All changes must go through pull requests
- Run tests locally before pushing: `go test ./...`
- Ensure code is formatted: `gofmt -s -w .`
- Fix all `go vet` warnings before submitting PR
- Build must succeed without errors
- All tests must pass (zero tolerance for failing tests)

### Git Workflow and Best Practices

**Branch Naming:**
- Feature branches: `feature/description` or `copilot/description`
- Bug fixes: `fix/issue-number-description`
- Documentation: `docs/description`

**Commit Messages:**
- Use clear, descriptive commit messages
- Start with a verb in imperative mood (Add, Fix, Update, Remove)
- Keep first line under 72 characters
- Add detailed description in commit body if needed
- Examples:
  - ✅ `Add validation for payment domain format`
  - ✅ `Fix race condition in watcher initialization`
  - ❌ `update code` (too vague)
  - ❌ `fixed stuff` (not descriptive)

**Pull Request Guidelines:**
- Provide clear PR title and description
- Reference related issues using `Fixes #123` or `Closes #123`
- Include test results in PR description
- Request review from appropriate team members
- Address all review comments before merging
- Keep PRs focused and reasonably sized (prefer multiple small PRs over one large PR)

**Code Review Checklist:**
1. All tests pass locally and in CI
2. Code follows Go idioms and project conventions
3. New code has appropriate test coverage
4. Documentation is updated if needed
5. No security vulnerabilities introduced
6. Performance impact is acceptable
7. Error handling is comprehensive

## Project-Specific Patterns

### Component Initialization
- Each internal package provides a `New()` constructor function
- Constructors should return fully initialized, ready-to-use instances
- Use functional options pattern for complex initialization when needed

### Package Structure
- `internal/aggregator`: Handles aggregation of validator data
- `internal/attest`: Provides attestation functionality
- `internal/bls`: BLS signature operations and cryptography
- `internal/farming`: Farming operations and payout management
- `internal/keystore`: Cryptographic key management
- `internal/payment`: Payment and paystring validation logic
- `internal/watcher`: Monitors validator state and blockchain events

**Detailed Package Responsibilities:**

**`internal/payment`** (Payment Processing)
- Validates paystring format (user$domain.com)
- Enforces paystring constraints (length, format, domain validation)
- Provides constructor `New(paystring)` with validation
- Getter method `GetPaystring()` for safe access
- Security: Input validation, length limits, format enforcement
- **Key Files**: `payment.go`, `payment_test.go`
- **Dependencies**: Standard library only (strings)

**`internal/farming`** (Farming Operations)
- Manages farming operations and payout scheduling for validators
- Handles onboarding process for farmers with paystring validation
- Supports multiple payout schedules: daily (default), weekly, monthly
- Thread-safe operations with mutex protection for concurrent access
- Provides constructor `New()` and `NewWithConfig()` for initialization
- Configuration management: payout schedule, minimum payout, enabled status
- **Key Files**: `farming.go`, `farming_test.go`, `template.go`, `template_test.go`
- **Dependencies**: `internal/payment` for paystring validation
- **Key Pattern**: Thread-safe state management with RWMutex
- **Testing**: Comprehensive test coverage including edge cases

**`internal/bls`** (BLS Cryptography)
- BLS signature operations for validator consensus
- Critical security component - handle with care
- All changes require security review
- Must use established cryptographic libraries
- **Key Pattern**: Immutable operations, no side effects
- **Testing**: Comprehensive test coverage required

**`internal/keystore`** (Key Management)
- Secure storage and retrieval of cryptographic keys
- Implements secure key lifecycle (generation, storage, deletion)
- Must zero sensitive data after use
- **Security**: Highest sensitivity - never log keys
- **Key Pattern**: Secure-by-default design

**`internal/aggregator`** (Data Aggregation)
- Aggregates validator data from multiple sources
- Performance-sensitive component
- Consider batching and caching strategies
- **Key Pattern**: Builder pattern for aggregation

**`internal/attest`** (Attestation)
- Provides attestation functionality for validators
- Works with BLS for signature verification
- **Key Pattern**: Stateless operations preferred

**`internal/watcher`** (State Monitoring)
- Monitors blockchain validator state changes
- Event-driven architecture
- **Key Pattern**: Observer pattern for state changes
- **Performance**: Efficient polling and event handling

**Package Dependencies:** Packages should minimize cross-dependencies. If a package needs functionality from another, consider:
1. Using interfaces to define contracts between packages
2. Creating a shared package for common types (avoid if possible)
3. Reviewing if functionality should be moved to a more appropriate package
4. Using dependency injection for testing
5. Keeping dependencies acyclic (no circular dependencies)

**Package Design Principles:**
- Each package should have a single, well-defined responsibility
- Packages should be independently testable
- Public API should be minimal and well-documented
- Internal implementation details should be unexported
- Use interfaces to define contracts, not implementations

## Common Tasks

### Adding a New Component
1. Create a new package under `internal/`
2. Add package documentation explaining its purpose
3. Implement a `New()` constructor function
4. Define clear interfaces for the component's behavior
5. Add comprehensive tests
6. Update `main.go` to initialize the component if needed

### Implementing Cryptographic Operations
- Use well-established cryptographic libraries
- Follow security best practices for key handling
- Add thorough testing including edge cases
- Document security considerations

### Modifying Existing Components
- Maintain backward compatibility where possible
- Update tests to cover new behavior
- Update documentation to reflect changes
- Consider impact on other components

## Dependencies
- Minimize external dependencies
- Prefer standard library when possible
- Carefully vet any new cryptographic or security-related dependencies
- Keep `go.mod` clean and up-to-date with `go mod tidy`

### Dependency Vetting Checklist

**Before adding any dependency, complete ALL steps:**

1. **Necessity Check**
   - [ ] Is this functionality available in Go standard library?
   - [ ] Can we implement this ourselves with reasonable effort?
   - [ ] Will this dependency provide significant value?
   - [ ] Is the added complexity justified?

2. **Maintenance and Activity**
   - [ ] Last commit within the past 6 months
   - [ ] Active response to issues and PRs
   - [ ] Clear versioning and release process
   - [ ] Responsive maintainers
   - [ ] Regular security updates

3. **Security Assessment**
   - [ ] Check for known CVEs using `go list -m -json all | nancy sleuth` or similar
   - [ ] Review the dependency's security policy
   - [ ] Check if the dependency has a security disclosure process
   - [ ] Scan for hardcoded credentials or secrets in the code
   - [ ] Verify cryptographic operations use standard libraries

4. **Code Quality**
   - [ ] Well-documented with examples
   - [ ] Has comprehensive test coverage
   - [ ] Follows Go best practices
   - [ ] Clear API and minimal breaking changes
   - [ ] No use of `unsafe` package (unless absolutely necessary)

5. **Transitive Dependencies**
   - [ ] Review all dependencies pulled in transitively
   - [ ] Check total dependency count impact
   - [ ] Verify licenses are compatible
   - [ ] Ensure no conflicting versions

6. **License Compatibility**
   - [ ] License is compatible with project (prefer MIT, Apache 2.0, BSD)
   - [ ] License allows commercial use
   - [ ] No copyleft requirements that conflict with project

7. **Community and Adoption**
   - [ ] Used by reputable projects
   - [ ] Has significant GitHub stars/community support
   - [ ] Active community discussions
   - [ ] Good reputation in the Go community

**Adding Dependencies - Step by Step:**

```bash
# Step 1: Add the dependency with specific version
go get github.com/example/package@v1.2.3

# Step 2: Verify the dependency was added correctly
go mod tidy
go mod verify

# Step 3: Check for vulnerabilities (if you have govulncheck)
govulncheck ./...

# Step 4: Review what was added to go.mod and go.sum
git diff go.mod go.sum

# Step 5: Run full test suite to ensure compatibility
go test ./...
go test -race ./...

# Step 6: Build to ensure no conflicts
go build -v .

# Step 7: Document why this dependency was added
# Add comment to go.mod or update documentation
```

**Dependency Update Strategy:**
```bash
# Check for available updates
go list -u -m all

# Update specific dependency
go get -u github.com/example/package

# Update all dependencies (use cautiously)
go get -u ./...

# Always verify after updates
go mod tidy
go test ./...
```

**Removing Dependencies:**
```bash
# Remove import from code first
# Then clean up go.mod
go mod tidy

# Verify dependency is gone
go list -m all | grep package-name
```

**Red Flags - Avoid Dependencies That:**
- Have no recent activity (abandoned projects)
- Lack documentation or examples
- Have many open security issues
- Use deprecated Go features
- Have excessive transitive dependencies (>20)
- Require `unsafe` package extensively
- Have unclear or restrictive licenses
- Show signs of poor code quality

## Documentation
- Keep the README.md up-to-date with project changes
- Document any setup or configuration requirements
- Provide examples for complex usage patterns
- Maintain inline code documentation for exported APIs

## Environment Setup

### Development Environment
1. **Go Installation:** Ensure Go 1.24.10 or later is installed
   ```bash
   go version  # Should show go1.24.10 or higher
   ```

2. **Clone and Setup:**
   ```bash
   git clone https://github.com/dreadwitdastacc-IFA/fluffy-umbrella.git
   cd fluffy-umbrella
   go mod download
   go build -v .
   ```
   
   **Note:** The repository is named `fluffy-umbrella`, but the Go module is `github.com/dreadwitdastacc-IFA/validatord`.

3. **Verify Setup:**
   ```bash
   go test ./...
   go vet ./...
   gofmt -l .
   ```

### Editor Configuration
- **VS Code:** Recommended extensions:
  - Go (by Go Team at Google)
  - GitHub Copilot
- **Settings:** Use workspace settings for consistent formatting
  ```json
  {
    "go.formatTool": "gofmt",
    "go.lintOnSave": "package",
    "editor.formatOnSave": true
  }
  ```

## Performance Considerations

### General Guidelines
- Avoid premature optimization; prioritize correctness first
- Profile before optimizing: use `go test -bench` and `pprof`
- Consider memory allocations in hot paths
- Use `sync.Pool` for frequently allocated objects when appropriate
- Prefer value types over pointers for small structs (< 3-4 fields)

**Example: Efficient String Building**
```go
// Bad: creates multiple intermediate strings
func buildMessage(parts []string) string {
    msg := ""
    for _, part := range parts {
        msg += part + ", "
    }
    return msg
}

// Good: uses strings.Builder
func buildMessage(parts []string) string {
    var builder strings.Builder
    for i, part := range parts {
        if i > 0 {
            builder.WriteString(", ")
        }
        builder.WriteString(part)
    }
    return builder.String()
}
```

### Cryptographic Performance
- BLS operations are computationally expensive; batch when possible
- Cache validated results when safe to do so
- Avoid unnecessary key generation; reuse when appropriate
- Consider using goroutines for independent crypto operations

## Debugging and Troubleshooting

### Common Issues and Solutions

**Build Failures:**
```bash
# Clear cache and rebuild
go clean -cache
go mod tidy
go build -v .
```

**Test Failures:**
```bash
# Run specific package tests with verbose output
go test -v ./internal/payment

# Run with race detector
go test -race ./...

# Run specific test
go test -v -run TestValidatePaystring ./internal/payment
```

**Import Issues:**
```bash
# Fix missing imports
go mod tidy

# Update a specific dependency
go get -u github.com/example/package
```

### Debugging Techniques
- Use `fmt.Printf` or `log.Printf` for quick debugging
- Use `go test -v` to see test output
- Use `t.Logf()` in tests instead of `fmt.Println`
- For complex issues, consider using `delve` debugger:
  ```bash
  # Install delve
  go install github.com/go-delve/delve/cmd/dlv@latest
  
  # Debug a test
  dlv test ./internal/payment -- -test.run TestValidatePaystring
  ```

### Diagnostic Tools
- **Scripts:** Use `./scripts/fluffy-payout-diagnostics.sh` for payment diagnostics
- **Logs:** Check system logs with `journalctl -u validatord` if running as service
- **Database:** Use `sqlite3` to inspect database if applicable

## Working with Copilot Agent

### Good Tasks for Copilot
Copilot coding agent works well for:
- Adding new internal packages following existing patterns
- Writing tests for existing code
- Refactoring small, well-defined functions
- Adding documentation and comments
- Implementing helper functions and utilities
- Fixing bugs with clear reproduction steps
- Adding validation and error handling
- Updating documentation (README, code comments)
- Creating diagnostic scripts

### Tasks Requiring Human Review
These tasks need extra scrutiny:
- Changes to cryptographic operations (BLS, keystore)
- Modifications to payment logic
- Security-sensitive code changes
- Large-scale refactoring across multiple packages
- Changes to public APIs
- Performance-critical code paths
- Database schema changes
- CI/CD configuration updates

### How to Get Best Results

**Crafting Effective Prompts:**
1. **Be Specific**: Provide clear, specific issue descriptions with acceptance criteria
2. **Include Examples**: Show expected behavior with input/output examples
3. **Reference Patterns**: Point to existing code patterns to follow
4. **Specify Scope**: Indicate which files or packages should be modified
5. **Set Constraints**: Mention performance, security, or compatibility requirements
6. **Break Down Complex Tasks**: Divide large tasks into smaller, focused subtasks
7. **Add Context**: Link to related issues, documentation, or relevant code sections

**Example of Good Prompt:**
```
Add validation to the payment package to ensure paystring domains 
have valid TLD (top-level domain). Follow the pattern in 
ValidatePaystring(). Add table-driven tests covering:
- Valid TLDs (.com, .org, .net)
- Invalid TLDs (.invalid, .test)
- Edge cases (empty TLD, multiple dots)
Performance requirement: validation should complete in <1ms.
```

**Example of Poor Prompt:**
```
Fix the payment thing
```

### Code Review Checklist

When reviewing Copilot-generated code, systematically verify:

**Functionality:**
- [ ] Code solves the stated problem completely
- [ ] All edge cases are handled
- [ ] Error cases are handled properly
- [ ] Tests cover success and failure scenarios
- [ ] No regression in existing functionality

**Code Quality:**
- [ ] Code follows Go idioms and project conventions
- [ ] Functions are focused and have single responsibility
- [ ] Variable and function names are descriptive
- [ ] Code is readable and maintainable
- [ ] No unnecessary complexity or over-engineering
- [ ] Comments explain "why" not "what"

**Security:**
- [ ] No security vulnerabilities introduced
- [ ] Input validation is comprehensive
- [ ] No secrets or sensitive data in code
- [ ] Cryptographic operations use approved methods
- [ ] Error messages don't leak sensitive information
- [ ] No timing attack vulnerabilities

**Testing:**
- [ ] Tests are comprehensive and meaningful
- [ ] Test names clearly describe what is being tested
- [ ] Both positive and negative cases are covered
- [ ] Edge cases are tested
- [ ] Tests are independent and deterministic
- [ ] No flaky tests introduced

**Dependencies:**
- [ ] Dependencies are necessary and vetted
- [ ] No unnecessary external dependencies added
- [ ] Transitive dependencies are acceptable
- [ ] License compatibility verified

**Performance:**
- [ ] Performance impact is acceptable
- [ ] No obvious performance anti-patterns
- [ ] Resource usage is reasonable
- [ ] No memory leaks or goroutine leaks

**Compatibility:**
- [ ] Backward compatibility is maintained
- [ ] Breaking changes are documented
- [ ] API changes are justified and necessary

**Documentation:**
- [ ] Documentation is clear and accurate
- [ ] All exported functions are documented
- [ ] Complex logic has explanatory comments
- [ ] README is updated if needed

### Feedback Loop and Iteration

**Providing Effective Feedback:**
1. **Be Specific**: Point to exact lines or sections that need changes
2. **Explain Why**: Don't just say what's wrong, explain why it's problematic
3. **Reference Standards**: Point to relevant sections of these instructions
4. **Provide Examples**: Show what you expected vs. what was generated
5. **Prioritize**: Indicate which issues are blocking vs. nice-to-have

**Feedback Template:**
```markdown
In [file.go:123], [specific issue]:
- Problem: [what's wrong]
- Why: [why it's problematic]
- Expected: [what should be done instead]
- Reference: [link to relevant guideline]
Example: [code example if applicable]
```

**Common Feedback Scenarios:**

**Security Issue:**
```
In payment.go:45, input validation is insufficient:
- Problem: No length check on paystring input
- Why: Could lead to DoS via large inputs
- Expected: Add max length validation (256 chars) before processing
- Reference: Security Considerations - Input size limits
```

**Pattern Mismatch:**
```
In aggregator.go:78, error handling doesn't follow project pattern:
- Problem: Using panic() for validation failure
- Why: Violates "What NOT to Do" - panic usage restriction
- Expected: Return descriptive error instead
- Reference: Error Handling section
Example: return nil, fmt.Errorf("validation failed: %w", err)
```

**Test Gap:**
```
In payment_test.go, missing edge case coverage:
- Problem: No test for empty domain after '$'
- Why: Edge cases must be tested per Testing Guidelines
- Expected: Add test case for "user$" input
- Reference: Testing Guidelines - Edge case coverage
```

### Iterative Improvement Process

1. **Initial Implementation**: Copilot provides first iteration
2. **Review**: Apply code review checklist systematically
3. **Feedback**: Provide specific, actionable feedback
4. **Refinement**: Copilot addresses feedback
5. **Verify**: Confirm issues are resolved
6. **Repeat**: Iterate until all quality bars are met

**Quality Gates:**
- All tests must pass
- Code coverage must be maintained or improved
- No security vulnerabilities
- Follows all coding guidelines
- Documentation is complete and accurate

### Learning and Adaptation

Help Copilot improve by:
- **Documenting Patterns**: Add new patterns to these instructions when established
- **Updating Examples**: Keep code examples current and relevant
- **Sharing Feedback**: Contribute learnings back to these instructions
- **Tracking Issues**: Note common mistakes to address in guidelines

## Quick Reference

### Essential Commands
```bash
# Development
go run main.go              # Run the application
go build -v .               # Build the project
go mod tidy                 # Clean dependencies

# Testing
go test ./...               # Run all tests
go test -v ./...            # Verbose test output
go test -race ./...         # Run with race detector
go test -cover ./...        # Show test coverage

# Code Quality
gofmt -s -w .               # Format all files
go vet ./...                # Run static analysis
go mod verify               # Verify dependencies

# Debugging
go test -v -run TestName    # Run specific test
./scripts/fluffy-payout-diagnostics.sh  # Payment diagnostics
```

### File Structure
```
fluffy-umbrella/
├── .github/
│   ├── copilot-instructions.md    # This file
│   └── agents/                     # Custom agent configs
├── internal/                       # Private packages
│   ├── aggregator/                # Data aggregation
│   ├── attest/                    # Attestation logic
│   ├── bls/                       # BLS cryptography
│   ├── farming/                   # Farming operations
│   ├── keystore/                  # Key management
│   ├── payment/                   # Payment handling
│   └── watcher/                   # State monitoring
├── scripts/                        # Utility scripts
├── main.go                        # Application entry point
├── go.mod                         # Dependencies
└── README.md                      # Project documentation
```

### Getting Help
- **Go Documentation:** https://go.dev/doc/effective_go
- **Project Issues:** Reference existing issues for context
- **Diagnostic Scripts:** Use scripts in `scripts/` directory
- **Code Examples:** Review tests for usage examples (`*_test.go` files)

---

**Note:** These instructions are living documentation. Update them when:
- New patterns or conventions are established
- New tools or dependencies are added
- Security requirements change
- Common pitfalls are identified
