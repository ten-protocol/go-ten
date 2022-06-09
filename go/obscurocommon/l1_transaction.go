package obscurocommon

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// L1Transaction represents how obscuro interprets Transactions that happen on the Layer 1
// These transactions can have different implementations
type L1Transaction interface{}

type L1RollupTx struct {
	Rollup EncodedRollup
}

type L1DepositTx struct {
	Amount        uint64          // Amount to be deposited
	To            *common.Address // Address the ERC20 Transfer was made to (always be the Management Contract Addr)
	Sender        *common.Address // Address that issued the ERC20, the token holder or tx.origin
	TokenContract *common.Address // Address of the ERC20 Contract address that was executed
}

type L1RespondSecretTx struct {
	Secret      []byte
	RequesterID common.Address
	AttesterID  common.Address
	AttesterSig []byte
}

// Sign signs the payload with a given private key
func (l *L1RespondSecretTx) Sign(privateKey *ecdsa.PrivateKey) *L1RespondSecretTx {
	var data []byte
	data = append(data, l.AttesterID.Bytes()...)
	data = append(data, l.RequesterID.Bytes()...)
	data = append(data, string(l.Secret)...)

	// form the data
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), string(data))
	// hash the data
	hashedData := crypto.Keccak256Hash([]byte(msg))
	// sign the hash
	signedHash, err := crypto.Sign(hashedData.Bytes(), privateKey)
	if err != nil {
		return nil
	}
	// remove ECDSA recovery id
	signedHash = signedHash[:len(signedHash)-1]
	l.AttesterSig = signedHash
	return l
}

type L1RequestSecretTx struct {
	Attestation EncodedAttestationReport
}

type L1InitializeSecretTx struct {
	AggregatorID  *common.Address
	InitialSecret []byte
}
