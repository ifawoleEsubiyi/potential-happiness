# Security Policy

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue, please report it responsibly.

### How to Report

**Do NOT open a public GitHub issue for security vulnerabilities.**

Instead, please report vulnerabilities by:

1. **Email**: Send details to the repository maintainers via the email listed in the repository contact information
2. **Private Disclosure**: Use GitHub's private vulnerability reporting feature if available

### What to Include

When reporting a vulnerability, please include:

- Description of the vulnerability
- Steps to reproduce the issue
- Potential impact assessment
- Any suggested fixes (if applicable)
- Your contact information for follow-up

### Response Timeline

- **Initial Response**: Within 48 hours of receiving your report
- **Status Update**: Within 7 days with our assessment
- **Resolution**: Security patches will be prioritized based on severity

### What to Expect

1. We will acknowledge your report within 48 hours
2. We will investigate and validate the vulnerability
3. We will work on a fix and coordinate disclosure timing with you
4. We will credit you for the discovery (unless you prefer to remain anonymous)

### Scope

This security policy applies to:

- The validatord daemon and all its components
- Any cryptographic operations (BLS, keystore)
- Payment and key handling functionality

### Security Best Practices

When contributing to this project:

- Never commit secrets, API keys, or credentials
- Use secure cryptographic practices
- Validate and sanitize all inputs
- Follow the principle of least privilege

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| main    | :white_check_mark: |

Thank you for helping keep validatord secure!
