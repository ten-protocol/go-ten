package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"
)

// LoadFlags for configuration and as required by the service/services
func LoadFlags(t TypeConfig) error {
	flagUsageMap := FlagUsageMap()
	fs := setupConfigFlags(t, flagUsageMap)
	err := setupFlagsByType(fs, t, flagUsageMap)
	if err != nil {
		return fmt.Errorf("could not setup flags: %s", err)
	}
	if err := fs.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}
	return nil
}

// setupConfigFlags creates a FlagSet with the default config file path
// the set will process both config and override
func setupConfigFlags(t TypeConfig, flagUsageMap map[string]string) *flag.FlagSet {
	fileMap := getTemplateFilePaths()

	// set the default config from file-map
	fs := flag.NewFlagSet("Config", flag.ExitOnError)
	fs.String(ConfigFlag, fileMap[t], flagUsageMap[ConfigFlag])
	fs.String(OverrideFlag, "", flagUsageMap[OverrideFlag])
	return fs
}

// setupFlagsByType propagates flags via the correct type association
func setupFlagsByType(fs *flag.FlagSet, t TypeConfig, flagUsageMap map[string]string) error {
	switch t {
	case Enclave:
		setupFlagsFromStruct(&EnclaveInputConfig{}, fs, flagUsageMap)
	case Host:
		setupFlagsFromStruct(&HostInputConfig{}, fs, flagUsageMap)
	case Node:
		{
			setupFlagsFromStruct(&HostInputConfig{}, fs, flagUsageMap)
			setupFlagsFromStruct(&EnclaveInputConfig{}, fs, flagUsageMap)
			setupFlagsFromStruct(&NodeConfig{}, fs, flagUsageMap)
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
func setupFlagsFromStruct[T Config](p *T, fs *flag.FlagSet, usageMap map[string]string) {
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
