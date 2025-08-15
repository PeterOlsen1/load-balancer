package receiver

import (
	"encoding/json"
	"fmt"
	b "load-balancer/pkg/balancer"
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/ws"
	"load-balancer/pkg/ws/input"
	"time"
)

func getBaseResponse(respType string) BaseResponse {
	return BaseResponse{
		Type:     respType,
		Time:     time.Now().Format(time.RFC3339),
		Response: true,
	}
}

func init() {
	ws.EventReciever.AddEventHandler("request_nodes", func(body []byte) ([]byte, error) {
		var nodeLiterals []node.Node
		for _, n := range b.LoadBalancer.Nodes {
			nodeLiterals = append(nodeLiterals, *n)
		}

		resp := RequestNodesResponse{
			BaseResponse: getBaseResponse("nodes"),
			Nodes:        nodeLiterals,
		}

		j, err := json.Marshal(resp)
		if err != nil {
			logger.Err("Marshalling node JSON", err)
			return nil, err
		}

		return j, nil
	})

	ws.EventReciever.AddEventHandler("node_stop", func(body []byte) ([]byte, error) {
		userRequest := input.NodeStop{}
		err := json.Unmarshal(body, &userRequest)
		if err != nil {
			logger.Err("Unmarshalling node_stop JSON", err)
			return nil, err
		}

		var address *string = nil
		b.LoadBalancer.Lock()
		defer b.LoadBalancer.Unlock()
		for _, n := range b.LoadBalancer.Nodes {
			if n.Address == userRequest.Address {
				address = &n.Address
				err := n.StopServer()
				if err != nil {
					return nil, err
				}
			}
		}

		resp := NodeStopResponse{
			BaseResponse: getBaseResponse("node_stop"),
			Message:      fmt.Sprintf("Could not locate node @ %s", userRequest.Address),
		}

		if address != nil {
			resp.Message = fmt.Sprintf("Successfully stopped node @ %s", address)
		}

		j, err := json.Marshal(resp)
		if err != nil {
			logger.Err("Marshalling node stop response", err)
		}

		return j, nil
	})

	ws.EventReciever.AddEventHandler("node_start", func(body []byte) ([]byte, error) {
		newNode, err := b.StartServer(b.PORT)
		if err != nil {
			return nil, err
		}
		b.PORT++
		b.LoadBalancer.AddNode(newNode)

		resp := NodeStartResponse{
			BaseResponse: getBaseResponse("node_start"),
			Message:      "Successfully started new node",
			ContainerID:  newNode.DockerInfo.Id,
			Address:      newNode.Address,
		}

		j, err := json.Marshal(resp)
		if err != nil {
			logger.Err("Marshalling node_start JSON", err)
			return nil, err
		}

		return j, nil
	})

	ws.EventReciever.AddEventHandler("node_pause", func(body []byte) ([]byte, error) {
		userRequest := input.NodePause{}
		if err := json.Unmarshal(body, &userRequest); err != nil {
			logger.Err("Failed to unmarshal node_pause data", err)
			return nil, err
		}

		var node *node.Node
		for _, n := range b.LoadBalancer.Nodes {
			if n.Address == userRequest.Address {
				node = n
				break
			}
		}

		if node == nil {
			logger.Err("Node not found", nil)
			return nil, fmt.Errorf("node not found")
		}

		node.Pause()
		logger.ContainerPause(node.DockerInfo.Id)

		resp := NodePauseResponse{
			BaseResponse: getBaseResponse("node_pause"),
			Address:      node.Address,
		}
		j, err := json.Marshal(resp)
		if err != nil {
			logger.Err("Marshalling node_pause JSON", err)
			return nil, err
		}
		return j, nil
	})

	ws.EventReciever.AddEventHandler("node_unpause", func(body []byte) ([]byte, error) {
		var nodeAddress string
		if err := json.Unmarshal(body, &nodeAddress); err != nil {
			logger.Err("Failed to unmarshal node address", err)
			return nil, err
		}

		var node *node.Node
		for _, n := range b.LoadBalancer.Nodes {
			if n.Address == nodeAddress {
				node = n
				break
			}
		}

		if node == nil {
			logger.Err("Node not found", nil)
			return nil, fmt.Errorf("node not found")
		}

		node.Unpause()
		logger.ContainerUnpause(node.DockerInfo.Id)

		resp := NodeUnpauseResponse{
			BaseResponse: getBaseResponse("unnode_pause"),
			Address:      node.Address,
		}
		j, err := json.Marshal(resp)
		if err != nil {
			logger.Err("Marshalling unnode_pause JSON", err)
			return nil, err
		}
		return j, nil
	})
}
