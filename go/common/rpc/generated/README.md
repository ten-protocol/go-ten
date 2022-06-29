The `.go` files in this folder are auto-generated from the `.proto` service definition using the `protoc` Protocol 
Buffer compiler.

Install Protobuf and Protoc-gen-go with:

    brew install protobuf

    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

Add this to the environment, because there are some binaries in ``$HOME/go``

    export GOROOT=/usr/local/go
    export GOPATH=$HOME/go
    export GOBIN=$GOPATH/bin
    export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN

The files were generated using the following command:

    cd go/common/rpc/generated
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative enclave.proto
