package rollupcontractlib

import (
	"testing"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
)

func TestSerialization(t *testing.T) {
	rol := datagenerator.RandomRollup()

	serializedRollup := EncodeToString(nodecommon.EncodeRollup(&rol))
	deserializedRollup := DecodeFromString(serializedRollup)
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

	compressedRollup, err := Compress(nodecommon.EncodeRollup(&rol))
	if err != nil {
		t.Fatal(err)
	}
	serializedRollup := EncodeToString(compressedRollup)
	deserializedRollup := DecodeFromString(serializedRollup)
	decompressedRollup, err := Decompress(deserializedRollup)
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
