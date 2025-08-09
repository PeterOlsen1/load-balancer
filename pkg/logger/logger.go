package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

var logfile string = fmt.Sprintf("app_%s_%s.log",
	time.Now().Format("2006-01-02"),
	uuid.New().String()[:8],
)

func LogErr(msg string, err error) {
	logLine := fmt.Sprintf("time=%s level=ERROR msg=\"%s\" error=\"%s\"", time.Now(), msg, err)
	writeToFile(logLine)
}

func Log(msg string) {
	logLine := fmt.Sprintf("time=%s level=INFO msg=\"%s\"", time.Now(), msg)
	writeToFile(logLine)
}

func LogContainerStart(containerID string) {
	logLine := fmt.Sprintf("time=%s level=CONTAINER_START containerID=\"%s\"", time.Now(), containerID)
	writeToFile(logLine)
}

func LogContainerStop(containerID string) {
	logLine := fmt.Sprintf("time=%s level=CONTAINER_STOP containerID=\"%s\"", time.Now(), containerID)
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
