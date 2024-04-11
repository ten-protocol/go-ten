package hostdb

import (
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/host/storage/init/sqlite"
)

const truncHash = 16

// An arbitrary number to put in the header
const batchNumber = 777

// truncTo16 checks if the leading half of the hash is filled with zeros and decides whether to truncate the first or last 16 bytes.
func truncTo16(hash gethcommon.Hash) []byte {
	hashBytes := hash.Bytes()
	// Check if the first half of the hash is all zeros
	if isLeadingHalfZeros(hashBytes) {
		return truncLastTo16(hashBytes)
	}
	return truncFirstTo16(hashBytes)
}

// isLeadingHalfZeros checks if the leading half of the hash is all zeros.
func isLeadingHalfZeros(bytes []byte) bool {
	halfLength := len(bytes) / 2
	for i := 0; i < halfLength; i++ {
		if bytes[i] != 0 {
			return false
		}
	}
	return true
}

// truncLastTo16 truncates the last 16 bytes of the hash.
func truncLastTo16(bytes []byte) []byte {
	if len(bytes) == 0 {
		return bytes
	}
	start := len(bytes) - truncHash
	if start < 0 {
		start = 0
	}
	b := bytes[start:]
	c := make([]byte, truncHash)
	copy(c, b)
	return c
}

// truncFirstTo16 truncates the first 16 bytes of the hash.
func truncFirstTo16(bytes []byte) []byte {
	if len(bytes) == 0 {
		return bytes
	}
	b := bytes[0:truncHash]
	c := make([]byte, truncHash)
	copy(c, b)
	return c
}

func createSQLiteDB(t *testing.T) (HostDB, error) {
	hostDB, err := sqlite.CreateTemporarySQLiteHostDB("", "mode=memory")
	if err != nil {
		t.Fatalf("unable to create temp sql db: %s", err)
	}
	return NewHostDB(hostDB, SQLiteSQLStatements())
}
