package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var Config ConfigType

// LoadConfig function to read YAML file and populate Config
func LoadConfig(configPath string) error {
	f, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Error reading config:", err)
		setDefaultConfig()
		return err
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&Config)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		setDefaultConfig()
		return err
	}

	nameMap := make(map[string]bool)
	for _, route := range Config.Routes {
		if nameMap[route.Name] {
			return fmt.Errorf("duplicate route name found: %s", route.Name)
		}
		nameMap[route.Name] = true
	}

	return nil
}

func setDefaultConfig() {
	Config = ConfigType{
		Server: ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
		Logging: LoggingConfig{
			Level:  "all",
			Folder: "./logs",
		},
		Routes: []RouteConfig{
			{
				Path:          "/*",
				Name:          "allServer",
				Strategy:      "round-robin",
				MaxNodes:      0,
				HealthTimeout: 5000,
				Docker: &DockerConfig{
					Image:                 "node-server",
					InternalPort:          3000,
					RequestScaleThreshold: 15,
				},
				Servers: []RouteServerConfig{},
			},
		},
	}
}
