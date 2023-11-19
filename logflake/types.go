package logflake

import "time"

// LogFlake struct
type LogFlake struct {
	Server            string
	AppKey            string
	Hostname          string
	EnableCompression bool
}

// LogLevel Log Level
type LogLevel int

const (
	LevelDebug     LogLevel = 0
	LevelInfo      LogLevel = 1
	LevelWarn      LogLevel = 2
	LevelError     LogLevel = 3
	LevelFatal     LogLevel = 4
	LevelException LogLevel = 5
)

// Log struct
type Log struct {
	Time        time.Time              `json:"time,omitempty"`
	Hostname    string                 `json:"hostname,omitempty"`
	Level       LogLevel               `json:"level,omitempty"`
	Correlation string                 `json:"correlation,omitempty"`
	Content     string                 `json:"content,omitempty"`
	Params      map[string]interface{} `json:"params,omitempty"`
}

// Performance struct
type Performance struct {
	Time     time.Time `json:"time,omitempty"`
	Label    string    `json:"label,omitempty"`
	Duration int64     `json:"duration,omitempty"`
}

// PerformanceCounter struct
type PerformanceCounter struct {
	Label    string
	start    time.Time
	instance *LogFlake
}
