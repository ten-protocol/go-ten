The `.go` files in this folder are auto-generated from the `.proto` service definition using the `protoc` Protocol 
Buffer compiler.

The files were generated using the following command:

    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative go/obscuronode/enclave/rpc/generated/enclave.proto