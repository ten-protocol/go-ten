package config

import (
	"os"
	"path"
	"reflect"
	"sort"
	"strconv"
	"testing"
	"time"
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
	os.Args = append(os.Args, "--"+p2pConnectionTimeoutSecsFlag, strconv.FormatInt(int64(p2pConnectionTimeout.Seconds()), 10))

	cfg, err := ParseHostConfig()
	if err != nil {
		t.Fatalf("could not parse config. Cause: %s", err)
	}
	if cfg.P2PConnectionTimeout != p2pConnectionTimeout {
		t.Fatalf("config file was not parsed from flags. Expected p2pConnectionTimeout of %d, got %d", p2pConnectionTimeout, cfg.P2PConnectionTimeout)
	}
}

func TestConfigFieldsMatchTomlConfigFields(t *testing.T) {
	// We get all the config fields.
	cfgReflection := reflect.TypeOf(HostConfig{})
	cfgFields := make([]string, cfgReflection.NumField())
	for i := 0; i < cfgReflection.NumField(); i++ {
		cfgFields[i] = cfgReflection.Field(i).Name
	}

	// We get all the .toml config fields.
	cfgTomlReflection := reflect.TypeOf(HostFileConfig{})
	cfgTomlFields := make([]string, cfgTomlReflection.NumField())
	for i := 0; i < cfgTomlReflection.NumField(); i++ {
		cfgTomlFields[i] = cfgTomlReflection.Field(i).Name
	}

	sort.Strings(cfgFields)
	sort.Strings(cfgTomlFields)

	nonIntersectingFields := make([]string, 0)
	i, j := 0, 0
	for i < len(cfgFields) && j < len(cfgTomlFields) {
		if cfgFields[i] > cfgTomlFields[j] {
			nonIntersectingFields = append(nonIntersectingFields, cfgTomlFields[j])
			j++
		} else if cfgFields[i] < cfgTomlFields[j] {
			nonIntersectingFields = append(nonIntersectingFields, cfgFields[i])
			i++
		} else {
			i++
			j++
		}
	}

	if len(nonIntersectingFields) > 0 {
		t.Fatalf("non-intersecting fields found: %s", nonIntersectingFields)
	}
}
