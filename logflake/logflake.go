package logflake

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

// New Returns new LogFlake instance
func New(appKey string) *LogFlake {
	hostname, _ := os.Hostname()
	return &LogFlake{
		Server:   "https://app-test.logflake.io",
		AppKey:   appKey,
		Hostname: hostname,
	}
}

// SendLog Sends log
func (i *LogFlake) SendLog(log Log) error {
	if len(log.Hostname) == 0 {
		log.Hostname = i.Hostname
	}
	return i.sendData("logs", log)
}

// SendPerformance Sends performance
func (i *LogFlake) SendPerformance(performance Performance) error {
	return i.sendData("perf", performance)
}

// HandleRecover Try to recover and send exception
func (i *LogFlake) HandleRecover() {
	if err := recover(); err != nil {
		pc, _, _, _ := runtime.Caller(2)
		fn := runtime.FuncForPC(pc)
		i.SendLog(Log{
			Level:   LevelException,
			Content: fmt.Sprintf("%s: %s\n%s", fn.Name(), err, string(debug.Stack())),
		})
	}
}

// MeasurePerformance Start performance counter
func (i *LogFlake) MeasurePerformance(label string) *PerformanceCounter {
	return &PerformanceCounter{
		Label:    label,
		start:    time.Now(),
		instance: i,
	}
}

// Stop stops performance counter and sends the result
func (p *PerformanceCounter) Stop() int64 {
	duration := time.Since(p.start).Milliseconds()
	p.instance.SendPerformance(Performance{
		Label:    p.Label,
		Duration: duration,
	})
	return duration
}

func (i *LogFlake) sendData(dataType string, data interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}
	url := i.Server + "/api/ingestion/" + i.AppKey + "/" + dataType
	_, err = http.Post(url, "application/json", bytes.NewBuffer(j))
	return err
}
