package config

type ConfigType struct {
	Server       ServerConfig       `yaml:"server"`
	LoadBalancer LoadBalancerConfig `yaml:"load_balancer"`
	Docker       DockerConfig       `yaml:"docker"`
	Logging      LoggingConfig      `yaml:"logging"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type LoadBalancerConfig struct {
	HealthInterval int    `yaml:"health_interval"`
	MaxNodes       int    `yaml:"max_nodes"`
	Strategy       string `yaml:"strategy"`
}

type DockerConfig struct {
	DockerImage  string `yaml:"docker_image"`
	InternalPort int    `yaml:"internal_port"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Folder string `yaml:"folder"`
}
