package logger

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"load-balancer/pkg/config"
	"load-balancer/pkg/types"

	"github.com/google/uuid"
)

var logfile string
var ll string = config.Config.Logging.Level
var logDir string = config.Config.Logging.Folder

func makeLogfile() (string, error) {
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		fmt.Println("error making logfile", err)
		return "", err
	}

	out := fmt.Sprintf("%s/app_%s_%s.log",
		logDir,
		time.Now().Format("2006-01-02"),
		uuid.New().String()[:8],
	)
	return out, nil
}

func init() {
	f, err := makeLogfile()
	if err != nil {
		return
	}
	logfile = f
}

func Err(msg string, err error) {
	if ll == "none" {
		return
	}

	logLine := fmt.Sprintf("time=%s type=ERROR msg=\"%s\" error=\"%s\"", time.Now().Format(time.RFC3339), msg, err)
	writeToFile(logLine)
}

func Info(msg string) {
	if ll == "error" || ll == "none" {
		return
	}
	logLine := fmt.Sprintf("time=%s type=INFO msg=\"%s\"", time.Now().Format(time.RFC3339), msg)
	writeToFile(logLine)
}

func ContainerStart(containerID string) {
	if ll != "all" {
		return
	}
	logLine := fmt.Sprintf("time=%s type=CONTAINER_START container_ID=\"%s\"", time.Now().Format(time.RFC3339), containerID)
	writeToFile(logLine)
}

func ContainerStop(containerID string) {
	if ll != "all" {
		return
	}
	logLine := fmt.Sprintf("time=%s type=CONTAINER_STOP container_ID=\"%s\"", time.Now().Format(time.RFC3339), containerID)
	writeToFile(logLine)
}

func ContainerPause(containerID string) {
	if ll != "all" {
		return
	}
	logLine := fmt.Sprintf("time=%s type=CONTAINER_PAUSE container_ID=\"%s\"", time.Now().Format(time.RFC3339), containerID)
	writeToFile(logLine)
}

func ContainerUnpause(containerID string) {
	if ll != "all" {
		return
	}
	logLine := fmt.Sprintf("time=%s type=CONTAINER_UNPAUSE container_ID=\"%s\"", time.Now().Format(time.RFC3339), containerID)
	writeToFile(logLine)
}

func Request(conn *types.Connection) {
	if ll != "all" {
		return
	}
	logLine := fmt.Sprintf("time=%s type=REQUEST ip=%s method=%s path=\"%s\" user_agent=\"%s\"", time.Now().Format(time.RFC3339), conn.Request.RemoteAddr, conn.Request.Method, conn.Request.URL.Path, conn.Request.UserAgent())
	writeToFile(logLine)
}

func WsRequest(body []byte, ip string) {
	if ll != "all" {
		return
	}
	logLine := fmt.Sprintf("time=%s type=WS_MESSAGE body=\"%s\" ip=\"%s\"", time.Now().Format(time.RFC3339), string(body), ip)
	writeToFile(logLine)
}

func WsConnect(req *http.Request) {
	if ll != "all" {
		return
	}
	logLine := fmt.Sprintf("time=%s type=WS_CONNECT ip=%s", time.Now().Format(time.RFC3339), req.RemoteAddr)
	writeToFile(logLine)
}

func WsClose(req *http.Request) {
	if ll != "all" {
		return
	}
	logLine := fmt.Sprintf("time=%s type=WS_CLOSE ip=%s", time.Now().Format(time.RFC3339), req.RemoteAddr)
	writeToFile(logLine)
}

func Health(status string, address string, respTime float32) {
	if ll != "all" {
		return
	}
	logLine := fmt.Sprintf("time=%s type=HEALTH status=%s address=\"%s\" response_time=%f", time.Now().Format(time.RFC3339), status, address, respTime)
	writeToFile(logLine)
}

func Proxy(path string, proxiedTo string, ip string) {
	if ll != "all" {
		return
	}
	logLine := fmt.Sprintf("time=%s type=PROXY ip=%s path=\"%s\" proxied_to=\"%s\"", time.Now().Format(time.RFC3339), ip, path, proxiedTo)
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
