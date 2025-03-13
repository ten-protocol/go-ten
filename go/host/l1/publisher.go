package l1

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"

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
	"github.com/ten-protocol/go-ten/go/wallet"
)

type Publisher struct {
	hostData         host.Identity
	hostWallet       wallet.Wallet // Wallet used to issue ethereum transactions
	ethClient        ethadapter.EthClient
	contractRegistry contractlib.ContractRegistryLib // Library to handle Management Contract lib operations
	storage          storage.Storage
	blobResolver     BlobResolver

	// cached map of important contract addresses (updated when we see a SetImportantContractsTx)
	importantContractAddresses map[string]gethcommon.Address
	// lock for the important contract addresses map
	importantAddressesMutex sync.RWMutex
	importantAddresses      *common.NetworkConfigAddresses

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
	contractRegistry contractlib.ContractRegistryLib,
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
		contractRegistry:          contractRegistry,
		repository:                repository,
		blobResolver:              blobResolver,
		hostStopper:               hostStopper,
		logger:                    logger,
		maxWaitForL1Receipt:       maxWaitForL1Receipt,
		retryIntervalForL1Receipt: retryIntervalForL1Receipt,
		storage:                   storage,

		importantAddressesMutex: sync.RWMutex{},
		importantAddresses:      &common.NetworkConfigAddresses{},

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
	initialiseSecretTx, err := p.contractRegistry.NetworkEnclaveLib().CreateInitializeSecret(l1tx)
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
	requestSecretTx, err := p.contractRegistry.NetworkEnclaveLib().CreateRequestSecret(l1tx)
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
	respondSecretTx, err := p.contractRegistry.NetworkEnclaveLib().CreateRespondSecret(l1tx, false)
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
		t, err := p.contractRegistry.NetworkEnclaveLib().DecodeTx(tx.Transaction)
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
	return p.ethClient.FetchLastBatchSeqNo(*p.contractRegistry.RollupLib().GetContractAddr())
}

func (p *Publisher) PublishBlob(result common.CreateRollupResult) {
	// Decode the rollup from the blobs
	rollupData, err := ethadapter.DecodeBlobs(result.Blobs)
	if err != nil {
		p.logger.Crit("could not decode rollup from blob.", log.ErrKey, err)
	}

	// Decode the rollup to get header info for logging
	extRollup, err := common.DecodeRollup(rollupData)
	if err != nil {
		p.logger.Crit("could not decode rollup.", log.ErrKey, err)
	}

	// Check if the signature is valid
	// This depends on how your signature verification works
	p.logger.Info("Signature validation", "is_valid")

	tx := &common.L1RollupTx{
		Rollup: rollupData,
	}
	p.logger.Info("Publishing rollup", "size", len(rollupData)/1024, log.RollupHashKey, extRollup.Hash())

	if p.logger.Enabled(context.Background(), gethlog.LevelTrace) {
		var headerLog string
		header, err := json.MarshalIndent(extRollup.Header, "", "   ")
		if err != nil {
			headerLog = err.Error()
		} else {
			headerLog = string(header)
		}

		p.logger.Trace("Sending transaction to publish rollup", "rollup_header", headerLog, log.RollupHashKey, extRollup.Header.Hash(), "batches_len", len(extRollup.BatchPayloads))
	}

	rollupBlobTx, err := p.contractRegistry.RollupLib().PopulateAddRollup(tx, result.Blobs, result.Signature)
	if err != nil {
		p.logger.Error("Could not create rollup blobs", log.RollupHashKey, extRollup.Hash(), log.ErrKey, err)
	}

	rollupBlockNum := extRollup.Header.CompressionL1Number
	// wait for the next block after the block that the rollup is bound to
	err = p.waitForBlockAfter(rollupBlockNum.Uint64())
	if err != nil {
		p.logger.Error("Failed waiting for block after rollup binding block number",
			"compression_block", rollupBlockNum,
			log.ErrKey, err)
	}

	err = p.publishTransaction(rollupBlobTx)
	if err != nil {
		var maxRetriesErr *MaxRetriesError
		if errors.As(err, &maxRetriesErr) {
			p.handleMaxRetriesFailure(maxRetriesErr, extRollup)
			return
		}

		p.logger.Error("Could not issue rollup tx",
			log.RollupHashKey, extRollup.Hash(),
			log.ErrKey, err)
	} else {
		p.logger.Info("Rollup included in L1", log.RollupHashKey, extRollup.Hash())
	}
}

func (p *Publisher) handleMaxRetriesFailure(err *MaxRetriesError, rollup *common.ExtRollup) {
	p.logger.Error("Blob transaction failed after max retries",
		"nonce", err.BlobTx.Nonce,
		log.RollupHashKey, rollup.Hash(),
		log.ErrKey, err.Error())
	// TODO store failed rollup details so we can easily remediate? ie send new tx with the same nonce
}

func (p *Publisher) PublishCrossChainBundle(_ *common.ExtCrossChainBundle, _ *big.Int, _ gethcommon.Hash) error {
	return nil
}

func (p *Publisher) GetImportantContracts() *common.NetworkConfigAddresses {
	p.importantAddressesMutex.RLock()
	defer p.importantAddressesMutex.RUnlock()
	return p.importantAddresses
}

// ResyncImportantContracts will fetch the latest contract addresses from the network config contract and update the cached map
// Note: this should be run in a goroutine as it makes L1 transactions in series and will block.
// Cache is not overwritten until it completes.
func (p *Publisher) ResyncImportantContracts() error {
	addresses, err := p.contractRegistry.NetworkConfigLib().GetContractAddresses()
	if err != nil {
		return fmt.Errorf("could not get contract addresses: %w", err)
	}

	p.importantAddressesMutex.Lock()
	defer p.importantAddressesMutex.Unlock()
	p.importantAddresses = addresses

	return nil
}

// publishTransaction will keep trying unless the L1 seems to be unavailable or the tx is otherwise rejected
// this method is guarded by a lock to ensure that only one transaction is attempted at a time to avoid nonce conflicts
// todo (@matt) this method should take a context so we can try to cancel if the tx is no longer required
func (p *Publisher) publishTransaction(tx types.TxData) error {
	p.sendingLock.Lock()
	defer p.sendingLock.Unlock()

	nonce, err := p.ethClient.Nonce(p.hostWallet.Address())
	if err != nil {
		return fmt.Errorf("could not get nonce for L1 tx: %w", err)
	}

	if _, ok := tx.(*types.BlobTx); ok {
		return p.publishBlobTxWithRetry(tx, nonce)
	}
	return p.publishDynamicTxWithRetry(tx, nonce)
}

func (p *Publisher) publishDynamicTxWithRetry(tx types.TxData, nonce uint64) error {
	retries := 0
	for !p.hostStopper.IsStopping() {
		if _, err := p.executeTransaction(tx, nonce, retries); err != nil {
			retries++
			continue
		}
		return nil
	}
	return errors.New("stopped while retrying transaction")
}

func (p *Publisher) publishBlobTxWithRetry(tx types.TxData, nonce uint64) error {
	const maxRetries = 5
	retries := 0

	for !p.hostStopper.IsStopping() && retries < maxRetries {
		pricedTx, err := p.executeTransaction(tx, nonce, retries)
		if pricedTx == nil {
			return fmt.Errorf("could not price transaction. Cause: %w", err)
		}
		if err != nil {
			if retries >= maxRetries-1 {
				blobTx, ok := pricedTx.(*types.BlobTx)
				if !ok {
					return &MaxRetriesError{
						Err:    fmt.Sprintf("unexpected tx type: %T", pricedTx),
						BlobTx: nil,
					}
				}
				return &MaxRetriesError{
					Err:    err.Error(),
					BlobTx: blobTx,
				}
			}
			retries++
			continue
		}
		return nil
	}
	return errors.New("stopped while retrying transaction")
}

// executeTransaction handles the common flow of pricing, signing, sending and waiting for receipt. Returns the priced
// transaction so we can log the values in the event we exceed the maximum number of retries.
func (p *Publisher) executeTransaction(tx types.TxData, nonce uint64, retryNum int) (types.TxData, error) {
	// Set gas prices and create transaction
	pricedTx, err := ethadapter.SetTxGasPrice(p.sendingContext, p.ethClient, tx, p.hostWallet.Address(), nonce, retryNum, p.logger)
	if err != nil {
		return pricedTx, errors.Wrap(err, "could not estimate gas/gas price for L1 tx")
	}

	// Sign and send
	signedTx, err := p.hostWallet.SignTransaction(pricedTx)
	if err != nil {
		return pricedTx, errors.Wrap(err, "could not sign L1 tx")
	}

	err = p.ethClient.SendTransaction(signedTx)
	if err != nil {
		p.logger.Warn("Failed to send transaction",
			"error", err,
			"nonce", signedTx.Nonce(),
			"txHash", signedTx.Hash())
		return pricedTx, errors.Wrap(err, "could not broadcast L1 tx")
	}

	// Wait for receipt
	receipt, err := p.waitForReceipt(signedTx)
	if err != nil || receipt.Status != types.ReceiptStatusSuccessful {
		return pricedTx, fmt.Errorf(signedTx.Hash().Hex()) // Return hash for MaxRetriesError
	}

	p.logger.Debug("L1 transaction successful receipt found.", log.TxKey, signedTx.Hash(),
		log.BlockHeightKey, receipt.BlockNumber, log.BlockHashKey, receipt.BlockHash)
	return pricedTx, nil
}

// Helper functions to reduce duplication
func (p *Publisher) waitForReceipt(signedTx *types.Transaction) (*types.Receipt, error) {
	var receipt *types.Receipt
	err := retry.Do(
		func() error {
			if p.hostStopper.IsStopping() {
				return retry.FailFast(errors.New("host is stopping or context canceled"))
			}
			var err error
			receipt, err = p.ethClient.TransactionReceipt(signedTx.Hash())
			if err != nil {
				return fmt.Errorf("could not get receipt publishing tx for L1 tx=%s: %w", signedTx.Hash(), err)
			}
			return err
		},
		retry.NewTimeoutStrategy(p.maxWaitForL1Receipt, p.retryIntervalForL1Receipt),
	)
	return receipt, err
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

// MaxRetriesError is a specific error type for handling max retries
type MaxRetriesError struct {
	Err    string
	BlobTx *types.BlobTx
}

func (e *MaxRetriesError) Error() string {
	return fmt.Sprintf("max retries reached for nonce %d  with BlobFeeCap: %d, GasTipCap: %d, GasFeeCap: %d, Gas: %d",
		e.BlobTx.Nonce, e.BlobTx.BlobFeeCap, e.BlobTx.GasTipCap, e.BlobTx.GasFeeCap, e.BlobTx.Gas)
}
