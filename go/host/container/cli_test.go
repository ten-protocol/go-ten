package container

import (
	"os"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/go/config"
)

const testToml = "/test.toml"

func TestConfigIsParsedFromTomlFileIfConfigFlagIsPresent(t *testing.T) {
	p2pConnectionTimeout := time.Duration(777000000000)
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	cfg, err := fileBasedConfig(path.Join(wd, testToml))
	if err != nil {
		t.Fatalf("could not parse config. Cause: %s", err)
	}
	if cfg.P2PConnectionTimeout != p2pConnectionTimeout {
		t.Fatalf("config file was not parsed from TOML. Expected P2PConnectionTimeout of %d, got %d", p2pConnectionTimeout, cfg.P2PConnectionTimeout)
	}
}

func TestConfigIsParsedFromCmdLineFlagsIfConfigFlagIsNotPresent(t *testing.T) {
	p2pConnectionTimeout := 6 * time.Second
	os.Args = append(os.Args, "--"+p2pConnectionTimeoutSecsName, strconv.FormatInt(int64(p2pConnectionTimeout.Seconds()), 10))

	cfg, err := ParseConfig()
	if err != nil {
		t.Fatalf("could not parse config. Cause: %s", err)
	}
	if cfg.P2PConnectionTimeout != p2pConnectionTimeout {
		t.Fatalf("config file was not parsed from flags. Expected p2pConnectionTimeout of %d, got %d", p2pConnectionTimeout, cfg.P2PConnectionTimeout)
	}
}

func TestConfigFieldsMatchTomlConfigFields(t *testing.T) {
	// We get all the config fields.
	cfgReflection := reflect.TypeOf(config.HostInputConfig{})
	cfgFields := make([]string, cfgReflection.NumField())
	for i := 0; i < cfgReflection.NumField(); i++ {
		cfgFields[i] = cfgReflection.Field(i).Name
	}

	// We get all the .toml config fields.
	cfgTomlReflection := reflect.TypeOf(HostConfigToml{})
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
