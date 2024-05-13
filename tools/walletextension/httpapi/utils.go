package httpapi

import (
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/status-im/keycard-go/hexutils"
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
func getUserID(conn UserConn) ([]byte, error) {
	// try getting userID (`token`) from query parameters and return it if successful
	userID, err := getQueryParameter(conn.ReadRequestParams(), common.EncryptedTokenQueryParameter)
	if err == nil {
		if len(userID) == common.MessageUserIDLenWithPrefix {
			return hexutils.HexToBytes(userID[2:]), nil
		} else if len(userID) == common.MessageUserIDLen {
			return hexutils.HexToBytes(userID), nil
		}

		return nil, fmt.Errorf(fmt.Sprintf("wrong length of userID from URL. Got: %d, Expected: %d od %d", len(userID), common.MessageUserIDLenWithPrefix, common.MessageUserIDLen))
	}

	return nil, fmt.Errorf("missing token field")
}

func handleError(conn UserConn, logger gethlog.Logger, err error) {
	logger.Warn("error processing request - Forwarding response to user", log.ErrKey, err)

	if err = conn.WriteResponse([]byte(err.Error())); err != nil {
		logger.Error("unable to write response back", log.ErrKey, err)
	}
}
