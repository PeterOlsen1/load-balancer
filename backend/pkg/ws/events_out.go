package ws

type BaseEvent struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

type OutputRequestEvent struct {
	BaseEvent
	IP        string `json:"ip"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	UserAgent string `json:"user_agent"`
}

type OutputProxyEvent struct {
	BaseEvent
	IP        string `json:"ip"`
	Path      string `json:"path"`
	ProxiedTo string `json:"address"`
}

type OutputHealthEvent struct {
	BaseEvent
	Status       string  `json:"status"`
	Address      string  `json:"address"`
	ResponseTime float32 `json:"response_time"`
}

type OutputContainerStartEvent struct {
	BaseEvent
	ContainerID string `json:"container_id"`
}

type OutputContainerStopEvent struct {
	BaseEvent
	ContainerID string `json:"container_id"`
}

type OutputInfoEvent struct {
	BaseEvent
	Message string `json:"message"`
}

type OutputErrorEvent struct {
	BaseEvent
	Message string `json:"message"`
	Error   error  `json:"error"`
}
