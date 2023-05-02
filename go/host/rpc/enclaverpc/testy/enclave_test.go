package testy

import (
	"context"
	"errors"
	"fmt"
	"github.com/obscuronet/go-obscuro/go/common/syserr"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"testing"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"

	"google.golang.org/grpc"

	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/go/host/rpc/enclaverpc/testy/generated"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

var globalError error

func TestName(t *testing.T) {
	err := StartServer()
	if err != nil {
		panic(err)
	}
	client := NewClient(gethlog.New())

	var clientSystemError error

	systemErr := syserr.New(fmt.Errorf("this an internal error"))
	globalError = status.Error(codes.Internal, systemErr.Error())
	_, err = client.ReplyWithStatusError(context.Background(), &generated.StatusRequest{})
	if err != nil {
		if status.Code(err) == codes.Internal {
			// It's an internal / systemErr it can be recreated on the clientside
			clientSystemError = syserr.New(err)
		}
		if !errors.Is(clientSystemError, &syserr.InternalError{}) {
			panic(fmt.Sprintf("err should be syserr.InternalError - err: %v", err))
		}
	}

	clientSystemError = nil

	globalError = status.Error(codes.Unknown, "This a connection error")
	_, err = client.ReplyWithStatusError(context.Background(), &generated.StatusRequest{})
	if err != nil {
		if status.Code(err) == codes.Unknown {
			// It's not systemErr it should not be recreated
			clientSystemError = syserr.New(err)
		}
		if errors.Is(err, &syserr.InternalError{}) {
			panic(fmt.Sprintf("err should be syserr.InternalError - err: %v", err))
		}
	}

	clientSystemError = nil
	globalError = status.Error(codes.Internal, "this an internal error -  but of type DBCRASH")
	details, err := status.Convert(globalError).WithDetails(
		&errdetails.ErrorInfo{
			Reason:   globalError.Error(),
			Domain:   "DBCRASH",
			Metadata: map[string]string{"whatIsThis": "this is some meta info"},
		})
	if err != nil {
		panic(err)
	}
	globalError = details.Err()

	_, err = client.ReplyWithStatusError(context.Background(), &generated.StatusRequest{})
	if err != nil {
		if status.Code(err) == codes.Internal {
			// It's an internal / systemErr it can be recreated on the clientside
			st, _ := status.FromError(err)
			fineDetails := st.Details()
			if fineDetails[0].(*errdetails.ErrorInfo).Domain == "DBCRASH" {
				fmt.Println("this is a DBCRASH system error")
				fmt.Println("Recreate it ")
			}
		}
	}

	clientSystemError = nil

	globalError = syserr.New(fmt.Errorf("this an internal error"))
	payload, err := client.ReplyWithErrorPayload(context.Background(), &generated.StatusRequest{})
	if err != nil {
		panic("no error expected")
	}

	if payload.Error != nil {
		switch payload.Error.ErrorCode {
		case 1:
			fmt.Println("this is the Internal Error")
		default:
			panic("this is can never happen")
		}
	}
}

func NewClient(logger gethlog.Logger) generated.TestProtoClient {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	connection, err := grpc.Dial("127.0.0.1:12345", opts...)
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

	return generated.NewTestProtoClient(connection)
}

func StartServer() error {
	s := &TestServer{}
	listenAddress := "127.0.0.1:12345"
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return fmt.Errorf("RPCServer could not listen on port: %w", err)
	}
	grpcServer := grpc.NewServer()

	generated.RegisterTestProtoServer(grpcServer, s)

	go func(lis net.Listener) {
		fmt.Printf("RPCServer listening on address %s.\n", listenAddress)
		err = grpcServer.Serve(lis)
		if err != nil {
			panic(err)
		}
	}(lis)

	return nil
}

type TestServer struct {
	generated.UnimplementedTestProtoServer
}

func (s *TestServer) ReplyWithStatusError(_ context.Context, _ *generated.StatusRequest) (*generated.StatusResponse, error) {
	err := globalError

	return &generated.StatusResponse{
		Status: 1,
	}, err
}

func (s *TestServer) ReplyWithErrorPayload(_ context.Context, _ *generated.StatusRequest) (*generated.StatusResponseWithError, error) {
	err := globalError
	errCode := int32(0)
	if errors.Is(err, &syserr.InternalError{}) {
		errCode = 1
	}

	return &generated.StatusResponseWithError{
		Status: 1,
		Error: &generated.Error{
			ErrString: err.Error(),
			ErrorCode: errCode,
		},
	}, nil
}
