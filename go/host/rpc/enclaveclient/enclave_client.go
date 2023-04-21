package enclaveclient

import (
	"fmt"
	"time"

	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/go/common/rpc/generated"
	"github.com/obscuronet/go-obscuro/go/config"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	gethlog "github.com/ethereum/go-ethereum/log"
	grpc "google.golang.org/grpc"
)

func NewEnclaveRPCClient(config *config.HostConfig, logger gethlog.Logger) common.Enclave {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	connection, err := grpc.Dial(config.EnclaveRPCAddress, opts...)
	if err != nil {
		logger.Crit("Failed to connect to enclave RPC service.", log.ErrKey, err)
	}
	connection.Connect()
	// perform an initial sleep because that Connect() method is not blocking and the retry immediately checks the status
	time.Sleep(500 * time.Millisecond)

	// We wait for the RPC connection to be ready.
	err = retry.Do(func() error {
		currState := connection.GetState()
		if currState != connectivity.Ready {
			logger.Info("retrying connection until enclave is available", "status", currState.String())
			connection.Connect()
			return fmt.Errorf("connection is not ready, status=%s", currState)
		}
		// connection is ready, break out of the loop
		return nil
	}, retry.NewBackoffAndRetryForeverStrategy([]time.Duration{500 * time.Millisecond, 1 * time.Second, 5 * time.Second}, 10*time.Second))

	if err != nil {
		// this should not happen as we retry forever...
		logger.Crit("failed to connect to enclave", log.ErrKey, err)
	}

	protoClient := generated.NewEnclaveProtoClient(connection)
	return &EnclaveRPCClient{
		// this casting guarantees the new structs are fulfilling the common.Enclave interface + the RPCClient struct
		EnclaveUserClient:   NewEnclaveUserClient(protoClient, connection, config, logger).(*EnclaveUserClient),
		EnclaveSystemClient: NewEnclaveSystemClient(protoClient, connection, config, logger).(*EnclaveSystemClient),
	}
}

type EnclaveRPCClient struct {
	*EnclaveUserClient
	*EnclaveSystemClient
}
