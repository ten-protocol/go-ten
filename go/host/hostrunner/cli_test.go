package hostrunner

import (
	"os"
	"path"
	"strconv"
	"testing"
	"time"
)

const testToml = "/test.toml"

func TestConfigIsParsedFromTomlFileIfConfigFlagIsPresent(t *testing.T) {
	expectedGossipRoundNanos := time.Duration(777)
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if config := fileBasedConfig(path.Join(wd, testToml)); config.GossipRoundDuration != expectedGossipRoundNanos {
		t.Fatalf("config file was not parsed from TOML. Expected GossipRoundNanos of %d, got %d", expectedGossipRoundNanos, config.GossipRoundDuration)
	}
}

func TestConfigIsParsedFromCmdLineFlagsIfConfigFlagIsNotPresent(t *testing.T) {
	expectedGossipRoundNanos := time.Duration(666)
	os.Args = append(os.Args, "--"+gossipRoundNanosName, strconv.FormatInt(expectedGossipRoundNanos.Nanoseconds(), 10))

	if config := ParseConfig(); config.GossipRoundDuration != expectedGossipRoundNanos {
		t.Fatalf("config file was not parsed from flags. Expected GossipRoundNanos of %d, got %d", expectedGossipRoundNanos, config.GossipRoundDuration)
	}
}
