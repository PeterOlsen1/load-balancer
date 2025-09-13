package node

import (
	"fmt"
	"io"
	"load-balancer/pkg/errors"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/types"
	"load-balancer/pkg/ws"
	"maps"
	"net/http"
	"time"
)

func (node *Node) processRequest(conn *types.Connection) {
	logger.Proxy(conn.Request.URL.Path, node.Address, conn.Request.RemoteAddr)
	ws.EventEmitter.Proxy(conn.Request.URL.Path, node.Address, conn.Request.RemoteAddr)

	backendURL := fmt.Sprintf("%s%s", node.Address, conn.Request.URL.Path)
	req, err := http.NewRequest(conn.Request.Method, backendURL, conn.Request.Body)
	if err != nil {
		logger.Err("Request creation failed", err)
		ws.EventEmitter.Error("Request creation failed", err)
		errors.Send500(conn, "Creating request to backend")
		return
	}

	maps.Copy(req.Header, conn.Request.Header)
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Err("Backend request failed", err)
		ws.EventEmitter.Error("Backend request failed", err)
		errors.Send500(conn, "Sending backend request")
		return
	}
	defer resp.Body.Close()

	conn.Response.WriteHeader(resp.StatusCode)
	_, err = io.Copy(conn.Response, resp.Body)
	if err != nil {
		logger.Err("Copying response", err)
		ws.EventEmitter.Error("Copying response", err)
		errors.Send500(conn, "Copying backend response")
		return
	}

	conn.Done <- true
}

// Send a request to the node backend to check the health
//
// If an OK status is returned, set node to healthy. Else, unhealthy
func (node *Node) CheckHealth() (string, error) {
	node.Metrics.mu.Lock()
	isPaused := node.Metrics.Health == "paused"
	node.Metrics.mu.Unlock()

	if isPaused {
		return "paused", nil
	}

	address := node.Address

	start := time.Now()
	resp, err := httpClient.Get(fmt.Sprintf("%s/health", address))
	duration := time.Since(start)

	if err != nil {
		logger.Err("Fetching node health", err)
		ws.EventEmitter.Error("Fetching node health", err)
		return "unhealthy", err
	}

	node.Metrics.mu.Lock()
	defer node.Metrics.mu.Unlock()

	respTime := float32(duration.Microseconds() / 1000)
	node.Metrics.ResponseTime = respTime

	health := "healthy"
	if resp.StatusCode != http.StatusOK {
		health = "unhealthy"
		go node.CloseQueue()
	} else {
		if !node.Queue.Open {
			go node.OpenQueue()
		}
	}

	// logger.Health(health, node.Address, respTime)
	ws.EventEmitter.Health(health, node.Address, respTime)

	node.Metrics.Health = health

	return health, nil
}

func (node *Node) Pause() {
	node.Metrics.mu.Lock()
	node.Metrics.Health = "paused"
	node.Metrics.mu.Unlock()

	node.CloseQueue()
}

func (node *Node) Unpause() {
	node.Metrics.mu.Lock()
	node.Metrics.Health = "unknown"
	node.Metrics.mu.Unlock()

	node.CheckHealth()
}

func (n *Node) Equals(other *Node) bool {
	return n.Address == other.Address && n.ContainerID == other.ContainerID
}
