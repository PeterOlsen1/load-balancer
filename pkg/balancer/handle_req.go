package balancer

import (
	"load-balancer/pkg/balancer/route"
	"load-balancer/pkg/errors"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/types"
	"path"
)

func (b *BalancerType) getRouteObject(conn *types.Connection) *route.Route {
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

// decide what route the request is going to, send it to that queue
func (b *BalancerType) HandleRequest(conn *types.Connection) {
	routeObject := b.getRouteObject(conn)
	if routeObject == nil {
		errors.Send400(conn, "Request does not match any defined routes")
		return
	}

	proxyNode := routeObject.GetProxyNode(conn.Request.RemoteAddr)
	if proxyNode != nil {
		proxyNode.Queue.Enqueue(conn)
	} else {
		routeObject.Queue.Enqueue(conn)
	}

}
