package env

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	"github.com/obscuronet/go-obscuro/integration"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/integration/networktest"
)

type testnetConnector struct {
	seqRPCAddress         string
	validatorRPCAddresses []string
	faucetHTTPAddress     string
	l1Host                string
	l1HttpPort            uint
}

func NewTestnetConnector(seqRPCAddr string, validatorRPCAddressses []string, faucetHTTPAddress string, l1Host string, l1Port uint) networktest.NetworkConnector {
	return &testnetConnector{
		seqRPCAddress:         seqRPCAddr,
		validatorRPCAddresses: validatorRPCAddressses,
		faucetHTTPAddress:     faucetHTTPAddress,
		l1Host:                l1Host,
		l1HttpPort:            l1Port,
	}
}

func (t *testnetConnector) ChainID() int64 {
	return integration.ObscuroChainID
}

func (t *testnetConnector) AllocateFaucetFunds(ctx context.Context, account gethcommon.Address) error {
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
	client, err := ethadapter.NewEthClient(t.l1Host, t.l1HttpPort, time.Minute, gethcommon.Address{}, testlog.Logger())
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
