package middleware

import (
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"

	"github.com/CloudPhoenix/logflake-client-go/logflake"
)

func NewRecoverer(instance *logflake.LogFlake) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil && err != http.ErrAbortHandler {
					pc, _, _, _ := runtime.Caller(2)
					fn := runtime.FuncForPC(pc)
					instance.SendLog(logflake.Log{
						Correlation: logflake.GetCorrelation(r),
						Level:       logflake.LevelException,
						Content:     fmt.Sprintf("%s: %s\n%s", fn.Name(), err, string(debug.Stack())),
					})
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
