package container

import (
	"github.com/stretchr/testify/assert"
	"github.com/ten-protocol/go-ten/go/config"
	"os"
	"strings"
	"testing"
)

// Helper function to set up environment for test
func setupEnv(key, value string) {
	err := os.Setenv(key, value)
	if err != nil {
		return
	}
}

// Helper function to cleanup environment
func cleanupEnv(envVars map[string]string) {
	for eFlag, _ := range envVars {
		targetEnvVar := "EDG_" + strings.ToUpper(eFlag)
		unSetHelper(targetEnvVar)
	}
}

func unSetHelper(enVar string) {
	err := os.Unsetenv(enVar)
	if err != nil {
		return
	}
}

// Test function for the retrieveOrSetEnclaveRestrictedFlags behavior in test mode
func TestRetrieveOrSetEnclaveRestrictedFlags_TestMode(t *testing.T) {
	// Setup
	setupEnv("EDG_TESTMODE", "true")
	defer unSetHelper("EDG_TESTMODE")

	cfg := &config.EnclaveInputConfig{} // Assume this struct is properly defined
	expectedCfg := &config.EnclaveInputConfig{}

	// Execution
	resultCfg, err := retrieveOrSetEnclaveRestrictedFlags(cfg)

	// Assertion - all EGD_<RESTRICTED> env_vars should be unset
	for eFlag, _ := range enclaveRestrictedFlags {
		targetEnvVar := "EDG_" + strings.ToUpper(eFlag)
		val := os.Getenv(targetEnvVar)
		assert.Equal(t, "", val, targetEnvVar+" should not be set.")
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedCfg, resultCfg)
}

// Test function for setting default configuration values when environment variables are not set
func TestRetrieveOrSetEnclaveRestrictedFlags_DefaultValues(t *testing.T) {
	// Setup
	cfg := &config.EnclaveInputConfig{
		L1ChainID:             123,
		TenChainID:            1337,
		TenGenesis:            "abc",
		UseInMemoryDB:         true,
		ProfilerEnabled:       false,
		DebugNamespaceEnabled: true,
	}
	expectedCfg := &config.EnclaveInputConfig{
		L1ChainID:             123,
		TenChainID:            1337,
		TenGenesis:            "abc",
		UseInMemoryDB:         true,
		ProfilerEnabled:       false,
		DebugNamespaceEnabled: true,
	}

	// Cleanup any relevant environment variable to simulate the scenario
	cleanupEnv(enclaveRestrictedFlags)       // run before
	defer cleanupEnv(enclaveRestrictedFlags) // after test

	// Execution
	resultCfg, err := retrieveOrSetEnclaveRestrictedFlags(cfg)

	// Assertion - all EGD_<RESTRICTED> env_vars should be set to the default values.
	for eFlag, _ := range enclaveRestrictedFlags {
		targetEnvVar := "EDG_" + strings.ToUpper(eFlag)
		val := os.Getenv(targetEnvVar)
		switch strings.ToUpper(eFlag) {
		case "L1CHAINID":
			assert.Equal(t, "123", val)
		case "TENCHAINID":
			assert.Equal(t, "1337", val)
		case "TENGENESIS":
			assert.Equal(t, "abc", val)
		case "USEINMEMORYDB":
			assert.Equal(t, "true", val)
		case "PROFILERENABLED":
			assert.Equal(t, "false", val)
		case "DEBUGNAMESPACEENABLED":
			assert.Equal(t, "true", val)
		}
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedCfg, resultCfg)
}

// Test function to ensure that environment variables are correctly set from an override EDG_<flag>
func TestRetrieveOrSetEnclaveRestrictedFlags_SetEnvVars(t *testing.T) {
	// Setup
	cfg := &config.EnclaveInputConfig{
		L1ChainID:             123,
		TenChainID:            1337,
		TenGenesis:            "abc",
		UseInMemoryDB:         true,
		ProfilerEnabled:       false,
		DebugNamespaceEnabled: true,
	}
	expectedCfg := &config.EnclaveInputConfig{
		L1ChainID:             123,
		TenChainID:            1337,
		TenGenesis:            "abc",
		UseInMemoryDB:         true,
		ProfilerEnabled:       false,
		DebugNamespaceEnabled: true,
	}

	// Cleanup any relevant environment variable to simulate the scenario
	cleanupEnv(enclaveRestrictedFlags) // run before

	// Overrides
	setupEnv("EDG_TENCHAINID", "888")
	setupEnv("EDG_USEINMEMORYDB", "false")

	defer cleanupEnv(enclaveRestrictedFlags) // after test

	// Execution
	resultCfg, _ := retrieveOrSetEnclaveRestrictedFlags(cfg)

	// Assertion - all EGD_<RESTRICTED> env_vars should be set with partial overrides.
	for eFlag, _ := range enclaveRestrictedFlags {
		targetEnvVar := "EDG_" + strings.ToUpper(eFlag)
		val := os.Getenv(targetEnvVar)
		switch strings.ToUpper(eFlag) {
		case "L1CHAINID":
			assert.Equal(t, "123", val)
		case "TENCHAINID":
			assert.Equal(t, "888", val) // here
		case "TENGENESIS":
			assert.Equal(t, "abc", val)
		case "USEINMEMORYDB":
			assert.Equal(t, "false", val) // here
		case "PROFILERENABLED":
			assert.Equal(t, "false", val)
		case "DEBUGNAMESPACEENABLED":
			assert.Equal(t, "true", val)
		}
	}
	assert.Equal(t, expectedCfg, resultCfg)
}

// Test function to ensure that ALL environment variables are set (either through EDB_ or a default)
func TestRetrieveOrSetEnclaveRestrictedFlags_AllRequired(t *testing.T) {
	cfg := &config.EnclaveInputConfig{}

	// Cleanup any relevant environment variable to simulate the scenario
	cleanupEnv(enclaveRestrictedFlags) // run before
	defer cleanupEnv(enclaveRestrictedFlags)

	// Overrides
	setupEnv("EDG_TENCHAINID", "888")
	setupEnv("EDG_USEINMEMORYDB", "false")

	// Expect the function to panic. Capture the panic
	defer func() {
		if r := recover(); r != nil {
			// r should contain the panic message
			assert.Contains(t, r, "Invalid default or EDG_ for", "The panic message should indicate invalid default")
		} else {
			t.Errorf("Expected a panic but did not occur")
		}
	}()

	// Execute the function that is expected to panic
	_, _ = retrieveOrSetEnclaveRestrictedFlags(cfg)
}
