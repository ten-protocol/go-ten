package config2

import (
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func Load() *TenConfig {
	// parse yaml file with viper
	v := viper.New()
	v.SetConfigFile("./go/config2/main/example-config.yml")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// apply any fields from the config overrides file
	v.SetConfigFile("./go/config2/main/example-override.yml")
	err = v.MergeInConfig()
	if err != nil {
		panic(err)
	}

	// Bind environment variables to config keys
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var tenCfg TenConfig
	err = v.Unmarshal(&tenCfg, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		mapstructure.TextUnmarshallerHookFunc(),
	)))
	if err != nil {
		panic(err)
	}

	tenCfg.PrettyPrint()
	return &tenCfg
}
