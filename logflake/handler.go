package logflake

import (
	"context"

	"log/slog"

	slogcommon "github.com/samber/slog-common"
)

type SlogOption struct {
	// log level (default: debug)
	Level slog.Leveler

	Instance *LogFlake

	// optional
	Converter Converter

	// optional: see slog.HandlerOptions
	AddSource   bool
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
}

func (o SlogOption) NewLogFlakeHandler() slog.Handler {
	if o.Level == nil {
		o.Level = slog.LevelDebug
	}
	if o.Instance == nil {
		panic("logflake: missing instance")
	}

	if o.Converter == nil {
		o.Converter = DefaultConverter
	}

	return &LogFlakeHandler{
		option: o,
		client: o.Instance,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

var _ slog.Handler = (*LogFlakeHandler)(nil)

type LogFlakeHandler struct {
	option SlogOption
	client *LogFlake
	attrs  []slog.Attr
	groups []string
}

func (h *LogFlakeHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
}

func (h *LogFlakeHandler) Handle(ctx context.Context, record slog.Record) error {
	newLog := h.option.Converter(h.option.AddSource, h.option.ReplaceAttr, h.attrs, h.groups, &record)
	go h.client.SendLog(newLog)
	return nil
}

func (h *LogFlakeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LogFlakeHandler{
		option: h.option,
		client: h.client,
		attrs:  slogcommon.AppendAttrsToGroup(h.groups, h.attrs, attrs...),
		groups: h.groups,
	}
}

func (h *LogFlakeHandler) WithGroup(name string) slog.Handler {
	// https://cs.opensource.google/go/x/exp/+/46b07846:slog/handler.go;l=247
	if name == "" {
		return h
	}

	return &LogFlakeHandler{
		option: h.option,
		client: h.client,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}
