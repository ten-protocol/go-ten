package buildhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

var CONNECTION_TIMEOUT = 15 * time.Second

// ethAPI struct
// based on https://eth.wiki/json-rpc/API
type ethAPI struct {
	apiClient *ethclient.Client
	retry     int
	timeout   time.Duration
	ipaddress string
	port      uint
	ctx       context.Context
}

func newEthAPI(ipaddress string, port uint) *ethAPI {
	return &ethAPI{
		timeout:   CONNECTION_TIMEOUT,
		ipaddress: ipaddress,
		port:      port,
		ctx:       context.Background(),
	}
}

func (e *ethAPI) connect() error {
	var err error
	for start := time.Now(); time.Since(start) < CONNECTION_TIMEOUT; time.Sleep(time.Second) {
		e.apiClient, err = ethclient.Dial(fmt.Sprintf("ws://%s:%d", e.ipaddress, e.port))
		if err == nil {
			break
		}
	}

	return err
}
