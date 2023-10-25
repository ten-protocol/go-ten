package faucet

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/obscuronet/go-obscuro/tools/walletextension/lib"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/obscuronet/go-obscuro/go/wallet"
)

// GatewayUser TODO (@ziga) refactor GatewayUser and integrate it with OGlib.
// GatewayUser is a struct that includes everything a gateway user has and uses (userID, wallets, http & ws addresses and client )
type GatewayUser struct {
	UserID            string
	Wallets           []wallet.Wallet
	HTTPClient        *ethclient.Client
	WSClient          *ethclient.Client
	ServerAddressHTTP string
	ServerAddressWS   string
}

func NewUser(wallets []wallet.Wallet, serverAddressHTTP string, serverAddressWS string) (*GatewayUser, error) {
	// automatically join OG
	userID, err := joinObscuroGateway(serverAddressHTTP)
	if err != nil {
		return nil, err
	}

	// create clients
	httpClient, err := ethclient.Dial(serverAddressHTTP + "/v1/" + "?u=" + userID)
	if err != nil {
		return nil, err
	}
	wsClient, err := ethclient.Dial(serverAddressWS + "/v1/" + "?u=" + userID)
	if err != nil {
		return nil, err
	}

	return &GatewayUser{
		UserID:            userID,
		Wallets:           wallets,
		HTTPClient:        httpClient,
		WSClient:          wsClient,
		ServerAddressHTTP: serverAddressHTTP,
		ServerAddressWS:   serverAddressWS,
	}, nil
}

func (u GatewayUser) RegisterAccounts() error {
	for _, w := range u.Wallets {
		response, err := registerAccount(u.ServerAddressHTTP, u.UserID, w.PrivateKey(), w.Address().Hex())
		if err != nil {
			return err
		}
		fmt.Printf("Successfully registered address %s for user: %s with response: %s \n", w.Address().Hex(), u.UserID, response)
	}

	return nil
}

func (u GatewayUser) PrintUserAccountsBalances() error {
	for _, w := range u.Wallets {
		balance, err := u.HTTPClient.BalanceAt(context.Background(), w.Address(), nil)
		if err != nil {
			return err
		}
		fmt.Println("Balance for account ", w.Address().Hex(), " - ", balance.String())
	}
	return nil
}

// TODO (@ziga) - use OGlib for registering accounts
func registerAccount(url string, userID string, pk *ecdsa.PrivateKey, hexAddress string) ([]byte, error) {
	payload := prepareRegisterPayload(userID, pk, hexAddress)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		url+"/v1/authenticate/?u="+userID,
		strings.NewReader(payload),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

func prepareRegisterPayload(userID string, pk *ecdsa.PrivateKey, hexAddress string) string {
	message := fmt.Sprintf("Register %s for %s", userID, strings.ToLower(hexAddress))
	prefixedMessage := fmt.Sprintf("\u0019Ethereum Signed Message:\n%d%s", len(message), message)
	messageHash := crypto.Keccak256([]byte(prefixedMessage))
	sig, err := crypto.Sign(messageHash, pk)
	if err != nil {
		fmt.Printf("Failed to sign message: %v\n", err)
	}
	sig[64] += 27
	signature := "0x" + hex.EncodeToString(sig)
	payload := fmt.Sprintf("{\"signature\": \"%s\", \"message\": \"%s\"}", signature, message)
	return payload
}

func joinObscuroGateway(url string) (string, error) {
	ogClient := lib.NewObscuroGatewayLibrary(url, url)
	err := ogClient.Join()
	if err != nil {
		return "", err
	}

	return ogClient.UserID(), nil
}
