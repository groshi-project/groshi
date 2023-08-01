// Package lookupcfg allows you to define and populate
// your configs from any kind of source using lookup function!
package lookupcfg

import (
	"fmt"
	"reflect"

	"github.com/jieggii/lookupcfg/internal"
)

// Field represents struct field after reading value from the source.
type Field struct {
	StructName string // name of the field in the struct
	SourceName string // name of the field in the source

	RawValue string // value of the field, read from source
	// (may be equal to "" if it was not provided by source)
	ExpectedValueType reflect.Type // expected type of the value
}

// IncorrectTypeField represents field of incorrect type.
type IncorrectTypeField struct {
	Field
	ConversionError error // error returned by type conversion function
}

// ConfigPopulationResult represents result of config population.
// Reports about mismatches between values provided by source and provided struct.
type ConfigPopulationResult struct {
	MissingFields       []Field              // slice of fields that are missing
	IncorrectTypeFields []IncorrectTypeField // slice of fields of incorrect type
}

// PopulateConfig fills the `v`'s fields with values read from the `source` using `lookupFunction`.
func PopulateConfig(
	source string,
	lookupFunction func(string) (string, bool),
	v any,
) *ConfigPopulationResult {
	result := &ConfigPopulationResult{}

	configType := reflect.Indirect(reflect.ValueOf(v)).Type()

	for i := 0; i < configType.NumField(); i++ { // iterating over struct fields
		field := configType.Field(i)
		fieldProperties, err := internal.ParseFieldTag(field.Tag)
		if err != nil {
			panic(fmt.Errorf("error parsing %v.%v's tag: %v", configType.Name(), field.Name, err))
		}
		if !fieldProperties.Participate {
			//skip fields which do not participate
			continue
		}
		fieldValue := reflect.ValueOf(v).Elem().Field(i)

		valueName, found := fieldProperties.ValueSources[source]
		if !found { // if `source` provided as function param is not present in the struct's field tag
			panic(
				fmt.Errorf(
					"%v.%v does not have tag \"%v\" (for the source \"%v\"). Use `%v` tag if you want to ignore this field",
					configType.Name(),
					field.Name,
					source,
					source,
					internal.IgnoranceTag,
				),
			)
		}
		rawValue, found := lookupFunction(valueName)
		if !found { // if value was not received from the provided source
			if !fieldProperties.DefaultValueWasSet { // if default value of the field was not indicated
				result.MissingFields = append(result.MissingFields, Field{
					StructName:        field.Name,
					SourceName:        valueName,
					RawValue:          rawValue,
					ExpectedValueType: field.Type,
				})
				continue
			}
			rawValue = fieldProperties.DefaultValue
		}
		convertedValue, err := internal.ParseValue(rawValue, field.Type)
		if err != nil {
			result.IncorrectTypeFields = append(
				result.IncorrectTypeFields, IncorrectTypeField{
					Field: Field{
						StructName:        field.Name,
						SourceName:        valueName,
						RawValue:          rawValue,
						ExpectedValueType: field.Type,
					},
					ConversionError: err,
				},
			)
			continue
		}
		fieldValue.Set(reflect.ValueOf(convertedValue))
	}
	return result
}
