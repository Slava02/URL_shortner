package sl

import "log/slog"

//  TODO: найти готовую замену в sl

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
