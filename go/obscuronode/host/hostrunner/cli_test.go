package hostrunner

import (
	"os"
	"path"
	"strconv"
	"testing"
	"time"
)

func TestConfigIsParsedFromTomlFileIfConfigFlagIsPresent(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	config := fileBasedConfig(path.Join(wd, "/test.toml"))

	expectedGossipRoundNanos := time.Duration(777)
	if config.GossipRoundDuration != expectedGossipRoundNanos {
		t.Fatalf("config file was not parsed from TOML. Expected GossipRoundNanos of %d, got %d", expectedGossipRoundNanos, config.GossipRoundDuration)
	}
}

func TestConfigIsParsedFromCmdLineRagsIfConfigFlagIsNotPresent(t *testing.T) {
	expectedGossipRoundNanos := time.Duration(666)
	os.Args = append(os.Args, "--"+gossipRoundNanosName, strconv.FormatInt(expectedGossipRoundNanos.Nanoseconds(), 10))

	if config := flagBasedConfig(); config.GossipRoundDuration != expectedGossipRoundNanos {
		t.Fatalf("config file was not parsed from flags. Expected GossipRoundNanos of %d, got %d", expectedGossipRoundNanos, config.GossipRoundDuration)
	}
}
