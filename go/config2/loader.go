package config2

import (
	"embed"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Any yaml files in the default config directory will be embedded into the binary
//
//go:embed default/*.yaml
var _baseConfig embed.FS

// Load reads and applies the config files and environment variables, returning a TenConfig struct
func Load(filePaths []string) (*TenConfig, error) {
	// parse yaml file with viper
	v := viper.New()
	var err error

	for i, filePath := range filePaths {
		// Check if the file exists in the embedded FS
		if content, err := _baseConfig.ReadFile(filePath); err == nil {
			// File found in embedded FS
			v.SetConfigType("yaml")
			if i == 0 {
				err = v.ReadConfig(strings.NewReader(string(content)))
			} else {
				err = v.MergeConfig(strings.NewReader(string(content)))
			}
		} else {
			// Otherwise, treat it as a file system path
			v.SetConfigFile(filePath)
			if i == 0 {
				err = v.ReadInConfig()
			} else {
				err = v.MergeInConfig()
			}
		}

		if err != nil {
			fmt.Println("Error reading config file: ", filePath)
			return nil, err
		}
	}

	// Bind environment variables to config keys
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// todo (@matt) for enclave processes apply signed configuration file **after** even the env variable overrides

	var tenCfg TenConfig
	err = v.Unmarshal(&tenCfg, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
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

func LoadTenConfigForEnv(env string) (*TenConfig, error) {
	return Load([]string{"default/0-base-config.yaml", fmt.Sprintf("default/1-env-%s.yaml", env)})
}
func DefaultTenConfig() (*TenConfig, error) {
	// load embedded base config
	return Load([]string{"default/0-base-config.yaml", "default/1-env-local.yaml"})
}
