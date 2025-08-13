package ws

type InputStopContainerEvent struct {
	BaseEvent,
	ContainerID string `json:"container_id"`
}

type InputStartContainerEvent struct {
	BaseEvent,
	Address string `json:"address"`
}

type InputStopServerEvent struct {
	BaseEvent
}

type InputStartServerEvent struct {
	BaseEvent
}

type InputRequestNodesEvent struct {
	BaseEvent
}