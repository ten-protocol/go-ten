package config

import (
	"flag"
	"fmt"
	"reflect"
)

const (
	OverrideFlag = "override"
	ConfigFlag   = "config"
)

// GetConfigFlagUsageMap is for overarching flags pointing to configuration files
func GetConfigFlagUsageMap() map[string]string {
	return map[string]string{
		OverrideFlag: "Additive config file to apply on top of default or -config",
		ConfigFlag:   "The path to the host's config file. Overrides all other flags",
	}
}

// SetupFlags iterates through a Config struct and usageMap and creates properly typed flags
// for each entry. Struct yaml-key must match a key in usageMap. No need to manually assign the
// flag value after parse, the flag links the associated struct pointer.
// Note: Flags will be assigned default value as EnvVar > parameter.
func SetupFlags[T Config](p T, fs *flag.FlagSet, usageMap map[string]string) {
	val := reflect.ValueOf(p).Elem()
	typ := val.Type()

	if val.Kind() != reflect.Struct {
		panic("SetupFlags only accepts struct types")
	}

	// Check for any missing flag definitions
	for i := 0; i < val.NumField(); i++ {
		yamlTag := typ.Field(i).Tag.Get("yaml")
		if yamlTag == "" {
			continue // skip fields without yaml tags
		}
		if _, exists := usageMap[yamlTag]; !exists {
			panic(fmt.Sprintf("Missing flag usage for yaml tag '%s'", yamlTag))
		}
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := field.Type()
		yamlTag := typ.Field(i).Tag.Get("yaml")
		if yamlTag == "" {
			continue // Skip fields without YAML tags
		}

		flagName := yamlTag
		flagUsage := usageMap[yamlTag]

		// Attach flag based on type - GetEnv<Type> makes the default prefer environment vars as priority
		switch fieldType.Kind() {
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
