package input

type BaseEvent struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

type StopContainer struct {
	BaseEvent,
	ContainerID string `json:"container_id"`
}

type StartContainer struct {
	BaseEvent,
	Address string `json:"address"`
}

type StopServer struct {
	BaseEvent
}

type StartServer struct {
	BaseEvent
}

type RequestNodes struct {
	BaseEvent
}
