package logflake

import (
	"errors"
	"log/slog"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/CloudPhoenix/logflake-client-go/logflake"
	"github.com/CloudPhoenix/logflake-client-go/middleware"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

var correlation string

func init() {
	correlation, _ = gonanoid.New()
}

func getInstance() *logflake.LogFlake {
	l := logflake.New(os.Getenv("LOGFLAKE_TEST"))
	l.Server = "https://app-test.logflake.io"
	return l
}

func TestLogs(t *testing.T) {
	i := getInstance()

	err := i.SendLog(logflake.Log{
		Content:     "Test Log (No Level)",
		Correlation: correlation,
	})
	if err != nil {
		t.Error(err)
	}

	err = i.SendLog(logflake.Log{
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
		Correlation: correlation,
	})
	if err != nil {
		t.Error(err)
	}

	err = i.SendLog(logflake.Log{
		Content:     "Test Log (Warn Level)",
		Level:       logflake.LevelWarn,
		Correlation: correlation,
	})
	if err != nil {
		t.Error(err)
	}

	err = i.SendLog(logflake.Log{
		Content:     "Test Log (Error Level)",
		Level:       logflake.LevelError,
		Correlation: correlation,
	})
	if err != nil {
		t.Error(err)
	}

	err = i.SendLog(logflake.Log{
		Content:     "Test Log (Fatal Level)",
		Level:       logflake.LevelFatal,
		Correlation: correlation,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestException(t *testing.T) {
	i := getInstance()
	defer i.HandleRecover(correlation)
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
	for range [5]int{} {
		if _, err := http.Get(svr.URL); err != nil {
			t.Error(err)
			return
		}
	}
}

func TestPerformance(t *testing.T) {
	i := getInstance()

	err := i.SendPerformance(logflake.Performance{
		Label:    "Test",
		Duration: 100,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestPerformanceCounter(t *testing.T) {
	i := getInstance()
	p := i.MeasurePerformance("Counter")
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // Simulate work
	p.Stop()
}

func TestSLog(t *testing.T) {
	startTime := time.Now()
	i := getInstance()
	logger := slog.New(logflake.SlogOption{Level: slog.LevelDebug, Instance: i}.NewLogFlakeHandler())
	slog.SetDefault(logger)
	slog.Debug("SLog test debug")
	slog.Info("SLog test info")
	slog.Warn("SLog test warn")
	slog.Error("SLog test error")
	slog.Info(
		"SLog test info with params",
		slog.String("correlation", correlation),
		slog.Duration("testDuration", time.Since(startTime)),
		slog.Group("testGroup", slog.String("testString", "Test"), slog.Int("testInt", 123)),
	)
	slog.Info("slog test with correlation",
		slog.String("correlation", correlation),
	)
}
