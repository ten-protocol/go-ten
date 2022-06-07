package obscuroscan

import (
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

const expectedNonce = 777

func TestCanDecryptRollup(t *testing.T) {
	contractABI, err := abi.JSON(strings.NewReader(mgmtcontractlib.MgmtContractABI))
	if err != nil {
		panic(err)
	}

	rollupJSON, err := decryptRollup(generateEncryptedRollupHex(), contractABI)
	if err != nil {
		t.Fatalf("rollup decryption failed. Cause: %s", err)
	}
	var rollupJSONMap map[string]interface{}
	err = json.Unmarshal(rollupJSON, &rollupJSONMap)
	if err != nil {
		t.Fatalf("rollup JSON unmarshalling failed. Cause: %s", err)
	}

	// We use the nonce as an indicator of whether the entire rollup was successfully decrypted.
	recoveredNonce := rollupJSONMap["Header"].(map[string]interface{})["Nonce"].(float64)
	if recoveredNonce != expectedNonce {
		t.Fatalf("rollup JSON did not contain correct nonce")
	}
}

func TestThrowsIfEncryptedRollupIsInvalid(t *testing.T) {
	contractABI, err := abi.JSON(strings.NewReader(mgmtcontractlib.MgmtContractABI))
	if err != nil {
		panic(err)
	}

	_, err = decryptRollup([]byte("invalid_rollup"), contractABI)
	if err == nil {
		t.Fatal("did not error on invalid rollup")
	}
}

// Generates an encrypted rollup in hex encoding.
func generateEncryptedRollupHex() []byte {
	rollup := core.NewRollup(
		common.BigToHash(big.NewInt(0)),
		nil,
		obscurocommon.L2GenesisHeight,
		common.HexToAddress("0x0"),
		[]nodecommon.L2Tx{},
		[]nodecommon.Withdrawal{},
		expectedNonce,
		common.BigToHash(big.NewInt(0)),
	)
	rollupTx := &obscurocommon.L1RollupTx{
		Rollup: nodecommon.EncodeRollup(rollup.ToExtRollup(core.NewTransactionBlobCryptoImpl()).ToRollup()),
	}

	mgmtContractAddress := common.BigToAddress(big.NewInt(0))
	mgmtContractLib := mgmtcontractlib.NewMgmtContractLib(&mgmtContractAddress)
	rollupTxData := mgmtContractLib.CreateRollup(rollupTx, 0)

	prvKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	signedRollupTx, err := types.SignNewTx(prvKey, types.NewEIP155Signer(big.NewInt(0)), rollupTxData)
	if err != nil {
		panic(err)
	}

	encryptedRollupHex := common.Bytes2Hex(signedRollupTx.Data())
	return []byte(encryptedRollupHex)
}
