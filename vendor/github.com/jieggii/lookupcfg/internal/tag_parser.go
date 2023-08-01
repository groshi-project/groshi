package internal

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const IgnoranceTag = "lookupcfg:\"ignore\""

type FieldProperties struct {
	Participate bool // indicates if struct field participates in all stuff that this lib does

	ValueSources       map[string]string // map of sources of value. E.g {"env": "HOST", "json": "host"}
	DefaultValue       string            // default value is stored as string because we parse it from string
	DefaultValueWasSet bool              // indicates if default value was set. (needed to check if default value was set if default value was set to "")
}

func ParseFieldTag(fieldTag reflect.StructTag) (*FieldProperties, error) {
	fieldTagString := string(fieldTag)

	fieldProperties := &FieldProperties{Participate: true}
	fieldProperties.ValueSources = make(map[string]string)

	if len(fieldTagString) == 0 || strings.Contains(fieldTagString, IgnoranceTag) {
		// skips fields without (or with empty) tags and fields with `lookupcfg:"ignore"` tag

		// todo: think about length check. Maybe it is not necessary and panic must be
		// triggered even on empty tags

		fieldProperties.Participate = false
		return fieldProperties, nil
	}

	tags := strings.Fields(fieldTagString)
	for _, tag := range tags {
		parts := strings.Split(tag, ":")
		if len(parts) != 2 {
			return nil, errors.New("invalid tag format")
		}
		key := parts[0]
		value := strings.Trim(parts[1], "\"")
		if key == "$default" {
			if fieldProperties.DefaultValueWasSet { // if default value was already set before
				return nil, errors.New("default value for this field has already been set")
			}
			// todo: check if type of default value matches field type
			fieldProperties.DefaultValue = value
			fieldProperties.DefaultValueWasSet = true
		} else {
			if _, found := fieldProperties.ValueSources[key]; found {
				return nil, fmt.Errorf("source \"%v\" for this field has already been set", key)
			}
			fieldProperties.ValueSources[key] = value
		}
	}
	return fieldProperties, nil
}
