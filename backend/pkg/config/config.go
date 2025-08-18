package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var Config = ConfigType{
	Server: ServerConfig{
		Port: 8080,
		Host: "localhost",
	},
	LoadBalancer: LoadBalancerConfig{
		HealthInterval: 5000,
		MaxNodes:       10,
		Strategy:       "round-robin",
	},
	Docker: DockerConfig{
		DockerImage:  "node-server",
		InternalPort: 3000,
	},
	Logging: LoggingConfig{
		Level:  "all",
		Folder: "./logs",
	},
}

func LoadConfig(configPath string) error {
	f, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Error reading config:", err)
		return err
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&Config)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return err
	}

	return nil
}
