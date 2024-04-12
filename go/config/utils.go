package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"reflect"
)

type TypeConfig int

type Config interface{}

const (
	enclave = "Enclave"
	host    = "Host"
	network = "Network"
)

const (
	Enclave TypeConfig = iota
	Host
	Network
)

func (t TypeConfig) String() string {
	return [...]string{"Enclave", "Host", "Network"}[t]
}

func ToTypeConfig(s string) (TypeConfig, error) {
	switch s {
	case enclave:
		return Enclave, nil
	case host:
		return Host, nil
	case network:
		return Network, nil
	default:
		panic("string " + s + " cannot be converted to TypeConfig.")
	}
}

func getFileMap() map[TypeConfig]string {
	return map[TypeConfig]string{
		Enclave: "./go/config/templates/default_enclave_config.yaml",
		Host:    "./go/config/templates/default_host_config.yaml",
		Network: "./go/config/templates/ITN_node_network.yaml",
	}
}

// LoadDefaultInputConfig parses optional or default configuration file and returns struct.
func LoadDefaultInputConfig(t TypeConfig) (Config, error) {
	fileMap := getFileMap()
	flagUsageMap := getFlagUsageMap()

	// set the default config from file-map
	configPath := flag.String("config", fileMap[t], flagUsageMap["configFlag"])
	overridePath := flag.String("override", "", flagUsageMap["overrideFlag"])

	// Parse only once capturing all necessary flags
	flag.Parse()

	var err error
	conf, err := loadConfigFromFile(t, *configPath)
	if err != nil {
		panic(err)
	}

	// Apply overrides if the override path is provided
	if *overridePath != "" {
		overridesConf, err := loadConfigFromFile(t, *overridePath)
		if err != nil {
			panic(err)
		}

		ApplyOverrides(conf, overridesConf)
	}

	return conf, nil
}

// loadConfigFromFile reads configuration from a file and environment variables
func loadConfigFromFile(t TypeConfig, configPath string) (Config, error) {
	var defaultConfig Config
	switch t {
	case Enclave:
		defaultConfig = &EnclaveInputConfig{}
	case Host:
		defaultConfig = &HostInputConfig{}
	case Network:
		defaultConfig = &NetworkInputConfig{}
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
