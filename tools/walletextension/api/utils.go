package api

import (
	"encoding/json"
	"fmt"
	"strings"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/tools/walletextension/common"
	"github.com/obscuronet/go-obscuro/tools/walletextension/userconn"
)

func parseRequest(body []byte) (*common.RPCRequest, error) {
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

	return &common.RPCRequest{
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

func handleEthError(req *common.RPCRequest, conn userconn.UserConn, logger gethlog.Logger, err error) {
	var method string
	id := json.RawMessage("1")
	if req != nil {
		method = req.Method
		id = req.ID
	}

	errjson := &common.JSONError{
		Code:    0,
		Message: err.Error(),
		Data:    nil,
	}

	jsonRPRCError := common.JSONRPCMessage{
		Version: "2.0",
		ID:      id,
		Method:  method,
		Params:  nil,
		Error:   errjson,
		Result:  nil,
	}

	if evmError, ok := err.(errutil.EVMSerialisableError); ok { //nolint: errorlint
		jsonRPRCError.Error.Data = evmError.Reason
		jsonRPRCError.Error.Code = evmError.ErrorCode()
	}

	errBytes, err := json.Marshal(jsonRPRCError)
	if err != nil {
		logger.Error("unable to marshal error - %w", log.ErrKey, err)
		return
	}

	logger.Info(fmt.Sprintf("Forwarding %s error response from Obscuro node: %s", method, errBytes))

	if err = conn.WriteResponse(errBytes); err != nil {
		logger.Error("unable to write response back - %w", log.ErrKey, err)
	}
}

func handleError(conn userconn.UserConn, logger gethlog.Logger, err error) {
	logger.Error("error processing request - Forwarding response to user", log.ErrKey, err)

	if err = conn.WriteResponse([]byte(err.Error())); err != nil {
		logger.Error("unable to write response back", log.ErrKey, err)
	}
}
