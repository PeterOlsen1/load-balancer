package types

import (
	"bytes"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Response http.ResponseWriter
	Request  *http.Request
	lock     sync.Mutex
}

type LockedConnection struct {
	Conn *websocket.Conn
	Lock sync.Mutex
}

func (conn *Connection) Body() (string, error) {
	defer conn.lock.Unlock()
	conn.lock.Lock()

	bodyText, err := io.ReadAll(conn.Request.Body)
	if err != nil {
		return "", err
	}
	conn.Request.Body = io.NopCloser(bytes.NewBuffer(bodyText))
	return string(bodyText), err
}
