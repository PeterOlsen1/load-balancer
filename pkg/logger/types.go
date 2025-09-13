package logger

import "sync"

type LogLevel uint8

const (
	All LogLevel = iota
	NoRequests
	NoProxy
	NoInfo
	None
)

type Logger struct {
	mu           sync.Mutex
	maxLines     uint32
	linesWritten uint32
	logFile      string
	logLevel     uint8
	logDir       string
}
