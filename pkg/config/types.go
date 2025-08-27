package config

type ConfigType struct {
	Server  ServerConfig  `yaml:"server"`
	Logging LoggingConfig `yaml:"logging"`
	Routes  []RouteConfig `yaml:"routes"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Folder string `yaml:"folder"`
}

// Image: the name of the given docker image to scale
// InternalPort: the port on which the server runs
// RequestScaleThreshold: the number of requests at a time necesasry to start a new container
// NoRequestsTimeout: the number of ms to wait for no requests before
type DockerConfig struct {
	Image                 string `yaml:"image"`
	InternalPort          int    `yaml:"internal_port"`
	RequestScaleThreshold int    `yaml:"request_scale_threshold"`
	NoRequestsTimeout     int    `yaml:"no_requests_timeout_ms"`
	InitialContainers     int    `yaml:"initial_containers"`
}

type RouteServerConfig struct {
	URL string `yaml:"url"`
}

type RouteConfig struct {
	Path            string              `yaml:"path"`
	Name            string              `yaml:"name"`
	Strategy        string              `yaml:"strategy"`
	HealthTimeout   int                 `yaml:"health_timeout_ms"`
	InactiveTimeout int                 `yaml:"inactive_timeout_ms"`
	RequestLimit    int                 `yaml:"node_request_limit"`
	Docker          *DockerConfig       `yaml:"docker"`
	Servers         []RouteServerConfig `yaml:"servers"`
}
