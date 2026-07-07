package resolver

import (
	"encoding/hex"
	"testing"
)

func TestNewSimpleQuery(t *testing.T) {
	expectedHex := "00160100000100000000000003646e7306676f6f676c6503636f6d0000010001"
	expectedBytes, err := hex.DecodeString(expectedHex)
	if err != nil {
		t.Fatalf("Failed to decode expected hex string: %v", err)
	}

	// ID is 22 (0x0016), domain is "dns.google.com"
	query := NewSimpleQuery(22, "dns.google.com")

	if string(query) != string(expectedBytes) {
		t.Errorf("Query mismatch.\nExpected: %x\nGot:      %x", expectedBytes, query)
	}
}
