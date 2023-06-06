package datagenerator

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/common/viewingkey"
)

func RandomViewingKey() (*viewingkey.ViewingKey, error) {
	vkKeyPair, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate VK KeyPair - %w", err)
	}
	vkPubKey := crypto.CompressPubkey(&vkKeyPair.PublicKey)

	userKeyPair, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate User KeyPair - %w", err)
	}
	address := crypto.PubkeyToAddress(userKeyPair.PublicKey)

	vkSignature, err := viewingkey.Sign(userKeyPair, vkPubKey)
	if err != nil {
		return nil, fmt.Errorf("unable to sign the vk with the user pk - %w", err)
	}
	return &viewingkey.ViewingKey{
		Account:    &address,
		PrivateKey: ecies.ImportECDSA(vkKeyPair),
		PublicKey:  vkPubKey,
		SignedKey:  vkSignature, // the Viewing key is signed with the users Private Key
	}, err
}
