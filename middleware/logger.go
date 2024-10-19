package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CloudPhoenix/logflake-client-go/logflake"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func NewLogger(instance *logflake.LogFlake, tracePerformance bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqTime := time.Now()
			correlation, _ := gonanoid.New()
			ww := NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				level := logflake.LevelDebug
				if ww.Status() >= 300 {
					level = logflake.LevelInfo
				}
				if ww.Status() >= 400 {
					level = logflake.LevelWarn
				}
				if ww.Status() >= 500 {
					level = logflake.LevelError
				}
				respTime := time.Now()
				d := respTime.Sub(reqTime)
				_ = instance.SendLog(logflake.Log{
					Level:       level,
					Time:        respTime,
					Correlation: correlation,
					Content:     fmt.Sprintf("%s %s status %d in %s", r.Method, r.RequestURI, ww.Status(), d.String()),
					Params: map[string]interface{}{
						"RequestMethod":   r.Method,
						"RequestURL":      r.RequestURI,
						"RequestHeaders":  r.Header,
						"ResponseStatus":  fmt.Sprintf("%d %s", ww.Status(), http.StatusText(ww.Status())),
						"ResponseHeaders": ww.Header(),
						"ResponseBytes":   ww.BytesWritten(),
					},
				})
				if tracePerformance && ww.Status() != 404 {
					_ = instance.SendPerformance(logflake.Performance{
						Time:     respTime,
						Label:    "HTTP Response",
						Duration: d.Milliseconds(),
					})
				}
			}()

			next.ServeHTTP(ww, logflake.WithCorrelationKey(r, correlation))
		}
		return http.HandlerFunc(fn)
	}
}
