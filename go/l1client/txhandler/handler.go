package txhandler

type TxHandler interface {
}

type EthTxHandler struct {
}

func NewEthTxHandler() TxHandler {
	return &EthTxHandler{}
}
