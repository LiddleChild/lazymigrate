package log

import "log/slog"

type AttributeKey string

const (
	AttributeKeySecondary AttributeKey = "secondary"
)

func AttributeSecondary() slog.Attr {
	return slog.Bool(string(AttributeKeySecondary), true)
}

func Attribute(attrs ...slog.Attr) slog.Attr {
	return slog.GroupAttrs("attributes", attrs...)
}
