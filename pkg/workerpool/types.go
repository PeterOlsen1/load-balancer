package workerpool

type WorkerPool[T any] struct {
	// The number of workers to be initialized in the pool.
	//
	// Each of these numWorkers corresponds to 1 goroutine,
	// so keep that in mind when choosing what to put for this value
	numWorkers uint16

	// The channel on which events are sent
	//
	// Capacity of this chan is static, 1000.
	// Keep this number in mind when deciding how many worker threads to spawn
	eventChan chan T

	// The function to run when events are dequeued from the eventChan
	//
	// Should log its own errors
	eventHandler func(T)
}
