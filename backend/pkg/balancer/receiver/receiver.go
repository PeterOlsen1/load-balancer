package receiver

import (
	"encoding/json"
	"fmt"
	b "load-balancer/pkg/balancer"
	"load-balancer/pkg/balancer/docker"
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/balancer/route"
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

		for _, n := range b.Balancer.NodeTable {
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
		for _, route := range b.Balancer.Routes {
			for _, n := range route.Nodes {
				if n.ContainerID == userRequest.ContainerID {
					address = &n.Address
					err := docker.StopContainer(n.ContainerID)
					if err != nil {
						return nil, err
					}
					break
				}
			}
		}

		resp := NodeStopResponse{
			BaseResponse: getBaseResponse("node_stop"),
			Message:      fmt.Sprintf("Could not locate node with ID %s", userRequest.ContainerID),
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
		userRequest := input.NodeStart{}
		err := json.Unmarshal(body, &userRequest)
		if err != nil {
			logger.Err("Unmarshalling node_start JSON", err)
			return nil, err
		}

		var routeObject *route.Route = nil
		for _, route := range b.Balancer.Routes {
			if route.Name == userRequest.RouteName {
				routeObject = route
				break
			}
		}

		newNode, err := routeObject.Scale()
		if err != nil {
			return nil, err
		}
		b.Balancer.NodeTable[newNode.ContainerID] = newNode

		resp := NodeStartResponse{
			BaseResponse: getBaseResponse("node_start"),
			Message:      "Successfully started new node",
			ContainerID:  newNode.ContainerID,
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

		node := b.Balancer.NodeTable[userRequest.ContainerID]
		if node == nil {
			logger.Err("Node not found", nil)
			return nil, fmt.Errorf("node not found")
		}

		node.Pause()
		logger.ContainerPause(node.ContainerID)

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
		userRequest := input.NodeUnpause{}
		if err := json.Unmarshal(body, &userRequest); err != nil {
			logger.Err("Failed to unmarshal node address", err)
			return nil, err
		}

		node := b.Balancer.NodeTable[userRequest.ContainerID]
		if node == nil {
			logger.Err("Node not found", nil)
			return nil, fmt.Errorf("node not found")
		}

		node.Unpause()
		logger.ContainerUnpause(node.ContainerID)

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

	ws.EventReciever.AddEventHandler("request_routes", func(body []byte) ([]byte, error) {
		return []byte("TODO"), nil
	})
}
