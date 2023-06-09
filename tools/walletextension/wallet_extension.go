package walletextension

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/obscuronet/go-obscuro/tools/walletextension/useraccountmanager"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/obscuronet/go-obscuro/go/common/stopcontrol"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/tools/walletextension/accountmanager"
	"github.com/obscuronet/go-obscuro/tools/walletextension/common"
	"github.com/obscuronet/go-obscuro/tools/walletextension/storage"
	"github.com/obscuronet/go-obscuro/tools/walletextension/userconn"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

var ErrSubscribeFailHTTP = fmt.Sprintf("received an %s request but the connection does not support subscriptions", rpc.Subscribe)

// WalletExtension handles the management of viewing keys and the forwarding of Ethereum JSON-RPC requests.
type WalletExtension struct {
	hostAddr           string // The address on which the Obscuro host can be reached.
	userAccountManager *useraccountmanager.UserAccountManager
	unsignedVKs        map[gethcommon.Address]*rpc.ViewingKey // Map temporarily holding VKs that have been generated but not yet signed
	storage            *storage.Storage
	logger             gethlog.Logger
	stopControl        *stopcontrol.StopControl
}

func New(
	hostAddr string,
	userAccountManager *useraccountmanager.UserAccountManager,
	storage *storage.Storage,
	stopControl *stopcontrol.StopControl,
	logger gethlog.Logger,
) *WalletExtension {
	return &WalletExtension{
		hostAddr:           hostAddr,
		userAccountManager: userAccountManager,
		unsignedVKs:        map[gethcommon.Address]*rpc.ViewingKey{},
		storage:            storage,
		logger:             logger,
		stopControl:        stopControl,
	}
}

// IsStopping returns whether the WE is stopping
func (w *WalletExtension) IsStopping() bool {
	return w.stopControl.IsStopping()
}

// Logger returns the WE set logger
func (w *WalletExtension) Logger() gethlog.Logger {
	return w.logger
}

// ProxyEthRequest proxys an incoming user request to the enclave
func (w *WalletExtension) ProxyEthRequest(request *accountmanager.RPCRequest, conn userconn.UserConn) (map[string]interface{}, error) {
	response := map[string]interface{}{}
	// all responses must contain the request id. Both successful and unsuccessful.
	response[common.JSONKeyRPCVersion] = jsonrpc.Version
	response[common.JSONKeyID] = request.ID

	// proxyRequest will find the correct client to proxy the request (or try them all if appropriate)
	var rpcResp interface{}

	// todo (@ziga) Remove this code after implementatio for all userIDs is done
	defaultAccManager, err := w.userAccountManager.GetUserAccountManager(common.DefaultUser)
	if err != nil {
		return nil, err
	}
	err = defaultAccManager.ProxyRequest(request, &rpcResp, conn)

	if err != nil && !errors.Is(err, rpc.ErrNilResponse) {
		response = common.CraftErrorResponse(err)
	} else if errors.Is(err, rpc.ErrNilResponse) {
		// if err was for a nil response then we will return an RPC result of null to the caller (this is a valid "not-found" response for some methods)
		response[common.JSONKeyResult] = nil
	} else {
		response[common.JSONKeyResult] = rpcResp

		// todo (@ziga) - fix this upstream on the decode
		// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-658.md
		adjustStateRoot(rpcResp, response)
	}
	return response, nil
}

// GenerateViewingKey generates the user viewing key and waits for signature
func (w *WalletExtension) GenerateViewingKey(addr gethcommon.Address) (string, error) {
	viewingKeyPrivate, err := crypto.GenerateKey()
	if err != nil {
		return "", fmt.Errorf("unable to generate a new keypair - %w", err)
	}

	viewingPublicKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	viewingPrivateKeyEcies := ecies.ImportECDSA(viewingKeyPrivate)

	w.unsignedVKs[addr] = &rpc.ViewingKey{
		Account:    &addr,
		PrivateKey: viewingPrivateKeyEcies,
		PublicKey:  viewingPublicKeyBytes,
		SignedKey:  nil, // we await a signature from the user before we can set up the EncRPCClient
	}

	// We return the hex of the viewing key's public key for MetaMask to sign over.
	viewingKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	return hex.EncodeToString(viewingKeyBytes), nil
}

// SubmitViewingKey checks a signed vieweing key and forwards it to the enclave
func (w *WalletExtension) SubmitViewingKey(address gethcommon.Address, signature []byte) error {
	vk, found := w.unsignedVKs[address]
	if !found {
		return fmt.Errorf(fmt.Sprintf("no viewing key found to sign for acc=%s, please call %s to generate key before sending signature", address, common.PathGenerateViewingKey))
	}

	// We transform the V from 27/28 to 0/1. This same change is made in Geth internals, for legacy reasons to be able
	// to recover the address: https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	signature[64] -= 27

	vk.SignedKey = signature
	// create an encrypted RPC client with the signed VK and register it with the enclave
	// todo (@ziga) - Create the clients lazily, to reduce connections to the host.
	client, err := rpc.NewEncNetworkClient(w.hostAddr, vk, w.logger)
	if err != nil {
		return fmt.Errorf("failed to create encrypted RPC client for account %s - %w", address, err)
	}
	defaultAccountManager, err := w.userAccountManager.GetUserAccountManager(common.DefaultUser)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("error getting default user account manager: %s", err))
	}

	defaultAccountManager.AddClient(address, client)

	err = w.storage.AddAccount([]byte(common.DefaultUser), vk.Account.Bytes(), vk.SignedKey)
	if err != nil {
		return fmt.Errorf("error saving account") // todo (@ziga) - improve error messages!
	}

	if err != nil {
		return fmt.Errorf("error saving viewing key to the database: %w", err)
	}

	// finally we remove the VK from the pending 'unsigned VKs' map now the client has been created
	delete(w.unsignedVKs, address)

	return nil
}

func adjustStateRoot(rpcResp interface{}, respMap map[string]interface{}) {
	if resultMap, ok := rpcResp.(map[string]interface{}); ok {
		if val, foundRoot := resultMap[common.JSONKeyRoot]; foundRoot {
			if val == "0x" {
				respMap[common.JSONKeyResult].(map[string]interface{})[common.JSONKeyRoot] = nil
			}
		}
	}
}
