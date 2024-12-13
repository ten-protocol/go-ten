package config

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Any yaml files in the default config directory will be embedded into the binary
//
//go:embed defaults/*.yaml defaults/**/*.yaml
var _baseConfig embed.FS

const _defaultBaseConfig = "defaults/0-base-config.yaml"

// LoadTenConfig reads the base config file and applying additional files provided as well as any env variables,
// returns a TenConfig struct
func LoadTenConfig(files ...string) (*TenConfig, error) {
	configFiles := []string{_defaultBaseConfig}
	configFiles = append(configFiles, files...)
	return load(configFiles)
}

// load reads and applies the config files and environment variables, returning a TenConfig struct
// This method is not publicly available as we want callers to always use the base config file to avoid gotchas.
func load(filePaths []string) (*TenConfig, error) {
	// parse yaml file with viper
	v := viper.New()

	// Bind environment variables to config keys to override yaml files
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	for i, filePath := range filePaths {
		var content []byte
		var err error
		// Check if the file exists in the embedded FS
		if content, err = _baseConfig.ReadFile(filePath); err == nil {
			// File found in embedded FS
			v.SetConfigType("yaml")
			if i == 0 { // only first file is read, the rest are 'merged'
				fmt.Println("reading embedded config file: ", filePath)
				err = v.ReadConfig(strings.NewReader(string(content)))
			} else {
				fmt.Println("merging embedded config file: ", filePath)
				err = v.MergeConfig(strings.NewReader(string(content)))
			}
		} else {
			// Otherwise, check if it exists as a file in the filesystem
			_, err = os.Stat(filePath)
			if os.IsNotExist(err) {
				fmt.Println("Config file not found: ", filePath)
				return nil, err
			}

			v.SetConfigFile(filePath)
			if i == 0 { // only first file is read, the rest are 'merged'
				fmt.Println("reading config file: ", filePath)
				err = v.ReadInConfig()
			} else {
				fmt.Println("merging config file: ", filePath)
				err = v.MergeInConfig()
			}
		}

		if err != nil {
			fmt.Println("Error reading config file: ", filePath)
			return nil, err
		}
	}

	// todo (@matt) for enclave processes apply signed configuration file **after** even the env variable overrides

	var tenCfg TenConfig
	err := v.Unmarshal(&tenCfg, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(), // handle string -> time.Duration
		mapstructure.StringToSliceHookFunc(","),     // handle string -> []string
		mapstructure.TextUnmarshallerHookFunc(),     // handle all types that implement encoding.TextUnmarshaler
		bigIntHookFunc(),                            // handle int values -> big.Int fields
	)))
	if err != nil {
		fmt.Println("Error unmarshalling config: ", err)
		return nil, err
	}

	fmt.Println("Successfully loaded Ten config.")
	tenCfg.PrettyPrint()
	return &tenCfg, nil
}
