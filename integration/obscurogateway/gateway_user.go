package faucet

import (
	"context"
	"fmt"

	"github.com/obscuronet/go-obscuro/tools/walletextension/lib"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/obscuronet/go-obscuro/go/wallet"
)

// GatewayUser TODO (@ziga) refactor GatewayUser and integrate it with OGlib.
// GatewayUser is a struct that includes everything a gateway user has and uses (userID, wallets, http & ws addresses and client )
type GatewayUser struct {
	Wallets    []wallet.Wallet
	HTTPClient *ethclient.Client
	WSClient   *ethclient.Client
	ogClient   *lib.OGLib
}

func NewUser(wallets []wallet.Wallet, serverAddressHTTP string, serverAddressWS string) (*GatewayUser, error) {
	ogClient := lib.NewObscuroGatewayLibrary(serverAddressHTTP, serverAddressWS)

	// automatically join
	ogClient.Join()

	// create clients
	httpClient, err := ethclient.Dial(serverAddressHTTP + "/v1/" + "?u=" + ogClient.UserID())
	if err != nil {
		return nil, err
	}
	wsClient, err := ethclient.Dial(serverAddressWS + "/v1/" + "?u=" + ogClient.UserID())
	if err != nil {
		return nil, err
	}

	return &GatewayUser{
		Wallets:    wallets,
		HTTPClient: httpClient,
		WSClient:   wsClient,
		ogClient:   ogClient,
	}, nil
}

func (u GatewayUser) RegisterAccounts() error {
	for _, w := range u.Wallets {
		u.ogClient.RegisterAccount(w.PrivateKey(), w.Address())
		fmt.Printf("Successfully registered address %s for user: %s.\n", w.Address().Hex(), u.ogClient.UserID())
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
