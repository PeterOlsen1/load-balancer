package input

type BaseEvent struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

type NodeEvent struct {
	BaseEvent,
	ContainerID string `json:"container_id"`
}

type NodeStop = NodeEvent
type NodeStart struct {
	NodeEvent
	RouteName string `json:route_name`
}
type NodePause = NodeEvent
type NodeUnpause = NodeEvent

type StopServer struct {
	BaseEvent
}

type RequestNodes struct {
	BaseEvent
}
