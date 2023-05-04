package errutil

// SerialisableError is an API error that encompasses an EVM error with a code and a reason
type SerialisableError struct {
	Err    string
	Reason interface{}
	Code   int
}

func (e SerialisableError) Error() string {
	return e.Err
}

func (e SerialisableError) ErrorCode() int {
	return e.Code
}

func (e SerialisableError) ErrorData() interface{} {
	return e.Reason
}
