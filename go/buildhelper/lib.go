package buildhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

var CONNECTION_TIMEOUT = 15 * time.Second

// EthAPI struct
// based on https://eth.wiki/json-rpc/API
type EthAPI struct {
	ApiClient *ethclient.Client
	retry     int
	timeout   time.Duration
	ipaddress string
	port      uint
	ctx       context.Context
}

func NewEthAPI(ipaddress string, port uint) *EthAPI {
	return &EthAPI{
		timeout:   CONNECTION_TIMEOUT,
		ipaddress: ipaddress,
		port:      port,
		ctx:       context.Background(),
	}
}

func (e *EthAPI) Connect() error {
	var err error
	for start := time.Now(); time.Since(start) < CONNECTION_TIMEOUT; time.Sleep(time.Second) {
		e.ApiClient, err = ethclient.Dial(fmt.Sprintf("ws://%s:%d", e.ipaddress, e.port))
		if err == nil {
			break
		}
	}

	return err
}
