#!/usr/bin/env bash
set -euo pipefail
OUTDIR="/tmp/fluffy-check"
mkdir -p "$OUTDIR"
REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
echo "Out: $OUTDIR"
git -C "$REPO_ROOT" grep -nE -C2 "payout|ExecutePayout|RunPayout|ProcessPayout|SendPayments|PayoutHandler|paypointer|paystring" 2>/dev/null | tee "$OUTDIR/repo-payout-grep.txt" || true
git -C "$REPO_ROOT" grep -nE -C2 "paystring|paypointer|txid|transaction|payments? sent|payment to|payout_id" 2>/dev/null | tee "$OUTDIR/repo-logs-grep.txt" || true
git -C "$REPO_ROOT" grep -nE -C2 "INSERT INTO|db.Exec|sqlx|gorm|sqlite3|bbolt|bolt|WriteFile|json.Marshal|Put\\(|Bucket\\(" 2>/dev/null | tee "$OUTDIR/repo-db-grep.txt" || true
if command -v journalctl >/dev/null 2>&1; then
  journalctl -u validatord --no-pager --since "30 days ago" | grep -iE "payout|payment|payments? sent|payment to|txid|payout_id|paypointer" | tee "$OUTDIR/journal-payout-lines.txt" || true
  tail -n 200 "$OUTDIR/journal-payout-lines.txt" > "$OUTDIR/journal-tail-200.txt" || true
  if [ -f "$OUTDIR/journal-tail-200.txt" ]; then
    sed -E 's/([A-Za-z0-9._-]+)\$([A-ZMIN_TEXT_KEYWORDS='bitcoin|btc|mining|min(er|ed)|stratum|getblocktemplate|getwork|cgminer|bfgminer|cpuminer|xmrig|pool|wallet|coinbase|minergate|solo-miner'
LARGE_FILE_SIZE="+100k" # find threshold for "large" files
BINARY_SIZE_THRESHOLD="+200k" # additional candidates for binary inspectiona-z0-9._-]+)/\1$REDACTED/g' "$OUTDIR/journal-tail-200.txt" > "$OUTDIR/journal-tail-200.redacted.txt" || true
  fi
else
  echo "journalctl not available" > "$OUTDIR/journal-payout-lines.txt"
fi
SQLITE_DB="${VALIDATORD_DB_PATH:-/var/lib/validatord/ledger.db}"
if [ -f "$SQLITE_DB" ] && command -v sqlite3 >/dev/null 2>&1; then
  sqlite3 "$SQLITE_DB" ".tables" > "$OUTDIR/sqlite-tables.txt" || true
  sqlite3 "$SQLITE_DB" "PRAGMA table_info(payments);" > "$OUTDIR/sqlite-payments-schema.txt" || true
  sqlite3 "$SQLITE_DB" "SELECT * FROM payments ORDER BY created_at DESC LIMIT 20;" > "$OUTDIR/sqlite-payments-sample.txt" || true
else
  echo "sqlite DB not found at $SQLITE_DB or sqlite3 not installed" > "$OUTDIR/sqlite-info.txt"
fi
echo "Saved to $OUTDIR"
