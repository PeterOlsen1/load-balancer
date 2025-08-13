package ws

type BaseEvent struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

type RequestEvent struct {
	BaseEvent
	IP        string `json:"ip"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	UserAgent string `json:"user_agent"`
}

type ProxyEvent struct {
	BaseEvent
	IP        string `json:"ip"`
	Path      string `json:"path"`
	ProxiedTo string `json:"address"`
}

type HealthEvent struct {
	BaseEvent
	Status       string  `json:"status"`
	Address      string  `json:"address"`
	ResponseTime float32 `json:"response_time"`
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
