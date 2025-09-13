package workerpool

func InitWorkerPool[T any](numWorkers uint16, eventHandler func(e T)) WorkerPool[T] {
	pool := WorkerPool[T]{
		numWorkers:   numWorkers,
		eventChan:    make(chan T, 1000),
		eventHandler: eventHandler,
	}

	// initialize worker threads
	for range numWorkers {
		go func() {
			for e := range pool.eventChan {
				pool.eventHandler(e)
			}
		}()
	}

	return pool
}
