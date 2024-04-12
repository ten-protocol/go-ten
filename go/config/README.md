# Configurations

- `templates`: contains default static templates in yaml format
- `test`: basic testing for default configs and flags
- `enclave/host_config.go`
  - structs for yaml unmarshalling and run-time configuration
- `utils.go`
  - basic generic tool for overriding configuration fields with reduced-set
 yamls using `-override <filename>.yaml`

Basic behaviour of configuration applied:
1. default template is first applied as base regardless - user can customize the base configuration with `-config <filename>.yaml`
2. sub-templates can additively be applied either to default template (no flag) or an initial flagged config by using `-override <filename>.yaml` the behavior is only keys matching the override will replace the existing default value
3. any single or multiple config value may be overridden with its corresponding `-<flagName> <flagValue>` at runtime.