# Obscuroscan

A webapp that connects to a running Obscuro network, displays network statistics, and decodes encrypted rollups on 
testnet.

## Usage

* Get the client server address for an Obscuro host on a running Obscuro network

* Run `obscuroscan/main/main()` with the following flags to start Obscuroscan:

  ```--rpcServerAddress=<client server address>```

  This will create an Obscuroscan instance listening on `http://localhost:3000/`

* Visit `http://localhost:3000/`
