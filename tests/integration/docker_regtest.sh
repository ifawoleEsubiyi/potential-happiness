#!/usr/bin/env bash
set -e

# Start bitcoind regtest in Docker
if ! docker ps | grep -q bitcoind-test; then
  docker run -d --name bitcoind-test \
    -p 18443:18443 \
    -e BITCOIND_RPCUSER=regtest \
    -e BITCOIND_RPCPASSWORD=regtest \
    ruimarinho/bitcoin-core:24.0.1 \
    -regtest -printtoconsole -rpcbind=0.0.0.0 -rpcallowip=0.0.0.0/0
fi

# Wait for bitcoind to be healthy
for i in {1..30}; do
  if docker run --rm --network host ruimarinho/bitcoin-core:24.0.1 bitcoin-cli \
    -regtest -rpcconnect=127.0.0.1 -rpcport=18443 -rpcuser=regtest -rpcpassword=regtest getblockchaininfo >/dev/null 2>&1; then
    echo "bitcoind ready"
    break
  fi
  sleep 2
done

# Build Go binary
export CGO_ENABLED=1
go build ./cmd/validatord

# Run integration tests
export BITCOIND_RPCURL="http://regtest:regtest@127.0.0.1:18443"
./tests/integration/run_regtest.sh
