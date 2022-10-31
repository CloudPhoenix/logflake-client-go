package logflake

import (
	"errors"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/CloudPhoenix/logflake-client-go/logflake"
	"github.com/CloudPhoenix/logflake-client-go/middleware"
)

func getInstance() *logflake.LogFlake {
	return logflake.New(os.Getenv("LOGFLAKE_TEST"))
}

func TestLogs(t *testing.T) {
	i := getInstance()

	i.SendLog(logflake.Log{
		Content:     "Test Log (No Level)",
		Correlation: "test",
	})

	i.SendLog(logflake.Log{
		Content: "Test Log (Info with Params)",
		Level:   logflake.LevelInfo,
		Params: map[string]interface{}{
			"String": "Test",
			"Number": 123,
			"Bool":   true,
			"JSON": map[string]interface{}{
				"a": 1,
				"b": "c",
				"d": true,
			},
		},
		Correlation: "test",
	})

	i.SendLog(logflake.Log{
		Content: "Test Log (Warn Level)",
		Level:   logflake.LevelWarn,
	})

	i.SendLog(logflake.Log{
		Content: "Test Log (Error Level)",
		Level:   logflake.LevelError,
	})

	i.SendLog(logflake.Log{
		Content: "Test Log (Fatal Level)",
		Level:   logflake.LevelFatal,
	})
}

func TestException(t *testing.T) {
	i := getInstance()
	defer i.HandleRecover()

	panic("Testing")
}

func TestMiddleware(t *testing.T) {
	i := getInstance()
	response := 100
	svr := httptest.NewServer(
		middleware.NewLogger(i, true)(
			middleware.NewRecoverer(i)(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // Simulate work
					if response >= 500 {
						panic(errors.New("example panic"))
					}
					w.WriteHeader(response)
					response += 100
				}))))
	defer svr.Close()
	for r := 0; r < 5; r++ {
		if _, err := http.Get(svr.URL); err != nil {
			t.Error(err)
			return
		}
	}
}

func TestPerformance(t *testing.T) {
	i := getInstance()

	i.SendPerformance(logflake.Performance{
		Label:    "Test",
		Duration: 100,
	})
}

func TestPerformanceCounter(t *testing.T) {
	i := getInstance()

	p := i.MeasurePerformance("Counter")
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // Simulate work
	p.Stop()
}
