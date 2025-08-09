package logger

import (
	"fmt"
	"os"
	"time"

	"load-balancer/pkg/types"

	"github.com/google/uuid"
)

var logfile string = fmt.Sprintf("logs/app_%s_%s.log",
	time.Now().Format("2006-01-02"),
	uuid.New().String()[:8],
)

func LogErr(msg string, err error) {
	logLine := fmt.Sprintf("time=%s type=ERROR msg=\"%s\" error=\"%s\"", time.Now().Format(time.RFC3339), msg, err)
	writeToFile(logLine)
}

func Log(msg string) {
	logLine := fmt.Sprintf("time=%s type=INFO msg=\"%s\"", time.Now().Format(time.RFC3339), msg)
	writeToFile(logLine)
}

func LogContainerStart(containerID string) {
	logLine := fmt.Sprintf("time=%s type=CONTAINER_START container_ID=\"%s\"", time.Now().Format(time.RFC3339), containerID)
	writeToFile(logLine)
}

func LogContainerStop(containerID string) {
	logLine := fmt.Sprintf("time=%s type=CONTAINER_STOP container_ID=\"%s\"", time.Now().Format(time.RFC3339), containerID)
	writeToFile(logLine)
}

func LogRequest(conn *types.Connection) {
	logLine := fmt.Sprintf("time=%s type=REQUEST method=%s path=\"%s\" user_agent=\"%s\"", time.Now().Format(time.RFC3339), conn.Request.Method, conn.Request.URL.Path, conn.Request.UserAgent())
	writeToFile(logLine)
}

func LogStatusCheck(status string, address string) {
	logLine := fmt.Sprintf("time=%s type=HEALTH status=%s address=\"%s\"", time.Now().Format(time.RFC3339), status, address)
	writeToFile(logLine)
}

func LogProxy(path string, proxiedTo string) {
	logLine := fmt.Sprintf("time=%s type=PROXY path=\"%s\" proxied_to=\"%s\"", time.Now().Format(time.RFC3339), path, proxiedTo)
	writeToFile(logLine)
}

func writeToFile(logLine string) {
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(logLine + "\n"); err != nil {
		fmt.Printf("Failed to write to log file: %v\n", err)
	}
}
