// Package transaction provides functionality for parsing and handling blockchain transactions.
package transaction

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

// TxIn represents a transaction input.
type TxIn struct {
	PreviousTxID    [32]byte // 32-byte hash of the previous transaction
	PreviousOutIdx  uint32   // Index of the output in the previous transaction
	ScriptSigLength uint64   // Length of the signature script
	ScriptSig       []byte   // Signature script
	Sequence        uint32   // Sequence number
}

// TxOut represents a transaction output.
type TxOut struct {
	Value           uint64 // Value in satoshis
	ScriptPubKeyLen uint64 // Length of the public key script
	ScriptPubKey    []byte // Public key script
}

// TxWitness represents witness data for a transaction input.
type TxWitness struct {
	Items [][]byte // Witness stack items
}

// Transaction represents a parsed Bitcoin transaction.
type Transaction struct {
	Version    int32       // Transaction version
	Marker     byte        // Marker (0x00 for SegWit transactions)
	Flag       byte        // Flag (0x01 for SegWit transactions)
	TxInCount  uint64      // Number of inputs
	TxIns      []TxIn      // Transaction inputs
	TxOutCount uint64      // Number of outputs
	TxOuts     []TxOut     // Transaction outputs
	Witnesses  []TxWitness // Witness data (for SegWit transactions)
	LockTime   uint32      // Lock time
	IsSegWit   bool        // Whether this is a SegWit transaction
}

// ParseError represents an error that occurred during transaction parsing.
type ParseError struct {
	Field   string
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("failed to parse %s: %s", e.Field, e.Message)
}

// readVarInt reads a variable-length integer from a byte slice and returns the value
// along with the number of bytes consumed.
func readVarInt(data []byte, offset int) (uint64, int, error) {
	if offset >= len(data) {
		return 0, 0, errors.New("unexpected end of data")
	}

	first := data[offset]
	switch {
	case first < 0xFD:
		return uint64(first), 1, nil
	case first == 0xFD:
		if offset+3 > len(data) {
			return 0, 0, errors.New("unexpected end of data for varint16")
		}
		return uint64(binary.LittleEndian.Uint16(data[offset+1 : offset+3])), 3, nil
	case first == 0xFE:
		if offset+5 > len(data) {
			return 0, 0, errors.New("unexpected end of data for varint32")
		}
		return uint64(binary.LittleEndian.Uint32(data[offset+1 : offset+5])), 5, nil
	default: // 0xFF
		if offset+9 > len(data) {
			return 0, 0, errors.New("unexpected end of data for varint64")
		}
		return binary.LittleEndian.Uint64(data[offset+1 : offset+9]), 9, nil
	}
}

// ParseHex parses a hex-encoded transaction string and returns a Transaction struct.
func ParseHex(hexStr string) (*Transaction, error) {
	if hexStr == "" {
		return nil, &ParseError{Field: "transaction", Message: "empty hex string"}
	}

	// Decode hex string to bytes
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, &ParseError{Field: "hex", Message: fmt.Sprintf("invalid hex encoding: %v", err)}
	}

	return Parse(data)
}

// Parse parses raw transaction bytes and returns a Transaction struct.
func Parse(data []byte) (*Transaction, error) {
	if len(data) < 10 { // Minimum: version(4) + marker/flag or input count + locktime(4)
		return nil, &ParseError{Field: "transaction", Message: "data too short"}
	}

	tx := &Transaction{}
	offset := 0

	// Parse version (4 bytes, little-endian)
	if offset+4 > len(data) {
		return nil, &ParseError{Field: "version", Message: "unexpected end of data"}
	}
	tx.Version = int32(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4

	// Check for SegWit marker and flag
	if offset+2 <= len(data) && data[offset] == 0x00 && data[offset+1] == 0x01 {
		tx.IsSegWit = true
		tx.Marker = data[offset]
		tx.Flag = data[offset+1]
		offset += 2
	}

	// Parse input count
	inputCount, bytesRead, err := readVarInt(data, offset)
	if err != nil {
		return nil, &ParseError{Field: "input_count", Message: err.Error()}
	}
	tx.TxInCount = inputCount
	offset += bytesRead

	// Parse inputs
	tx.TxIns = make([]TxIn, inputCount)
	for i := uint64(0); i < inputCount; i++ {
		txIn, bytesUsed, err := parseTxIn(data, offset)
		if err != nil {
			return nil, &ParseError{Field: fmt.Sprintf("input[%d]", i), Message: err.Error()}
		}
		tx.TxIns[i] = txIn
		offset += bytesUsed
	}

	// Parse output count
	outputCount, bytesRead, err := readVarInt(data, offset)
	if err != nil {
		return nil, &ParseError{Field: "output_count", Message: err.Error()}
	}
	tx.TxOutCount = outputCount
	offset += bytesRead

	// Parse outputs
	tx.TxOuts = make([]TxOut, outputCount)
	for i := uint64(0); i < outputCount; i++ {
		txOut, bytesUsed, err := parseTxOut(data, offset)
		if err != nil {
			return nil, &ParseError{Field: fmt.Sprintf("output[%d]", i), Message: err.Error()}
		}
		tx.TxOuts[i] = txOut
		offset += bytesUsed
	}

	// Parse witness data if SegWit
	if tx.IsSegWit {
		tx.Witnesses = make([]TxWitness, inputCount)
		for i := uint64(0); i < inputCount; i++ {
			witness, bytesUsed, err := parseWitness(data, offset)
			if err != nil {
				return nil, &ParseError{Field: fmt.Sprintf("witness[%d]", i), Message: err.Error()}
			}
			tx.Witnesses[i] = witness
			offset += bytesUsed
		}
	}

	// Parse locktime (4 bytes, little-endian)
	if offset+4 > len(data) {
		return nil, &ParseError{Field: "locktime", Message: "unexpected end of data"}
	}
	tx.LockTime = binary.LittleEndian.Uint32(data[offset : offset+4])

	return tx, nil
}

// parseTxIn parses a transaction input from raw bytes.
func parseTxIn(data []byte, offset int) (TxIn, int, error) {
	txIn := TxIn{}
	startOffset := offset

	// Previous transaction hash (32 bytes)
	if offset+32 > len(data) {
		return txIn, 0, errors.New("unexpected end of data for previous tx hash")
	}
	copy(txIn.PreviousTxID[:], data[offset:offset+32])
	offset += 32

	// Previous output index (4 bytes)
	if offset+4 > len(data) {
		return txIn, 0, errors.New("unexpected end of data for previous output index")
	}
	txIn.PreviousOutIdx = binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4

	// Script length (varint)
	scriptLen, bytesRead, err := readVarInt(data, offset)
	if err != nil {
		return txIn, 0, fmt.Errorf("failed to read script length: %w", err)
	}
	txIn.ScriptSigLength = scriptLen
	offset += bytesRead

	// Script signature
	if offset+int(scriptLen) > len(data) {
		return txIn, 0, errors.New("unexpected end of data for script signature")
	}
	txIn.ScriptSig = make([]byte, scriptLen)
	copy(txIn.ScriptSig, data[offset:offset+int(scriptLen)])
	offset += int(scriptLen)

	// Sequence (4 bytes)
	if offset+4 > len(data) {
		return txIn, 0, errors.New("unexpected end of data for sequence")
	}
	txIn.Sequence = binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4

	return txIn, offset - startOffset, nil
}

// parseTxOut parses a transaction output from raw bytes.
func parseTxOut(data []byte, offset int) (TxOut, int, error) {
	txOut := TxOut{}
	startOffset := offset

	// Value (8 bytes, little-endian)
	if offset+8 > len(data) {
		return txOut, 0, errors.New("unexpected end of data for value")
	}
	txOut.Value = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8

	// Script length (varint)
	scriptLen, bytesRead, err := readVarInt(data, offset)
	if err != nil {
		return txOut, 0, fmt.Errorf("failed to read script length: %w", err)
	}
	txOut.ScriptPubKeyLen = scriptLen
	offset += bytesRead

	// Script public key
	if offset+int(scriptLen) > len(data) {
		return txOut, 0, errors.New("unexpected end of data for script public key")
	}
	txOut.ScriptPubKey = make([]byte, scriptLen)
	copy(txOut.ScriptPubKey, data[offset:offset+int(scriptLen)])
	offset += int(scriptLen)

	return txOut, offset - startOffset, nil
}

// parseWitness parses witness data from raw bytes.
func parseWitness(data []byte, offset int) (TxWitness, int, error) {
	witness := TxWitness{}
	startOffset := offset

	// Number of witness items (varint)
	itemCount, bytesRead, err := readVarInt(data, offset)
	if err != nil {
		return witness, 0, fmt.Errorf("failed to read witness item count: %w", err)
	}
	offset += bytesRead

	// Parse each witness item
	witness.Items = make([][]byte, itemCount)
	for i := uint64(0); i < itemCount; i++ {
		// Item length (varint)
		itemLen, bytesRead, err := readVarInt(data, offset)
		if err != nil {
			return witness, 0, fmt.Errorf("failed to read witness item %d length: %w", i, err)
		}
		offset += bytesRead

		// Item data
		if offset+int(itemLen) > len(data) {
			return witness, 0, fmt.Errorf("unexpected end of data for witness item %d", i)
		}
		witness.Items[i] = make([]byte, itemLen)
		copy(witness.Items[i], data[offset:offset+int(itemLen)])
		offset += int(itemLen)
	}

	return witness, offset - startOffset, nil
}

// GetTxIDHex returns the transaction ID of a previous input as a hex string.
// Note: Bitcoin displays transaction IDs in reversed byte order.
func (t *TxIn) GetTxIDHex() string {
	// Reverse the bytes for display
	reversed := make([]byte, 32)
	for i := 0; i < 32; i++ {
		reversed[i] = t.PreviousTxID[31-i]
	}
	return hex.EncodeToString(reversed)
}

// GetScriptPubKeyHex returns the output script as a hex string.
func (t *TxOut) GetScriptPubKeyHex() string {
	return hex.EncodeToString(t.ScriptPubKey)
}

// GetValueBTC returns the output value in BTC.
func (t *TxOut) GetValueBTC() float64 {
	return float64(t.Value) / 100000000.0
}
