package l1

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto/kzg4844"

	"github.com/ten-protocol/go-ten/go/common/gethutil"

	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	"github.com/ten-protocol/go-ten/go/host/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/host"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"
	"github.com/ten-protocol/go-ten/go/wallet"
)

type Publisher struct {
	hostData        host.Identity
	hostWallet      wallet.Wallet // Wallet used to issue ethereum transactions
	ethClient       ethadapter.EthClient
	mgmtContractLib mgmtcontractlib.MgmtContractLib // Library to handle Management Contract lib operations
	storage         storage.Storage
	blobResolver    BlobResolver

	// cached map of important contract addresses (updated when we see a SetImportantContractsTx)
	importantContractAddresses map[string]gethcommon.Address
	// lock for the important contract addresses map
	importantAddressesMutex sync.RWMutex

	repository host.L1DataService
	logger     gethlog.Logger

	hostStopper *stopcontrol.StopControl

	maxWaitForL1Receipt       time.Duration
	retryIntervalForL1Receipt time.Duration

	// we only allow one transaction in-flight at a time to avoid nonce conflicts
	// We also have a context to cancel the tx if host stops
	sendingLock      sync.Mutex
	sendingContext   context.Context
	sendingCtxCancel context.CancelFunc
}

func NewL1Publisher(
	hostData host.Identity,
	hostWallet wallet.Wallet,
	client ethadapter.EthClient,
	mgmtContract mgmtcontractlib.MgmtContractLib,
	repository host.L1DataService,
	blobResolver BlobResolver,
	hostStopper *stopcontrol.StopControl,
	logger gethlog.Logger,
	maxWaitForL1Receipt time.Duration,
	retryIntervalForL1Receipt time.Duration,
	storage storage.Storage,
) *Publisher {
	sendingCtx, cancelSendingCtx := context.WithCancel(context.Background())
	return &Publisher{
		hostData:                  hostData,
		hostWallet:                hostWallet,
		ethClient:                 client,
		mgmtContractLib:           mgmtContract,
		repository:                repository,
		blobResolver:              blobResolver,
		hostStopper:               hostStopper,
		logger:                    logger,
		maxWaitForL1Receipt:       maxWaitForL1Receipt,
		retryIntervalForL1Receipt: retryIntervalForL1Receipt,
		storage:                   storage,

		importantContractAddresses: map[string]gethcommon.Address{},
		importantAddressesMutex:    sync.RWMutex{},

		sendingLock:      sync.Mutex{},
		sendingContext:   sendingCtx,
		sendingCtxCancel: cancelSendingCtx,
	}
}

func (p *Publisher) Start() error {
	go func() {
		// Do an initial read of important contract addresses when service starts up
		err := p.ResyncImportantContracts()
		if err != nil {
			p.logger.Error("Could not load important contract addresses", log.ErrKey, err)
		}
	}()
	return nil
}

func (p *Publisher) Stop() error {
	p.sendingCtxCancel()
	return nil
}

func (p *Publisher) HealthStatus(context.Context) host.HealthStatus {
	// todo (@matt) do proper health status based on failed transactions or something
	errMsg := ""
	if p.hostStopper.IsStopping() {
		errMsg = "not running"
	}
	return &host.BasicErrHealthStatus{ErrMsg: errMsg}
}

func (p *Publisher) InitializeSecret(attestation *common.AttestationReport, encSecret common.EncryptedSharedEnclaveSecret) error {
	encodedAttestation, err := common.EncodeAttestation(attestation)
	if err != nil {
		return errors.Wrap(err, "could not encode attestation")
	}
	l1tx := &common.L1InitializeSecretTx{
		EnclaveID:     &attestation.EnclaveID,
		Attestation:   encodedAttestation,
		InitialSecret: encSecret,
	}
	initialiseSecretTx, err := p.mgmtContractLib.CreateInitializeSecret(l1tx)
	if err != nil {
		return err
	}
	// we block here until we confirm a successful receipt. It is important this is published before the initial rollup.
	return p.publishTransaction(initialiseSecretTx)
}

func (p *Publisher) RequestSecret(attestation *common.AttestationReport) (gethcommon.Hash, error) {
	encodedAttestation, err := common.EncodeAttestation(attestation)
	if err != nil {
		return gethutil.EmptyHash, errors.Wrap(err, "could not encode attestation")
	}
	l1tx := &common.L1RequestSecretTx{
		Attestation: encodedAttestation,
	}
	// record the L1 head height before we submit the secret request, so we know which block to watch from
	l1Head, err := p.ethClient.FetchHeadBlock()
	if err != nil {
		err = p.ethClient.ReconnectIfClosed()
		if err != nil {
			panic(errors.Wrap(err, "could not reconnect to eth client"))
		}
		l1Head, err = p.ethClient.FetchHeadBlock()
		if err != nil {
			panic(errors.Wrap(err, "could not fetch head block"))
		}
	}
	requestSecretTx, err := p.mgmtContractLib.CreateRequestSecret(l1tx)
	if err != nil {
		return gethutil.EmptyHash, err
	}

	// we wait until the secret req transaction has succeeded before we start polling for the secret
	err = p.publishTransaction(requestSecretTx)
	if err != nil {
		return gethutil.EmptyHash, err
	}

	return l1Head.Hash(), nil
}

func (p *Publisher) PublishSecretResponse(secretResponse *common.ProducedSecretResponse) error {
	l1tx := &common.L1RespondSecretTx{
		Secret:      secretResponse.Secret,
		RequesterID: secretResponse.RequesterID,
		AttesterID:  secretResponse.AttesterID,
	}
	// todo (#1624) - l1tx.Sign(a.attestationPubKey) doesn't matter as the waitSecret will process a tx that was reverted
	respondSecretTx, err := p.mgmtContractLib.CreateRespondSecret(l1tx, false)
	if err != nil {
		return err
	}
	p.logger.Info("Broadcasting secret response L1 tx.", "requester", secretResponse.RequesterID)

	// fire-and-forget (track the receipt asynchronously)
	go func() {
		err := p.publishTransaction(respondSecretTx)
		if err != nil {
			p.logger.Error("Could not broadcast secret response L1 tx", log.ErrKey, err)
		}
	}()

	return nil
}

// FindSecretResponseTx will attempt to decode the transactions passed in
func (p *Publisher) FindSecretResponseTx(processed []*common.L1TxData) []*common.L1RespondSecretTx {
	secretRespTxs := make([]*common.L1RespondSecretTx, 0)

	for _, tx := range processed {
		t, err := p.mgmtContractLib.DecodeTx(tx.Transaction)
		if err != nil {
			p.logger.Error("Could not decode transaction", log.ErrKey, err)
			continue
		}
		if t == nil {
			continue
		}
		if scrtTx, ok := t.(*common.L1RespondSecretTx); ok {
			secretRespTxs = append(secretRespTxs, scrtTx)
			continue
		}
	}
	return secretRespTxs
}

func (p *Publisher) FetchLatestSeqNo() (*big.Int, error) {
	return p.ethClient.FetchLastBatchSeqNo(*p.mgmtContractLib.GetContractAddr())
}

func (p *Publisher) PublishBlob(producedRollup *common.ExtRollup, blobs []*kzg4844.Blob) {
	encRollup, err := common.EncodeRollup(producedRollup)
	if err != nil {
		p.logger.Crit("could not encode rollup.", log.ErrKey, err)
	}
	tx := &common.L1RollupTx{
		Rollup: encRollup,
	}
	p.logger.Info("Publishing rollup", "size", len(encRollup)/1024, log.RollupHashKey, producedRollup.Hash())

	if p.logger.Enabled(context.Background(), gethlog.LevelTrace) {
		var headerLog string
		header, err := json.MarshalIndent(producedRollup.Header, "", "   ")
		if err != nil {
			headerLog = err.Error()
		} else {
			headerLog = string(header)
		}

		p.logger.Trace("Sending transaction to publish rollup", "rollup_header", headerLog, log.RollupHashKey, producedRollup.Header.Hash(), "batches_len", len(producedRollup.BatchPayloads))
	}

	rollupBlobTx, err := p.mgmtContractLib.PopulateAddRollup(tx, blobs)
	if err != nil {
		p.logger.Error("Could not create rollup blobs", log.RollupHashKey, producedRollup.Hash(), log.ErrKey, err)
	}

	// wait for the next block after the block that the rollup is bound to
	err = p.waitForBlockAfter(producedRollup.Header.CompressionL1Number.Uint64())
	if err != nil {
		p.logger.Error("Failed waiting for block after rollup binding block number",
			"compression_block", producedRollup.Header.CompressionL1Number,
			log.ErrKey, err)
	}

	err = p.publishTransaction(rollupBlobTx)
	if err != nil {
		p.logger.Error("Could not issue rollup tx", log.RollupHashKey, producedRollup.Hash(), log.ErrKey, err)
	} else {
		p.logger.Info("Rollup included in L1", log.RollupHashKey, producedRollup.Hash())
	}
	// TODO publish rollup to archive service if not already done
}

func (p *Publisher) PublishCrossChainBundle(_ *common.ExtCrossChainBundle, _ *big.Int, _ gethcommon.Hash) error {
	return nil
}

func (p *Publisher) GetImportantContracts() map[string]gethcommon.Address {
	p.importantAddressesMutex.RLock()
	defer p.importantAddressesMutex.RUnlock()
	return p.importantContractAddresses
}

// ResyncImportantContracts will fetch the latest important contracts from the management contract and update the cached map
// Note: this should be run in a goroutine as it makes L1 transactions in series and will block.
// Cache is not overwritten until it completes.
func (p *Publisher) ResyncImportantContracts() error {
	getKeysCallMsg, err := p.mgmtContractLib.GetImportantContractKeysMsg()
	if err != nil {
		return fmt.Errorf("could not build callMsg for important contracts: %w", err)
	}
	keysResp, err := p.ethClient.CallContract(getKeysCallMsg)
	if err != nil {
		return fmt.Errorf("could not fetch important contracts: %w", err)
	}

	importantContracts, err := p.mgmtContractLib.DecodeImportantContractKeysResponse(keysResp)
	if err != nil {
		return fmt.Errorf("could not decode important contracts resp: %w", err)
	}

	contractsMap := make(map[string]gethcommon.Address)

	for _, contract := range importantContracts {
		getAddressCallMsg, err := p.mgmtContractLib.GetImportantAddressCallMsg(contract)
		if err != nil {
			return fmt.Errorf("could not build callMsg for important contract=%s: %w", contract, err)
		}
		addrResp, err := p.ethClient.CallContract(getAddressCallMsg)
		if err != nil {
			return fmt.Errorf("could not fetch important contract=%s: %w", contract, err)
		}
		contractAddress, err := p.mgmtContractLib.DecodeImportantAddressResponse(addrResp)
		if err != nil {
			return fmt.Errorf("could not decode important contract=%s resp: %w", contract, err)
		}
		contractsMap[contract] = contractAddress
	}

	p.importantAddressesMutex.Lock()
	defer p.importantAddressesMutex.Unlock()
	p.importantContractAddresses = contractsMap

	return nil
}

// publishTransaction will keep trying unless the L1 seems to be unavailable or the tx is otherwise rejected
// this method is guarded by a lock to ensure that only one transaction is attempted at a time to avoid nonce conflicts
// todo (@matt) this method should take a context so we can try to cancel if the tx is no longer required
func (p *Publisher) publishTransaction(tx types.TxData) error {
	// this log message seems superfluous but is useful to debug deadlock issues, we expect 'Host issuing l1 tx' soon
	// after unless we're stuck blocking.
	p.logger.Info("Host preparing to issue L1 tx")

	p.sendingLock.Lock()
	defer p.sendingLock.Unlock()

	retries := -1

	nonce, err := p.ethClient.Nonce(p.hostWallet.Address())
	if err != nil {
		return fmt.Errorf("could not get nonce for L1 tx: %w", err)
	}

	// while the publisher service is still alive we keep trying to get the transaction into the L1
	for !p.hostStopper.IsStopping() {
		retries++ // count each attempt so we can increase gas price

		// update the tx gas price before each attempt
		tx, err := ethadapter.SetTxGasPrice(p.sendingContext, p.ethClient, tx, p.hostWallet.Address(), nonce, retries)
		if err != nil {
			return errors.Wrap(err, "could not estimate gas/gas price for L1 tx")
		}

		signedTx, err := p.hostWallet.SignTransaction(tx)
		if err != nil {
			return errors.Wrap(err, "could not sign L1 tx")
		}
		p.logger.Info("Host issuing L1 tx", log.TxKey, signedTx.Hash(), "size", signedTx.Size()/1024, "retries", retries)
		err = p.ethClient.SendTransaction(signedTx)
		if err != nil {
			return errors.Wrap(err, "could not broadcast L1 tx")
		}
		p.logger.Info("Successfully submitted tx to L1", "txHash", signedTx.Hash())

		var receipt *types.Receipt
		// retry until receipt is found or context is canceled
		err = retry.Do(
			func() error {
				if p.hostStopper.IsStopping() {
					return retry.FailFast(errors.New("host is stopping or context canceled"))
				}
				receipt, err = p.ethClient.TransactionReceipt(signedTx.Hash())
				if err != nil {
					return fmt.Errorf("could not get receipt publishing tx for L1 tx=%s: %w", signedTx.Hash(), err)
				}
				return err
			},
			retry.NewTimeoutStrategy(p.maxWaitForL1Receipt, p.retryIntervalForL1Receipt),
		)
		if err != nil {
			p.logger.Info("Receipt not found for transaction, we will re-attempt", log.ErrKey, err)
			continue // try again with updated gas price
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			return fmt.Errorf("unsuccessful receipt found for published L1 transaction, status=%d", receipt.Status)
		}

		p.logger.Debug("L1 transaction successful receipt found.", log.TxKey, signedTx.Hash(),
			log.BlockHeightKey, receipt.BlockNumber, log.BlockHashKey, receipt.BlockHash)
		break
	}
	return nil
}

// waitForBlockAfter waits until the current block number is greater than the target block number
func (p *Publisher) waitForBlockAfter(targetBlock uint64) error {
	err := retry.Do(
		func() error {
			if p.hostStopper.IsStopping() {
				return retry.FailFast(errors.New("host is stopping"))
			}

			currentBlock, err := p.ethClient.BlockNumber()
			if err != nil {
				return fmt.Errorf("failed to get current block number: %w", err)
			}

			if currentBlock <= targetBlock {
				return fmt.Errorf("waiting for block after %d (current: %d)", targetBlock, currentBlock)
			}

			return nil
		},
		retry.NewTimeoutStrategy(p.maxWaitForL1Receipt, p.retryIntervalForL1Receipt),
	)

	if err != nil {
		return fmt.Errorf("timeout waiting for block after %d: %w", targetBlock, err)
	}

	return nil
}
