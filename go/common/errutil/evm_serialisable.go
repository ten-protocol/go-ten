package errutil

import "fmt"

// DataError is an API error that encompasses an EVM error with a code and a reason
type DataError struct {
	Code   int         `json:"code"`
	Err    string      `json:"message"`
	Reason interface{} `json:"data,omitempty"`
}

func (e DataError) Error() string {
	return e.Err
}

func (e DataError) ErrorCode() int {
	return e.Code
}

func (e DataError) ErrorData() interface{} {
	return e.Reason
}

func (e DataError) String() string {
	return fmt.Sprintf("Data Error. Message: %s, Data: %v", e.Err, e.Reason)
}
