package env

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ten-protocol/go-ten/go/wallet"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/integration/networktest/userwallet"

	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/integration/common/testlog"

	"github.com/ten-protocol/go-ten/integration"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/integration/networktest"
)

var _defaultFaucetAmount = big.NewInt(5_000_000_000_000_000)

type testnetConnector struct {
	seqRPCAddress         string
	validatorRPCAddresses []string
	faucetHTTPAddress     string
	l1RPCURL              string
	tenGatewayURL         string
	tenGatewayWSURL       string
	faucetWallet          userwallet.User
}

func newTestnetConnector(seqRPCAddr string, validatorRPCAddressses []string, faucetHTTPAddress string, l1WSURL string, tenGatewayURL string, tenGatewayWSURL string) *testnetConnector {
	return &testnetConnector{
		seqRPCAddress:         seqRPCAddr,
		validatorRPCAddresses: validatorRPCAddressses,
		faucetHTTPAddress:     faucetHTTPAddress,
		l1RPCURL:              l1WSURL,
		tenGatewayURL:         tenGatewayURL,
		tenGatewayWSURL:       tenGatewayWSURL,
	}
}

func newTestnetConnectorWithFaucetAccount(seqRPCAddr string, validatorRPCAddressses []string, faucetPK string, l1RPCAddress string, tenGatewayURL string) *testnetConnector {
	ecdsaKey, err := crypto.HexToECDSA(faucetPK)
	if err != nil {
		panic(err)
	}
	wal := wallet.NewInMemoryWalletFromPK(big.NewInt(integration.TenChainID), ecdsaKey, testlog.Logger())
	return &testnetConnector{
		seqRPCAddress:         seqRPCAddr,
		validatorRPCAddresses: validatorRPCAddressses,
		faucetWallet:          userwallet.NewUserWallet(wal, validatorRPCAddressses[0], testlog.Logger()),
		l1RPCURL:              l1RPCAddress,
		tenGatewayURL:         tenGatewayURL,
	}
}

func (t *testnetConnector) ChainID() int64 {
	return integration.TenChainID
}

func (t *testnetConnector) AllocateFaucetFunds(ctx context.Context, account gethcommon.Address) error {
	if t.faucetWallet != nil {
		// if we have a faucet wallet available, use it to send funds - if not we'll use the http faucet
		return t.AllocateFaucetFundsWithWallet(ctx, account)
	}
	payload := map[string]string{"address": account.Hex()}
	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.faucetHTTPAddress, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("unable to make http faucet request - %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing http faucet request - %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status of http faucet request, code=%d status=%s", resp.StatusCode, resp.Status)
	}
	return nil
}

func (t *testnetConnector) SequencerRPCAddress() string {
	return t.seqRPCAddress
}

func (t *testnetConnector) ValidatorRPCAddress(idx int) string {
	return t.validatorRPCAddresses[idx]
}

func (t *testnetConnector) NumValidators() int {
	return len(t.validatorRPCAddresses)
}

func (t *testnetConnector) GetL1Client() (ethadapter.EthClient, error) {
	client, err := ethadapter.NewEthClientFromURL(t.l1RPCURL, time.Minute, testlog.Logger())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (t *testnetConnector) GetSequencerNode() networktest.NodeOperator {
	panic("node operators cannot be accessed for testnets")
}

func (t *testnetConnector) GetValidatorNode(_ int) networktest.NodeOperator {
	panic("node operators cannot be accessed for testnets")
}

func (t *testnetConnector) AllocateFaucetFundsWithWallet(ctx context.Context, account gethcommon.Address) error {
	txHash, err := t.faucetWallet.SendFunds(ctx, account, _defaultFaucetAmount)
	if err != nil {
		return err
	}

	receipt, err := t.faucetWallet.AwaitReceipt(ctx, txHash)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("faucet transaction receipt status not successful - %v", receipt.Status)
	}
	return nil
}

func (t *testnetConnector) GetMCOwnerWallet() (wallet.Wallet, error) {
	return nil, errors.New("testnet connector environments cannot access the MC owner wallet")
}

func (t *testnetConnector) GetGatewayClient() (ethadapter.EthClient, error) {
	if t.tenGatewayURL == "" {
		return nil, errors.New("gateway client not set for this environment")
	}
	return ethadapter.NewEthClientFromURL(t.tenGatewayURL, time.Minute, testlog.Logger())
}

func (t *testnetConnector) GetGatewayURL() (string, error) {
	return t.tenGatewayURL, nil
}

func (t *testnetConnector) GetGatewayWSURL() (string, error) {
	return t.tenGatewayWSURL, nil
}
