package balancer

import (
	"fmt"
	"io"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/types"
	"load-balancer/pkg/ws"
	"maps"
	"net/http"
	"path"
)

func (b *BalancerType) ProxyRequest(conn *types.Connection) {
	routeObject := b.getRouteObject(conn)
	if routeObject == nil {
		send500(conn)
		return
	}

	node := routeObject.RoundRobin()
	if node == nil {
		logger.Err("Failed to find node for proxy", fmt.Errorf("failed to find node for proxy"))
		send500(conn)
		return
	}

	node.Metrics.Lock.Lock()
	node.Metrics.Connections++

	if node.Metrics.Connections > 30 {
		node, err := StartServer(PORT, routeObject.Docker)
		if err != nil {
			send500(conn)
			return
		}
		routeObject.AddNode(node)
	} else if node.Metrics.Connections == 1 && len(routeObject.Nodes) > 1 {
		routeObject.lock.Lock()
		routeObject.RemoveNode(node)
		routeObject.lock.Unlock()
		node.StopServer()
	}
	node.Metrics.Lock.Unlock()

	defer func() {
		node.Metrics.Lock.Lock()
		node.Metrics.Connections--
		node.Metrics.Lock.Unlock()
	}()

	go logger.Proxy(conn.Request.URL.Path, node.Address, conn.Request.RemoteAddr)
	go ws.EventEmitter.Proxy(conn.Request.URL.Path, node.Address, conn.Request.RemoteAddr)

	backendURL := fmt.Sprintf("%s%s", node.Address, conn.Request.URL.Path)
	req, err := http.NewRequest(conn.Request.Method, backendURL, conn.Request.Body)
	if err != nil {
		go logger.Err("Request creation failed", err)
		go ws.EventEmitter.Error("Request creation failed", err)
		send500(conn)
		return
	}

	maps.Copy(req.Header, conn.Request.Header)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		go logger.Err("Backend request failed", err)
		go ws.EventEmitter.Error("Backend request failed", err)
		send500(conn)
		return
	}
	defer resp.Body.Close()

	conn.Response.WriteHeader(resp.StatusCode)
	_, err = io.Copy(conn.Response, resp.Body)
	if err != nil {
		go logger.Err("Copying response", err)
		go ws.EventEmitter.Error("Copying response", err)
		send500(conn)
		return
	}
}

func (b *BalancerType) getRouteObject(conn *types.Connection) *Route {
	for _, route := range b.Routes {
		routePath := route.Path
		if routePath == "/*" {
			routePath = "/"
		}

		matched, err := path.Match(routePath, conn.Request.URL.Path)
		if err != nil {
			go logger.Err("Route matching failed", err)
			continue
		}

		if matched {
			return route
		}
	}

	return nil
}
