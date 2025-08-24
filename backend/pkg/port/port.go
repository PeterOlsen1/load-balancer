package port

import "sync"

var port int = 3000
var portMutex sync.Mutex

func ConsumePort() int {
	var ret int
	portMutex.Lock()
	ret = port
	port++
	portMutex.Unlock()
	return ret
}

func ConsumeMultiplePorts(numPorts int) []int {
	var ret []int
	portMutex.Lock()
	for range numPorts {
		ret = append(ret, port)
		port++
	}
	portMutex.Unlock()
	return ret
}
