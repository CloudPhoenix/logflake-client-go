package logflake

import (
	"context"
	"net/http"
)

type CtxKeys string

const LogFlakeCorrelationCtxKey CtxKeys = "LogFlakeCorrelation"

func GetCorrelation(r *http.Request) string {
	correlation, _ := r.Context().Value(LogFlakeCorrelationCtxKey).(string)
	return correlation
}

func WithCorrelationKey(r *http.Request, correlation string) *http.Request {
	r = r.WithContext(context.WithValue(r.Context(), LogFlakeCorrelationCtxKey, correlation))
	return r
}
