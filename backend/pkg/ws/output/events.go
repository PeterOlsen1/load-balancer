package output

type BaseEvent struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

type Request struct {
	BaseEvent
	IP        string `json:"ip"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	UserAgent string `json:"user_agent"`
}

type Proxy struct {
	BaseEvent
	IP        string `json:"ip"`
	Path      string `json:"path"`
	ProxiedTo string `json:"address"`
}

type Health struct {
	BaseEvent
	Status       string  `json:"status"`
	Address      string  `json:"address"`
	ResponseTime float32 `json:"response_time"`
}

type ContainerStart struct {
	BaseEvent
	ContainerID string `json:"container_id"`
}

type ContainerStop struct {
	BaseEvent
	ContainerID string `json:"container_id"`
}

type Info struct {
	BaseEvent
	Message string `json:"message"`
}

type Error struct {
	BaseEvent
	Message string `json:"message"`
	Error   error  `json:"error"`
}
