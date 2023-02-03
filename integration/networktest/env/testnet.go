package env

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/obscuronet/go-obscuro/integration"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/integration/networktest"
)

type testnetConnector struct {
	seqRPCAddress         string
	validatorRPCAddresses []string
	faucetHTTPAddress     string
}

func NewTestnetConnector(seqRPCAddr string, validatorRPCAddressses []string, faucetHTTPAddress string) networktest.NetworkConnector {
	return &testnetConnector{
		seqRPCAddress:         seqRPCAddr,
		validatorRPCAddresses: validatorRPCAddressses,
		faucetHTTPAddress:     faucetHTTPAddress,
	}
}

func (t *testnetConnector) ChainID() int64 {
	return integration.ObscuroChainID
}

func (t *testnetConnector) AllocateFaucetFunds(account gethcommon.Address) error {
	req := map[string]string{"address": account.Hex()}
	jsonPayload, _ := json.Marshal(req)
	resp, err := http.Post(t.faucetHTTPAddress, "application/json", bytes.NewBuffer(jsonPayload)) //nolint:noctx
	if err != nil {
		return fmt.Errorf("unable to make http faucet request - %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
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

func (t *testnetConnector) GetSequencerNode() networktest.NodeOperator {
	panic("node operators cannot be accessed for testnets")
}

func (t *testnetConnector) GetValidatorNode(_ int) networktest.NodeOperator {
	panic("node operators cannot be accessed for testnets")
}
