All go processes in the TEN ecosystem will use the config mechanism here to load their configuration.

Ten config has a hierarchical structure, so every field has a unique field name. Some fields are used by multiple processes, 
while others are only used by a single process.

For example, `network.chainId` is a field that might be used by multiple processes, but `host.rpc.httpPort` is only used by the host processes.

The fields are defined in a single place, in the Config structs declared in `config.go`. These declarations include their type,
the mapstructure annotation that gives their yaml field name and any comments about their usage.

This library provides a loading mechanism powered by viper, it reads 0-base-config.yaml to initialise the config and then
overrides fields with any other yaml files provided and finally overrides with environment variables.

Environment variable keys are the same as the yaml keys but are uppercased and have dots replaced with underscores.
For example, `network.chainId` would be `NETWORK_CHAINID`, and `host.rpc.httpPort` would be `HOST_RPC_HTTPPORT`.

The ability to set them with environment variables is important, it allows for easy configuration of the system in a docker environment.


## Loading the config

The loading method allows an ordered list of yaml files to be provided, values in later files will override values in earlier files.

The loading method will also apply environment variables, overriding any values set in the yaml files.

In practice for production deployments, we expect that the 0-base-config.yaml will be the only file that is always provided 
and config will be provided by the orchestration engine as env variables.


## Defining a new field
Defining a new field requires adding it in 2 or 3 places:
1. The Config struct in config.go
2. The 0-base-config.yaml file
3. The enclave(-test).json file (if the field is used by the enclave process)

The config struct (1) is the source of truth for the field, it defines the type and the yaml field name. (2) and (3) are required because:

- Environment variable values are only applied by viper when the field is defined in at least one of the yaml files. For this reason,
  we want to make sure that all fields are defined in the base config file even if just with trivial or empty values.

- The ego enclave restricts the environment variables that can be accessed by the enclave process. This means that the enclave
  process will not be able to access environment variables that are not whitelisted. This is a security feature of the ego enclave.

  The enclave.json file used to produce the signed ego artifact allows environment variables to be specified, either with a hardcoded value
  which we will use for fixed constants that are not allowed to change or with a 'fromHost' flag which will allow the enclave to access the
  environment variable from the host process. This is useful for configuration values that are allowed to change between deployments.

  So any configuration value that is expected to be set by an environment variable should be whitelisted in the enclave.json file.
