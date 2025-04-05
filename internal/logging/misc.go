package logging

import "log/slog"

func ErrAttr(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func ComponentAttr(component string) slog.Attr {
	return slog.String("component", component)
}
