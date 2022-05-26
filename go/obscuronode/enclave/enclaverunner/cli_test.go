package enclaverunner

import (
	"os"
	"path"
	"strconv"
	"testing"
)

func TestConfigIsParsedFromTomlFileIfConfigFlagIsPresent(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	config := fileBasedConfig(path.Join(wd, "/test.toml"))

	if expectedChainID := int64(777); config.ChainID != expectedChainID {
		t.Fatalf("config file was not parsed from TOML. Expected ChainID of %d, got %d", expectedChainID, config.ChainID)
	}
}

func TestConfigIsParsedFromCmdLineFlagsIfConfigFlagIsNotPresent(t *testing.T) {
	expectedChainID := int64(777)
	os.Args = append(os.Args, "--"+chainIDName, strconv.FormatInt(expectedChainID, 10))

	if config := flagBasedConfig(); config.ChainID != expectedChainID {
		t.Fatalf("config file was not parsed from flags. Expected ChainID of %d, got %d", expectedChainID, config.ChainID)
	}
}
