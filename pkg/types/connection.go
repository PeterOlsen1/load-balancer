package types

import (
	"bytes"
	"io"
	"net/http"
	"sync"
)

type Connection struct {
	Writer  http.ResponseWriter
	Request *http.Request
	lock    sync.Mutex
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
