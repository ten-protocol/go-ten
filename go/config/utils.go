package config

import (
	"embed"
	"encoding/base64"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	ConfigEnvKey   = "CONFIG_YAML_BASE64"
	OverrideEnvKey = "OVERRIDE_YAML_BASE64"
)

// Embedding the default YAML files into the binary.
//
//go:embed templates/*
var templateFS embed.FS

// HostEnvs alias
type HostEnvs = map[string]string

// EncEnvs alias
type EncEnvs = map[string]string

// Config represents structs for Input with associated flag FlagUsageMap
type Config interface{}

// getTemplateFilePaths returns a map of the default static config per TypeConfig
func getTemplateFilePaths() map[TypeConfig]string {
	return map[TypeConfig]string{
		Enclave:    "templates/default_enclave_config.yaml",
		Host:       "templates/default_host_config.yaml",
		Network:    "templates/ITN_network.yaml",
		Node:       "templates/default_node.yaml",
		Testnet:    "templates/default_testnet.yaml",
		L1Deployer: "templates/supplemental/default_l1_deployer_config.yaml",
	}
}

// LoadDefaultInputConfig parses optional or default configuration file and returns interface.
func LoadDefaultInputConfig(t TypeConfig, paths RunParams) (Config, error) {
	configPath := paths[ConfigFlag]
	overridePath := paths[OverrideFlag]
	var err error
	conf, err := LoadConfigFromFile(t, configPath)
	if err != nil {
		panic(err)
	}

	// Apply overrides if the override path is provided
	if overridePath != "" {
		overridesConf, err := LoadConfigFromFile(t, overridePath)
		if err != nil {
			panic(err)
		}

		ApplyOverrides(conf, overridesConf)
	}

	return conf, nil
}

// LoadConfigFromFile reads configuration from a file and returns as interface
func LoadConfigFromFile(t TypeConfig, configPath string) (Config, error) {
	var defaultConfig Config
	switch t {
	case Enclave:
		defaultConfig = &EnclaveInputConfig{}
	case Host:
		defaultConfig = &HostInputConfig{}
	case Network:
		defaultConfig = &NetworkInputConfig{}
	case Node:
		defaultConfig = &NodeConfig{}
	case Testnet:
		defaultConfig = &TestnetConfig{}
	case L1Deployer:
		defaultConfig = &L1ContractDeployerConfig{}
	default:
		return nil, fmt.Errorf("invalid TypeConfig %s", t)
	}

	// Read YAML configuration, Attempt to read from embedded file system first
	data, err := templateFS.ReadFile(configPath)
	if err != nil { // If not found in embedded FS, try reading from local file system
		data, err = os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("file not found in embedded FS and local FS: %s", err)
		}
	}
	err = yaml.Unmarshal(data, defaultConfig)
	if err != nil {
		return nil, err
	}

	return defaultConfig, nil
}

// WriteConfigToFile for serializing a `_InputConfig` structs. Note, using the yaml
// marshall for structs without `yaml:"key"` definitions will lose the casing, however,
// the persistence will successfully unmarshall again.
func WriteConfigToFile(c Config, filePath string) error {
	yamlStr, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// create read-only file
	err = os.WriteFile(filePath, yamlStr, 0o644) //nolint:gosec
	if err != nil {
		return err
	}
	return nil
}

// WriteConfigToEnv writes the configuration to the environment variables as serialized YAML.
// This is useful for passing configuration to docker containers (see Dockerfiles)
func WriteConfigToEnv(c Config, envs map[string]string, key string) error {
	yamlStr, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// base64 to capture newlines (docker doesn't support in env vars)
	configBase64 := base64.StdEncoding.EncodeToString(yamlStr)

	envs[strings.ToUpper(key)] = configBase64
	return nil
}

// ApplyOverrides applies the overrides FROM the 'o' struct to the 'c' struct.
func ApplyOverrides(c, o interface{}) {
	cVal := reflect.ValueOf(c).Elem()
	oVal := reflect.ValueOf(o).Elem()

	applyFieldOverrides(cVal, oVal)
}

func applyFieldOverrides(cVal, oVal reflect.Value) {
	for i := 0; i < oVal.NumField(); i++ {
		oField := oVal.Field(i)
		oFieldType := oVal.Type().Field(i)

		cField := cVal.FieldByName(oFieldType.Name)
		cFieldType, ok := cVal.Type().FieldByName(oFieldType.Name)

		// Ensure the field exists and has the same type.
		if !ok || cFieldType.Type != oFieldType.Type {
			continue
		}

		// Check if the field is a struct and not a primitive type.
		if oField.Kind() == reflect.Struct {
			// Recursively apply overrides on struct fields.
			applyFieldOverrides(cField, oField)
		} else {
			// Apply override if the field in 'o' is set (non-zero for simplicity).
			if isFieldSet(oField) {
				cField.Set(oField)
			}
		}
	}
}

// Example of isFieldSet function.
func isFieldSet(v reflect.Value) bool {
	zero := reflect.Zero(v.Type()).Interface()
	return !reflect.DeepEqual(v.Interface(), zero)
}

// GetEnvString returns key as string or fallback
func GetEnvString(key, fallback string) string {
	if value, exists := os.LookupEnv(strings.ToUpper(key)); exists {
		return value
	}
	return fallback
}

// GetEnvBool returns key as bool or fallback
func GetEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(strings.ToUpper(key)); exists {
		parsed, err := strconv.ParseBool(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

// GetEnvInt returns key as int or fallback
func GetEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(strings.ToUpper(key)); exists {
		parsed, err := strconv.Atoi(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

// GetEnvInt64 returns key as int64 or fallback
func GetEnvInt64(key string, fallback int64) int64 {
	if value, exists := os.LookupEnv(strings.ToUpper(key)); exists {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

// GetEnvUint64 returns key as uint64 or fallback
func GetEnvUint64(key string, fallback uint64) uint64 {
	if value, exists := os.LookupEnv(strings.ToUpper(key)); exists {
		parsed, err := strconv.ParseUint(value, 10, 64)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

// GetEnvUint returns key as uint or fallback
func GetEnvUint(key string, fallback uint) uint {
	if value, exists := os.LookupEnv(strings.ToUpper(key)); exists {
		parsed, err := strconv.ParseUint(value, 10, 32)
		if err == nil {
			return uint(parsed)
		}
	}
	return fallback
}

// GetEnvStringSlice returns key as string slice or fallback
func GetEnvStringSlice(key string, fallback []string) []string {
	if value, exists := os.LookupEnv(strings.ToUpper(key)); exists {
		return strings.Split(value, ",")
	}
	return fallback
}

// MergeEnvMaps takes in two maps and returns one, map2 is canonical
func MergeEnvMaps(map1, map2 map[string]string) map[string]string {
	mergedMap := make(map[string]string)
	for key, value := range map1 {
		mergedMap[strings.ToUpper(key)] = value
	}
	for key, value := range map2 {
		mergedMap[strings.ToUpper(key)] = value
	}
	return mergedMap
}
