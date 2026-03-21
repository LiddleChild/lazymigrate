package log

import (
	"encoding/json"
	"log/slog"
)

type LogAttribute struct {
	Secondary bool `json:"secondary,omitempty"`
}

type logAttributeSetter func(attr LogAttribute) LogAttribute

func AttributeSecondary() logAttributeSetter {
	return func(attr LogAttribute) LogAttribute {
		attr.Secondary = true
		return attr
	}
}

func Attributes(setters ...logAttributeSetter) []any {
	var attribute LogAttribute
	for _, setter := range setters {
		attribute = setter(attribute)
	}

	buffer, _ := json.Marshal(attribute)

	var attributeMap map[string]any
	_ = json.Unmarshal(buffer, &attributeMap)

	var attrs []any
	for key, value := range attributeMap {
		attrs = append(attrs, slog.Any(key, value))
	}

	return attrs
}
