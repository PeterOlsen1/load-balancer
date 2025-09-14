package logger

import (
	"load-balancer/pkg/batch"
	"load-balancer/pkg/workerpool"
	"os"
	"sync"
)

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
	logFileRef   *os.File
	logLevel     uint8
	logDir       string
	logBatch     *batch.Batch[string]
	workerPool   *workerpool.WorkerPool[string]
}
