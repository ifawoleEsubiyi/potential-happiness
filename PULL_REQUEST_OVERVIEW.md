# Pull Request Overview

## PR #104: Pin Amplify Runner Action to v0.4.0

This PR improved the stability and reproducibility of the Amplify Security workflow by pinning the runner action to a specific release version (v0.4.0) instead of tracking the unstable `develop` branch. This is a best practice for production workflows as it prevents unexpected breaking changes from being automatically introduced.

### Key Changes

- Pinned `amplify-security/runner-action` from `@develop` to `@v0.4.0` for predictable and stable security scans

### Benefits

- **Stability**: Using a fixed version prevents unexpected behavior changes
- **Reproducibility**: Security scans produce consistent results across runs
- **Security**: Avoids potential supply chain risks from tracking a development branch
- **Predictability**: Updates to the security scanning workflow are intentional and reviewed

### Files Changed

- `.github/workflows/amplify.yml`: Updated the runner action version from `@develop` to `@v0.4.0`

### Related Links

- Original PR: [#104](https://github.com/dreadwitdastacc-IFA/fluffy-umbrella-dd7ebd84/pull/104)
