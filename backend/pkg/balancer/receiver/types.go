package receiver

import "load-balancer/pkg/balancer/node"

type BaseResponse struct {
	Type     string `json:"type"`
	Time     string `json:"time"`
	Response bool   `json:"response"`
}

type RequestNodesResponse struct {
	BaseResponse
	Nodes []node.Node `json:"nodes"`
}

type NodeStopResponse struct {
	BaseResponse
	Message string `json:"message"`
}

type NodeStartResponse struct {
	BaseResponse
	Address     string `json:"address"`
	ContainerID string `json:"container_id"`
	Message     string `json:"message"`
}

type NodePauseResponse = NodeStartResponse

type NodeUnpauseResponse = NodeStartResponse
