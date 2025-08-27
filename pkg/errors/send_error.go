package errors

import (
	"fmt"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/types"
	"load-balancer/pkg/ws"
)

func Send500(conn *types.Connection, reason string) {
	message := fmt.Sprintf("500 Internal Server Error: %s", reason)
	conn.Response.Header().Set("Content-Type", "text/plain")
	conn.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(message)))
	conn.Response.WriteHeader(500)

	// Write the response body once
	_, err := conn.Response.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing 500 response:", err)

		logger.Err("Writing 500 response", err)
		ws.EventEmitter.Error("Writing 500 response", err)
	}

	conn.Done <- true
}

func Send400(conn *types.Connection, reason string) {
	message := fmt.Sprintf("400 Bad request: %s", reason)
	conn.Response.Header().Set("Content-Type", "text/plain")
	conn.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(message)))
	conn.Response.WriteHeader(400)

	// Write the response body once
	_, err := conn.Response.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing 400 response:", err)

		logger.Err("Writing 400 response", err)
		ws.EventEmitter.Error("Writing 400 response", err)
	}

	conn.Done <- true
}
