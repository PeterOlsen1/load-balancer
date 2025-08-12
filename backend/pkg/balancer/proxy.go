package balancer

import (
	"fmt"
	"io"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/types"
	"maps"
	"net/http"
)

func (b *Balancer) ProxyRequest(conn *types.Connection) {
	node := b.RoundRobin()
	if node == nil {
		return
	}

	node.Metrics.Lock.Lock()
	node.Metrics.Connections++
	node.Metrics.Lock.Unlock()
	defer func() {
		node.Metrics.Lock.Lock()
		node.Metrics.Connections--
		node.Metrics.Lock.Unlock()
	}()

	go logger.LogProxy(conn.Request.URL.Path, node.Address)

	backendURL := fmt.Sprintf("%s%s", node.Address, conn.Request.URL.Path)
	req, err := http.NewRequest(conn.Request.Method, backendURL, conn.Request.Body)
	if err != nil {
		logger.LogErr("Request creation failed", err)
		send500(conn)
		return
	}

	maps.Copy(req.Header, conn.Request.Header)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.LogErr("Backend request failed", err)
		send500(conn)
		return
	}
	defer resp.Body.Close()

	conn.Response.WriteHeader(resp.StatusCode)
	_, err = io.Copy(conn.Response, resp.Body)
	if err != nil {
		logger.LogErr("Copying response", err)
		send500(conn)
		return
	}
}
