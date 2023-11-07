package accountmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/obscuronet/go-obscuro/tools/walletextension/subscriptions"

	"github.com/obscuronet/go-obscuro/go/common/gethencoding"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/go/common"

	wecommon "github.com/obscuronet/go-obscuro/tools/walletextension/common"

	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/tools/walletextension/userconn"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	ethCallPaddedArgLen = 64
	ethCallAddrPadding  = "000000000000000000000000"

	ErrNoViewingKey = "method %s cannot be called with an unauthorised client - no signed viewing keys found"
)

// AccountManager provides a single location for code that helps wallet extension in determining the appropriate
// account to use to send a request when multiple are registered
type AccountManager struct {
	unauthedClient rpc.Client
	// todo (@ziga) - create two types of clients - WS clients, and HTTP clients - to not create WS clients unnecessarily.
	accountsMutex        sync.RWMutex
	accountClients       map[gethcommon.Address]*rpc.EncRPCClient // An encrypted RPC client per registered account
	subscriptionsManager *subscriptions.SubscriptionManager
	logger               gethlog.Logger
}

func NewAccountManager(unauthedClient rpc.Client, logger gethlog.Logger) *AccountManager {
	return &AccountManager{
		unauthedClient:       unauthedClient,
		accountClients:       make(map[gethcommon.Address]*rpc.EncRPCClient),
		subscriptionsManager: subscriptions.New(logger),
		logger:               logger,
	}
}

// AddClient adds a client to the list of clients, keyed by account address.
func (m *AccountManager) AddClient(address gethcommon.Address, client *rpc.EncRPCClient) {
	m.accountsMutex.Lock()
	defer m.accountsMutex.Unlock()
	m.accountClients[address] = client
}

// ProxyRequest tries to identify the correct EncRPCClient to proxy the request to the Obscuro node, or it will attempt
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
			return fmt.Errorf("subscriptionID needs to be a string. Got: %d", rpcReq.Params[0])
		}
		m.subscriptionsManager.HandleUnsubscribe(subscriptionID, rpcResp)
		return nil
	}
	return m.executeCall(rpcReq, rpcResp)
}

const emptyFilterCriteria = "[]" // This is the value that gets passed for an empty filter criteria.

// determine the client based on the topics
// if none is found use all clients from current user
func (m *AccountManager) suggestSubscriptionClient(rpcReq *wecommon.RPCRequest) ([]rpc.Client, error) {
	m.accountsMutex.RLock()
	defer m.accountsMutex.RUnlock()

	clients := make([]rpc.Client, 0, len(m.accountClients))

	// by default, if no client is identified as a candidate, then subscribe to all accounts
	for _, c := range m.accountClients {
		clients = append(clients, c)
	}

	if len(rpcReq.Params) < 2 {
		return clients, nil
	}

	// The filter is the second parameter
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

	// Go through each topic filter and look for registered addresses
	for i, topicCondition := range filterCriteria.Topics {
		// the first topic is always the signature of the event, so it can't be an address
		if i == 0 {
			continue
		}
		for _, topic := range topicCondition {
			potentialAddr := common.ExtractPotentialAddress(topic)
			m.logger.Info(fmt.Sprintf("Potential address (%s) found for the request %s", potentialAddr, rpcReq))
			if potentialAddr != nil {
				cl, found := m.accountClients[*potentialAddr]
				if found {
					m.logger.Info("Client found for potential address: ", potentialAddr)
					return []rpc.Client{cl}, nil
				}
				m.logger.Info("Potential address does not have a client", potentialAddr)
			}
		}
	}

	return clients, nil
}

func (m *AccountManager) executeCall(rpcReq *wecommon.RPCRequest, rpcResp *interface{}) error {
	m.accountsMutex.RLock()
	defer m.accountsMutex.RUnlock()
	// for obscuro RPC requests it is important we know the sender account for the viewing key encryption/decryption
	suggestedClient := m.suggestAccountClient(rpcReq, m.accountClients)

	switch {
	case suggestedClient != nil: // use the suggested client if there is one
		// todo (@ziga) - if we have a suggested client, should we still loop through the other clients if it fails?
		// 		The call data guessing won't often be wrong but there could be edge-cases there
		return submitCall(suggestedClient, rpcReq, rpcResp)

	case len(m.accountClients) > 0: // try registered clients until there's a successful execution
		m.logger.Info(fmt.Sprintf("appropriate client not found, attempting request with up to %d clients", len(m.accountClients)))
		var err error
		for _, client := range m.accountClients {
			err = submitCall(client, rpcReq, rpcResp)
			if err == nil || errors.Is(err, rpc.ErrNilResponse) {
				// request didn't fail, we don't need to continue trying the other clients
				return nil
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

	paramsMap, err := parseParams(req.Params)
	if err != nil {
		// no further info to deduce calling client
		return nil
	}

	if req.Method == rpc.Call {
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
