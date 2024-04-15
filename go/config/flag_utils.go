package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"
)

type Action = string
type ConfPaths = map[string]string
type CliFlagSet = map[string]interface{}

type NodeFlagStringSet = map[string]string

// LoadFlagStrings calls LoadFlags and converts all flag values to strings.
func LoadFlagStrings(t TypeConfig) (Action, ConfPaths, NodeFlagStringSet, error) {
	action, cPaths, nodeFlags, err := LoadFlags(t, false)
	if err != nil {
		return "", nil, nil, err
	}

	// Convert CliFlagSet to NodeFlagStringSet
	flagStrings := make(NodeFlagStringSet)
	for key, value := range nodeFlags {
		// Convert each value to a string, assuming CliFlagSet is map[string]interface{}
		switch v := value.(type) {
		case string:
			flagStrings[key] = v
		case int, int64, uint, uint64, float32, float64, bool:
			flagStrings[key] = fmt.Sprintf("%v", v)
		default:
			// Handle complex types that may not be directly representable as a string
			// This could be adjusted based on what types are expected in your flags
			flagStrings[key] = fmt.Sprintf("%#v", v)
		}
	}

	return action, cPaths, flagStrings, nil
}

// LoadFlags parses flags and returns a map of flag values.
func LoadFlags(t TypeConfig, withDefaults bool) (Action, ConfPaths, CliFlagSet, error) {
	fs := SetupConfigFlags(t)

	err := setupFlagsByType(fs, t)
	if err != nil {
		return "", nil, nil, fmt.Errorf("could not setup flags: %s", err)
	}

	// Parse the Flags
	if err := fs.Parse(os.Args[1:]); err != nil {
		return "", nil, nil, fmt.Errorf("error parsing flags: %w", err)
	}

	// Only capture flags that were explicitly set
	flagValues := make(map[string]interface{})
	if withDefaults {
		fs.VisitAll(func(f *flag.Flag) {
			var value interface{}
			// Use reflection to get the actual data type of the flag
			switch f.Value.(type) {
			case flag.Getter:
				// If the value implements flag.Getter, we can safely retrieve the underlying value
				value = f.Value.(flag.Getter).Get()
			default:
				// Fallback to storing as a string if type is unknown or complex
				value = f.Value.String()
			}
			flagValues[f.Name] = value
		})
	}
	fs.Visit(func(f *flag.Flag) {
		var value interface{}
		// Use reflection to get the actual data type of the flag
		switch f.Value.(type) {
		case flag.Getter:
			// If the value implements flag.Getter, we can safely retrieve the underlying value
			value = f.Value.(flag.Getter).Get()
		default:
			// Fallback to storing as a string if type is unknown or complex
			value = f.Value.String()
		}
		flagValues[f.Name] = value
	})

	// add config paths even though not explicitly added
	cPaths := ConfPaths{
		ConfigFlag:   fs.Lookup(ConfigFlag).Value.String(),
		OverrideFlag: fs.Lookup(OverrideFlag).Value.String(),
	}

	// capture action
	action := fs.Arg(0)

	return action, cPaths, flagValues, nil
}

// SetupConfigFlags creates a FlagSet with the default config file path
// the set will process both config and override
func SetupConfigFlags(t TypeConfig) *flag.FlagSet {
	flagUsageMap := FlagUsageMap()
	fileMap := getTemplateFilePaths()

	// set the default config from file-map; ContinueOnError allows two stage parsing
	fs := flag.NewFlagSet("Config", flag.ContinueOnError)
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
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := field.Type()
		yamlTag := typ.Field(i).Tag.Get("yaml")
		if yamlTag == "" {
			continue // Skip fields without YAML tags
		}

		if _, exists := usageMap[yamlTag]; !exists {
			println(fmt.Sprintf("Missing flag usage for yaml tag '%s'", yamlTag))
		}

		if fieldType.Kind() == reflect.Struct {
			// Recurse into nested struct
			setupStructFlags(field, fs, usageMap)
		} else {
			// Set up flag for field
			flagName := yamlTag

			if isFlagSet(fs, flagName) {
				continue
			}

			flagUsage := usageMap[yamlTag]
			switch fieldType.Kind() {
			case reflect.Slice:
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
				fmt.Printf("Unsupported field type %s for field %s\n", fieldType, flagName)
			}
		}
	}
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
