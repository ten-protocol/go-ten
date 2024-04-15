# Configurations

- `templates`: contains default static templates in yaml format
- `test`: basic testing for default configs and flags
- `node_config.go / host_config.go / enclave_config.go`
  - structs for yaml unmarshalling and run-time configuration
- `flag_usage.go`
  - basic flag usage for runtime configuration
  - assigns flags to corresponding configuration services
- `flag_utils.go`
  - processing and assignment of flags to configuration services
- `utils.go`
  - basic generic tool for overriding configuration fields with reduced-set
 yamls using `-override <filename>.yaml`

Basic behaviour of configuration applied:
1. Default `templates` are first applied as base in all scenarios. The user can customize the base configuration with `-config <filename>.yaml`
2. Sub-templates can additively be applied to a base `template` by using `-override <filename>.yaml`. The difference between `-config` and `-override` is that the former
will only be loaded once, thus if you `-config myfile.yaml` the `templates/default_<service>.yaml` will not be loaded. However, 
with `-override` the `-config` will be loaded first and then any keys in the override file will apply on-top of the `-config` file.
3. Any static `-config` or `-override` value(s) may be overridden with its corresponding `-<flagName> <flagValue>` or environment variables
`UPPERCASE_FLAGNAME=<flagValue>`. The flag values are processed in the following order: `template` -> `override` -> `flag` -> `env`