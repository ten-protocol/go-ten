The `.go` files in this folder are auto-generated from the `.proto` service definition using the `protoc` Protocol 
Buffer compiler.

The files were generated using the following command:

    cd path/to/generated
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative enclave.proto`


Install Protobuf with:

    brew install protobuf

Install Protoc-gen-go with:
    
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest