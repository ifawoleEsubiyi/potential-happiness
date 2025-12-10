package transaction

import (
	"encoding/hex"
	"strings"
	"testing"
)

// Sample SegWit transaction from the issue
const sampleTxHex = "02000000000101fb721e5ebcb80ded7e58f6eee2662e88670d52843f7b6ede91762bf695d9469e0900000000feffffff0b10270000000000001976a9142cacb2169c6255696a9d6f60ad948ccf723393a488ac10270000000000001976a914b8bcb22b1493094da0050319ba36ae4973c3bd2988ac10270000000000001600146aacb07738d18f83c829629438a19fb834a70ea650c300000000000017a914442a0d577b30ed57f8d366f530e2cf0f1de111768710270000000000001600142b9c0a7b79154761338c875f04866ad4e79c756410270000000000001976a91487f63000d6271fb09377b2e339c50cdadea8a76e88ac10270000000000001976a914bd87749eeca207fa3ad24133d78587f43f4c8d3388ac1027000000000000160014c7517af2fb087ac5d0f4eec0077fad2cc251e7ef50c30000000000001976a914e9b4cc38c8e85d156a1d391dfaa5ad39e7114b1c88ac1027000000000000160014c7fc2813c1bc402ee6ec7fa3d1d631274f2c5363680c090000000000160014b93ea02346e65fb5f9594e18673e4f378aee0035024730440220386aa649986e3645fab58cd7d536a4e1b3ee290f125986ff0dec559c61d84ece0220114a90f84d9f5835ff70a2eb0723caf26b0b287c790745e95c38eed729dd10b90121030b9f3a29ed9860e4b8e5885fdef17824e15d955bb8e8e205a77836392bf3c37770002e00"

func TestParseHex(t *testing.T) {
	tests := []struct {
		name        string
		hexStr      string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid segwit transaction",
			hexStr:  sampleTxHex,
			wantErr: false,
		},
		{
			name:        "empty hex string",
			hexStr:      "",
			wantErr:     true,
			errContains: "empty hex string",
		},
		{
			name:        "invalid hex characters",
			hexStr:      "zzzz",
			wantErr:     true,
			errContains: "invalid hex encoding",
		},
		{
			name:        "hex string too short",
			hexStr:      "0100",
			wantErr:     true,
			errContains: "data too short",
		},
		{
			name:        "odd length hex string",
			hexStr:      "010",
			wantErr:     true,
			errContains: "invalid hex encoding",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := ParseHex(tt.hexStr)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseHex() expected error, got nil")
					return
				}
				if tt.errContains != "" {
					if pe, ok := err.(*ParseError); ok {
						if pe.Message == "" || (pe.Message != "" && !strings.Contains(pe.Error(), tt.errContains)) {
							t.Errorf("ParseHex() error = %v, want error containing %q", err, tt.errContains)
						}
					} else if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("ParseHex() error = %v, want error containing %q", err, tt.errContains)
					}
				}
				return
			}
			if err != nil {
				t.Errorf("ParseHex() unexpected error = %v", err)
				return
			}
			if tx == nil {
				t.Errorf("ParseHex() returned nil transaction")
			}
		})
	}
}

func TestParseSegWitTransaction(t *testing.T) {
	tx, err := ParseHex(sampleTxHex)
	if err != nil {
		t.Fatalf("ParseHex() error = %v", err)
	}

	// Check version
	if tx.Version != 2 {
		t.Errorf("Version = %d, want 2", tx.Version)
	}

	// Check SegWit flag
	if !tx.IsSegWit {
		t.Error("IsSegWit = false, want true")
	}

	// Check input count
	if tx.TxInCount != 1 {
		t.Errorf("TxInCount = %d, want 1", tx.TxInCount)
	}

	// Check output count
	if tx.TxOutCount != 11 {
		t.Errorf("TxOutCount = %d, want 11", tx.TxOutCount)
	}

	// Check that inputs and outputs slices match counts
	if len(tx.TxIns) != int(tx.TxInCount) {
		t.Errorf("len(TxIns) = %d, want %d", len(tx.TxIns), tx.TxInCount)
	}
	if len(tx.TxOuts) != int(tx.TxOutCount) {
		t.Errorf("len(TxOuts) = %d, want %d", len(tx.TxOuts), tx.TxOutCount)
	}

	// Check witness data exists
	if len(tx.Witnesses) != int(tx.TxInCount) {
		t.Errorf("len(Witnesses) = %d, want %d", len(tx.Witnesses), tx.TxInCount)
	}

	// Check locktime - last 4 bytes "70002e00" in little-endian = 0x002e0070 = 3014768
	expectedLockTime := uint32(0x002e0070)
	if tx.LockTime != expectedLockTime {
		t.Errorf("LockTime = %d, want %d", tx.LockTime, expectedLockTime)
	}
}

func TestParseInputDetails(t *testing.T) {
	tx, err := ParseHex(sampleTxHex)
	if err != nil {
		t.Fatalf("ParseHex() error = %v", err)
	}

	if len(tx.TxIns) == 0 {
		t.Fatal("No inputs parsed")
	}

	input := tx.TxIns[0]

	// Check previous output index
	if input.PreviousOutIdx != 9 {
		t.Errorf("PreviousOutIdx = %d, want 9", input.PreviousOutIdx)
	}

	// Check sequence - hex "feffffff" in little-endian = 0xfffffffe
	expectedSequence := uint32(0xfffffffe)
	if input.Sequence != expectedSequence {
		t.Errorf("Sequence = 0x%08x, want 0x%08x", input.Sequence, expectedSequence)
	}

	// Check script sig length (0 for SegWit)
	if input.ScriptSigLength != 0 {
		t.Errorf("ScriptSigLength = %d, want 0 for SegWit", input.ScriptSigLength)
	}
}

func TestParseOutputDetails(t *testing.T) {
	tx, err := ParseHex(sampleTxHex)
	if err != nil {
		t.Fatalf("ParseHex() error = %v", err)
	}

	// Check first output value (10000 satoshis = 0x2710)
	if len(tx.TxOuts) == 0 {
		t.Fatal("No outputs parsed")
	}

	expectedValue := uint64(10000) // 0x2710 in little-endian
	if tx.TxOuts[0].Value != expectedValue {
		t.Errorf("TxOuts[0].Value = %d, want %d", tx.TxOuts[0].Value, expectedValue)
	}

	// Check output script exists
	if len(tx.TxOuts[0].ScriptPubKey) == 0 {
		t.Error("TxOuts[0].ScriptPubKey is empty")
	}
}

func TestParseWitnessDetails(t *testing.T) {
	tx, err := ParseHex(sampleTxHex)
	if err != nil {
		t.Fatalf("ParseHex() error = %v", err)
	}

	if len(tx.Witnesses) == 0 {
		t.Fatal("No witnesses parsed")
	}

	witness := tx.Witnesses[0]

	// Check witness item count (should be 2 for a P2WPKH input)
	if len(witness.Items) != 2 {
		t.Errorf("len(Witnesses[0].Items) = %d, want 2", len(witness.Items))
	}

	// First item should be the signature (71 bytes in this case)
	if len(witness.Items[0]) != 71 {
		t.Errorf("Witnesses[0].Items[0] length = %d, want 71", len(witness.Items[0]))
	}

	// Second item should be the public key (33 bytes for compressed pubkey)
	if len(witness.Items[1]) != 33 {
		t.Errorf("Witnesses[0].Items[1] length = %d, want 33", len(witness.Items[1]))
	}
}

func TestGetTxIDHex(t *testing.T) {
	tx, err := ParseHex(sampleTxHex)
	if err != nil {
		t.Fatalf("ParseHex() error = %v", err)
	}

	if len(tx.TxIns) == 0 {
		t.Fatal("No inputs parsed")
	}

	txIDHex := tx.TxIns[0].GetTxIDHex()
	// The txid should be 64 hex characters (32 bytes)
	if len(txIDHex) != 64 {
		t.Errorf("GetTxIDHex() length = %d, want 64", len(txIDHex))
	}

	// Verify it's valid hex
	_, err = hex.DecodeString(txIDHex)
	if err != nil {
		t.Errorf("GetTxIDHex() returned invalid hex: %v", err)
	}
}

func TestGetScriptPubKeyHex(t *testing.T) {
	tx, err := ParseHex(sampleTxHex)
	if err != nil {
		t.Fatalf("ParseHex() error = %v", err)
	}

	if len(tx.TxOuts) == 0 {
		t.Fatal("No outputs parsed")
	}

	scriptHex := tx.TxOuts[0].GetScriptPubKeyHex()
	if scriptHex == "" {
		t.Error("GetScriptPubKeyHex() returned empty string")
	}

	// Verify it's valid hex
	_, err = hex.DecodeString(scriptHex)
	if err != nil {
		t.Errorf("GetScriptPubKeyHex() returned invalid hex: %v", err)
	}
}

func TestGetValueBTC(t *testing.T) {
	tx, err := ParseHex(sampleTxHex)
	if err != nil {
		t.Fatalf("ParseHex() error = %v", err)
	}

	if len(tx.TxOuts) == 0 {
		t.Fatal("No outputs parsed")
	}

	// First output is 10000 satoshis = 0.0001 BTC
	btcValue := tx.TxOuts[0].GetValueBTC()
	expectedBTC := 0.0001
	if btcValue != expectedBTC {
		t.Errorf("GetValueBTC() = %f, want %f", btcValue, expectedBTC)
	}
}

func TestReadVarInt(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		expected  uint64
		bytesRead int
		wantErr   bool
	}{
		{
			name:      "single byte value",
			data:      []byte{0x42},
			expected:  0x42,
			bytesRead: 1,
		},
		{
			name:      "max single byte",
			data:      []byte{0xFC},
			expected:  0xFC,
			bytesRead: 1,
		},
		{
			name:      "two byte value",
			data:      []byte{0xFD, 0x00, 0x01},
			expected:  0x0100,
			bytesRead: 3,
		},
		{
			name:      "four byte value",
			data:      []byte{0xFE, 0x01, 0x02, 0x03, 0x04},
			expected:  0x04030201,
			bytesRead: 5,
		},
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "incomplete two byte",
			data:    []byte{0xFD, 0x00},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, bytesRead, err := readVarInt(tt.data, 0)
			if tt.wantErr {
				if err == nil {
					t.Error("readVarInt() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("readVarInt() unexpected error = %v", err)
				return
			}
			if value != tt.expected {
				t.Errorf("readVarInt() value = %d, want %d", value, tt.expected)
			}
			if bytesRead != tt.bytesRead {
				t.Errorf("readVarInt() bytesRead = %d, want %d", bytesRead, tt.bytesRead)
			}
		})
	}
}

func TestParseError(t *testing.T) {
	err := &ParseError{
		Field:   "test_field",
		Message: "test message",
	}

	expected := "failed to parse test_field: test message"
	if err.Error() != expected {
		t.Errorf("ParseError.Error() = %q, want %q", err.Error(), expected)
	}
}

func TestParse_InvalidData(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		errContains string
	}{
		{
			name:        "nil data",
			data:        nil,
			errContains: "data too short",
		},
		{
			name:        "empty data",
			data:        []byte{},
			errContains: "data too short",
		},
		{
			name:        "truncated version",
			data:        []byte{0x01, 0x02, 0x03},
			errContains: "data too short",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.data)
			if err == nil {
				t.Error("Parse() expected error, got nil")
				return
			}
			if !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("Parse() error = %v, want error containing %q", err, tt.errContains)
			}
		})
	}
}

// Benchmark for parsing
func BenchmarkParseHex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ParseHex(sampleTxHex)
		if err != nil {
			b.Fatalf("ParseHex() error = %v", err)
		}
	}
}
