package logger

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"load-balancer/pkg/batch"
	"load-balancer/pkg/config"
	"load-balancer/pkg/types"
	"load-balancer/pkg/workerpool"

	"github.com/google/uuid"
)

var logger Logger

func makeLogfile() (string, error) {
	err := os.MkdirAll(logger.logDir, os.ModePerm)
	if err != nil {
		fmt.Println("error making logfile", err)
		return "", err
	}

	out := fmt.Sprintf("%s/app_%s_%s.log",
		logger.logDir,
		time.Now().Format("2006-01-02"),
		uuid.New().String()[:8],
	)
	return out, nil
}

func InitLogger() {
	logger.logLevel = config.Config.Logging.Level
	logger.logDir = config.Config.Logging.Folder

	f, err := makeLogfile()
	if err != nil {
		return
	}

	//setup batch and callback function
	batch := batch.InitBatch(100, time.Millisecond*500, func(b []string) {
		logger.mu.Lock()
		defer logger.mu.Unlock()

		fRef, err := os.OpenFile(logger.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Failed to open log file: %v\n", err)
			return
		}

		for _, line := range b {
			if _, err := fRef.WriteString(line + "\n"); err != nil {
				continue
			}

			logger.linesWritten++
			if logger.linesWritten > logger.maxLines {
				newLogFile, err := makeLogfile()
				if err != nil {
					fmt.Printf("Failed to create new log file: %v\n", err)
					return
				}
				fmt.Printf("50000 log lines, starting a new logfile: %s\n", newLogFile)
				fmt.Fprintf(fRef, "%d log lines written, moving to new file %s", logger.maxLines, newLogFile)

				logger.logFile = newLogFile
				logger.linesWritten = 0

				fRef, err = os.OpenFile(logger.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Printf("Failed to open log file: %v\n", err)
					return
				}
			}
		}
	})

	//setup pool and callback function
	pool := workerpool.InitWorkerPool(10, func(logLine string) {
		batch.Add(logLine)
	})

	logCfg := config.Config.Logging
	logger = Logger{
		maxLines:   logCfg.MaxLines,
		logFile:    f,
		logLevel:   logCfg.Level,
		logDir:     logCfg.Folder,
		logBatch:   batch,
		workerPool: pool,
	}
}

func CleanupLogger() {
	logger.logBatch.Flush()
}

func Err(msg string, err error) {
	if logger.logLevel == 4 {
		return
	}

	logLine := fmt.Sprintf("time=%s type=ERROR msg=\"%s\" error=\"%s\"", time.Now().Format(time.RFC3339), msg, err)
	writeToFile(logLine)
}

func Info(msg string) {
	if logger.logLevel >= 3 {
		return
	}
	logLine := fmt.Sprintf("time=%s type=INFO msg=\"%s\"", time.Now().Format(time.RFC3339), msg)
	writeToFile(logLine)
}

func ContainerStart(containerID string) {
	if logger.logLevel >= 3 {
		return
	}
	logLine := fmt.Sprintf("time=%s type=CONTAINER_START container_ID=\"%s\"", time.Now().Format(time.RFC3339), containerID)
	writeToFile(logLine)
}

func ContainerStop(containerID string) {
	if logger.logLevel >= 3 {
		return
	}
	logLine := fmt.Sprintf("time=%s type=CONTAINER_STOP container_ID=\"%s\"", time.Now().Format(time.RFC3339), containerID)
	writeToFile(logLine)
}

func ContainerPause(containerID string) {
	if logger.logLevel >= 3 {
		return
	}
	logLine := fmt.Sprintf("time=%s type=CONTAINER_PAUSE container_ID=\"%s\"", time.Now().Format(time.RFC3339), containerID)
	writeToFile(logLine)
}

func ContainerUnpause(containerID string) {
	if logger.logLevel >= 3 {
		return
	}
	logLine := fmt.Sprintf("time=%s type=CONTAINER_UNPAUSE container_ID=\"%s\"", time.Now().Format(time.RFC3339), containerID)
	writeToFile(logLine)
}

func Request(conn *types.Connection) {
	if logger.logLevel >= 1 {
		return
	}
	logLine := fmt.Sprintf("time=%s type=REQUEST ip=%s method=%s path=\"%s\" user_agent=\"%s\"", time.Now().Format(time.RFC3339), conn.Request.RemoteAddr, conn.Request.Method, conn.Request.URL.Path, conn.Request.UserAgent())
	writeToFile(logLine)
}

func WsRequest(body []byte, ip string) {
	if logger.logLevel >= 1 {
		return
	}
	logLine := fmt.Sprintf("time=%s type=WS_MESSAGE body=\"%s\" ip=\"%s\"", time.Now().Format(time.RFC3339), string(body), ip)
	writeToFile(logLine)
}

func WsConnect(req *http.Request) {
	if logger.logLevel >= 1 {
		return
	}
	logLine := fmt.Sprintf("time=%s type=WS_CONNECT ip=%s", time.Now().Format(time.RFC3339), req.RemoteAddr)
	writeToFile(logLine)
}

func WsClose(req *http.Request) {
	if logger.logLevel >= 1 {
		return
	}
	logLine := fmt.Sprintf("time=%s type=WS_CLOSE ip=%s", time.Now().Format(time.RFC3339), req.RemoteAddr)
	writeToFile(logLine)
}

func Health(status string, address string, respTime float32) {
	if logger.logLevel >= 1 {
		return
	}
	logLine := fmt.Sprintf("time=%s type=HEALTH status=%s address=\"%s\" response_time=%f", time.Now().Format(time.RFC3339), status, address, respTime)
	writeToFile(logLine)
}

func Proxy(path string, proxiedTo string, ip string) {
	if logger.logLevel >= 2 {
		return
	}
	logLine := fmt.Sprintf("time=%s type=PROXY ip=%s path=\"%s\" proxied_to=\"%s\"", time.Now().Format(time.RFC3339), ip, path, proxiedTo)
	writeToFile(logLine)
}

func PoolSize(active int, inactive int) {
	if logger.logLevel >= 2 {
		return
	}

	logLine := fmt.Sprintf("time=%s type=INFO active=%d inactive=%d", time.Now().Format(time.RFC3339), active, inactive)
	writeToFile(logLine)
}

func writeToFile(logLine string) {
	logger.workerPool.Event(logLine)
}
