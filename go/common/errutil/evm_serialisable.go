package errutil

// EVMSerialisableError is an API error that encompasses an EVM error with a code and a reason
type EVMSerialisableError struct {
	Err    string
	Reason interface{}
	Code   int
}

func (e EVMSerialisableError) Error() string {
	return e.Err
}

func (e EVMSerialisableError) ErrorCode() int {
	return e.Code
}

func (e EVMSerialisableError) ErrorData() interface{} {
	return e.Reason
}
