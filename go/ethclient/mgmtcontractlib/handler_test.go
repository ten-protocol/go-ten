package mgmtcontractlib

import (
	"testing"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
)

func TestSerialization(t *testing.T) {
	rol := datagenerator.RandomRollup()

	serializedRollup := base64EncodeToString(nodecommon.EncodeRollup(&rol))
	deserializedRollup := base64DecodeFromString(serializedRollup)
	newRollup, err := nodecommon.DecodeRollup(deserializedRollup)
	if err != nil {
		t.Fatal(err)
	}

	if rol.Hash() != newRollup.Hash() {
		t.Errorf("unexpected hashes when converting")
	}
}

func TestCompression(t *testing.T) {
	rol := datagenerator.RandomRollup()

	compressedRollup, err := compress(nodecommon.EncodeRollup(&rol))
	if err != nil {
		t.Fatal(err)
	}
	serializedRollup := base64EncodeToString(compressedRollup)
	deserializedRollup := base64DecodeFromString(serializedRollup)
	decompressedRollup, err := decompress(deserializedRollup)
	if err != nil {
		t.Fatal(err)
	}
	newRollup, err := nodecommon.DecodeRollup(decompressedRollup)
	if err != nil {
		t.Fatal(err)
	}

	if rol.Hash() != newRollup.Hash() {
		t.Errorf("unexpected hashes when converting")
	}
}
