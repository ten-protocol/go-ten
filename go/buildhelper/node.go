package buildhelper

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum"

	"github.com/obscuronet/obscuro-playground/go/buildhelper/buildconstants"
	"github.com/obscuronet/obscuro-playground/go/buildhelper/helpertypes"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

var nonceLock sync.RWMutex

type EthNode struct {
	port      uint
	ipaddress string
	apiClient *EthAPI
	id        common.Address
}

// TODO remove the id common.Address
func NewEthNode(id common.Address, ipaddress string, port uint) (obscurocommon.L1Node, error) {
	apiClient := NewEthAPI(ipaddress, port)
	if err := apiClient.Connect(); err != nil {
		return nil, fmt.Errorf("unable to connect to the eth node - %w", err)
	}

	log.Log(fmt.Sprintf("Initializing eth node at contract: %s", buildconstants.ContractAddress))
	return &EthNode{
		ipaddress: ipaddress,
		port:      port,
		apiClient: apiClient,
		id:        id,
	}, nil
}

func (e *EthNode) RPCBlockchainFeed() []*types.Block {
	var availBlocks []*types.Block

	block, err := e.apiClient.APIClient.BlockByNumber(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	availBlocks = append(availBlocks, block)

	for {
		parentHash := block.ParentHash()
		// todo set this to genesis hash
		if parentHash.Hex() == "0x0000000000000000000000000000000000000000000000000000000000000000" {
			break
		}

		block, err = e.apiClient.APIClient.BlockByHash(context.Background(), block.ParentHash())
		if err != nil {
			fmt.Printf("ERROR %v\n", err)
		}

		availBlocks = append(availBlocks, block)
	}

	// TODO double check the list is ordered [genesis, 1, 2, 3, 4, ..., last]
	// TODO It's pretty ugly but it avoids creating a new slice
	// TODO The approach of feeding all the blocks should change from all-blocks-in-memory to a stream
	for i, j := 0, len(availBlocks)-1; i < j; i, j = i+1, j-1 {
		availBlocks[i], availBlocks[j] = availBlocks[j], availBlocks[i]
	}
	return availBlocks
}

func (e *EthNode) BroadcastTx(t obscurocommon.EncodedL1Tx) {
	nonceLock.Lock()
	defer nonceLock.Unlock()

	privateKey := Addr1PK()
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := e.apiClient.APIClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		panic(err)
	}

	l1tx, err := t.Decode()
	if err != nil {
		panic(err)
	}

	l1txData := obscurocommon.TxData(&l1tx)

	ethTx := &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: big.NewInt(20000000000),
		Gas:      1024_000_000,
		To:       &buildconstants.ContractAddress,
	}

	contractABI, err := abi.JSON(strings.NewReader(buildconstants.ContractAbi))
	if err != nil {
		panic(err)
	}

	// TODO each of these cases should be a function:
	// TODO like: func createRollupTx() or func createDepositTx()
	// TODO And then eventually, these functions would be called directly, when we get rid of our special format. (we'll have to change the mock thing as well for that)
	switch l1txData.TxType {
	case obscurocommon.DepositTx:
		ethTx.Value = big.NewInt(int64(l1txData.Amount))
		data, err := contractABI.Pack("Deposit", l1txData.Dest)
		if err != nil {
			panic(err)
		}
		ethTx.Data = data
		log.Log(fmt.Sprintf("BROADCAST TX: Issuing DepositTx - Addr: %s deposited %d to %s ",
			fromAddress, l1txData.Amount, l1txData.Dest))

	case obscurocommon.RollupTx:
		r, err := nodecommon.DecodeRollup(l1txData.Rollup)
		if err != nil {
			panic(err)
		}
		zipped := helpertypes.Compress(l1txData.Rollup)
		encRollupData := helpertypes.EncodeToString(zipped)
		data, err := contractABI.Pack("AddRollup", encRollupData)
		if err != nil {
			panic(err)
		}

		ethTx.Data = data
		derolled, _ := nodecommon.DecodeRollup(l1txData.Rollup)

		msg := ethereum.CallMsg{
			From:     Addr1(),
			To:       &common.Address{},
			GasPrice: big.NewInt(20000000000),
			Gas:      1024_000_000,
			Data:     data,
		}

		estimate, err := e.apiClient.APIClient.EstimateGas(context.Background(), msg)
		if err != nil {
			panic(err)
		}
		log.Log(fmt.Sprintf("Estimated Cost of: %d\n", estimate))
		log.Log(fmt.Sprintf("BROADCAST TX - Issuing Rollup: %s - %d txs - datasize: %d - gas: %d \n", r.Hash(), len(derolled.Transactions), len(data), ethTx.Gas))

	case obscurocommon.StoreSecretTx:
		data, err := contractABI.Pack("StoreSecret", helpertypes.EncodeToString(l1txData.Secret))
		if err != nil {
			panic(err)
		}
		ethTx.Data = data
		log.Log(fmt.Sprintf("BROADCAST TX: Issuing StoreSecretTx: encoded as %s", helpertypes.EncodeToString(l1txData.Secret)))
	case obscurocommon.RequestSecretTx:
		data, err := contractABI.Pack("RequestSecret")
		if err != nil {
			panic(err)
		}
		ethTx.Data = data
		log.Log("BROADCAST TX: Issuing RequestSecret")
	}

	signedTx, err := types.SignNewTx(privateKey, types.NewEIP155Signer(big.NewInt(1337)), ethTx)
	if err != nil {
		panic(err)
	}

	err = e.apiClient.APIClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
	}
}

func (e *EthNode) BlockListener() chan *types.Header {
	ch := make(chan *types.Header, 1)
	subs, err := e.apiClient.APIClient.SubscribeNewHead(context.Background(), ch)
	if err != nil {
		panic(err)
	}
	// we should hook the subs to cleanup
	fmt.Println(subs)

	return ch
}

func (e *EthNode) FetchBlockByNumber(n *big.Int) (*types.Block, error) {
	return e.apiClient.APIClient.BlockByNumber(context.Background(), n)
}

func (e *EthNode) FetchBlock(hash common.Hash) (*types.Block, error) {
	return e.apiClient.APIClient.BlockByHash(context.Background(), hash)
}
