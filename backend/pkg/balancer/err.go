package balancer

import (
	"fmt"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/types"
)

func send500(conn *types.Connection) {
	message := "500 Internal Server Error"
	conn.Response.Header().Set("Content-Type", "text/plain")
	conn.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(message)))
	conn.Response.WriteHeader(500)

	// Write the response body once
	_, err := conn.Response.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing 500 response:", err)
		logger.LogErr("Writing 500 response", err)
	}
}
