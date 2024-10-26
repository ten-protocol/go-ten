All go processes in the TEN ecosystem will use the config mechanism here to load their configuration.

We have a TenConfig hierarchical structure that contains typed configuration fields used across all our systems.

They are defined in a single place, in the Config structs declared in `config.go`. These declarations include their type,
the mapstructure annotation that gives their yaml field name and any comments about their usage.

For example, `network.chainId` is a field that might be used by multiple processes, but `host.rpc.httpPort` is only used by the host processes.

This library provides a loading mechanism powered by viper, it reads 0-base-config.yaml to initialise the config and then
overrides fields with any other yaml files provided and finally overrides with environment variables.

Environment variable keys are the same as the yaml keys but are uppercased and have dots replaced with underscores.
For example, `network.chainId` would be `NETWORK_CHAINID`, and `host.rpc.httpPort` would be `HOST_RPC_HTTPPORT`.

## Defining a new field
When defining a new field in the Config struct, there are a few things to keep in mind:

- Environment variable values are only applied when the field is defined in at least one of the yaml files. For this reason, 
  we want to make sure that all fields are defined in the base config file even if just with trivial or empty values.

- The ego enclave restricts the environment variables that can be accessed by the enclave process. This means that the enclave
  process will not be able to access environment variables that are not whitelisted. This is a security feature of the ego enclave.

  The enclave.json file used to produce the signed ego artifact allows environment variables to be specified, either with a hardcoded value
  which we will use for fixed constants that are not allowed to change or with a 'fromHost' flag which will allow the enclave to access the
  environment variable from the host process. This is useful for configuration values that are allowed to change between deployments.

  So any configuration value that is expected to be set by an environment variable should be whitelisted in the enclave.json file.

The upshot of this is that defining a new field requires adding it in 2 or 3 places.
The nested Config struct declared from config.go, the 0-base-config.yaml file and, if field is used by the enclave, the enclave.json file.

## Loading the config

The loading method allows an ordered list of yaml files to be provided, values in later files will override values in earlier files.

In practice for production deployments, we expect that the 0-base-config.yaml will be the only file that is always provided, and the rest will be optional.

As a convention, we're using a prefixes for the files when lists are provided:
`0-base-config.yaml` - the base configuration that is always provided
`1-env-<env-name>.yaml` - environment-specific configuration
`2-node-<node-name>.yaml` - node-specific configuration (e.g. sequencer, validator in integration tests)
`3-network.yaml` - network-specific configuration (e.g. chainId, networkId)

The loader is not performing any verification on these conventions though, callers can provide files with any names and they are applied in the order provided.