package lookupcfg

import (
	"github.com/jieggii/lookupcfg/internal"
	"reflect"
)

// Field represents struct field after reading value from the source
type Field struct {
	StructName string // name of the field in the struct
	SourceName string // name of the field in the source

	RawValue string // value of the field, read from source
	// (may be equal to "" if it was not provided by source)
	ExpectedValueType reflect.Type // expected type of the value
}

type IncorrectTypeField struct {
	Field
	ConversionError error // error returned by type conversion function
}

type ConfigPopulationResult struct {
	MissingFields       []Field              // list of fields that are missing
	IncorrectTypeFields []IncorrectTypeField // array of fields of incorrect type
}

// PopulateConfig fills the `object`'s fields with values read from the `source` using `lookupFunction`
func PopulateConfig(
	source string,
	lookupFunction func(string) (string, bool),
	object interface{},
) *ConfigPopulationResult {
	result := &ConfigPopulationResult{}

	configType := reflect.Indirect(reflect.ValueOf(object)).Type()

	for i := 0; i < configType.NumField(); i++ { // iterating over struct fields
		field := configType.Field(i)
		err, fieldMeta := internal.ParseFieldTag(field.Tag)
		if err != nil {
			internal.Panicf("Error parsing %v.%v's tag: %v", configType.Name(), field.Name, err)
		}
		if !fieldMeta.Participate {
			//skip fields which do not participate
			continue
		}
		fieldValue := reflect.ValueOf(object).Elem().Field(i)

		valueName, ok := fieldMeta.ValueSources[source]
		if !ok { // if `source` provided as function param is not present in the struct's field tag
			internal.Panicf(
				"%v.%v does not have tag \"%v\" (for the source \"%v\"). Use `%v` tag if you want to ignore this field.",
				configType.Name(),
				field.Name,
				source,
				source,
				internal.IgnoranceTag,
			)
		}
		value, ok := lookupFunction(valueName)
		if !ok { // if value was not received from the provided source
			if fieldMeta.DefaultValue == "" { // if default value of the field was not indicated
				result.MissingFields = append(result.MissingFields, Field{
					StructName:        field.Name,
					SourceName:        valueName,
					RawValue:          value,
					ExpectedValueType: field.Type,
				})
				continue
			}
			value = fieldMeta.DefaultValue
		}
		convertedValue, err := internal.Parse(value, field.Type)
		if err != nil {
			result.IncorrectTypeFields = append(
				result.IncorrectTypeFields, IncorrectTypeField{
					Field: Field{
						StructName:        field.Name,
						SourceName:        valueName,
						RawValue:          value,
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
