package internal

import (
	"errors"
	"reflect"
	"strings"
)

const IgnoranceTag = "lookupcfg:\"ignore\""

type FieldMeta struct {
	Participate bool // indicates if this field participates in all stuff that this lib does

	ValueSources       map[string]string // map of sources of value. E.g {"env": "HOST", "json": "host"}
	DefaultValue       string            // default value is stored as string because we parse it from string
	DefaultValueWasSet bool              // indicates if default value was set. (needed to check if default value was set if default value was set to "")
}

func ParseFieldTag(fieldTag reflect.StructTag) (error, *FieldMeta) {
	fieldTagString := string(fieldTag)

	fieldMeta := &FieldMeta{Participate: true}
	fieldMeta.ValueSources = make(map[string]string)

	if len(fieldTagString) == 0 || strings.Contains(fieldTagString, IgnoranceTag) {
		// skips fields without (or with empty) tags and fields with `lookupcfg:"ignore"` tag

		// todo: think about length check. Maybe it is not necessary and panic must be
		// triggered even on empty tags

		fieldMeta.Participate = false
		return nil, fieldMeta
	}

	tags := strings.Split(fieldTagString, " ")
	for _, tag := range tags {
		parts := strings.Split(tag, ":")
		if len(parts) != 2 {
			return errors.New("invalid tag format"), nil
		}
		key := parts[0]
		value := strings.Trim(parts[1], "\"")
		if key == "$default" {
			// todo: check if user tries to set default value multiple times
			// todo: check if type of default value matches field type
			fieldMeta.DefaultValue = value
			fieldMeta.DefaultValueWasSet = true
		} else {
			// todo: check if key already exists
			fieldMeta.ValueSources[key] = value
		}
	}
	return nil, fieldMeta
}
