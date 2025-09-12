package config

type ConfigType struct {
	// Server configuration info. See ServerConfig type for more
	Server ServerConfig `yaml:"server"`

	//Logger configuration info. See LoggingConfig type for more
	Logging LoggingConfig `yaml:"logging"`

	//WS emitter configuration info. See EmitterConfig type for more
	Emitter EmitterConfig `yaml:"emitter"`

	//Route configuration info. See RouteConfig for more
	Routes []RouteConfig `yaml:"routes"`
}

type EmitterConfig struct {
	// A flag to signal whether or not the event emitter is enabled.
	// If true, events will be sent out on a ws connection to {Path}
	Enabled bool `yaml:"enabled"`

	// The url path where we will initialize a websocket connection
	// "wss://{host}:{port}{Path}"
	//
	// The leading `/` is not included in the path, you must include it
	Path string `yaml:"path"`
}
type ServerConfig struct {
	// The port on which the server will run. Default is 8080
	Port int `yaml:"port"`

	// The host on which the server will run. Default is localhost
	Host string `yaml:"host"`
}

type LoggingConfig struct {
	// The desired level of logging. Options:
	//  - All: all logging statements are written, including REQUEST, PROXY, and INFO
	//	- Error: only error statements are written
	//	- None: no logs are written
	Level uint `yaml:"level"`

	// Path to the folder where logs will be stored
	Folder string `yaml:"folder"`
}

type DockerConfig struct {
	// The name of the docker image which your server will run on.
	// Should be pre-built before running the load balancer
	Image string `yaml:"image"`

	// The port which is exposed within your image.
	// For example,
	//	import { express } from "express";
	//
	//	const app = express();
	//	app.listen(3000);
	// The `internal_port` variable would be 3000
	InternalPort int `yaml:"internal_port"`
}

type PoolConfig struct {
	// The number of containers to keep warm for spikes in requests.
	// These containers do not recieve reqeusts until moved to active
	InactiveSize int `yaml:"inactive_size"`

	// The minimum number of containers to keep active.
	// More containers may be pulled from the inactive pool if necessary
	//
	// These containers recieve reqeusts, but are moved to inactive if a
	// reqeust to /health fails
	ActiveSize int `yaml:"active_size"`

	// The maximum number of active containers we can have
	//
	// This means that the total number of max containers we can have is
	// max_active + inactive_size
	MaxActive int `yaml:"max_active"`

	// The amount of time between each node activation.
	// A node activation refers to moving a node from "inactive" to "active"
	// while the balancer is under load.
	//
	// If this is set too low, too many containers will be created under load
	ActivationInterval int `yaml:"activation_interval_ms"`

	// The amount of time between node cleanup checks.
	// Every n ms, a goroutine will check the load.
	// If load < 10%, a node will be paused
	CleanupInterval int `yaml:"cleanup_interval_ms"`
}

type RouteServerConfig struct {
	// URL of a pre-running server that will be able to handle requests.
	//
	// It must have a /health route or else the load balancer will fail
	URL string `yaml:"url"`

	// Weight of the server. TODO: implement server weights?
	Weight int `yaml:"weight"`
}

type RouteConfig struct {
	// The base route of the path. For all requests, pass /*
	//
	// For a route that only deals with /api/..., pass /api/*
	Path string `yaml:"path"`

	// The name of the route, used as an identifier, primarily for the frontend
	Name string `yaml:"name"`

	// The load balancing strategy of the route. Options are:
	//	- Round robin
	//	- IP hash
	//	- Least connections
	Strategy string `yaml:"strategy"`

	// The number of ms in between every check to /heath
	HealthTimeout int `yaml:"health_timeout_ms"`

	// The max number of requests that a node can have in its queue at a time.
	//
	// This number is also used to calculate the load % of a given route
	NodeQueueSize int `yaml:"node_queue_size"`

	// The max size of the route queue where connections go before being sent to nodes.
	// It is recommended to keep this high since requests dropped from the route queue will send a 500
	RouteQueueSize uint `yaml:"route_queue_size"`

	// The number of worker threads to be wathing a queue at a time
	WorkerThreads uint `yaml:"worker_threads"`
	
	// Docker configuration, see DockerConfig type for more info
	Docker *DockerConfig `yaml:"docker"`

	// Node pool configuration, see PoolConfig type for more info
	Pool PoolConfig `yaml:"pool"`

	// List of pre-running servers to proxy requests to. See RouteServerConfig type for more info
	Servers []RouteServerConfig `yaml:"servers"`

}
