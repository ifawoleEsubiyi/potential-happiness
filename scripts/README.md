# Utility Scripts
# Scripts Directory

This directory contains utility scripts and JavaScript utilities for the validatord project.

## JavaScript Utilities

### utilities.js

A collection of JavaScript utility functions for web and date operations.

#### Functions

##### `calculateDaysBetweenDates(begin, end)`

Calculates the number of days between two dates.

**Parameters:**
- `begin` (Date|string): The start date
- `end` (Date|string): The end date

**Returns:** `number` - The number of days between the two dates

**Examples:**
```javascript
// Using date strings
calculateDaysBetweenDates('2024-01-01', '2024-01-10'); // returns 9

// Using Date objects
calculateDaysBetweenDates(new Date('2024-01-01'), new Date('2024-01-10')); // returns 9

// Negative result when end is before begin
calculateDaysBetweenDates('2024-01-10', '2024-01-01'); // returns -9
```

**Features:**
- Accepts both Date objects and date strings
- Validates input dates
- Returns negative values when end date is before begin date
- Accounts for leap years

##### `highlightImagesWithoutAlt([borderStyle])`

Finds all images without alternate text and applies a red border to highlight them for accessibility review.

**Parameters:**
- `borderStyle` (string, optional): The CSS border style to apply (default: '3px solid red')

**Returns:** `Array` - The collection of images that were highlighted

**Examples:**
```javascript
// Highlight with default red border
const images = highlightImagesWithoutAlt();
console.log(`Found ${images.length} images without alt text`);

// Highlight with custom border style
highlightImagesWithoutAlt('5px dashed orange');
```

**Features:**
- Finds images missing the alt attribute
- Finds images with empty alt attributes
- Finds images with whitespace-only alt attributes
- Applies customizable border styling
- Returns array of affected images

##### `removeImageHighlighting()`

Removes highlighting from all images (utility to undo `highlightImagesWithoutAlt`).

**Example:**
```javascript
removeImageHighlighting(); // Removes all borders added by highlightImagesWithoutAlt
```

#### Testing

Run the test suite:
```bash
node scripts/utilities.test.js
```

#### Demo

Open `scripts/demo.html` in a web browser to see an interactive demonstration of all utilities.

---

## Cleanup Script

### clean.sh

A dedicated cleanup script for the scripts directory that removes temporary files, logs, and build artifacts.

#### Purpose

This script cleans:
- Log files (`*.log`)
- Temporary files (`*.tmp`, `*.temp`)
- Editor backup files (`*~`)
- JavaScript dependencies (`node_modules/`)
- Package lock files (`package-lock.json`)

#### Usage

From the scripts directory:
```bash
./scripts/clean.sh
```

Or use the Makefile target from the repository root:
```bash
make clean-scripts
```

#### Features

- Safe execution with error checking
- Informative output showing what's being cleaned
- Only removes files that exist (no errors for missing files)
- Can be run multiple times safely

---

## Diagnostics Scripts

## clean.sh

A dedicated cleanup script that removes all build artifacts, test files, temporary files, and script outputs from the validatord project.

### Purpose

This script cleans:
- Build artifacts (`validatord` binary, `coverage.out`)
- Test artifacts (`*.test`, `*.out`)
- Temporary files (`*.tmp`, `*.temp`, `*.log`)
- Profiling files (`*.prof`, `*.pprof`)
- Script outputs (`/tmp/fluffy-check/`)

### Usage

From the repository root:

```bash
./scripts/clean.sh
```

Or use the Makefile targets:

```bash
make clean              # Clean all artifacts
make clean-scripts      # Clean only scripts directory outputs
```

### Features

- Automatically detects repository root directory
- Provides progress messages for each cleanup step
- Safe to run multiple times (idempotent)
- Matches `.gitignore` patterns for consistency

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
