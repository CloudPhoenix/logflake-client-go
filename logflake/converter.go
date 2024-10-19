package logflake

import (
	"log/slog"

	slogcommon "github.com/samber/slog-common"
)

var SourceKey = "source"

type Converter func(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) Log

func DefaultConverter(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) Log {
	newLog := Log{
		Time:    record.Time,
		Level:   slogLevelToLogFlake(record.Level),
		Content: record.Message,
	}
	// aggregate all attributes
	attrs := slogcommon.AppendRecordAttrsToAttrs(loggerAttr, groups, record)
	// developer formatters
	attrs = slogcommon.ReplaceAttrs(replaceAttr, []string{}, attrs...)
	attrs = slogcommon.RemoveEmptyAttrs(attrs)
	// handler formatter
	newLog.Correlation, newLog.Params = attrToLogFlakeParams("", attrs)
	return newLog
}

func slogLevelToLogFlake(level slog.Level) LogLevel {
	switch level {
	case slog.LevelDebug:
		return LevelDebug
	case slog.LevelInfo:
		return LevelInfo
	case slog.LevelWarn:
		return LevelWarn
	case slog.LevelError:
		return LevelError
	default:
		return LevelDebug
	}
}

func attrToLogFlakeParams(base string, attrs []slog.Attr) (string, map[string]interface{}) {
	result := map[string]interface{}{}
	correlation := ""
	for i := range attrs {
		attr := attrs[i]
		k := base + attr.Key
		v := attr.Value
		kind := attr.Value.Kind()

		if kind == slog.KindGroup {
			_, groupAttrs := attrToLogFlakeParams(k+".", v.Group())
			for k, v := range groupAttrs {
				result[k] = v
			}
		} else {
			if k == "correlation" {
				correlation = slogcommon.ValueToString(v)
				continue
			}
			result[k] = slogcommon.ValueToString(v)
		}
	}
	return correlation, result
}
