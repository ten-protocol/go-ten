package env

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/integration/networktest/userwallet"

	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	"github.com/obscuronet/go-obscuro/integration"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/integration/networktest"
)

var _defaultFaucetAmount = big.NewInt(750_000_000_000_000)

type testnetConnector struct {
	seqRPCAddress         string
	validatorRPCAddresses []string
	faucetHTTPAddress     string
	l1RPCAddress          string
	faucetWallet          *userwallet.UserWallet
}

func NewTestnetConnector(seqRPCAddr string, validatorRPCAddressses []string, faucetHTTPAddress string, l1RPCAddress string) networktest.NetworkConnector {
	return &testnetConnector{
		seqRPCAddress:         seqRPCAddr,
		validatorRPCAddresses: validatorRPCAddressses,
		faucetHTTPAddress:     faucetHTTPAddress,
		l1RPCAddress:          l1RPCAddress,
	}
}

func NewTestnetConnectorWithFaucetAccount(seqRPCAddr string, validatorRPCAddressses []string, faucetPK string, l1RPCAddress string) networktest.NetworkConnector {
	ecdsaKey, err := crypto.HexToECDSA(faucetPK)
	if err != nil {
		panic(err)
	}
	return &testnetConnector{
		seqRPCAddress:         seqRPCAddr,
		validatorRPCAddresses: validatorRPCAddressses,
		faucetWallet:          userwallet.NewUserWallet(ecdsaKey, validatorRPCAddressses[0], testlog.Logger(), userwallet.WithChainID(big.NewInt(integration.ObscuroChainID))),
		l1RPCAddress:          l1RPCAddress,
	}
}

func (t *testnetConnector) ChainID() int64 {
	return integration.ObscuroChainID
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
	client, err := ethadapter.NewEthClientFromAddress(t.l1RPCAddress, time.Minute, gethcommon.Address{}, testlog.Logger())
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
