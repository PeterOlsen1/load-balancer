package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var Config ConfigType

var defaultConfig = ConfigType{
	Server: ServerConfig{
		Port: 8080,
		Host: "localhost",
	},
	Logging: LoggingConfig{
		Level:    4,
		Folder:   "./logs",
		MaxLines: 50000,
	},
	Emitter: EmitterConfig{
		Enabled: true,
		Path:    "/events",
	},
	Routes: []RouteConfig{
		{
			Path:           "/*",
			Name:           "allServer",
			Strategy:       "round-robin",
			HealthTimeout:  5000,
			RouteQueueSize: 50,
			NodeQueueSize:  1000,
			WorkerThreads:  10,
			Docker: &DockerConfig{
				Image:        "node-server",
				InternalPort: 3000,
			},
			Pool: PoolConfig{
				InactiveSize:       3,
				ActiveSize:         2,
				MaxActive:          10,
				ActivationInterval: 500,
				CleanupInterval:    5000,
			},
			Servers: []RouteServerConfig{},
		},
	},
}

// LoadConfig function to read YAML file and populate Config
func LoadConfig(configPath string) error {
	f, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Error reading config:", err)
		return err
	}
	defer f.Close()

	// set default config before
	Config = defaultConfig
	err = yaml.NewDecoder(f).Decode(&Config)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return err
	}

	nameMap := make(map[string]bool)
	for _, route := range Config.Routes {
		if nameMap[route.Name] {
			return fmt.Errorf("duplicate route name found: %s", route.Name)
		}
		nameMap[route.Name] = true
	}

	if Config.Logging.Level > 4 {
		Config.Logging.Level = 4
	}

	return nil
}
