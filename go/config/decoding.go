package config

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// bigIntHookFunc is a mapstructure HookFunc to handle both ints and strings for big.Int
func bigIntHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		// Check if we are trying to map into a *big.Int
		if t == reflect.TypeOf(&big.Int{}) {
			switch data := data.(type) {
			case int:
				// Convert int to *big.Int
				return big.NewInt(int64(data)), nil
			case int64:
				// Convert int64 to *big.Int
				return big.NewInt(data), nil
			case float64:
				// Convert float64 to *big.Int (if the YAML parser gives float types)
				return big.NewInt(int64(data)), nil
			case string:
				// Convert string to *big.Int
				bi := new(big.Int)
				_, ok := bi.SetString(data, 10) // Assuming base 10 for string conversion
				if !ok {
					return nil, fmt.Errorf("cannot convert %s to big.Int", data)
				}
				return bi, nil
			}
		}
		return data, nil
	}
}
