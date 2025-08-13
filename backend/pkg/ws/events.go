package ws

type BaseEvent struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

type RequestEvent struct {
	BaseEvent
	Method    string `json:"method"`
	Path      string `json:"path"`
	UserAgent string `json:"user_agent"`
}

type ProxyEvent struct {
	BaseEvent
	Path      string `json:"path"`
	ProxiedTo string `json:"address"`
}

type HealthEvent struct {
	BaseEvent
	Status  string `json:"status"`
	Address string `json:"address"`
}

type ContainerStartEvent struct {
	BaseEvent
	ContainerID string `json:"container_id"`
}

type ContainerStopEvent struct {
	BaseEvent
	ContainerID string `json:"container_id"`
}

type InfoEvent struct {
	BaseEvent
	Message string `json:"message"`
}

type ErrorEvent struct {
	BaseEvent
	Message string `json:"message"`
	Error   error  `json:"error"`
}
