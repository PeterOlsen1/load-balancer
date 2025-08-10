package balancer

import (
	"fmt"
	"io"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/types"
	"net/http"
	"net/url"
)

func (b *Balancer) ProxyRequest(conn *types.Connection) {
	node := b.RoundRobin()

	node.Metrics.Lock.Lock()
	node.Metrics.Connections++
	node.Metrics.Lock.Unlock()

	go logger.LogProxy(conn.Request.URL.Path, node.Address)
	go node.CheckHealth()

	backendURL := fmt.Sprintf("%s%s", node.Address, conn.Request.URL.Path)

	//url of the service here
	proxyURL, err := url.Parse("http://localhost:3000")
	if err != nil {
		logger.LogErr("Parsing proxy URL", err)
		return
	}

	tr := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{
		Transport: tr,
	}

	req, err := http.NewRequest(conn.Request.Method, backendURL, conn.Request.Body)
	if err != nil {
		logger.LogErr("Creating proxy request", err)
		return
	}

	//could be faster to copy headers if we just change pointers, possible implications
	for key, values := range conn.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.LogErr("Forwarding request to backend", err)
		return
	}
	defer resp.Body.Close()

	conn.Response.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	for key, values := range resp.Header {
		for _, value := range values {
			conn.Response.Header().Add(key, value)
		}
	}

	// Copy the response body to the client
	_, err = io.Copy(conn.Response, resp.Body)
	if err != nil {
		logger.LogErr("Copying response", err)
	}

	node.Metrics.Lock.Lock()
	node.Metrics.Connections--
	node.Metrics.Lock.Unlock()
}
