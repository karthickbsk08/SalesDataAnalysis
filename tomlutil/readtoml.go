package tomlutil

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/BurntSushi/toml"
)

func ReadTomlConfig(fileName string) (f interface{}) {
	// Read the TOML file
	_, err := toml.DecodeFile(fileName, &f)
	if err != nil {
		log.Println(err)
		return f
	}
	return f
}

func ReadTomlMapinAnyType(fileName string, pType *interface{}) {
	// Read the TOML file
	_, err := toml.DecodeFile(fileName, &pType)
	if err != nil {
		log.Println(err)
	}
}

// DecodeTOMLWithTypeCheck decodes TOML file into 'out' if it is a supported pointer type.
func DecodeTOMLWithTypeCheck(filepath string, out interface{}) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("out argument must be a non-nil pointer")
	}

	// Use a type switch on the concrete pointer type.
	switch out.(type) {
	case *struct{}, *map[string]interface{}:
		// These two cover most generic uses (struct or map)
		// Decode directly:
		if _, err := toml.DecodeFile(filepath, out); err != nil {
			return err
		}
		return nil

	// Add cases for pointer to specific types if needed, e.g.,
	case *string, *int, *bool, *float64:
		if _, err := toml.DecodeFile(filepath, out); err != nil {
			return err
		}
		return nil

	default:
		// Unsupported type
		return fmt.Errorf("unsupported type: %T", out)
	}
}
