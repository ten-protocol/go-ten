package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// HostEnvs alias
type HostEnvs = map[string]string

// EncEnvs alias
type EncEnvs = map[string]string

// Config represents structs for Input with associated flag FlagUsageMap
type Config interface{}

// IsItemInSet helper Function to check if an item is in the set
func IsItemInSet(set map[string]bool, item string) bool {
	return set[item]
}

// getTemplateFilePaths returns a map of the default static config per TypeConfig
func getTemplateFilePaths() map[TypeConfig]string {
	return map[TypeConfig]string{
		Enclave: "./go/config/templates/default_enclave_config.yaml",
		Host:    "./go/config/templates/default_host_config.yaml",
		Network: "./go/config/templates/ITN_network.yaml",
		Node:    "./go/config/templates/default_node.yaml",
	}
}

// LoadDefaultInputConfig parses optional or default configuration file and returns interface.
func LoadDefaultInputConfig(t TypeConfig, paths ConfPaths) (Config, error) {
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
	default:
		return nil, fmt.Errorf("invalid TypeConfig %s", t)
	}

	// Read YAML configuration
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
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

// ApplyOverrides is a generic function that applies non-zero value fields from the override struct 'o' to 'c'.
func ApplyOverrides[T any](c, o T) {
	cVal := reflect.ValueOf(c).Elem()
	oVal := reflect.ValueOf(o).Elem()

	// Iterate over each field in the override struct.
	for i := 0; i < oVal.NumField(); i++ {
		oField := oVal.Field(i)
		cField := cVal.Field(i)

		// Apply override if the field in 'o' is set.
		if isFieldSet(oField) {
			cField.Set(oField)
		}
	}
}

// isFieldSet determines whether the provided reflect.Value holds a non-default value.
func isFieldSet(field reflect.Value) bool {
	// Handle based on the field kind.
	switch field.Kind() {
	case reflect.Slice:
		return !field.IsNil() && field.Len() > 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() != 0
	case reflect.String:
		return field.String() != ""
	case reflect.Bool:
		return field.Bool()
	default:
		panic("unhandled default case")
	}

	// For struct or other complex types, you might need a more sophisticated approach.
	return false
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

// MergeFlagSets takes any number of FlagSets and merges them into a new FlagSet.
// Later FlagSets override earlier ones.
func MergeFlagSets(flagSets ...*flag.FlagSet) *flag.FlagSet {
	result := flag.NewFlagSet("merged", flag.ExitOnError)

	tempMap := make(map[string]*flag.Flag)

	for _, fs := range flagSets {
		fs.VisitAll(func(f *flag.Flag) {
			tempMap[f.Name] = f
		})
	}

	// Register flags in the new FlagSet based on the final values in tempMap
	for _, f := range tempMap {
		if getter, ok := f.Value.(flag.Getter); ok {
			switch v := getter.Get().(type) {
			case int:
				result.Int(f.Name, v, f.Usage)
			case string:
				result.String(f.Name, v, f.Usage)
			case bool:
				result.Bool(f.Name, v, f.Usage)
			}
		}
	}

	return result
}
