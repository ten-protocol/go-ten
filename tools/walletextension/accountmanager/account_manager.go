package accountmanager

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ten-protocol/go-ten/tools/walletextension/storage"

	"github.com/ten-protocol/go-ten/tools/walletextension/subscriptions"

	"github.com/ten-protocol/go-ten/go/common/gethencoding"

	gethlog "github.com/ethereum/go-ethereum/log"

	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"

	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/tools/walletextension/userconn"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	ethCallPaddedArgLen = 64
	ethCallAddrPadding  = "000000000000000000000000"

	ErrNoViewingKey = "method %s cannot be called with an unauthorised client - no signed viewing keys found"
)

// AccountManager provides a single location for code that helps the gateway in determining the appropriate
// account to use to send a request for selected user when multiple accounts are registered
type AccountManager struct {
	userID               string
	unauthedClient       rpc.Client
	accountsMutex        sync.RWMutex
	accountClientsHTTP   map[gethcommon.Address]*rpc.EncRPCClient // An encrypted RPC http client per registered account
	hostRPCBindAddrWS    string
	subscriptionsManager *subscriptions.SubscriptionManager
	storage              storage.Storage
	logger               gethlog.Logger
}

func NewAccountManager(userID string, unauthedClient rpc.Client, hostRPCBindAddressWS string, storage storage.Storage, logger gethlog.Logger) *AccountManager {
	return &AccountManager{
		userID:               userID,
		unauthedClient:       unauthedClient,
		accountClientsHTTP:   make(map[gethcommon.Address]*rpc.EncRPCClient),
		hostRPCBindAddrWS:    hostRPCBindAddressWS,
		subscriptionsManager: subscriptions.New(logger),
		storage:              storage,
		logger:               logger,
	}
}

// GetAllAddressesWithClients returns a list of addresses which already have clients (are in accountClients map)
func (m *AccountManager) GetAllAddressesWithClients() []string {
	m.accountsMutex.RLock()
	defer m.accountsMutex.RUnlock()

	addresses := make([]string, 0, len(m.accountClientsHTTP))
	for address := range m.accountClientsHTTP {
		addresses = append(addresses, address.Hex())
	}
	return addresses
}

// AddClient adds a client to the list of clients, keyed by account address.
func (m *AccountManager) AddClient(address gethcommon.Address, client *rpc.EncRPCClient) {
	m.accountsMutex.Lock()
	defer m.accountsMutex.Unlock()
	m.accountClientsHTTP[address] = client
}

// ProxyRequest tries to identify the correct EncRPCClient to proxy the request to the Ten node, or it will attempt
// the request with all clients until it succeeds
func (m *AccountManager) ProxyRequest(rpcReq *wecommon.RPCRequest, rpcResp *interface{}, userConn userconn.UserConn) error {
	// We need to handle a special case for subscribing and unsubscribing from events,
	// because we need to handle multiple accounts with a single user request
	if rpcReq.Method == rpc.Subscribe {
		clients, err := m.suggestSubscriptionClient(rpcReq)
		if err != nil {
			return err
		}
		err = m.subscriptionsManager.HandleNewSubscriptions(clients, rpcReq, rpcResp, userConn)
		if err != nil {
			m.logger.Error("Error subscribing to multiple clients")
			return err
		}
		return nil
	}
	if rpcReq.Method == rpc.Unsubscribe {
		if len(rpcReq.Params) != 1 {
			return fmt.Errorf("one parameter (subscriptionID) expected, %d parameters received", len(rpcReq.Params))
		}
		subscriptionID, ok := rpcReq.Params[0].(string)
		if !ok {
			return fmt.Errorf("subscriptionID needs to be a string. Got: %v", rpcReq.Params[0])
		}
		m.subscriptionsManager.HandleUnsubscribe(subscriptionID, rpcResp)
		return nil
	}
	return m.executeCall(rpcReq, rpcResp)
}

const emptyFilterCriteria = "[]" // This is the value that gets passed for an empty filter criteria.

// suggestSubscriptionClient returns clients that should be used for the subscription request.
// For other requests we use http clients, but for subscriptions ws clients are required, that is the reason for
// creating ws clients here.
// We only want to have the connections open for the duration of the subscription, so we create the clients here and
// don't store them in the accountClients map.
func (m *AccountManager) suggestSubscriptionClient(rpcReq *wecommon.RPCRequest) ([]rpc.Client, error) {
	m.accountsMutex.RLock()
	defer m.accountsMutex.RUnlock()

	userIDBytes, err := wecommon.GetUserIDbyte(m.userID)
	if err != nil {
		return nil, fmt.Errorf("error decoding string (%s), %w", m.userID, err)
	}

	accounts, err := m.storage.GetAccounts(userIDBytes)
	if err != nil {
		return nil, fmt.Errorf("error getting accounts for user: %s, %w", m.userID, err)
	}

	userPrivateKey, err := m.storage.GetUserPrivateKey(userIDBytes)
	if err != nil {
		return nil, fmt.Errorf("error getting private key for user: %s, %w", m.userID, err)
	}

	if len(rpcReq.Params) > 1 {
		filteredAccounts, err := m.filterAccounts(rpcReq, accounts)
		if err != nil {
			return nil, err
		}
		// return filtered clients if we found any
		if len(filteredAccounts) > 0 {
			accounts = filteredAccounts
		}
	}
	// create clients for all accounts if we didn't find any clients that match the filter or if no topics were provided
	return m.createClientsForAccounts(accounts, userPrivateKey)
}

// filterClients checks if any of the accounts match the filter criteria and returns those accounts
func (m *AccountManager) filterAccounts(rpcReq *wecommon.RPCRequest, accounts []wecommon.AccountDB) ([]wecommon.AccountDB, error) {
	var filteredAccounts []wecommon.AccountDB
	filterCriteriaJSON, err := json.Marshal(rpcReq.Params[1])
	if err != nil {
		return nil, fmt.Errorf("could not marshal filter criteria to JSON. Cause: %w", err)
	}
	filterCriteria := filters.FilterCriteria{}
	if string(filterCriteriaJSON) != emptyFilterCriteria {
		err = filterCriteria.UnmarshalJSON(filterCriteriaJSON)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal filter criteria from the following JSON: `%s`. Cause: %w", string(filterCriteriaJSON), err)
		}
	}

	for _, topicCondition := range filterCriteria.Topics {
		for _, topic := range topicCondition {
			potentialAddr := common.ExtractPotentialAddress(topic)
			m.logger.Info(fmt.Sprintf("Potential address (%s) found for the request %s", potentialAddr, rpcReq))
			if potentialAddr != nil {
				for _, account := range accounts {
					// if we find a match, we append the account to the list of filtered accounts
					if bytes.Equal(account.AccountAddress, potentialAddr.Bytes()) {
						filteredAccounts = append(filteredAccounts, account)
					}
				}
			}
		}
	}

	return filteredAccounts, nil
}

// createClientsForAllAccounts creates ws clients for all accounts for given user and returns them
func (m *AccountManager) createClientsForAccounts(accounts []wecommon.AccountDB, userPrivateKey []byte) ([]rpc.Client, error) {
	clients := make([]rpc.Client, 0, len(accounts))
	for _, account := range accounts {
		encClient, err := wecommon.CreateEncClient(m.hostRPCBindAddrWS, account.AccountAddress, userPrivateKey, account.Signature, account.SignatureType, m.logger)
		if err != nil {
			m.logger.Error(fmt.Errorf("error creating new client, %w", err).Error())
			continue
		}
		clients = append(clients, encClient)
	}
	return clients, nil
}

// todo - better way
const notAuthorised = "not authorised"

var platformAuthorisedCalls = map[string]bool{
	rpc.GetBalance: true,
	// rpc.GetCode, //todo
	rpc.GetTransactionCount:   true,
	rpc.GetTransactionReceipt: true,
	rpc.GetLogs:               true,
}

func (m *AccountManager) executeCall(rpcReq *wecommon.RPCRequest, rpcResp *interface{}) error {
	m.accountsMutex.RLock()
	defer m.accountsMutex.RUnlock()
	// for Ten RPC requests, it is important we know the sender account for the viewing key encryption/decryption
	suggestedClient := m.suggestAccountClient(rpcReq, m.accountClientsHTTP)

	switch {
	case suggestedClient != nil: // use the suggested client if there is one
		// todo (@ziga) - if we have a suggested client, should we still loop through the other clients if it fails?
		// 		The call data guessing won't often be wrong but there could be edge-cases there
		return submitCall(suggestedClient, rpcReq, rpcResp)

	case len(m.accountClientsHTTP) > 0: // try registered clients until there's a successful execution
		m.logger.Info(fmt.Sprintf("appropriate client not found, attempting request with up to %d clients", len(m.accountClientsHTTP)))
		var err error
		for _, client := range m.accountClientsHTTP {
			err = submitCall(client, rpcReq, rpcResp)
			if err == nil || errors.Is(err, rpc.ErrNilResponse) {
				// request didn't fail, we don't need to continue trying the other clients
				return nil
			}
			// platform calls return a standard error for calls that are not authorised.
			// any other error can be returned early
			if platformAuthorisedCalls[rpcReq.Method] && err.Error() != notAuthorised {
				return err
			}
		}
		// every attempt errored
		return err

	default: // no clients registered, use the unauthenticated one
		if rpc.IsSensitiveMethod(rpcReq.Method) {
			return fmt.Errorf(ErrNoViewingKey, rpcReq.Method)
		}
		return m.unauthedClient.Call(rpcResp, rpcReq.Method, rpcReq.Params...)
	}
}

// suggestAccountClient works through various methods to try and guess which available client to use for a request, returns nil if none found
func (m *AccountManager) suggestAccountClient(req *wecommon.RPCRequest, accClients map[gethcommon.Address]*rpc.EncRPCClient) *rpc.EncRPCClient {
	if len(accClients) == 1 {
		for _, client := range accClients {
			// return the first (and only) client
			return client
		}
	}
	switch req.Method {
	case rpc.Call, rpc.EstimateGas:
		return m.handleEthCall(req, accClients)
	case rpc.GetBalance:
		return extractAddress(0, req.Params, accClients)
	case rpc.GetLogs:
		return extractAddress(1, req.Params, accClients)
	case rpc.GetTransactionCount:
		return extractAddress(0, req.Params, accClients)
	default:
		return nil
	}
}

func extractAddress(pos int, params []interface{}, accClients map[gethcommon.Address]*rpc.EncRPCClient) *rpc.EncRPCClient {
	if len(params) < pos+1 {
		return nil
	}
	requestedAddress, err := gethencoding.ExtractAddress(params[pos])
	if err == nil {
		return accClients[*requestedAddress]
	}
	return nil
}

func (m *AccountManager) handleEthCall(req *wecommon.RPCRequest, accClients map[gethcommon.Address]*rpc.EncRPCClient) *rpc.EncRPCClient {
	paramsMap, err := parseParams(req.Params)
	if err != nil {
		// no further info to deduce calling client
		return nil
	}
	// check if request params had a "from" address and if we had a client for that address
	fromClient, found := checkForFromField(paramsMap, accClients)
	if found {
		return fromClient
	}

	// Otherwise, we search the `data` field for an address matching a registered viewing key.
	addr, err := searchDataFieldForAccount(paramsMap, accClients)
	if err == nil {
		return accClients[*addr]
	}
	return nil
}

// Many eth RPC requests provide params as first argument in a json map with similar fields (e.g. a `from` field)
func parseParams(args []interface{}) (map[string]interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no params found to unmarshal")
	}

	// only interested in trying first arg
	params, ok := args[0].(map[string]interface{})
	if !ok {
		callParamsJSON, ok := args[0].([]byte)
		if !ok {
			return nil, fmt.Errorf("first arg was not a byte array")
		}

		err := json.Unmarshal(callParamsJSON, &params)
		if err != nil {
			return nil, fmt.Errorf("first arg couldn't be unmarshaled into a params map")
		}
	}

	return params, nil
}

func checkForFromField(paramsMap map[string]interface{}, accClients map[gethcommon.Address]*rpc.EncRPCClient) (*rpc.EncRPCClient, bool) {
	fromVal, found := paramsMap[wecommon.JSONKeyFrom]
	if !found {
		return nil, false
	}

	fromStr, ok := fromVal.(string)
	if !ok {
		return nil, false
	}

	fromAddr := gethcommon.HexToAddress(fromStr)
	client, found := accClients[fromAddr]
	return client, found
}

// Extracts the arguments from the request's `data` field. If any of them, after removing padding, match the viewing
// key address, we return that address. Otherwise, we return nil.
func searchDataFieldForAccount(callParams map[string]interface{}, accClients map[gethcommon.Address]*rpc.EncRPCClient) (*gethcommon.Address, error) {
	// We ensure that the `data` field is present.
	data := callParams[wecommon.JSONKeyData]
	if data == nil {
		return nil, fmt.Errorf("eth_call request did not have its `data` field set")
	}
	dataString, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("eth_call request's `data` field was not of the expected type `string`")
	}

	// We check that the data field is long enough before removing the leading "0x" (1 bytes/2 chars) and the method ID
	// (4 bytes/8 chars).
	if len(dataString) < 10 {
		return nil, fmt.Errorf("data field is not long enough - no known account found in data bytes")
	}
	dataString = dataString[10:]

	// We split up the arguments in the `data` field.
	var dataArgs []string
	for i := 0; i < len(dataString); i += ethCallPaddedArgLen {
		if i+ethCallPaddedArgLen > len(dataString) {
			break
		}
		dataArgs = append(dataArgs, dataString[i:i+ethCallPaddedArgLen])
	}

	// We iterate over the arguments, looking for an argument that matches a viewing key address
	for _, dataArg := range dataArgs {
		// If the argument doesn't have the correct padding, it's not an address.
		if !strings.HasPrefix(dataArg, ethCallAddrPadding) {
			continue
		}

		maybeAddress := gethcommon.HexToAddress(dataArg[len(ethCallAddrPadding):])
		if _, ok := accClients[maybeAddress]; ok {
			return &maybeAddress, nil
		}
	}

	return nil, fmt.Errorf("no known account found in data bytes")
}

func submitCall(client *rpc.EncRPCClient, req *wecommon.RPCRequest, resp *interface{}) error {
	if req.Method == rpc.Call || req.Method == rpc.EstimateGas {
		// Never modify the original request, as it might be reused.
		req = req.Clone()

		// Any method using an ethereum.CallMsg is a sensitive method that requires a viewing key lookup but the 'from' field is not mandatory
		// and is often not included from metamask etc. So we ensure it is populated here.
		account := client.Account()
		var err error
		req.Params, err = setFromFieldIfMissing(req.Params, *account)
		if err != nil {
			return err
		}
	}

	if req.Method == rpc.GetLogs {
		// Never modify the original request, as it might be reused.
		req = req.Clone()

		// We add the account to the list of arguments, so we know which account to use to filter the logs and encrypt
		// the result.
		req.Params = append(req.Params, client.Account().Hex())
	}

	return client.Call(resp, req.Method, req.Params...)
}

// The enclave requires the `from` field to be set so that it can encrypt the response, but sources like MetaMask often
// don't set it. So we check whether it's present; if absent, we walk through the arguments in the request's `data`
// field, and if any of the arguments match our viewing key address, we set the `from` field to that address.
func setFromFieldIfMissing(args []interface{}, account gethcommon.Address) ([]interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no params found to unmarshal")
	}

	callMsg, err := gethencoding.ExtractEthCallMapString(args[0])
	if err != nil {
		return nil, fmt.Errorf("unable to marshal callMsg - %w", err)
	}

	// We only modify `eth_call` requests where the `from` field is not set.
	if callMsg[gethencoding.CallFieldFrom] != gethcommon.HexToAddress("0x0").Hex() {
		return args, nil
	}

	// override the existing args
	callMsg[gethencoding.CallFieldFrom] = account.Hex()

	// do not modify other existing arguments
	request := []interface{}{callMsg}
	for i := 1; i < len(args); i++ {
		request = append(request, args[i])
	}

	return request, nil
}
