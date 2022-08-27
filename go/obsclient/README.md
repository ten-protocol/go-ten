This package is analogous to the ethclient package in go-ethereum.

It provides a higher level, standard way to interact with an Obscuro network programmatically.

It aims to provide all the same methods that the geth ethclient provides for compatibility/familiarity, as well as obscuro-specific methods.

There are two clients, `ObsClient` and `AuthObsClient`

`ObsClient` just requires a NodeClient and provides access to general Obscuro functionality that doesn't require viewing keys.

`AuthObsClient` requires a NodeClient, an account and a Viewing Key with signature. It provides full Obscuro functionality,
authenticating with the node and encrypting/decrypting sensitive requests.