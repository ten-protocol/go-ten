package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

// TestnetConfigCLI represents the configurations passed into the testnet over CLI
type TestnetConfigCLI struct {
	validatorEnclaveDockerImage string
	validatorEnclaveDebug       bool
	sequencerEnclaveDockerImage string
	sequencerEnclaveDebug       bool
	contractDeployerDockerImage string
	contractDeployerDebug       bool
	isSGXEnabled                bool
	logLevel                    int
}

// ParseConfigCLI returns a NodeConfigCLI based the cli params and defaults.
func ParseConfigCLI() *TestnetConfigCLI {
	cfg := &TestnetConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	validatorEnclaveDockerImage := flag.String(validatorEnclaveDockerImageFlag, getEnvOrDefault(validatorEnclaveDockerImageFlag, "testnetobscuronet.azurecr.io/obscuronet/enclave:latest"), flagUsageMap[validatorEnclaveDockerImageFlag])
	validatorEnclaveDebug := flag.Bool(validatorEnclaveDebugFlag, getEnvOrDefaultBool(validatorEnclaveDebugFlag, false), flagUsageMap[validatorEnclaveDebugFlag])
	sequencerEnclaveDockerImage := flag.String(sequencerEnclaveDockerImageFlag, getEnvOrDefault(sequencerEnclaveDockerImageFlag, "testnetobscuronet.azurecr.io/obscuronet/enclave:latest"), flagUsageMap[sequencerEnclaveDockerImageFlag])
	sequencerEnclaveDebug := flag.Bool(sequencerEnclaveDebugFlag, getEnvOrDefaultBool(validatorEnclaveDebugFlag, false), flagUsageMap[sequencerEnclaveDebugFlag])
	contractDeployerDockerImage := flag.String(contractDeployerDockerImageFlag, getEnvOrDefault(contractDeployerDockerImageFlag, "testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest"), flagUsageMap[contractDeployerDockerImageFlag])
	contractDeployerDebug := flag.Bool(contractDeployerDebugFlag, getEnvOrDefaultBool(validatorEnclaveDebugFlag, false), flagUsageMap[contractDeployerDebugFlag])
	isSGXEnabled := flag.Bool(isSGXEnabledFlag, getEnvOrDefaultBool(validatorEnclaveDebugFlag, false), flagUsageMap[isSGXEnabledFlag])
	logLevel := flag.Int(logLevelFlag, getEnvOrDefaultInt(logLevelFlag, 4), flagUsageMap[logLevelFlag])
	flag.Parse()

	cfg.validatorEnclaveDockerImage = *validatorEnclaveDockerImage
	cfg.sequencerEnclaveDockerImage = *sequencerEnclaveDockerImage
	cfg.validatorEnclaveDebug = *validatorEnclaveDebug
	cfg.sequencerEnclaveDebug = *sequencerEnclaveDebug
	cfg.contractDeployerDebug = *contractDeployerDebug
	cfg.contractDeployerDockerImage = *contractDeployerDockerImage
	cfg.isSGXEnabled = *isSGXEnabled
	cfg.logLevel = *logLevel

	return cfg
}

// getEnvOrDefault tries to get an associated environment variable; if not present, returns a default value.
func getEnvOrDefault(flagName string, defaultValue string) string {
	envName := strings.ToUpper(strings.ReplaceAll(flagName, "-", "_"))
	if value, exists := os.LookupEnv(envName); exists {
		return value
	}
	return defaultValue
}

// getEnvOrDefaultBool tries to get an associated environment variable as boolean; if not present, returns a default value.
func getEnvOrDefaultBool(flagName string, defaultValue bool) bool {
	envName := strings.ToUpper(strings.ReplaceAll(flagName, "-", "_"))
	if value, exists := os.LookupEnv(envName); exists {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			// Log the error
			log.Printf("Error parsing environment variable %s=%s as boolean: %v", envName, value, err)
			return defaultValue
		}
		return boolValue
	}
	return defaultValue
}

// getEnvOrDefault tries to get an associated environment variable as integer; if not present, returns a default value.
func getEnvOrDefaultInt(flagName string, defaultValue int) int {
	envName := strings.ToUpper(strings.ReplaceAll(flagName, "-", "_"))
	if value, exists := os.LookupEnv(envName); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			// Log the error
			log.Printf("Error parsing environment variable %s=%s as integer: %v", envName, value, err)
			return defaultValue
		}
		return intValue
	}
	return defaultValue
}
