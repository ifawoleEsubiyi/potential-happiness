# Diagnostics Scripts

This directory contains utility scripts for the validatord project.

## fluffy-payout-diagnostics.sh

A comprehensive diagnostics script for troubleshooting payout systems, logs, and database operations.

### Purpose

This script collects diagnostic information about:
- Payout-related code patterns in the repository
- Payment and transaction logs
- Database operations and schema
- System logs from journalctl (if running as a service)
- SQLite database contents (if available)

### Usage

From the repository root:

```bash
./scripts/fluffy-payout-diagnostics.sh
```

Or, to install system-wide (requires sudo):

```bash
sudo cp scripts/fluffy-payout-diagnostics.sh /usr/local/bin/fluffy-payout-diagnostics.sh
sudo chmod +x /usr/local/bin/fluffy-payout-diagnostics.sh
sudo /usr/local/bin/fluffy-payout-diagnostics.sh
```

#### Environment Variables

- `VALIDATORD_DB_PATH` - Path to the SQLite database (default: `/var/lib/validatord/ledger.db`)

Example:
```bash
VALIDATORD_DB_PATH=/custom/path/ledger.db ./scripts/fluffy-payout-diagnostics.sh
```

### Output

All diagnostic data is saved to `/tmp/fluffy-check/` with the following files:

- `repo-payout-grep.txt` - Code grep for payout-related patterns
- `repo-logs-grep.txt` - Code grep for payment/transaction log patterns
- `repo-db-grep.txt` - Code grep for database operations
- `journal-payout-lines.txt` - System logs related to payouts (if journalctl available)
- `journal-tail-200.txt` - Last 200 lines of payout-related logs
- `journal-tail-200.redacted.txt` - Redacted version with paystrings anonymized
- `sqlite-tables.txt` - SQLite database tables (if DB exists)
- `sqlite-payments-schema.txt` - Payments table schema (if exists)
- `sqlite-payments-sample.txt` - Sample payment records (if exists)
- `sqlite-info.txt` - SQLite status message (if DB not found)

### Privacy

The script automatically redacts sensitive paystring information in the format `user$domain.com` to `user$REDACTED` in the journal output files to protect privacy.

### Requirements

- Git (for repository operations)
- Bash 4.0+ (with `set -euo pipefail` support)
- Optional: `journalctl` (for system log collection)
- Optional: `sqlite3` (for database diagnostics)

### Exit Codes

The script uses `|| true` for grep and optional commands to ensure it completes successfully even when:
- No matches are found in repository greps
- journalctl is not available
- SQLite database doesn't exist
- sqlite3 is not installed

This ensures the script always exits with code 0 and generates useful diagnostic output.
