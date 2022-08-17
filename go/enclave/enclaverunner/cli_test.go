package enclaverunner

import (
	"os"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/obscuronet/go-obscuro/go/config"
)

const (
	testToml        = "/test.toml"
	expectedChainID = int64(1377)
)

func TestConfigIsParsedFromTomlFileIfConfigFlagIsPresent(t *testing.T) {
	panic("test fails")
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

func TestConfigFlagsMatchConfigFields(t *testing.T) {
	t.Skip("TODO - Reenable test when it's less disruptive to rename the CLI flags for consistency.")

	// We get all the config fields.
	cfgReflection := reflect.TypeOf(config.HostConfig{})
	cfgFields := make([]string, cfgReflection.NumField())
	for i := 0; i < cfgReflection.NumField(); i++ {
		cfgFields[i] = strings.ToLower(cfgReflection.Field(i).Name)
	}

	// We get all the CLI flags via the usages map.
	flagUsageMap := getFlagUsageMap()
	i := 0
	cliFlags := make([]string, len(flagUsageMap))
	for key := range flagUsageMap {
		cliFlags[i] = strings.ToLower(key)
		i++
	}

	sort.Strings(cfgFields)
	sort.Strings(cliFlags)

	if !reflect.DeepEqual(cfgFields, cliFlags) {
		t.Fatalf("config file supports the following fields: %s, but there are CLI flags for the following fields: %s", cfgFields, cliFlags)
	}
}
