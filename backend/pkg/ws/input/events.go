package input

type BaseEvent struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

type ContainerStop struct {
	BaseEvent,
	ContainerID string `json:"container_id"`
}

type ContainerStart struct {
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