package config

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// TenConfig is the top-level configuration struct for all TEN services.
// Fields are organised hierarchically, not all fields and structures will be used by all services.
//
// The default config is defined in yaml files formatted hierarchically like:
// ```
// network:
//
//	chainId: 1234
//	l1Contracts:
//	  managementContract: "0x1234567890abcdef1234567890abcdef12345678"
//	  ...
//
// node:
//
//	nodeType: "validator"
//	...
//
// ```
// The config can be overridden by other yaml files formatted the same.
// But fields can also be overridden directly per yaml spec, e.g.:
// ```
// network.chainId: 5678
// node.nodeType: "sequencer"
// ```
//
// Fields can also be overridden by environment variables, with the key being the flattened path of the field, like:
// ```
// export NETWORK_CHAINID=5678
// export NODE_NODETYPE=sequencer
// ```
// For ease of reading only the top-level struct is defined in this file, the nested structs are defined in their own files.
type TenConfig struct {
	Network *NetworkConfig `mapstructure:"network"`
	Node    *NodeConfig    `mapstructure:"node"`
	Host    *HostConfig    `mapstructure:"host"`
	Enclave *EnclaveConfig `mapstructure:"enclave"`
}

//
// Note: Just TenConfig utility functions below here.
// All the top-level nested config structs are defined in their own files.
//

func (t *TenConfig) PrettyPrint() {
	output, err := yaml.Marshal(t)
	if err != nil {
		fmt.Printf("Error printing config as YAML: %v\n", err)
		return
	}
	fmt.Println(string(output))
}

// ToEnvironmentVariables converts the config structure into environment variables map
func (t *TenConfig) ToEnvironmentVariables() map[string]string {
	return structToEnvMap("", t)
}

// ToEnvironmentVariablesRecursive recursively converts the config structure into environment variables map
func structToEnvMap(prefix string, cfg interface{}) map[string]string {
	envMap := make(map[string]string)
	value := reflect.ValueOf(cfg)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return envMap
	}

	valType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := valType.Field(i)
		mapKey := fieldType.Tag.Get("mapstructure")

		if mapKey == "" {
			mapKey = fieldType.Name
		}

		// Convert mapstructure key to uppercase with underscores
		envKey := strings.ToUpper(mapKey)
		if prefix != "" {
			envKey = prefix + "_" + envKey
		}

		// Handle *big.Int explicitly before the switch
		if field.Type() == reflect.TypeOf(new(big.Int)) {
			if !field.IsNil() {
				envMap[envKey] = field.Interface().(*big.Int).String()
			}
			continue
		} else if field.Type() == reflect.TypeOf(big.Int{}) {
			ptrBigInt := field.Addr().Interface().(*big.Int)
			envMap[envKey] = ptrBigInt.String()
			continue
		}

		switch field.Kind() {
		case reflect.Struct:
			// Recursively handle nested structures
			nestedMap := structToEnvMap(envKey, field.Interface())
			for k, v := range nestedMap {
				envMap[k] = v
			}
		case reflect.Slice:
			// Handle string slices as comma-separated strings
			if field.Type().Elem().Kind() == reflect.String {
				strSlice := field.Interface().([]string)
				envMap[envKey] = strings.Join(strSlice, ",")
			} else {
				// Handle other slice types, if needed
				envMap[envKey] = fmt.Sprintf("%v", field.Interface())
			}
		case reflect.Ptr:
			if !field.IsNil() {
				// Handle pointer types
				nestedMap := structToEnvMap(envKey, field.Interface())
				for k, v := range nestedMap {
					envMap[k] = v
				}
			}
		default:
			// Handle basic types
			if !field.CanInterface() {
				fmt.Println("Field cannot be interfaced:", prefix, envKey)
			}
			envMap[envKey] = fmt.Sprintf("%v", field.Interface())
		}
	}
	return envMap
}
