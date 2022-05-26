package enclaverunner

import (
	"os"
	"path"
	"strconv"
	"testing"
)

const testToml = "/test.toml"

func TestConfigIsParsedFromTomlFileIfConfigFlagIsPresent(t *testing.T) {
	expectedChainID := int64(777) //nolint:ifshort
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if config := fileBasedConfig(path.Join(wd, testToml)); config.ChainID != expectedChainID {
		t.Fatalf("config file was not parsed from TOML. Expected ChainID of %d, got %d", expectedChainID, config.ChainID)
	}
}

func TestConfigIsParsedFromCmdLineFlagsIfConfigFlagIsNotPresent(t *testing.T) {
	expectedChainID := int64(777)
	os.Args = append(os.Args, "--"+chainIDName, strconv.FormatInt(expectedChainID, 10))

	if config := ParseConfig(); config.ChainID != expectedChainID {
		t.Fatalf("config file was not parsed from flags. Expected ChainID of %d, got %d", expectedChainID, config.ChainID)
	}
}
