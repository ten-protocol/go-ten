package config

import (
	"flag"
	enclavecontainer "github.com/ten-protocol/go-ten/go/enclave/container"
	hostcontainer "github.com/ten-protocol/go-ten/go/host/container"
	"os"
	"testing"
)

const defaultHost = "/../templates/default_host_config.yaml"
const defaultEnclave = "/../templates/default_enclave_config.yaml"
const overrideConfig = "/partial.yaml"

// Same mechanism for host and enclave
func TestConfigIsParsedFromYamlFileIfConfigFlagIsPresent(t *testing.T) {
	resetFlagSet()

	l1WebsocketURL := "ws://127.0.0.1:8546"
	logLevel := 3

	// Back up the original os.Args to be available after unit test
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Mock os.Args
	os.Args = []string{"your-program", "-config", wd + defaultHost}

	cfg, err := hostcontainer.ParseConfig()
	if err != nil {
		t.Fatalf("could not parse config. Cause: %s", err)
	}
	if cfg.L1WebsocketURL != l1WebsocketURL || cfg.LogLevel != logLevel {
		t.Fatalf("config file was not parsed from YAML. Expected l1WebsockerURL of %s"+
			"and logLevel %d, got %s and %d", l1WebsocketURL, logLevel, cfg.L1WebsocketURL, cfg.LogLevel)
	}
}

// The default config will set the regular values including logLevel 3, override will swap logLevel
func TestOverrideAdditiveReplacementOfDefaultConfig(t *testing.T) {
	resetFlagSet()

	gasBatchExecutionLimit := uint64(300_000_000_000)
	logLevel := 2

	// Back up the original os.Args to be available after unit test
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Mock os.Args
	os.Args = []string{
		"your-program",
		"-config", wd + defaultEnclave,
		"-override", wd + overrideConfig,
	}

	cfg, err := enclavecontainer.ParseConfig()
	if err != nil {
		t.Fatalf("could not parse config. Cause: %s", err)
	}
	if cfg.GasBatchExecutionLimit != gasBatchExecutionLimit {
		t.Fatalf("config file was not parsed from YAML. Expected gasBatchExecutionLimit of %d, got %d", gasBatchExecutionLimit, cfg.GasBatchExecutionLimit)
	}
	if cfg.LogLevel != logLevel {
		t.Fatalf("override failed, logLevel of default was 3 but should have overriden to %d", logLevel)
	}
}

// needed for subsequent runs testing flags
func resetFlagSet() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
