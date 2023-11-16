This package is analogous to the ethclient package in go-ethereum.

It provides a higher level, standard way to interact with an Ten network programmatically.

It aims to provide all the same methods that the geth ethclient provides for compatibility/familiarity, as well as ten-specific methods.

There are two clients, `ObsClient` and `AuthObsClient`

`ObsClient` just requires a Client and provides access to general Ten functionality that doesn't require viewing keys.

`AuthObsClient` requires a EncRPCClient, which is an RPC client with an account and a signed Viewing Key for authentication.
It provides full Ten functionality, authenticating with the node and encrypting/decrypting sensitive requests.