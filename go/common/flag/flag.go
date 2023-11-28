package flag

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// WrappedFlag is a construct that allows to have go flags while obeying to a new set of restrictiveness rules
type WrappedFlag struct {
	flagType     string
	ptr          any
	defaultValue any
	description  string
}

// GetString returns the flag current value cast to string
func (f WrappedFlag) GetString() string {
	return *f.ptr.(*string)
}

// GetInt64 returns the flag current value cast to int64
func (f WrappedFlag) GetInt64() int64 {
	return *f.ptr.(*int64)
}

// GetBool returns the flag current value cast to bool
func (f WrappedFlag) GetBool() bool {
	return *f.ptr.(*bool)
}

// singletonFlagger is the singleton instance of the loaded flags
var singletonFlagger = map[string]*WrappedFlag{}

// String directly uses the go flag package
func String(flagName, defaultValue, description string) *string {
	return flag.String(flagName, defaultValue, description)
}

// RestrictedString wraps the go flag package depending on the restriction mode
func RestrictedString(flagName, defaultValue, description string) *WrappedFlag {
	prtVal := new(string)
	singletonFlagger[flagName] = &WrappedFlag{
		flagType:     "string",
		ptr:          prtVal,
		defaultValue: defaultValue,
		description:  description,
	}
	return singletonFlagger[flagName]
}

// Int64 directly uses the go flag package
func Int64(flagName string, defaultValue int64, description string) *int64 {
	return flag.Int64(flagName, defaultValue, description)
}

// RestrictedInt64 wraps the go flag package depending on the restriction mode
func RestrictedInt64(flagName string, defaultValue int64, description string) *WrappedFlag {
	prtVal := new(int64)
	singletonFlagger[flagName] = &WrappedFlag{
		flagType:     "int64",
		ptr:          prtVal,
		defaultValue: defaultValue,
		description:  description,
	}
	return singletonFlagger[flagName]
}

// Bool directly uses the go flag package
func Bool(flagName string, defaultValue bool, description string) *bool {
	return flag.Bool(flagName, defaultValue, description)
}

// RestrictedBool wraps the go flag package depending on the restriction mode
func RestrictedBool(flagName string, defaultValue bool, description string) *WrappedFlag {
	prtVal := new(bool)
	singletonFlagger[flagName] = &WrappedFlag{
		flagType:     "bool",
		ptr:          prtVal,
		defaultValue: defaultValue,
		description:  description,
	}
	return singletonFlagger[flagName]
}

// Int directly uses the go flag package
func Int(flagName string, defaultValue int, description string) *int {
	return flag.Int(flagName, defaultValue, description)
}

// Uint64 directly uses the go flag package
func Uint64(flagName string, defaultValue uint64, description string) *uint64 {
	return flag.Uint64(flagName, defaultValue, description)
}

// Parse ensures the restricted mode is applied only to restricted flags
// Restricted Mode - Flags can only be inputted via ENV Vars via the enclave.json
// Non-Restricted Mode - Flags can only be inputted via normal CLI command line
func Parse() error {
	mandatoryEnvFlags := false
	val := os.Getenv("EDG_RESTRICTED")
	if val == "true" {
		fmt.Println("Using mandatory signed configurations.")
		mandatoryEnvFlags = true
	}

	for flagName, wflag := range singletonFlagger {
		// parse restricted flags if in restricted mode
		if mandatoryEnvFlags {
			err := parseMandatoryFlags(flagName, wflag)
			if err != nil {
				return fmt.Errorf("unable to parse mandatory flag: %s - %w", flagName, err)
			}
		} else {
			err := parseNonMandatoryFlag(flagName, wflag)
			if err != nil {
				return fmt.Errorf("unable to parse flag: %s - %w", flagName, err)
			}
		}
	}

	// parse all remaining flags
	flag.Parse()
	return nil
}

func parseNonMandatoryFlag(flagName string, wflag *WrappedFlag) error {
	switch wflag.flagType {
	case "string":
		wflag.ptr = flag.String(flagName, wflag.defaultValue.(string), wflag.description)
	case "int64":
		wflag.ptr = flag.Int64(flagName, wflag.defaultValue.(int64), wflag.description)
	case "bool":
		wflag.ptr = flag.Bool(flagName, wflag.defaultValue.(bool), wflag.description)
	default:
		return fmt.Errorf("unexpected type: %s", wflag.flagType)
	}
	return nil
}

func parseMandatoryFlags(flagName string, wflag *WrappedFlag) error {
	val := os.Getenv("EDG_" + strings.ToUpper(flagName))
	if val == "" {
		return fmt.Errorf("mandatory restricted flag not available - %s", flagName)
	}

	switch wflag.flagType {
	case "string":
		wflag.ptr = &val
	case "int64":
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fmt.Errorf("unable to parse flag %s - %w", flagName, err)
		}
		wflag.ptr = &i
	case "bool":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("unable to parse flag %s - %w", flagName, err)
		}
		wflag.ptr = &b
	default:
		return fmt.Errorf("unexpected mandatory type: %s", wflag.flagType)
	}
	return nil
}
