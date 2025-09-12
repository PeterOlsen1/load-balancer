package port

import "sync"

var port uint16 = 3000
var portMutex sync.Mutex

func ConsumePort() uint16 {
	var ret uint16
	portMutex.Lock()
	ret = port
	port++
	portMutex.Unlock()
	return ret
}

func ConsumeMultiplePorts(numPorts uint16) []uint16 {
	var ret []uint16
	portMutex.Lock()
	for range numPorts {
		ret = append(ret, port)
		port++
	}
	portMutex.Unlock()
	return ret
}
