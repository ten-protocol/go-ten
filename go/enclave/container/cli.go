package container

import (
	"flag"
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func ParseConfig() (*config.EnclaveConfig, error) {
	inputCfg, err := config.LoadDefaultInputConfig(config.Enclave)
	if err != nil {
		return nil, fmt.Errorf("issues loading default and override config from file: %w", err)
	}
	cfg := inputCfg.(*config.EnclaveInputConfig) // assert

	fs := flag.NewFlagSet("enclave", flag.ExitOnError)
	usageMap := config.FlagUsageMap()
	config.SetupFlags(cfg, fs, usageMap)

	// Parse command-line flags
	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("error parsing flags: %w", err)
	}

	enclaveConfig, err := cfg.ToEnclaveConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to convert EnclaveInputConfig to EnclaveConfig: %w", err)
	}
	return enclaveConfig, nil
}

// retrieveOrSetEnclaveRestrictedFlags, ensures relevant flags are able to pass into the enclave in scenarios where
// an `ego sign` procedure isn't enabled - and no env[] array used. In this case, it will take the EDG_ env vars as
// first or the default config values as fall-back.
func retrieveOrSetEnclaveRestrictedFlags(cfg *config.EnclaveInputConfig) (*config.EnclaveInputConfig, error) {
	val := os.Getenv("EDG_TESTMODE")
	if val == "true" {
		fmt.Println("Using test mode flags")
		return cfg, nil
	} else {
		fmt.Println("Using mandatory signed configurations.")
	}

	v := reflect.ValueOf(cfg).Elem() // Get the reflect.Value of the struct

	for eFlag, flagType := range config.EnclaveRestrictedFlags {
		eFlag = capitalizeFirst(eFlag)
		targetEnvVar := "EDG_" + strings.ToUpper(eFlag)
		val := os.Getenv(targetEnvVar)
		if val == "" {
			fieldVal := v.FieldByName(eFlag) // Access the struct field by name
			if !fieldVal.IsValid() {
				panic("No valid field found for flag " + eFlag)
			}

			var strVal string
			switch flagType {
			case "int64":
				strVal = strconv.FormatInt(fieldVal.Int(), 10)
			case "string":
				strVal = fieldVal.String()
			case "bool":
				strVal = strconv.FormatBool(fieldVal.Bool())
			default:
				panic("Unsupported type for field " + eFlag)
			}

			if strVal == "" {
				panic("Invalid default or EDG_ for " + eFlag)
			}

			if err := os.Setenv(targetEnvVar, strVal); err != nil {
				panic("Failed to set environment variable " + targetEnvVar)
			}
			fmt.Printf("Set %s to %s from default configuration.\n", targetEnvVar, strVal)
		}
	}
	return cfg, nil
}

// capitalizeFirst capitalizes the first letter of the given string. handles mismatch between flag and config struct
func capitalizeFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
