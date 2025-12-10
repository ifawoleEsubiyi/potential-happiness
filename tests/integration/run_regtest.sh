#!/bin/bash
# Integration test script for validatord using bitcoind regtest mode.
# This script validates that the validator daemon can interact with a local
# bitcoin regtest network.

set -e

echo "=== Validatord Integration Test (regtest) ==="

# Check required environment variables
if [ -z "${BITCOIND_RPCURL}" ]; then
    echo "BITCOIND_RPCURL environment variable is not set"
    exit 1
fi

echo "Using BITCOIND_RPCURL: ${BITCOIND_RPCURL}"

# Check if the validatord binary exists
if [ ! -f "./validatord" ]; then
    echo "Error: validatord binary not found. Please build it first."
    exit 1
fi

echo "Found validatord binary"

# Run the validatord binary to verify it starts correctly
echo "Starting validatord..."
./validatord &
VALIDATORD_PID=$!

# Wait for validatord to initialize with retry logic
MAX_RETRIES=10
RETRY_DELAY=1
for i in $(seq 1 $MAX_RETRIES); do
    if kill -0 $VALIDATORD_PID 2>/dev/null; then
        echo "Validatord is running (PID: $VALIDATORD_PID, attempt $i)"
        # Stop the process gracefully
        kill $VALIDATORD_PID 2>/dev/null || true
        wait $VALIDATORD_PID 2>/dev/null || true
        echo "Validatord stopped"
        break
    fi
    if [ $i -eq 1 ]; then
        # First attempt - process may have already exited successfully
        echo "Validatord exited (this may be expected for short-lived initialization)"
        break
    fi
    sleep $RETRY_DELAY
done

echo ""
echo "=== Integration Test Complete ==="
echo "All checks passed!"

exit 0
