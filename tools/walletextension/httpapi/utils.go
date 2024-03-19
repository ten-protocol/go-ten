package httpapi

import (
	"fmt"
	"strings"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
)

func getQueryParameter(params map[string]string, selectedParameter string) (string, error) {
	value, exists := params[selectedParameter]
	if !exists {
		return "", fmt.Errorf("parameter '%s' is not in the query params", selectedParameter)
	}

	return value, nil
}

// getUserID returns userID from query params / url of the URL
// it always first tries to get userID from a query parameter `u` or `token` (`u` parameter will become deprecated)
// if it fails to get userID from a query parameter it tries to get it from the URL and it needs position as the second parameter
func getUserID(conn UserConn, userIDPosition int) (string, error) {
	// try getting userID (`token`) from query parameters and return it if successful
	userID, err := getQueryParameter(conn.ReadRequestParams(), common.EncryptedTokenQueryParameter)
	if err == nil {
		if len(userID) != common.MessageUserIDLen {
			return "", fmt.Errorf(fmt.Sprintf("wrong length of userID from URL. Got: %d, Expected: %d", len(userID), common.MessageUserIDLen))
		}
		return userID, err
	}

	// try getting userID(`u`) from query parameters and return it if successful
	userID, err = getQueryParameter(conn.ReadRequestParams(), common.UserQueryParameter)
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

func handleError(conn UserConn, logger gethlog.Logger, err error) {
	logger.Warn("error processing request - Forwarding response to user", log.ErrKey, err)

	if err = conn.WriteResponse([]byte(err.Error())); err != nil {
		logger.Error("unable to write response back", log.ErrKey, err)
	}
}
