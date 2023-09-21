package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/obscuronet/go-obscuro/tools/walletextension/accountmanager"
	"github.com/obscuronet/go-obscuro/tools/walletextension/common"
	"github.com/obscuronet/go-obscuro/tools/walletextension/userconn"
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

func getQueryParameter(params map[string]string, selectedParameter string) (string, error) {
	value, exists := params[selectedParameter]
	if !exists {
		return "", fmt.Errorf("parameter '%s' is not in the query params", selectedParameter)
	}

	return value, nil
}

func getUserID(conn userconn.UserConn, userIDPosition int) (string, error) {
	// try getting userID from query parameters and return it if successful
	userID, err := getQueryParameter(conn.ReadRequestParams(), common.UserQueryParameter)
	if err == nil {
		if len(userID) != common.MessageUserIDLen {
			return "", fmt.Errorf(fmt.Sprintf("wrong length of userID from URL. Got: %d, Expected: %d", len(userID), common.MessageUserIDLen))
		}
		return userID, err
	}

	// Alternatively, try to get it from URL path
	// This is a temporary hack to work around hardhat bug which causes hardhat to ignore query parameters.
	// It is unsafe because https encrypts query parameters,
	// but not URL itself and will be removed once hardhat bug is resolved.
	path := conn.GetHTTPRequest().URL.Path
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	// our URLs, which require userID, have following pattern: <version>/<endpoint (*optional)>/<userID (*optional)>
	// userID can be only on second or third position
	if len(parts) != userIDPosition+1 {
		return "", fmt.Errorf("URL structure of the request looks wrong")
	}
	userID = parts[userIDPosition]

	// Check if userID has the correct length
	if len(userID) != common.MessageUserIDLen {
		return "", fmt.Errorf(fmt.Sprintf("wrong length of userID from URL. Got: %d, Expected: %d", len(userID), common.MessageUserIDLen))
	}

	return userID, nil
}
