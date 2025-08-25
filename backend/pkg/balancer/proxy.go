package balancer

import (
	"fmt"
	"load-balancer/pkg/balancer/route"
	"load-balancer/pkg/errors"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/types"
	"path"
	"time"
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

func (b *BalancerType) ProxyRequest(conn *types.Connection) {
	routeObject := b.getRouteObject(conn)
	if routeObject == nil {
		errors.Send500(conn, "Failed to find route match")
		logger.Err("Failed to find route match", fmt.Errorf("finding route match"))
		return
	}

	node := routeObject.GetProxyNode(conn.Request.RemoteAddr)
	if node == nil {
		logger.Err("Failed to find node for proxy", fmt.Errorf("failed to find node for proxy"))
		errors.Send500(conn, "Failed to find node for proxy")
		return
	}

	//add to queue here
	err := node.Queue.Enqueue(conn)
	if err != nil {
		errors.Send500(conn, "Failed to add connection to node queue")
		return
	}

	node.Metrics.Lock.Lock()
	node.Metrics.Connections++

	// add new node if we are above x connections
	// if we have one connection (slow) and more than one node, remove it
	// ^ could be improved upon,
	if !node.Metrics.CreatedNewNode && len(node.Queue.Queue) > routeObject.Docker.RequestScaleThreshold {
		node.Metrics.CreatedNewNode = true
		go func() {
			node, err := routeObject.Scale()
			if err != nil {
				errors.Send500(conn, "Failed starting server on connection threshhold")
				return
			}

			b.NodeTable[node.ContainerID] = node
		}()
	}
	node.Metrics.Lock.Unlock()

	defer func() {
		node.Metrics.Lock.Lock()
		node.Metrics.Connections--

		//if we are below 70% of connection threshold, it is okay to make a new node
		if len(node.Queue.Queue) < int(float64(routeObject.Docker.RequestScaleThreshold)*0.7) {
			node.Metrics.CreatedNewNode = false
		}

		node.Metrics.LastRequestTime = time.Now()
		node.Metrics.Lock.Unlock()
	}()
}
