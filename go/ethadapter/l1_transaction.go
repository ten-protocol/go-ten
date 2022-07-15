package ethadapter

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// L1Transaction is an abstraction that transforms an Ethereum transaction into a format that can be consumed more easily by Obscuro.
type L1Transaction interface{}

type L1RollupTx struct {
	Rollup common.EncodedRollup
}

type L1DepositTx struct {
	Amount        uint64              // Amount to be deposited
	To            *gethcommon.Address // Address the ERC20 Transfer was made to (always be the Management Contract Addr)
	Sender        *gethcommon.Address // Address that issued the ERC20, the token holder or tx.origin
	TokenContract *gethcommon.Address // Address of the ERC20 Contract address that was executed
}

type L1RespondSecretTx struct {
	Secret      []byte
	RequesterID gethcommon.Address
	AttesterID  gethcommon.Address
	AttesterSig []byte
	HostAddress string
}

// Sign signs the payload with a given private key
func (l *L1RespondSecretTx) Sign(privateKey *ecdsa.PrivateKey) *L1RespondSecretTx {
	var data []byte
	data = append(data, l.AttesterID.Bytes()...)
	data = append(data, l.RequesterID.Bytes()...)
	data = append(data, l.HostAddress...)
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
	Attestation common.EncodedAttestationReport
}

type L1InitializeSecretTx struct {
	AggregatorID  *gethcommon.Address
	InitialSecret []byte
	HostAddress   string
	Attestation   common.EncodedAttestationReport
}
