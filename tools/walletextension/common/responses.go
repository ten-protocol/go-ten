package common

import (
	"errors"

	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

func CraftErrorResponse(err error) map[string]interface{} {
	errMap := make(map[string]interface{})
	respMap := make(map[string]interface{})

	respMap[JSONKeyErr] = errMap
	errMap[JSONKeyMessage] = err.Error()

	var e gethrpc.Error
	ok := errors.As(err, &e)
	if ok {
		errMap[JSONKeyCode] = e.ErrorCode()
	}

	var de gethrpc.DataError
	ok = errors.As(err, &de)
	if ok {
		errMap[JSONKeyData] = de.ErrorData()
	}

	return respMap
}
