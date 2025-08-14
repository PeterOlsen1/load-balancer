package balancer

import (
	"encoding/json"
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/ws"
	"time"
)

type RequestNodesResponse struct {
	Type  string      `json:"type"`
	Time  string      `json:"time"`
	Nodes []node.Node `json:"nodes"`
}

func init() {
	ws.EventReciever.AddEventHandler("request_nodes", func(body []byte) ([]byte, error) {
		var nodeLiterals []node.Node
		for _, n := range LoadBalancer.nodes {
			nodeLiterals = append(nodeLiterals, *n)
		}

		resp := RequestNodesResponse{
			Type:  "nodes",
			Time:  time.Now().Format(time.RFC3339),
			Nodes: nodeLiterals,
		}

		j, err := json.Marshal(resp)
		if err != nil {
			logger.Err("Marshalling node JSON", err)
			return nil, err
		}

		return j, nil
	})

	// ws.EventReciever.AddEventHandler("container_stop", func(body []byte) ([]byte, error) {
	// 	return nil, nil
	// })
}
