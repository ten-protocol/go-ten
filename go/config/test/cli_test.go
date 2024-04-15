package config

import (
	"flag"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/config"
	enclavecontainer "github.com/ten-protocol/go-ten/go/enclave/container"
	hostcontainer "github.com/ten-protocol/go-ten/go/host/container"
	"os"
	"testing"
)

const defaultHost = "/../templates/default_host_config.yaml"
const defaultEnclave = "/../templates/default_enclave_config.yaml"
const overrideConfig = "/partial.yaml"

// Same mechanism for host and enclave
func TestHostConfigIsParsedFromYamlFileIfConfigFlagIsPresent(t *testing.T) {
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

	_, cPaths, _, err := config.LoadFlagStrings(config.Host)
	cfg, err := hostcontainer.ParseConfig(cPaths)
	if err != nil {
		t.Fatalf("could not parse config. Cause: %s", err)
	}
	if cfg.L1WebsocketURL != l1WebsocketURL || cfg.LogLevel != logLevel {
		t.Fatalf("config file was not parsed from YAML. Expected l1WebsockerURL of %s"+
			"and logLevel %d, got %s and %d", l1WebsocketURL, logLevel, cfg.L1WebsocketURL, cfg.LogLevel)
	}
}

// The default config will set the regular values including logLevel 3, override will swap logLevel
func TestEnclaveOverrideAdditiveReplacementOfDefaultConfig(t *testing.T) {
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

	_, cPaths, _, err := config.LoadFlagStrings(config.Enclave)
	cfg, err := enclavecontainer.ParseConfig(cPaths)
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

func TestHostFlagOverridesDefaultProperty(t *testing.T) {
	resetFlagSet()

	// Back up the original os.Args to be available after unit test
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	os.Args = []string{
		"your-program",
		"-config", wd + defaultHost,
		"-nodeType", "sequencer",
	}

	_, cPaths, _, err := config.LoadFlagStrings(config.Host)
	cfg, err := hostcontainer.ParseConfig(cPaths)
	if err != nil {
		t.Fatalf("could not parse config. Cause: %s", err)
	}

	// Assert that the flag value overrides the default configuration
	if cfg.NodeType != common.Sequencer {
		t.Fatalf("default config was not loaded. Expected nodeType of %s, got %s", common.Validator, cfg.NodeType)
	}
}

func TestEnclaveEnvVarOverridesDefaultConfigAndFlag(t *testing.T) {
	resetFlagSet()

	// Back up the original os.Args to be available after unit test
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Set environment variable which will override below flag for nodeType
	err := os.Setenv("NODETYPE", "validator")
	if err != nil {
		t.Fatalf("could not set environment variable. Cause: %s", err)
	}
	err = os.Setenv("LOGLEVEL", "2")
	if err != nil {
		t.Fatalf("could not set environment variable. Cause: %s", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	os.Args = []string{
		"your-program",
		"-config", wd + defaultEnclave,
		"-nodeType", "sequencer", // flag to be overridden by env var in space format
		"-logLevel=2",           // flag to be overridden by env var in = format
		"-logPath", "/tmp/logs", // keep override config because no envVar
	}

	_, cPaths, _, err := config.LoadFlagStrings(config.Enclave)
	cfg, err := enclavecontainer.ParseConfig(cPaths)
	if err != nil {
		t.Fatalf("could not parse config. Cause: %s", err)
	}

	// Assert that the flag value overrides the default configuration
	if cfg.NodeType != common.Validator {
		t.Fatalf("env override not successful. Expected nodeType of %s, got %s", common.Validator, cfg.NodeType)
	}
	if cfg.LogLevel != 2 {
		t.Fatalf("env override not successful. Expected logLevel of %d, got %d", 2, cfg.LogLevel)
	}
	if cfg.LogPath != "/tmp/logs" {
		t.Fatalf("flag override failed. Expected logPath of %s, got %s", "/tmp/logs", cfg.LogPath)
	}
}

// needed for subsequent runs testing flags
func resetFlagSet() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
