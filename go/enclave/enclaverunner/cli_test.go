package enclaverunner

import (
	"os"
	"path"
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

	if config := fileBasedConfig(path.Join(wd, testToml)); config.L1ChainID != expectedChainID {
		t.Fatalf("config file was not parsed from TOML. Expected L1ChainID of %d, got %d", expectedChainID, config.L1ChainID)
	}
}

func TestConfigIsParsedFromCmdLineFlagsIfConfigFlagIsNotPresent(t *testing.T) {
	os.Args = append(os.Args, "--"+l1ChainIDName, strconv.FormatInt(expectedChainID, 10))

	if config := ParseConfig(); config.L1ChainID != expectedChainID {
		t.Fatalf("config file was not parsed from flags. Expected L1ChainID of %d, got %d", expectedChainID, config.L1ChainID)
	}
}
