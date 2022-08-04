package enclaverunner

import (
	"github.com/obscuronet/go-obscuro/go/config"
	"os"
	"path"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

const (
	testToml        = "/test.toml"
	expectedChainID = int64(1377)
)

func TestConfigIsParsedFromTomlFileIfConfigFlagIsPresent(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if cfg := fileBasedConfig(path.Join(wd, testToml)); cfg.L1ChainID != expectedChainID {
		t.Fatalf("config file was not parsed from TOML. Expected L1ChainID of %d, got %d", expectedChainID, cfg.L1ChainID)
	}
}

func TestConfigIsParsedFromCmdLineFlagsIfConfigFlagIsNotPresent(t *testing.T) {
	os.Args = append(os.Args, "--"+l1ChainIDName, strconv.FormatInt(expectedChainID, 10))

	if cfg := ParseConfig(); cfg.L1ChainID != expectedChainID {
		t.Fatalf("config file was not parsed from flags. Expected L1ChainID of %d, got %d", expectedChainID, cfg.L1ChainID)
	}
}

func TestConfigFieldsMatchTomlConfigFields(t *testing.T) {
	// We get all the config fields.
	cfgReflection := reflect.TypeOf(config.EnclaveConfig{})
	cfgFields := make([]string, cfgReflection.NumField())
	for i := 0; i < cfgReflection.NumField(); i++ {
		cfgFields[i] = cfgReflection.Field(i).Name
	}

	// We get all the .toml config fields.
	cfgTomlReflection := reflect.TypeOf(EnclaveConfigToml{})
	cfgTomlFields := make([]string, cfgTomlReflection.NumField())
	for i := 0; i < cfgTomlReflection.NumField(); i++ {
		cfgTomlFields[i] = cfgTomlReflection.Field(i).Name
	}

	sort.Strings(cfgFields)
	sort.Strings(cfgTomlFields)

	if !reflect.DeepEqual(cfgFields, cfgTomlFields) {
		t.Fatalf("config file supports the following fields: %s, but .toml config file supports the following fields: %s", cfgFields, cfgTomlFields)
	}
}
