package input

type BaseEvent struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

type NodeEvent struct {
	BaseEvent,
	Address string `json:"address"`
}

type NodeStop = NodeEvent
type NodeStart = NodeEvent
type NodePause = NodeEvent
type NodeUnpause = NodeEvent

type StopServer struct {
	BaseEvent
}

type StartServer struct {
	BaseEvent
}

type RequestNodes struct {
	BaseEvent
}
