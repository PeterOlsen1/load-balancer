package types

import (
	"bytes"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Response   http.ResponseWriter
	Request    *http.Request
	Done       chan bool
	RetryCount uint8
}

type LockedConnection struct {
	Conn *websocket.Conn
	mu   sync.Mutex
}

// debugging purposes, just prints the body
func (conn *Connection) DebugBody() (string, error) {
	bodyText, err := io.ReadAll(conn.Request.Body)
	if err != nil {
		return "", err
	}
	conn.Request.Body = io.NopCloser(bytes.NewBuffer(bodyText))
	return string(bodyText), err
}

func (c *LockedConnection) WriteMessage(messageType int, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Conn.WriteMessage(messageType, data)
}
