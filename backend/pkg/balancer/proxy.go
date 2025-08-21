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
		send500(conn, "Failed to find route match")
		logger.Err("Failed to find route match", fmt.Errorf("finding route match"))
		return
	}

	node := routeObject.GetProxyNode(conn.Request.RemoteAddr)
	if node == nil {
		logger.Err("Failed to find node for proxy", fmt.Errorf("failed to find node for proxy"))
		send500(conn, "Failed to find node for proxy")
		return
	}

	node.Metrics.Lock.Lock()
	node.Metrics.Connections++

	fmt.Println("connections:", node.Metrics.Connections)

	// add new node if we are above x connections
	// if we have one connection (slow) and more than one node, remove it
	// ^ could be improved upon,
	if node.Metrics.Connections > 30 {
		go func() {
			node, err := StartServer(routeObject.Docker)
			if err != nil {
				send500(conn, "Failed starting server on connection threshhold")
				return
			}
			routeObject.AddNode(node)
		}()
	}
	//  else if node.Metrics.Connections == 1 && len(routeObject.Nodes) > 1 {
	// 	// close node once the proxy is done, this feels risky. re-evaluate how we want this to work
	// 	defer func() {
	// 		routeObject.lock.Lock()
	// 		routeObject.RemoveNode(node)
	// 		routeObject.lock.Unlock()
	// 		node.StopServer()
	// 	}()
	// }
	node.Metrics.Lock.Unlock()

	defer func() {
		node.Metrics.Lock.Lock()
		node.Metrics.Connections--
		node.Metrics.Lock.Unlock()
	}()

	logger.Proxy(conn.Request.URL.Path, node.Address, conn.Request.RemoteAddr)
	ws.EventEmitter.Proxy(conn.Request.URL.Path, node.Address, conn.Request.RemoteAddr)

	backendURL := fmt.Sprintf("%s%s", node.Address, conn.Request.URL.Path)
	req, err := http.NewRequest(conn.Request.Method, backendURL, conn.Request.Body)
	if err != nil {
		logger.Err("Request creation failed", err)
		ws.EventEmitter.Error("Request creation failed", err)
		send500(conn, "Creating request to backend")
		return
	}

	maps.Copy(req.Header, conn.Request.Header)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Err("Backend request failed", err)
		ws.EventEmitter.Error("Backend request failed", err)
		send500(conn, "Sending backend request")
		return
	}
	defer resp.Body.Close()

	conn.Response.WriteHeader(resp.StatusCode)
	_, err = io.Copy(conn.Response, resp.Body)
	if err != nil {
		logger.Err("Copying response", err)
		ws.EventEmitter.Error("Copying response", err)
		send500(conn, "Copying backend response")
		return
	}

}

func (b *BalancerType) getRouteObject(conn *types.Connection) *Route {
	for _, route := range b.Routes {
		matched, err := path.Match(route.Path, conn.Request.URL.Path)
		if err != nil {
			logger.Err("Route matching failed", err)
			continue
		}

		if matched {
			return route
		}
	}

	return nil
}
