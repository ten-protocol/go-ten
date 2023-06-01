package api

import (
	"encoding/json"
	"fmt"

	"github.com/obscuronet/go-obscuro/tools/walletextension/accountmanager"
	"github.com/obscuronet/go-obscuro/tools/walletextension/common"
)

func parseRequest(body []byte) (*accountmanager.RPCRequest, error) {
	// We unmarshal the JSON request
	var reqJSONMap map[string]json.RawMessage
	err := json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON-RPC request body to JSON: %s. "+
			"If you're trying to generate a viewing key, visit %s", err, common.PathViewingKeys)
	}

	reqID := reqJSONMap[common.JSONKeyID]
	var method string
	err = json.Unmarshal(reqJSONMap[common.JSONKeyMethod], &method)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal method string from JSON-RPC request body: %w", err)
	}

	// we extract the params into a JSON list
	var params []interface{}
	err = json.Unmarshal(reqJSONMap[common.JSONKeyParams], &params)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal params list from JSON-RPC request body: %w", err)
	}

	return &accountmanager.RPCRequest{
		ID:     reqID,
		Method: method,
		Params: params,
	}, nil
}
