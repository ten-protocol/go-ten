package accountmanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/obscuronet/go-obscuro/go/common"

	wecommon "github.com/obscuronet/go-obscuro/tools/walletextension/common"

	"github.com/go-kit/kit/transport/http/jsonrpc"

	"github.com/obscuronet/go-obscuro/go/common/gethenconding"

	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/tools/walletextension/userconn"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	methodEthSubscription = "eth_subscription"

	ethCallPaddedArgLen = 64
	ethCallAddrPadding  = "000000000000000000000000"

	ErrNoViewingKey = "method %s cannot be called with an unauthorised client - no signed viewing keys found"
)

// AccountManager provides a single location for code that helps wallet extension in determining the appropriate
// account to use to send a request when multiple are registered
type AccountManager struct {
	unauthedClient rpc.Client
	// TODO - Create two types of clients - WS clients, and HTTP clients - to not create WS clients unnecessarily.
	accountClients map[gethcommon.Address]*rpc.EncRPCClient // An encrypted RPC client per registered account
}

func NewAccountManager(unauthedClient rpc.Client) AccountManager {
	return AccountManager{
		unauthedClient: unauthedClient,
		accountClients: make(map[gethcommon.Address]*rpc.EncRPCClient),
	}
}

// AddClient adds a client to the list of clients, keyed by account address.
func (m *AccountManager) AddClient(address gethcommon.Address, client *rpc.EncRPCClient) {
	m.accountClients[address] = client
}

// ProxyRequest tries to identify the correct EncRPCClient to proxy the request to the Obscuro node, or it will attempt
// the request with all clients until it succeeds
func (m *AccountManager) ProxyRequest(rpcReq *RPCRequest, rpcResp *interface{}, userConn userconn.UserConn) error {
	// for obscuro RPC requests it is important we know the sender account for the viewing key encryption/decryption
	suggestedClient := suggestAccountClient(rpcReq, m.accountClients)

	switch {
	case suggestedClient != nil: // use the suggested client if there is one
		// todo: if we have a suggested client, should we still loop through the other clients if it fails?
		// 		The call data guessing won't often be wrong but there could be edge-cases there
		return performRequest(suggestedClient, rpcReq, rpcResp, userConn)

	case len(m.accountClients) > 0: // try registered clients until there's a successful execution
		log.Info("appropriate client not found, attempting request with up to %d clients", len(m.accountClients))
		var err error
		for _, client := range m.accountClients {
			err = performRequest(client, rpcReq, rpcResp, userConn)
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
func suggestAccountClient(req *RPCRequest, accClients map[gethcommon.Address]*rpc.EncRPCClient) *rpc.EncRPCClient {
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

	if req.Method == rpc.RPCCall {
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

	// todo: add other mechanisms for determining the correct account to use. E.g. we may want to start caching and
	// 	 	recent transaction hashes for accounts so that receipt lookups know which acc to use

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

func performRequest(client *rpc.EncRPCClient, req *RPCRequest, resp *interface{}, userConn userconn.UserConn) error {
	if req.Method == rpc.RPCSubscribe {
		return executeSubscribe(client, req, resp, userConn)
	}
	return executeCall(client, req, resp)
}

func executeSubscribe(client *rpc.EncRPCClient, req *RPCRequest, resp *interface{}, userConn userconn.UserConn) error {
	if len(req.Params) == 0 {
		return fmt.Errorf("could not subscribe as no subscription namespace was provided")
	}
	ch := make(chan common.IDAndLog)
	subscription, err := client.Subscribe(context.Background(), resp, rpc.RPCSubscribeNamespace, ch, req.Params...)
	if err != nil {
		return fmt.Errorf("could not call %s with params %v. Cause: %w", req.Method, req.Params, err)
	}

	// We listen for incoming messages on the subscription.
	go func() {
		for {
			select {
			case idAndLog := <-ch:
				if userConn.IsClosed() {
					log.Info("received log but websocket was closed")
					return
				}

				jsonResponse, err := prepareLogResponse(idAndLog)
				if err != nil {
					log.Error("could not marshal log response to JSON. Cause: %s", err)
					continue
				}

				log.Info("Forwarding log from Obscuro node: %s", jsonResponse)
				err = userConn.WriteResponse(jsonResponse)
				if err != nil {
					log.Error("could not write the JSON log to the websocket. Cause: %s", err)
					continue
				}

			case err = <-subscription.Err():
				// An error on this channel means the subscription has ended, so we exit the loop.
				userConn.HandleError(err.Error())
				return
			}
		}
	}()

	// We periodically check if the websocket is closed, and terminate the subscription.
	go func() {
		for {
			if userConn.IsClosed() {
				subscription.Unsubscribe()
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return nil
}

func executeCall(client *rpc.EncRPCClient, req *RPCRequest, resp *interface{}) error {
	if req.Method == rpc.RPCCall || req.Method == rpc.RPCEstimateGas {
		// Never modify the original request, as it might be reused.
		req = req.Clone()

		// Any method using an ethereum.CallMsg is a sensitive method that requires a viewing key lookup but the 'from' field is not mandatory
		// and is often not included from metamask etc. So we ensure it is populated here.
		account := client.Account()
		var err error
		req.Params, err = setCallFromFieldIfMissing(req.Params, *account)
		if err != nil {
			return err
		}
	}

	return client.Call(resp, req.Method, req.Params...)
}

// The enclave requires the `from` field to be set so that it can encrypt the response, but sources like MetaMask often
// don't set it. So we check whether it's present; if absent, we walk through the arguments in the request's `data`
// field, and if any of the arguments match our viewing key address, we set the `from` field to that address.
func setCallFromFieldIfMissing(args []interface{}, account gethcommon.Address) ([]interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no params found to unmarshal")
	}

	callMsg, err := gethenconding.ExtractEthCallMapString(args[0])
	if err != nil {
		return nil, fmt.Errorf("unable to marshal callMsg - %w", err)
	}

	// We only modify `eth_call` requests where the `from` field is not set.
	if callMsg[gethenconding.CallFieldFrom] != gethcommon.HexToAddress("0x0").Hex() {
		return args, nil
	}

	// override the existing args
	callMsg[gethenconding.CallFieldFrom] = account.Hex()

	// do not modify other existing arguments
	request := []interface{}{callMsg}
	for i := 1; i < len(args); i++ {
		request = append(request, args[i])
	}

	return request, nil
}

// Formats the log to be sent as an Eth JSON-RPC response.
func prepareLogResponse(idAndLog common.IDAndLog) ([]byte, error) {
	paramsMap := make(map[string]interface{})
	paramsMap[wecommon.JSONKeySubscription] = idAndLog.SubID
	paramsMap[wecommon.JSONKeyResult] = idAndLog.Log

	respMap := make(map[string]interface{})
	respMap[wecommon.JSONKeyRPCVersion] = jsonrpc.Version
	respMap[wecommon.JSONKeyMethod] = methodEthSubscription
	respMap[wecommon.JSONKeyParams] = paramsMap

	jsonResponse, err := json.Marshal(respMap)
	if err != nil {
		return nil, fmt.Errorf("could not marshal log response to JSON. Cause: %w", err)
	}
	return jsonResponse, nil
}

type RPCRequest struct {
	ID     json.RawMessage
	Method string
	Params []interface{}
}

// Clone returns a new instance of the *RPCRequest
func (r *RPCRequest) Clone() *RPCRequest {
	return &RPCRequest{
		ID:     r.ID,
		Method: r.Method,
		Params: r.Params,
	}
}
