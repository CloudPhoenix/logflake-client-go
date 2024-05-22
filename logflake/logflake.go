package logflake

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/golang/snappy"
)

// New Returns new LogFlake instance
func New(appKey string) *LogFlake {
	hostname, _ := os.Hostname()
	i := &LogFlake{
		Server:            "https://app.logflake.io",
		AppKey:            appKey,
		Hostname:          hostname,
		EnableCompression: true,
	}
	i.updateUA()
	return i
}

func (i *LogFlake) updateUA() string {
	i.userAgent = fmt.Sprintf("logflake-client-go/%s-%s/%s", runtime.GOOS, runtime.GOARCH, runtime.Version())
	return i.userAgent
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
		_ = i.SendLog(Log{
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
	_ = p.instance.SendPerformance(Performance{
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

	var req *http.Request

	if i.EnableCompression {
		// Encode to Base64
		var encoded bytes.Buffer
		encoder := base64.NewEncoder(base64.StdEncoding, &encoded)
		if _, err := encoder.Write(j); err != nil {
			return err
		}
		if err := encoder.Close(); err != nil {
			return err
		}
		// Compress with Snappy
		compressed := snappy.Encode(nil, encoded.Bytes())
		// Prepare compressed request
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(compressed))
		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "application/octet-stream")
	} else {
		// Prepare request
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(j))
		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("Accept", "application/json")
	if i.userAgent == "" {
		i.updateUA()
	}
	req.Header.Add("User-Agent", i.userAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
