package walletextension

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"testing"
)

const (
	dataFieldPrefix            = "0xmethodID"
	padding                    = "000000000000000000000000"
	viewingKeyAddressHex       = "71C7656EC7ab88b098defB751B7401B5f6d8976F"
	viewingKeyAddressHexPadded = padding + viewingKeyAddressHex
	otherAddressHexPadded      = padding + "71C7656EC7ab88b098defB751B7401B5f6d8976E" // Differs only in the final byte.
)

var (
	viewingKeyAddressOne = common.HexToAddress("0x" + viewingKeyAddressHex)
	viewingKeyAddressTwo = common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976D") // Not in the data field.
	viewingKeyMap        = map[common.Address]*ecies.PrivateKey{
		viewingKeyAddressOne: nil,
		viewingKeyAddressTwo: nil,
	}
)

func TestCanSearchDataFieldForFrom(t *testing.T) {
	callParams := map[string]interface{}{"data": dataFieldPrefix + otherAddressHexPadded + viewingKeyAddressHexPadded}
	address, err := searchDataFieldForFrom(callParams, viewingKeyMap)
	if err != nil {
		t.Fatalf("did not expect an error but got %s", err)
	}
	if *address != viewingKeyAddressOne {
		t.Fatal("did not find correct viewing key Address in `data` field")
	}
}

func TestCanSearchDataFieldWhenHasUnexpectedLength(t *testing.T) {
	incorrectLengthArg := "arg2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" // Only 31 bytes.
	callParams := map[string]interface{}{"data": dataFieldPrefix + otherAddressHexPadded + viewingKeyAddressHexPadded + incorrectLengthArg}
	address, err := searchDataFieldForFrom(callParams, viewingKeyMap)
	if err != nil {
		t.Fatalf("did not expect an error but got %s", err)
	}
	if *address != viewingKeyAddressOne {
		t.Fatal("did not find correct viewing key Address in `data` field")
	}
}

func TestErrorsWhenDataFieldIsMissing(t *testing.T) {
	_, err := searchDataFieldForFrom(make(map[string]interface{}), viewingKeyMap)

	if err == nil {
		t.Fatal("`data` field was missing but not error was thrown")
	}
}

func TestGracefulWhenDataFieldTooShort(t *testing.T) {
	callParams := map[string]interface{}{"data": "tooshort"}
	address, err := searchDataFieldForFrom(callParams, viewingKeyMap)
	if err != nil {
		t.Fatalf("did not expect an error but got %s", err)
	}
	if address != nil {
		t.Fatal("`data` field was too short but Address was found anyway")
	}
}
