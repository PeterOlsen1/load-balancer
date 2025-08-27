package balancer

import (
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/balancer/route"
	"load-balancer/pkg/config"
)

var Balancer = BalancerType{
	NodeTable: make(map[string]*node.Node),
}

func (b *BalancerType) InitBalancer() error {
	for _, r := range config.Config.Routes {
		route, err := route.InitRoute(r)
		if err != nil {
			return err
		}

		b.Routes = append(b.Routes, route)
	}

	return nil
}
