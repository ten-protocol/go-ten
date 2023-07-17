package l1

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/contracts/generated/ManagementContract"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/pkg/errors"
)

const (
	// Attempts to broadcast the rollup transaction to the L1. Worst-case, equates to 7 seconds, plus time per request.
	l1TxTriesRollup = 3
	// Attempts to send secret initialisation, request or response transactions to the L1. Worst-case, equates to 63 seconds, plus time per request.
	l1TxTriesSecret = 7

	// todo - these values have to be configurable
	maxWaitForL1Receipt       = 100 * time.Second
	retryIntervalForL1Receipt = 10 * time.Second
	maxWaitForSecretResponse  = 120 * time.Second
)

type Publisher struct {
	hostData        host.Identity
	hostWallet      wallet.Wallet // Wallet used to issue ethereum transactions
	ethClient       ethadapter.EthClient
	mgmtContractLib mgmtcontractlib.MgmtContractLib // Library to handle Management Contract lib operations

	repository host.L1BlockRepository
	logger     gethlog.Logger

	running atomic.Bool
}

func NewL1Publisher(hostData host.Identity, hostWallet wallet.Wallet, client ethadapter.EthClient, mgmtContract mgmtcontractlib.MgmtContractLib, repository host.L1BlockRepository, logger gethlog.Logger) *Publisher {
	return &Publisher{
		hostData:        hostData,
		hostWallet:      hostWallet,
		ethClient:       client,
		mgmtContractLib: mgmtContract,
		repository:      repository,
		logger:          logger,
	}
}

func (p *Publisher) Start() error {
	p.running.Store(true)
	return nil
}

func (p *Publisher) Stop() error {
	p.running.Store(false)
	return nil
}

func (p *Publisher) HealthStatus() host.HealthStatus {
	// todo (@matt) do proper health status based on failed transactions or something
	errMsg := ""
	if !p.running.Load() {
		errMsg = "not running"
	}
	return &host.BasicErrHealthStatus{ErrMsg: errMsg}
}

func (p *Publisher) InitializeSecret(attestation *common.AttestationReport, encSecret common.EncryptedSharedEnclaveSecret) error {
	encodedAttestation, err := common.EncodeAttestation(attestation)
	if err != nil {
		return errors.Wrap(err, "could not encode attestation")
	}
	l1tx := &ethadapter.L1InitializeSecretTx{
		AggregatorID:  &p.hostData.ID,
		Attestation:   encodedAttestation,
		InitialSecret: encSecret,
		HostAddress:   p.hostData.P2PPublicAddress,
	}
	initialiseSecretTx := p.mgmtContractLib.CreateInitializeSecret(l1tx, p.hostWallet.GetNonceAndIncrement())
	initialiseSecretTx, err = p.ethClient.EstimateGasAndGasPrice(initialiseSecretTx, p.hostWallet.Address())
	if err != nil {
		p.hostWallet.SetNonce(p.hostWallet.GetNonce() - 1)
		return err
	}
	// we block here until we confirm a successful receipt. It is important this is published before the initial rollup.
	err = p.signAndBroadcastL1Tx(initialiseSecretTx, l1TxTriesSecret, true)
	if err != nil {
		return err
	}
	return nil
}

func (p *Publisher) RequestSecret(attestation *common.AttestationReport) (gethcommon.Hash, error) {
	encodedAttestation, err := common.EncodeAttestation(attestation)
	if err != nil {
		return gethcommon.Hash{}, errors.Wrap(err, "could not encode attestation")
	}
	l1tx := &ethadapter.L1RequestSecretTx{
		Attestation: encodedAttestation,
	}
	// record the L1 head height before we submit the secret request, so we know which block to watch from
	l1Head, err := p.ethClient.FetchHeadBlock()
	if err != nil {
		err = p.ethClient.Reconnect()
		if err != nil {
			panic(errors.Wrap(err, "could not reconnect to eth client"))
		}
		l1Head, err = p.ethClient.FetchHeadBlock()
		if err != nil {
			panic(errors.Wrap(err, "could not fetch head block"))
		}
	}
	requestSecretTx := p.mgmtContractLib.CreateRequestSecret(l1tx, p.hostWallet.GetNonceAndIncrement())
	requestSecretTx, err = p.ethClient.EstimateGasAndGasPrice(requestSecretTx, p.hostWallet.Address())
	if err != nil {
		p.hostWallet.SetNonce(p.hostWallet.GetNonce() - 1)
		return gethcommon.Hash{}, err
	}
	// we wait until the secret req transaction has succeeded before we start polling for the secret
	err = p.signAndBroadcastL1Tx(requestSecretTx, l1TxTriesSecret, true)
	if err != nil {
		return gethcommon.Hash{}, err
	}

	return l1Head.Hash(), nil
}

func (p *Publisher) PublishSecretResponse(secretResponse *common.ProducedSecretResponse) error {
	l1tx := &ethadapter.L1RespondSecretTx{
		Secret:      secretResponse.Secret,
		RequesterID: secretResponse.RequesterID,
		AttesterID:  p.hostData.ID,
		HostAddress: secretResponse.HostAddress,
	}
	// todo (#1624) - l1tx.Sign(a.attestationPubKey) doesn't matter as the waitSecret will process a tx that was reverted
	respondSecretTx := p.mgmtContractLib.CreateRespondSecret(l1tx, p.hostWallet.GetNonceAndIncrement(), false)
	respondSecretTx, err := p.ethClient.EstimateGasAndGasPrice(respondSecretTx, p.hostWallet.Address())
	if err != nil {
		p.hostWallet.SetNonce(p.hostWallet.GetNonce() - 1)
		return err
	}
	p.logger.Info("Broadcasting secret response L1 tx.", "requester", secretResponse.RequesterID)
	// fire-and-forget (track the receipt asynchronously)
	err = p.signAndBroadcastL1Tx(respondSecretTx, l1TxTriesSecret, false)
	if err != nil {
		return errors.Wrap(err, "could not broadcast secret response L1 tx")
	}
	return nil
}

func (p *Publisher) ExtractSecretResponses(block *types.Block) []*ethadapter.L1RespondSecretTx {
	var secretRespTxs []*ethadapter.L1RespondSecretTx
	for _, tx := range block.Transactions() {
		t := p.mgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}
		if scrtTx, ok := t.(*ethadapter.L1RespondSecretTx); ok {
			secretRespTxs = append(secretRespTxs, scrtTx)
		}
	}
	return secretRespTxs
}

func (p *Publisher) FetchLatestSeqNo() (*big.Int, error) {
	contract, err := ManagementContract.NewManagementContract(*p.mgmtContractLib.GetContractAddr(), p.ethClient.EthClient())
	if err != nil {
		return nil, err
	}

	batchNo, err := contract.LastBatchSeqNo(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	return batchNo, nil
}

func (p *Publisher) PublishRollup(producedRollup *common.ExtRollup) {
	encRollup, err := common.EncodeRollup(producedRollup)
	if err != nil {
		p.logger.Crit("could not encode rollup.", log.ErrKey, err)
	}
	tx := &ethadapter.L1RollupTx{
		Rollup: encRollup,
	}
	p.logger.Info("Publishing rollup", "size", len(encRollup)/1024, log.RollupHashKey, producedRollup.Hash())

	p.logger.Trace("Sending transaction to publish rollup", "rollup_header",
		gethlog.Lazy{Fn: func() string {
			header, err := json.MarshalIndent(producedRollup.Header, "", "   ")
			if err != nil {
				return err.Error()
			}

			return string(header)
		}}, "rollup_hash", producedRollup.Header.Hash().Hex(), "batches_len", len(producedRollup.BatchPayloads))

	rollupTx := p.mgmtContractLib.CreateRollup(tx, p.hostWallet.GetNonceAndIncrement())
	rollupTx, err = p.ethClient.EstimateGasAndGasPrice(rollupTx, p.hostWallet.Address())
	if err != nil {
		// todo (#1624) - make rollup submission a separate workflow (design and implement the flow etc)
		p.hostWallet.SetNonce(p.hostWallet.GetNonce() - 1)
		p.logger.Error("could not estimate rollup tx", log.ErrKey, err)
		return
	}

	err = p.signAndBroadcastL1Tx(rollupTx, l1TxTriesRollup, true)
	if err != nil {
		p.logger.Error("could not issue rollup tx", log.ErrKey, err)
	} else {
		p.logger.Info("Rollup included in L1", "hash", producedRollup.Hash())
	}
}

func (p *Publisher) FetchLatestPeersList() ([]string, error) {
	msg, err := p.mgmtContractLib.GetHostAddresses()
	if err != nil {
		return nil, err
	}
	response, err := p.ethClient.CallContract(msg)
	if err != nil {
		return nil, err
	}
	decodedResponse, err := p.mgmtContractLib.DecodeCallResponse(response)
	if err != nil {
		return nil, err
	}
	hostAddresses := decodedResponse[0]

	// We remove any duplicate addresses and our own address from the retrieved peer list
	var filteredHostAddresses []string
	uniqueHostKeys := make(map[string]bool) // map to track addresses we've seen already
	for _, hostAddress := range hostAddresses {
		// We exclude our own address.
		if hostAddress == p.hostData.P2PPublicAddress {
			continue
		}
		if _, found := uniqueHostKeys[hostAddress]; !found {
			uniqueHostKeys[hostAddress] = true
			filteredHostAddresses = append(filteredHostAddresses, hostAddress)
		}
	}

	return filteredHostAddresses, nil
}

// `tries` is the number of times to attempt broadcasting the transaction.
// if awaitReceipt is true then this method will block and synchronously wait to check the receipt, otherwise it is fire
// and forget and the receipt tracking will happen in a separate go-routine
func (p *Publisher) signAndBroadcastL1Tx(tx types.TxData, tries uint64, awaitReceipt bool) error {
	var err error
	tx, err = p.ethClient.EstimateGasAndGasPrice(tx, p.hostWallet.Address())
	if err != nil {
		return errors.Wrap(err, "could not estimate gas/gas price for L1 tx")
	}

	signedTx, err := p.hostWallet.SignTransaction(tx)
	if err != nil {
		return err
	}

	p.logger.Info("Host issuing l1 tx", log.TxKey, signedTx.Hash(), "size", signedTx.Size()/1024)

	err = retry.Do(func() error {
		return p.ethClient.SendTransaction(signedTx)
	}, retry.NewDoublingBackoffStrategy(time.Second, tries)) // doubling retry wait (3 tries = 7sec, 7 tries = 63sec)
	if err != nil {
		return errors.Wrapf(err, "could not broadcast L1 tx after %d tries", tries)
	}
	p.logger.Info("Successfully submitted tx to L1", "txHash", signedTx.Hash())

	if awaitReceipt {
		// block until receipt is found and then return
		return p.waitForReceipt(signedTx.Hash())
	}

	// else just watch for receipt asynchronously and log if it fails
	go func() {
		// todo (#1624) - consider how to handle the various ways that L1 transactions could fail to improve node operator QoL
		err = p.waitForReceipt(signedTx.Hash())
		if err != nil {
			p.logger.Error("L1 transaction failed", log.ErrKey, err)
		}
	}()

	return nil
}

func (p *Publisher) waitForReceipt(txHash common.TxHash) error {
	var receipt *types.Receipt
	var err error
	err = retry.Do(
		func() error {
			receipt, err = p.ethClient.TransactionReceipt(txHash)
			if err != nil {
				// adds more info on the error
				return errors.Wrapf(err, "could not get receipt for L1 tx=%s", txHash)
			}
			return err
		},
		retry.NewTimeoutStrategy(maxWaitForL1Receipt, retryIntervalForL1Receipt),
	)
	if err != nil {
		return errors.Wrap(err, "receipt for L1 tx not found despite successful broadcast")
	}

	if err == nil && receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("unsuccessful receipt found for published L1 transaction, status=%d", receipt.Status)
	}
	p.logger.Debug("L1 transaction receipt found.", log.TxKey, txHash, log.BlockHeightKey, receipt.BlockNumber, log.BlockHashKey, receipt.BlockHash)
	return nil
}
