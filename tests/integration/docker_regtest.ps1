# PowerShell script for Go integration tests with bitcoind regtest
# Save as tests/integration/docker_regtest.ps1 and run in PowerShell

# Start bitcoind regtest in Docker if not already running
if (-not (docker ps | Select-String 'bitcoind-test')) {
    docker run -d --name bitcoind-test `
        -p 18443:18443 `
        -e BITCOIND_RPCUSER=regtest `
        -e BITCOIND_RPCPASSWORD=regtest `
        ruimarinho/bitcoin-core:24.0.1 `
        -regtest -printtoconsole -rpcbind=0.0.0.0 -rpcallowip=0.0.0.0/0
}

# Wait for bitcoind to be healthy
Write-Host "Waiting for bitcoind..."
for ($i=0; $i -lt 30; $i++) {
    $result = docker run --rm --network host ruimarinho/bitcoin-core:24.0.1 bitcoin-cli `
        -regtest -rpcconnect=127.0.0.1 -rpcport=18443 -rpcuser=regtest -rpcpassword=regtest getblockchaininfo 2>$null
    if ($result) {
        Write-Host "bitcoind ready"
        break
    }
    Start-Sleep -Seconds 2
}

# Build Go binary
$env:CGO_ENABLED = "1"
go build ./cmd/validatord

# Run integration tests
$env:BITCOIND_RPCURL = "http://regtest:regtest@127.0.0.1:18443"
./tests/integration/run_regtest.sh
