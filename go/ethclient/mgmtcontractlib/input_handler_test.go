package mgmtcontractlib

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
)

func TestRollupSerialization(t *testing.T) {
	rol := datagenerator.RandomRollup()

	serializedRollup := EncodeToString(nodecommon.EncodeRollup(&rol))
	deserializedRollup, _ := DecodeFromString(serializedRollup)
	newRollup, err := nodecommon.DecodeRollup(deserializedRollup)
	if err != nil {
		t.Fatal(err)
	}

	if rol.Hash() != newRollup.Hash() {
		t.Errorf("unexpected hashes when converting")
	}
}

func TestAttestationSerialization(t *testing.T) {
	att := obscurocommon.AttestationReport{
		Report: []byte("REPORT BYTES"),
		PubKey: []byte("PUBLIC KEY"),
		Owner:  common.Address{123},
	}

	serializedAttestation := EncodeToString(nodecommon.EncodeAttestation(&att))
	deserialized, err := DecodeFromString(serializedAttestation)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := nodecommon.DecodeAttestation(deserialized)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(decoded.Report, att.Report) ||
		!bytes.Equal(decoded.PubKey, att.PubKey) ||
		decoded.Owner != att.Owner {
		t.Errorf("unexpected hashes when converting attestation report")
	}
}

func TestCompression(t *testing.T) {
	rol := datagenerator.RandomRollup()

	compressedRollup, err := Compress(nodecommon.EncodeRollup(&rol))
	if err != nil {
		t.Fatal(err)
	}
	serializedRollup := EncodeToString(compressedRollup)
	deserializedRollup, _ := DecodeFromString(serializedRollup)
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
