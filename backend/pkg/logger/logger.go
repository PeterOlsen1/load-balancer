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

func Err(msg string, err error) {
	logLine := fmt.Sprintf("time=%s type=ERROR msg=\"%s\" error=\"%s\"", time.Now().Format(time.RFC3339), msg, err)
	writeToFile(logLine)
}

func Info(msg string) {
	logLine := fmt.Sprintf("time=%s type=INFO msg=\"%s\"", time.Now().Format(time.RFC3339), msg)
	writeToFile(logLine)
}

func ContainerStart(containerID string) {
	logLine := fmt.Sprintf("time=%s type=CONTAINER_START container_ID=\"%s\"", time.Now().Format(time.RFC3339), containerID)
	writeToFile(logLine)
}

func ContainerStop(containerID string) {
	logLine := fmt.Sprintf("time=%s type=CONTAINER_STOP container_ID=\"%s\"", time.Now().Format(time.RFC3339), containerID)
	writeToFile(logLine)
}

func Request(conn *types.Connection) {
	logLine := fmt.Sprintf("time=%s type=REQUEST method=%s path=\"%s\" user_agent=\"%s\"", time.Now().Format(time.RFC3339), conn.Request.Method, conn.Request.URL.Path, conn.Request.UserAgent())
	writeToFile(logLine)
}

func WsRequest(body []byte) {
	logLine := fmt.Sprintf("time=%s type=WS_MESSAGE body=\"%s\"", time.Now().Format(time.RFC3339), string(body))
	writeToFile(logLine)
}

func StatusCheck(status string, address string) {
	logLine := fmt.Sprintf("time=%s type=HEALTH status=%s address=\"%s\"", time.Now().Format(time.RFC3339), status, address)
	writeToFile(logLine)
}

func Proxy(path string, proxiedTo string) {
	logLine := fmt.Sprintf("time=%s type=PROXY path=\"%s\" proxied_to=\"%s\"", time.Now().Format(time.RFC3339), path, proxiedTo)
	writeToFile(logLine)
}

func writeToFile(logLine string) {
	os.MkdirAll("logs", os.ModePerm)
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
