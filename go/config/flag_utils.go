package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type RunParams = map[string]string
type CliFlagSet = map[string]interface{}
type CliFlagStringSet = map[string]string

func LoadFlagStrings(t TypeConfig) (RunParams, CliFlagStringSet, error) {
	rParams, nodeFlags, err := LoadFlags(t, false)
	if err != nil {
		return nil, nil, err
	}

	flagStrings := convertToFlagStringSet(nodeFlags)
	return rParams, flagStrings, nil
}

func LoadFlags(t TypeConfig, withDefaults bool) (RunParams, CliFlagSet, error) {
	fs := SetupConfigFlags(t)
	err := setupFlagsByType(fs, t)
	if err != nil {
		return nil, nil, err
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	flagValues := captureFlagValues(fs, withDefaults)
	rParams := getRunParams(fs)

	for key, _ := range rParams { // no need to keep the config and override paths in the flagValues
		if flagValues[key] != nil {
			delete(flagValues, key)
		}
	}

	os.Args = removeFlagsFromArgs(os.Args, rParams) // remove flags related to config from args

	action := os.Args[len(os.Args)-1] // last element is args
	rParams["action"] = action

	return rParams, flagValues, nil
}

// EnvOrFlag iterates across the program args, if any args in form `-<arg> val` or `-<arg>=val` are found the
// key <ARG> to upper is checked against existing environment variable keys. If there is a match, the
// arg and its associated value are removed from the args, the environment variable is left unmodified,
// however if the environment variable is not set, it will throw error.
func EnvOrFlag(args []string) ([]string, error) {
	for i := 0; i < len(args); i++ {
		eqDelimiter := false
		if strings.HasPrefix(args[i], "-") {
			arg := strings.TrimLeft(args[i], "-")
			if strings.Contains(arg, "=") {
				eqDelimiter = true
				arg = strings.Split(arg, "=")[0]
			}

			if val, ok := os.LookupEnv(strings.ToUpper(arg)); ok {
				if val == "" {
					return nil, fmt.Errorf("env var set with no value: %s", arg)
				}
				if eqDelimiter {
					args = append(args[:i], args[i+1:]...)
					i--
					continue
				}
				args = append(args[:i], args[i+2:]...)
				i--
			}
		}
	}
	return args, nil
}

func removeFlagsFromArgs(args []string, flagsToRemove map[string]string) []string {
	for key := range flagsToRemove {
		for i := 0; i < len(args); i++ {
			// Check if the argument is a flag with a value in the form `-<flag>=something`
			if strings.HasPrefix(args[i], "-"+key+"=") {
				args = append(args[:i], args[i+1:]...)
				i--
			} else if args[i] == "-"+key && i+1 < len(args) {
				// Check if the argument is a flag with a value in the form `-<flag> something`
				args = append(args[:i], args[i+2:]...)
				i--
			}
		}
	}
	return args
}

// convertToFlagStringSet helper for LoadFlagStrings
func convertToFlagStringSet(nodeFlags CliFlagSet) CliFlagStringSet {
	flagStrings := make(CliFlagStringSet)
	for key, value := range nodeFlags {
		flagStrings[key] = fmt.Sprintf("%v", value)
	}
	return flagStrings
}

// getRunParams adds the config and override paths to a map
func getRunParams(fs *flag.FlagSet) RunParams {
	return RunParams{
		DryRunFlag:   fs.Lookup(DryRunFlag).Value.String(),
		ConfigFlag:   fs.Lookup(ConfigFlag).Value.String(),
		OverrideFlag: fs.Lookup(OverrideFlag).Value.String(),
	}
}

// captureFlagValues helper for LoadFlags
func captureFlagValues(fs *flag.FlagSet, withDefaults bool) CliFlagSet {
	flagValues := make(map[string]interface{})
	visitor := func(f *flag.Flag) {
		value := getFlagValue(f)
		flagValues[f.Name] = value
	}

	if withDefaults {
		fs.VisitAll(visitor)
	} else {
		fs.Visit(visitor)
	}

	return flagValues
}

// getFlagValue helper for captureFlagValues
func getFlagValue(f *flag.Flag) interface{} {
	switch v := f.Value.(type) {
	case flag.Getter:
		return v.Get()
	default:
		return f.Value.String()
	}
}

// SetupConfigFlags creates a FlagSet with the default config file path
// the set will process both config and override
func SetupConfigFlags(t TypeConfig) *flag.FlagSet {
	flagUsageMap := FlagUsageMap()
	fileMap := getTemplateFilePaths()

	// set the default config from file-map; ContinueOnError allows two stage parsing
	fs := flag.NewFlagSet("Config", flag.ContinueOnError)
	fs.Bool(DryRunFlag, false, flagUsageMap[DryRunFlag])
	fs.String(ConfigFlag, fileMap[t], flagUsageMap[ConfigFlag])
	fs.String(OverrideFlag, "", flagUsageMap[OverrideFlag])
	return fs
}

// setupFlagsByType propagates flags via the correct type association
func setupFlagsByType(fs *flag.FlagSet, t TypeConfig) error {
	flagUsageMap := FlagUsageMap()
	switch t {
	case Enclave:
		SetupFlagsFromStruct(&EnclaveInputConfig{}, fs, flagUsageMap)
	case Host:
		SetupFlagsFromStruct(&HostInputConfig{}, fs, flagUsageMap)
	case Node:
		{
			SetupFlagsFromStruct(&HostInputConfig{}, fs, flagUsageMap)
			SetupFlagsFromStruct(&EnclaveInputConfig{}, fs, flagUsageMap)
			SetupFlagsFromStruct(&NodeConfig{}, fs, flagUsageMap)
		}
	case Testnet:
		{
			SetupFlagsFromStruct(&HostInputConfig{}, fs, flagUsageMap)
			SetupFlagsFromStruct(&EnclaveInputConfig{}, fs, flagUsageMap)
			SetupFlagsFromStruct(&NodeConfig{}, fs, flagUsageMap)
			SetupFlagsFromStruct(&Eth2NetworkConfig{}, fs, flagUsageMap)
			SetupFlagsFromStruct(&L1ContractDeployerConfig{}, fs, flagUsageMap)
			SetupFlagsFromStruct(&L2ContractDeployerConfig{}, fs, flagUsageMap)

		}
	case L1Deployer:
		SetupFlagsFromStruct(&L1ContractDeployerConfig{}, fs, flagUsageMap)
	default:
		return fmt.Errorf("unknown TypeConfig %s", t.String())
	}
	return nil
}

// SetupFlagsFromStruct iterates through a Config struct and usageMap and creates properly typed flags
// for each entry. Struct yaml-key must match a key in usageMap. No need to manually assign the
// flag value after parse, the flag links the associated struct pointer.
// Note: Flags will be assigned default value as EnvVar > parameter.
func SetupFlagsFromStruct[T Config](p *T, fs *flag.FlagSet, usageMap map[string]string) {
	val := reflect.ValueOf(p).Elem()
	if val.Kind() != reflect.Struct {
		panic("SetupFlagsFromStruct only accepts struct types")
	}
	setupStructFlags(val, fs, usageMap)
}

func setupStructFlags(val reflect.Value, fs *flag.FlagSet, usageMap map[string]string) {
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		yamlTag := val.Type().Field(i).Tag.Get("yaml")
		if yamlTag == "" || isFlagSet(fs, yamlTag) {
			continue
		}

		if _, exists := usageMap[yamlTag]; !exists {
			// exclude the following tags from the check
			if yamlTag != "networkConfig" && yamlTag != "nodeDetails" && yamlTag != "nodeSettings" && yamlTag != "nodeImages" {
				println(fmt.Sprintf("Missing flag usage for yaml tag '%s'", yamlTag))
			}
		}

		if field.Type().Kind() == reflect.Struct {
			setupStructFlags(field, fs, usageMap)
		} else {
			setupFlag(field, fs, yamlTag, usageMap[yamlTag])
		}
	}
}

// setupFlag assigns a flag to a field based on its type prefers environment variable over default
func setupFlag(field reflect.Value, fs *flag.FlagSet, flagName, flagUsage string) {
	switch field.Type().Kind() {
	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.String {
			fs.Var(newStringSliceValue(GetEnvStringSlice(flagName, field.Interface().([]string)), field.Addr().Interface().(*[]string)), flagName, flagUsage)
		} else {
			fmt.Printf("Unsupported slice type %s for field %s\n", field.Type(), flagName)
		}
	case reflect.String:
		fs.StringVar(field.Addr().Interface().(*string), flagName, GetEnvString(flagName, field.Interface().(string)), flagUsage)
	case reflect.Int:
		fs.IntVar(field.Addr().Interface().(*int), flagName, GetEnvInt(flagName, field.Interface().(int)), flagUsage)
	case reflect.Int64:
		fs.Int64Var(field.Addr().Interface().(*int64), flagName, GetEnvInt64(flagName, field.Interface().(int64)), flagUsage)
	case reflect.Uint:
		fs.UintVar(field.Addr().Interface().(*uint), flagName, GetEnvUint(flagName, field.Interface().(uint)), flagUsage)
	case reflect.Uint64:
		fs.Uint64Var(field.Addr().Interface().(*uint64), flagName, GetEnvUint64(flagName, field.Interface().(uint64)), flagUsage)
	case reflect.Bool:
		fs.BoolVar(field.Addr().Interface().(*bool), flagName, GetEnvBool(flagName, field.Interface().(bool)), flagUsage)
	default:
		fmt.Printf("Unsupported field type %s for field %s\n", field.Type(), flagName)
	}
}

type stringSliceValue []string

func newStringSliceValue(val []string, p *[]string) *stringSliceValue {
	*p = val
	return (*stringSliceValue)(p)
}

func (s *stringSliceValue) Set(val string) error {
	*s = stringSliceValue(strings.Split(val, ","))
	return nil
}

func (s *stringSliceValue) String() string {
	return strings.Join(*s, ",")
}

// isFlagSet is used to check if a flag has been defined (incl. before parse)
func isFlagSet(fs *flag.FlagSet, fName string) bool {
	found := false
	fs.VisitAll(func(fl *flag.Flag) {
		if fl.Name == fName {
			found = true
		}
	})
	return found
}
