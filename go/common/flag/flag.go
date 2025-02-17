package flag

import (
	"flag"
	"fmt"
)

type TenFlag struct {
	Name         string
	Value        any
	FlagType     string
	Description  string
	DefaultValue any
}

func NewStringFlag(name, defaultValue, description string) *TenFlag {
	return &TenFlag{
		Name:         name,
		Value:        "",
		FlagType:     "string",
		Description:  description,
		DefaultValue: defaultValue,
	}
}

func NewIntFlag(name string, defaultValue int, description string) *TenFlag {
	return &TenFlag{
		Name:         name,
		Value:        0,
		FlagType:     "int",
		Description:  description,
		DefaultValue: defaultValue,
	}
}

func NewBoolFlag(name string, defaultValue bool, description string) *TenFlag {
	return &TenFlag{
		Name:         name,
		Value:        false,
		FlagType:     "bool",
		Description:  description,
		DefaultValue: defaultValue,
	}
}

func NewInt64Flag(name string, defaultValue int64, description string) *TenFlag {
	return &TenFlag{
		Name:         name,
		Value:        false,
		FlagType:     "int64",
		Description:  description,
		DefaultValue: defaultValue,
	}
}

func NewUint64Flag(name string, defaultValue uint64, description string) *TenFlag {
	return &TenFlag{
		Name:         name,
		Value:        false,
		FlagType:     "uint64",
		Description:  description,
		DefaultValue: defaultValue,
	}
}

func (f TenFlag) String() string {
	if ptrVal, ok := f.Value.(*string); ok {
		return *ptrVal
	}
	return f.Value.(string)
}

func (f TenFlag) Int() int {
	if ptrVal, ok := f.Value.(*int); ok {
		return *ptrVal
	}
	return f.Value.(int)
}

func (f TenFlag) Int64() int64 {
	if ptrVal, ok := f.Value.(*int64); ok {
		return *ptrVal
	}
	return f.Value.(int64)
}

func (f TenFlag) Uint64() uint64 {
	if ptrVal, ok := f.Value.(*uint64); ok {
		return *ptrVal
	}
	return f.Value.(uint64)
}

func (f TenFlag) Bool() bool {
	if ptrVal, ok := f.Value.(*bool); ok {
		return *ptrVal
	}
	return f.Value.(bool)
}

func (f TenFlag) IsSet() bool {
	found := false
	flag.Visit(func(fl *flag.Flag) {
		if fl.Name == f.Name {
			found = true
		}
	})
	return found
}

func CreateCLIFlags(flags map[string]*TenFlag) error {
	for _, tflag := range flags {
		switch tflag.FlagType {
		case "string":
			tflag.Value = flag.String(tflag.Name, tflag.DefaultValue.(string), tflag.Description)
		case "bool":
			tflag.Value = flag.Bool(tflag.Name, tflag.DefaultValue.(bool), tflag.Description)
		case "int":
			tflag.Value = flag.Int(tflag.Name, tflag.DefaultValue.(int), tflag.Description)
		case "int64":
			tflag.Value = flag.Int64(tflag.Name, tflag.DefaultValue.(int64), tflag.Description)
		case "uint64":
			tflag.Value = flag.Uint64(tflag.Name, tflag.DefaultValue.(uint64), tflag.Description)
		default:
			return fmt.Errorf("unexpected flag type %s", tflag.FlagType)
		}
	}
	return nil
}

func Parse() {
	flag.Parse()
}
